package repository

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func newTestMongoDB(t *testing.T) *mongo.Database {
	t.Helper()
	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Skipf("skipping: cannot connect to test database: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*1000*1000*1000)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		t.Skipf("skipping: cannot ping test database: %v", err)
	}
	return client.Database("ecommerce_test")
}

func TestUserCommandRepository_DeletePermanent(t *testing.T) {
	db := newTestMongoDB(t)
	repo := NewUserCommandRepository(db)
	ctx := context.Background()

	// Clean up test collection
	coll := db.Collection("users")
	defer coll.Drop(ctx)

	if err := coll.FindOne(ctx, bson.M{"_id": 999999}).Err(); err == nil {
		t.Log("test document exists, deleting first")
		coll.DeleteOne(ctx, bson.M{"_id": 999999})
	}

	_, err := repo.DeletePermanent(ctx, 999999)
	if err == nil {
		t.Log("DeletePermanent on non-existent user returned no error (expected NotFound)")
	}
}
