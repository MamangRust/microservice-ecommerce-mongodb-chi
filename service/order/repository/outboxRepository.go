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

func (r *outboxRepository) Create(ctx context.Context, topic, key string, payload []byte) (*models.OutboxEvent, error) {
	now := time.Now()
	event := &models.OutboxEvent{
		OutboxID:      time.Now().UnixNano(),
		Topic:         topic,
		EventKey:      key,
		Payload:       payload,
		Status:        "pending",
		NextAttemptAt: now,
		CreatedAt:     &now,
		UpdatedAt:     &now,
	}
	_, err := r.db.Collection("outbox_events").InsertOne(ctx, event)
	return event, err
}

func (r *outboxRepository) GetPending(ctx context.Context, limit int) ([]*models.OutboxEvent, error) {
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
	return events, nil
}

func (r *outboxRepository) Claim(ctx context.Context, limit int, leaseUntil time.Time) ([]*models.OutboxEvent, error) {
	events, err := r.GetPending(ctx, limit)
	if err != nil {
		return nil, err
	}

	ids := make([]int64, len(events))
	for i, e := range events {
		ids[i] = e.OutboxID
	}
	if len(ids) > 0 {
		now := time.Now()
		_, err = r.db.Collection("outbox_events").UpdateMany(ctx,
			bson.M{"_id": bson.M{"$in": ids}},
			bson.M{"$set": bson.M{"next_attempt_at": leaseUntil, "updated_at": now}},
		)
		if err != nil {
			return nil, err
		}
		for _, e := range events {
			e.NextAttemptAt = leaseUntil
			e.UpdatedAt = &now
		}
	}
	return events, nil
}

func (r *outboxRepository) MarkDelivered(ctx context.Context, outboxID int64) (*models.OutboxEvent, error) {
	now := time.Now()
	var event models.OutboxEvent
	err := r.db.Collection("outbox_events").FindOneAndUpdate(ctx,
		bson.M{"_id": outboxID, "status": "pending"},
		bson.M{"$set": bson.M{"status": "delivered", "updated_at": now}},
	).Decode(&event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *outboxRepository) MarkFailed(ctx context.Context, outboxID int64, nextAttemptAt time.Time) (*models.OutboxEvent, error) {
	now := time.Now()
	var event models.OutboxEvent
	err := r.db.Collection("outbox_events").FindOneAndUpdate(ctx,
		bson.M{"_id": outboxID, "status": "pending"},
		bson.M{"$inc": bson.M{"attempts": 1}, "$set": bson.M{"next_attempt_at": nextAttemptAt, "updated_at": now}},
	).Decode(&event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *outboxRepository) MarkDead(ctx context.Context, outboxID int64) (*models.OutboxEvent, error) {
	now := time.Now()
	var event models.OutboxEvent
	err := r.db.Collection("outbox_events").FindOneAndUpdate(ctx,
		bson.M{"_id": outboxID, "status": "pending"},
		bson.M{"$set": bson.M{"status": "dead", "updated_at": now}},
	).Decode(&event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *outboxRepository) DeleteOld(ctx context.Context, cutoff time.Time) (int64, error) {
	result, err := r.db.Collection("outbox_events").DeleteMany(ctx,
		bson.M{"status": bson.M{"$in": []string{"delivered", "dead"}}, "updated_at": bson.M{"$lt": cutoff}},
	)
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}
