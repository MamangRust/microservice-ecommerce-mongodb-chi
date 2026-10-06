package outbox

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"go.uber.org/zap"
)

// Outbox constants control the durable retry behavior of the outbox relay
// (Phase 6 — Transactional Outbox). They mirror the ecommerce transaction
// service so behavior is uniform across repositories.
const (
	OutboxMaxAttempts    = 5
	OutboxBackoff        = 30 * time.Second
	OutboxRelayInterval  = 5 * time.Second
	OutboxRelayBatchSize = 100
	// OutboxClaimLease is how long a relay worker owns a claimed event. If the
	// worker dies after claiming but before marking the event delivered, the
	// lease expires and another relay instance re-claims and retries it.
	OutboxClaimLease = 1 * time.Minute
	// OutboxRetention is how long delivered/dead events are kept before the
	// relay purges them as part of the retention policy.
	OutboxRetention = 7 * 24 * time.Hour

	// OutboxRetentionEveryTicks runs the retention purge every N relay ticks so
	// the DELETE scan does not run on every relay cycle.
	OutboxRetentionEveryTicks = 60
)

// OutboxEvent mirrors the generated outbox_events row shape. It is defined here
// (not in a generated schema package) so the relay contract stays
// schema-agnostic: every service adapts its own generated db.Queries to it.
type OutboxEvent struct {
	OutboxID      int64
	Topic         string
	EventKey      string
	Payload       []byte
	Status        string
	Attempts      int32
	NextAttemptAt time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type CreateOutboxEventParams struct {
	Topic    string
	EventKey string
	Payload  []byte
}

type ClaimPendingOutboxEventsParams struct {
	Limit         int32
	NextAttemptAt time.Time
}

type MarkOutboxEventFailedParams struct {
	OutboxID      int64
	NextAttemptAt time.Time
}

// OutboxQuerier is the subset of the per-service generated db.Queries that the
// outbox relay needs. Each service's generated Queries does not satisfy this
// interface structurally (its params/rows are generated types), so services
// pass a small adapter (see service/<name>/apps).
type OutboxQuerier interface {
	CreateOutboxEvent(ctx context.Context, params CreateOutboxEventParams) (OutboxEvent, error)
	ClaimPendingOutboxEvents(ctx context.Context, params ClaimPendingOutboxEventsParams) ([]OutboxEvent, error)
	MarkOutboxEventFailed(ctx context.Context, params MarkOutboxEventFailedParams) (OutboxEvent, error)
	MarkOutboxEventDelivered(ctx context.Context, outboxID int64) (OutboxEvent, error)
	MarkOutboxEventDead(ctx context.Context, outboxID int64) (OutboxEvent, error)
	DeleteOldOutboxEvents(ctx context.Context, cutoff time.Time) (int64, error)
}

// OutboxPublisher is the minimal Kafka producer surface the relay needs.
// *kafka.Kafka satisfies this interface via SendMessage.
type OutboxPublisher interface {
	SendMessage(ctx context.Context, topic string, key string, value []byte) error
}

// OutboxService persists email events durably and relays them to Kafka with
// retry and dead-letter semantics. Producers enqueue the event inside the same
// database transaction as the business write (EnqueueInTx) so a crash between
// the two cannot lose the event; the relay then guarantees delivery.
type OutboxService struct {
	queries   OutboxQuerier
	publisher OutboxPublisher
	logger    logger.LoggerInterface
}

// NewOutboxService builds the outbox service. The queries may be nil (e.g. a
// service whose database has no outbox_events table) — Enqueue and the relay
// then no-op; the publisher may also be nil (local dev without Kafka) — the
// relay then drains the queue without sending.
func NewOutboxService(queries OutboxQuerier, publisher OutboxPublisher, log logger.LoggerInterface) *OutboxService {
	return &OutboxService{queries: queries, publisher: publisher, logger: log}
}

// EnqueueInTx persists a pending event inside the given database transaction so
// the caller can commit the business write and the event atomically. q must be
// a querier bound to the transaction (e.g. schema.New(tx) wrapped in the
// service's adapter). This is the production path: the event survives the
// commit and is published by the relay.
func (s *OutboxService) EnqueueInTx(ctx context.Context, q OutboxQuerier, topic, key string, payload []byte) error {
	_, err := q.CreateOutboxEvent(ctx, CreateOutboxEventParams{
		Topic:    topic,
		EventKey: key,
		Payload:  payload,
	})
	if err != nil {
		return err
	}
	s.logger.Info("outbox event enqueued", zap.String("topic", topic), zap.String("key", key))
	return nil
}

// Enqueue persists a pending event AFTER the business transaction has already
// committed. It is the NON-ATOMIC fallback path: a crash between the commit and
// this insert silently loses the event, so it must not be used where the
// business write is local. It exists for aggregator services whose business
// write happens in another service over gRPC (best-effort guarantee).
func (s *OutboxService) Enqueue(ctx context.Context, topic, key string, payload []byte) error {
	if s.queries == nil {
		return nil
	}
	_, err := s.queries.CreateOutboxEvent(ctx, CreateOutboxEventParams{
		Topic:    topic,
		EventKey: key,
		Payload:  payload,
	})
	if err != nil {
		return err
	}
	s.logger.Info("outbox event enqueued", zap.String("topic", topic), zap.String("key", key))
	return nil
}

// PublishPending claims up to limit pending events whose retry window has
// elapsed, publishes each to Kafka, and marks it delivered. Claiming is atomic
// (FOR UPDATE SKIP LOCKED + lease), so concurrent relay instances never publish
// the same event twice. It returns the number of events successfully delivered.
func (s *OutboxService) PublishPending(ctx context.Context, limit int) (int, error) {
	if s.queries == nil || s.publisher == nil {
		return 0, nil
	}
	events, err := s.queries.ClaimPendingOutboxEvents(ctx, ClaimPendingOutboxEventsParams{
		Limit:         int32(limit),
		NextAttemptAt: time.Now().Add(OutboxClaimLease),
	})
	if err != nil {
		return 0, err
	}

	delivered := 0
	for _, event := range events {
		if err := s.publisher.SendMessage(ctx, event.Topic, event.EventKey, event.Payload); err != nil {
			s.logger.Error("failed to publish outbox event, scheduling retry",
				zap.Error(err),
				zap.Int64("outbox_id", event.OutboxID),
				zap.String("topic", event.Topic),
				zap.Int32("attempts", event.Attempts),
			)
			if int(event.Attempts)+1 >= OutboxMaxAttempts {
				if _, deadErr := s.queries.MarkOutboxEventDead(ctx, event.OutboxID); deadErr != nil {
					s.logger.Error("failed to dead-letter outbox event", zap.Error(deadErr), zap.Int64("outbox_id", event.OutboxID))
				}
				continue
			}
			nextAttempt := time.Now().Add(OutboxBackoff * time.Duration(event.Attempts+1))
			if _, failErr := s.queries.MarkOutboxEventFailed(ctx, MarkOutboxEventFailedParams{
				OutboxID:      event.OutboxID,
				NextAttemptAt: nextAttempt,
			}); failErr != nil {
				s.logger.Error("failed to record outbox retry", zap.Error(failErr), zap.Int64("outbox_id", event.OutboxID))
			}
			continue
		}
		if _, err := s.queries.MarkOutboxEventDelivered(ctx, event.OutboxID); err != nil {
			s.logger.Error("failed to mark outbox event delivered", zap.Error(err), zap.Int64("outbox_id", event.OutboxID))
			continue
		}
		delivered++
	}
	return delivered, nil
}

// Start runs the relay loop until ctx is cancelled.
func (s *OutboxService) Start(ctx context.Context, interval time.Duration, limit int) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	tickCount := 0
	for {
		select {
		case <-ctx.Done():
			s.logger.Info("outbox relay stopped")
			return
		case <-ticker.C:
			if _, err := s.PublishPending(ctx, limit); err != nil {
				s.logger.Error("outbox relay cycle failed", zap.Error(err))
			}
			tickCount++
			// Retention runs periodically (not every tick) to avoid scanning the
			// outbox table on every relay cycle; it purges delivered/dead events
			// whose terminal state is older than the retention window.
			if s.queries != nil && tickCount%OutboxRetentionEveryTicks == 0 {
				if removed, err := s.queries.DeleteOldOutboxEvents(ctx, time.Now().Add(-OutboxRetention)); err != nil {
					s.logger.Error("outbox retention cleanup failed", zap.Error(err))
				} else if removed > 0 {
					s.logger.Info("outbox retention cleanup", zap.Int64("removed", removed))
				}
			}
		}
	}
}
