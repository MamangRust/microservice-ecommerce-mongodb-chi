package inbox

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var ErrInvalidInboxKey = errors.New("invalid consumer inbox key")

func inboxKey(consumerName, eventKey string) string {
	return consumerName + ":" + eventKey
}

// Reserve claims an event for a consumer. It returns false when the event was
// already processed. An expired processing lease may be reclaimed after a
// consumer crashes.
func Reserve(ctx context.Context, db *mongo.Database, consumerName, eventKey, topic string, partition int32, offset int64) (bool, bool, int64, error) {
	if db == nil || consumerName == "" || eventKey == "" {
		return false, false, 0, ErrInvalidInboxKey
	}

	coll := db.Collection("consumer_inbox")
	id := inboxKey(consumerName, eventKey)
	now := time.Now()
	leaseUntil := now.Add(1 * time.Minute)

	// Atomically claim an existing unprocessed/expired document.
	var claimed bson.M
	err := coll.FindOneAndUpdate(ctx,
		bson.M{
			"_id":         id,
			"status":      bson.M{"$ne": "processed"},
			"lease_until": bson.M{"$lte": now},
		},
		bson.M{
			"$set": bson.M{"status": "processing", "lease_until": leaseUntil, "last_error": "", "topic": topic, "partition_id": partition, "message_offset": offset},
			"$inc": bson.M{"attempts": 1, "reservation_version": 1},
		},
	).Decode(&claimed)

	if err == nil {
		return true, false, toInt64(claimed["reservation_version"]), nil
	}
	if err != mongo.ErrNoDocuments {
		return false, false, 0, err
	}

	// No claimable doc. Check whether it already exists.
	var existing bson.M
	err = coll.FindOne(ctx, bson.M{"_id": id}).Decode(&existing)
	if err == nil {
		// Exists but processed or still leased.
		if s, _ := existing["status"].(string); s == "processed" {
			return false, true, 0, nil
		}
		return false, false, 0, nil
	}
	if err != mongo.ErrNoDocuments {
		return false, false, 0, err
	}

	// Insert a fresh reservation.
	doc := bson.M{
		"_id":                 id,
		"consumer_name":       consumerName,
		"event_key":           eventKey,
		"topic":               topic,
		"partition_id":        partition,
		"message_offset":      offset,
		"status":              "processing",
		"attempts":            int32(1),
		"reservation_version": int64(1),
		"lease_until":         leaseUntil,
		"last_error":          "",
		"processed_at":        nil,
	}
	if _, err := coll.InsertOne(ctx, doc); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return false, false, 0, nil
		}
		return false, false, 0, err
	}
	return true, false, 1, nil
}

func MarkProcessed(ctx context.Context, db *mongo.Database, consumerName, eventKey string, reservationVersion int64) error {
	if db == nil || consumerName == "" || eventKey == "" {
		return ErrInvalidInboxKey
	}
	now := time.Now()
	_, err := db.Collection("consumer_inbox").UpdateOne(ctx,
		bson.M{"_id": inboxKey(consumerName, eventKey), "status": "processing", "reservation_version": reservationVersion},
		bson.M{"$set": bson.M{"status": "processed", "processed_at": now, "lease_until": now, "last_error": ""}},
	)
	return err
}

func Release(ctx context.Context, db *mongo.Database, consumerName, eventKey string, reservationVersion int64, processingErr error) error {
	if db == nil || consumerName == "" || eventKey == "" {
		return ErrInvalidInboxKey
	}
	lastError := "consumer processing failed"
	if processingErr != nil {
		lastError = processingErr.Error()
	}
	now := time.Now()
	_, err := db.Collection("consumer_inbox").UpdateOne(ctx,
		bson.M{"_id": inboxKey(consumerName, eventKey), "status": "processing", "reservation_version": reservationVersion},
		bson.M{"$set": bson.M{"status": "pending", "lease_until": now, "last_error": lastError}},
	)
	return err
}

func toInt64(v interface{}) int64 {
	switch t := v.(type) {
	case int64:
		return t
	case int32:
		return int64(t)
	case float64:
		return int64(t)
	default:
		return 0
	}
}

// MongoInbox adapts a MongoDB database to the outbox.ConsumerInbox contract.
type MongoInbox struct {
	db *mongo.Database
}

func NewMongoInbox(db *mongo.Database) (*MongoInbox, error) {
	if db == nil {
		return nil, errors.New("inbox db is nil")
	}
	return &MongoInbox{db: db}, nil
}

func (i *MongoInbox) Reserve(ctx context.Context, consumerName, eventKey, topic string, partition int32, offset int64) (bool, bool, int64, error) {
	return Reserve(ctx, i.db, consumerName, eventKey, topic, partition, offset)
}

func (i *MongoInbox) MarkProcessed(ctx context.Context, consumerName, eventKey string, reservationVersion int64) error {
	return MarkProcessed(ctx, i.db, consumerName, eventKey, reservationVersion)
}

func (i *MongoInbox) Release(ctx context.Context, consumerName, eventKey string, reservationVersion int64, processingErr error) error {
	return Release(ctx, i.db, consumerName, eventKey, reservationVersion, processingErr)
}
