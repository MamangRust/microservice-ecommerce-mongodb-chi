package service

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-pkg/outbox"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// mongoOutboxQuerier adapts *mongo.Database to satisfy the outbox.OutboxQuerier interface.
type mongoOutboxQuerier struct {
	db *mongo.Database
}

// NewGormOutboxQuerier builds the outbox adapter from a MongoDB database.
func NewGormOutboxQuerier(db *mongo.Database) outbox.OutboxQuerier {
	if db == nil {
		return nil
	}
	return &mongoOutboxQuerier{db: db}
}

func (q *mongoOutboxQuerier) CreateOutboxEvent(ctx context.Context, params outbox.CreateOutboxEventParams) (outbox.OutboxEvent, error) {
	now := time.Now()
	event := models.OutboxEvent{
		OutboxID:      time.Now().UnixNano(),
		Topic:         params.Topic,
		EventKey:      params.EventKey,
		Payload:       params.Payload,
		Status:        "pending",
		NextAttemptAt: now,
		CreatedAt:     &now,
		UpdatedAt:     &now,
	}

	coll := q.db.Collection("outbox_events")
	_, err := coll.InsertOne(ctx, event)
	if err != nil {
		return outbox.OutboxEvent{}, err
	}

	return outboxEventToDomain(event), nil
}

func (q *mongoOutboxQuerier) ClaimPendingOutboxEvents(ctx context.Context, params outbox.ClaimPendingOutboxEventsParams) ([]outbox.OutboxEvent, error) {
	coll := q.db.Collection("outbox_events")

	// Find pending events and update them atomically
	filter := bson.M{
		"status":          "pending",
		"next_attempt_at": bson.M{"$lte": time.Now()},
	}
	update := bson.M{
		"$set": bson.M{"next_attempt_at": params.NextAttemptAt},
	}

	// Find first, then update (MongoDB doesn't have SKIP LOCKED)
	var events []models.OutboxEvent
	cursor, err := coll.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "_id", Value: 1}}).SetLimit(int64(params.Limit)))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &events); err != nil {
		return nil, err
	}

	// Update claimed events
	ids := make([]int64, len(events))
	for i, e := range events {
		ids[i] = e.OutboxID
	}
	if len(ids) > 0 {
		_, err = coll.UpdateMany(ctx, bson.M{"_id": bson.M{"$in": ids}}, update)
		if err != nil {
			return nil, err
		}
	}

	result := make([]outbox.OutboxEvent, 0, len(events))
	for _, e := range events {
		result = append(result, outboxEventToDomain(e))
	}
	return result, nil
}

func (q *mongoOutboxQuerier) MarkOutboxEventFailed(ctx context.Context, params outbox.MarkOutboxEventFailedParams) (outbox.OutboxEvent, error) {
	coll := q.db.Collection("outbox_events")

	filter := bson.M{"_id": params.OutboxID, "status": "pending"}
	now := time.Now()
	update := bson.M{
		"$inc": bson.M{"attempts": 1},
		"$set": bson.M{
			"next_attempt_at": params.NextAttemptAt,
			"updated_at":      now,
		},
	}

	var event models.OutboxEvent
	err := coll.FindOneAndUpdate(ctx, filter, update).Decode(&event)
	if err != nil {
		return outbox.OutboxEvent{}, err
	}
	return outboxEventToDomain(event), nil
}

func (q *mongoOutboxQuerier) MarkOutboxEventDelivered(ctx context.Context, outboxID int64) (outbox.OutboxEvent, error) {
	coll := q.db.Collection("outbox_events")

	filter := bson.M{"_id": outboxID, "status": "pending"}
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"status":     "delivered",
			"updated_at": now,
		},
	}

	var event models.OutboxEvent
	err := coll.FindOneAndUpdate(ctx, filter, update).Decode(&event)
	if err != nil {
		return outbox.OutboxEvent{}, fmt.Errorf("outbox event not found or already processed: %w", err)
	}
	return outboxEventToDomain(event), nil
}

func (q *mongoOutboxQuerier) MarkOutboxEventDead(ctx context.Context, outboxID int64) (outbox.OutboxEvent, error) {
	coll := q.db.Collection("outbox_events")

	filter := bson.M{"_id": outboxID, "status": "pending"}
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"status":     "dead",
			"updated_at": now,
		},
	}

	var event models.OutboxEvent
	err := coll.FindOneAndUpdate(ctx, filter, update).Decode(&event)
	if err != nil {
		return outbox.OutboxEvent{}, fmt.Errorf("outbox event not found or already processed: %w", err)
	}
	return outboxEventToDomain(event), nil
}

func (q *mongoOutboxQuerier) DeleteOldOutboxEvents(ctx context.Context, cutoff time.Time) (int64, error) {
	coll := q.db.Collection("outbox_events")

	filter := bson.M{
		"status":     bson.M{"$in": []string{"delivered", "dead"}},
		"updated_at": bson.M{"$lt": cutoff},
	}

	result, err := coll.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

func outboxEventToDomain(e models.OutboxEvent) outbox.OutboxEvent {
	var createdAt, updatedAt time.Time
	if e.CreatedAt != nil {
		createdAt = *e.CreatedAt
	}
	if e.UpdatedAt != nil {
		updatedAt = *e.UpdatedAt
	}
	return outbox.OutboxEvent{
		OutboxID:      e.OutboxID,
		Topic:         e.Topic,
		EventKey:      e.EventKey,
		Payload:       e.Payload,
		Status:        e.Status,
		Attempts:      e.Attempts,
		NextAttemptAt: e.NextAttemptAt,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}
