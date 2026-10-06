package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type outboxRepository struct {
	db *mongo.Database
}

func NewOutboxRepository(db *mongo.Database) OutboxRepository {
	return &outboxRepository{db: db}
}

func toOutboxResult(e *models.OutboxEvent) *OutboxEventResult {
	if e == nil {
		return nil
	}
	return &OutboxEventResult{
		OutboxID:      e.OutboxID,
		Topic:         e.Topic,
		EventKey:      e.EventKey,
		Payload:       e.Payload,
		Status:        e.Status,
		Attempts:      e.Attempts,
		NextAttemptAt: e.NextAttemptAt,
		CreatedAt:     e.CreatedAt,
		UpdatedAt:     e.UpdatedAt,
	}
}

func (r *outboxRepository) Create(ctx context.Context, topic, key string, payload []byte) (*OutboxEventResult, error) {
	return r.create(ctx, r.db, topic, key, payload)
}

func (r *outboxRepository) CreateInTx(ctx context.Context, tx *mongo.Database, topic, key string, payload []byte) (*OutboxEventResult, error) {
	return r.create(ctx, tx, topic, key, payload)
}

func (r *outboxRepository) create(ctx context.Context, db *mongo.Database, topic, key string, payload []byte) (*OutboxEventResult, error) {
	now := time.Now()
	event := &models.OutboxEvent{
		OutboxID:      now.UnixNano(),
		Topic:         topic,
		EventKey:      key,
		Payload:       payload,
		Status:        "pending",
		NextAttemptAt: now,
		CreatedAt:     &now,
		UpdatedAt:     &now,
	}
	if _, err := db.Collection("outbox_events").InsertOne(ctx, event); err != nil {
		return nil, err
	}
	return toOutboxResult(event), nil
}

func (r *outboxRepository) GetPending(ctx context.Context, limit int) ([]*OutboxEventResult, error) {
	filter := bson.M{"status": "pending", "next_attempt_at": bson.M{"$lte": time.Now()}}
	opts := options.Find().SetSort(bson.D{{Key: "_id", Value: 1}}).SetLimit(int64(limit))
	cursor, err := r.db.Collection("outbox_events").Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var events []*models.OutboxEvent
	if err := cursor.All(ctx, &events); err != nil {
		return nil, err
	}
	var results []*OutboxEventResult
	for _, e := range events {
		results = append(results, toOutboxResult(e))
	}
	return results, nil
}

func (r *outboxRepository) Claim(ctx context.Context, limit int, leaseUntil time.Time) ([]*OutboxEventResult, error) {
	results, err := r.GetPending(ctx, limit)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	ids := make([]int64, len(results))
	for i, e := range results {
		ids[i] = e.OutboxID
	}
	if len(ids) > 0 {
		if _, err := r.db.Collection("outbox_events").UpdateMany(ctx,
			bson.M{"_id": bson.M{"$in": ids}},
			bson.M{"$set": bson.M{"next_attempt_at": leaseUntil, "updated_at": now}},
		); err != nil {
			return nil, err
		}
		for _, e := range results {
			e.NextAttemptAt = leaseUntil
			e.UpdatedAt = &now
		}
	}
	return results, nil
}

func (r *outboxRepository) MarkDelivered(ctx context.Context, outboxID int64) (*OutboxEventResult, error) {
	return r.mark(ctx, outboxID, bson.M{"$set": bson.M{"status": "delivered", "updated_at": time.Now()}})
}

func (r *outboxRepository) MarkFailed(ctx context.Context, outboxID int64, nextAttemptAt time.Time) (*OutboxEventResult, error) {
	return r.mark(ctx, outboxID, bson.M{
		"$inc": bson.M{"attempts": 1},
		"$set": bson.M{"next_attempt_at": nextAttemptAt, "updated_at": time.Now()},
	})
}

func (r *outboxRepository) MarkDead(ctx context.Context, outboxID int64) (*OutboxEventResult, error) {
	return r.mark(ctx, outboxID, bson.M{"$set": bson.M{"status": "dead", "updated_at": time.Now()}})
}

func (r *outboxRepository) mark(ctx context.Context, outboxID int64, update bson.M) (*OutboxEventResult, error) {
	var event models.OutboxEvent
	err := r.db.Collection("outbox_events").FindOneAndUpdate(ctx,
		bson.M{"_id": outboxID, "status": "pending"}, update,
	).Decode(&event)
	if err != nil {
		return nil, err
	}
	return toOutboxResult(&event), nil
}

func (r *outboxRepository) DeleteOld(ctx context.Context, cutoff time.Time) (int64, error) {
	result, err := r.db.Collection("outbox_events").DeleteMany(ctx, bson.M{
		"status":     bson.M{"$in": []string{"delivered", "dead"}},
		"updated_at": bson.M{"$lt": cutoff},
	})
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}
