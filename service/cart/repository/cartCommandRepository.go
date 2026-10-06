package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/errors/cart_errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type cartCommandRepository struct {
	db *mongo.Database
}

func NewCartCommandRepository(db *mongo.Database) CartCommandRepository {
	return &cartCommandRepository{db: db}
}

func (r *cartCommandRepository) CreateCart(ctx context.Context, req *requests.CartCreateRecord) (*CartCreateResult, error) {
	id, err := database.GetNextID(ctx, r.db, "carts")
	if err != nil {
		return nil, cart_errors.ErrCreateCart.WithInternal(err)
	}
	now := time.Now()
	cart := models.Cart{
		CartID:    id,
		UserID:    int32(req.UserID),
		ProductID: int32(req.ProductID),
		Name:      req.Name,
		Price:     int32(req.Price),
		Image:     req.ImageProduct,
		Quantity:  int32(req.Quantity),
		Weight:    int32(req.Weight),
		CreatedAt: &now,
		UpdatedAt: &now,
	}
	if _, err := r.db.Collection("carts").InsertOne(ctx, cart); err != nil {
		return nil, cart_errors.ErrCreateCart.WithInternal(err)
	}
	return &CartCreateResult{
		CartID:    id,
		UserID:    int32(req.UserID),
		ProductID: int32(req.ProductID),
		Name:      req.Name,
		Price:     int32(req.Price),
		Image:     req.ImageProduct,
		Quantity:  int32(req.Quantity),
		Weight:    int32(req.Weight),
		CreatedAt: now.Format("2006-01-02 15:04:05"),
		UpdatedAt: now.Format("2006-01-02 15:04:05"),
	}, nil
}

func (r *cartCommandRepository) DeletePermanent(ctx context.Context, req *requests.DeleteCartRequest) (bool, error) {
	res, err := r.db.Collection("carts").DeleteOne(ctx, bson.M{"_id": req.CartID, "user_id": req.UserID})
	if err != nil {
		return false, cart_errors.ErrDeleteCartPermanent.WithInternal(err)
	}
	if res.DeletedCount == 0 {
		return false, cart_errors.ErrFindAllCarts
	}
	return true, nil
}

func (r *cartCommandRepository) DeleteAllPermanently(ctx context.Context, req *requests.DeleteAllCartRequest) (bool, error) {
	ids := make([]int32, len(req.CartIds))
	for i, id := range req.CartIds {
		ids[i] = int32(id)
	}
	_, err := r.db.Collection("carts").DeleteMany(ctx, bson.M{"_id": bson.M{"$in": ids}, "user_id": req.UserID})
	if err != nil {
		return false, cart_errors.ErrDeleteAllCarts.WithInternal(err)
	}
	return true, nil
}
