package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	sharedErrors "github.com/MamangRust/microservice-ecommerce-shared/errors"
	resettoken_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/reset_token_errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type resetTokenRepository struct {
	db *mongo.Database
}

func NewResetTokenRepository(db *mongo.Database) *resetTokenRepository {
	return &resetTokenRepository{db: db}
}

func (r *resetTokenRepository) FindByToken(ctx context.Context, code string) (*models.ResetToken, error) {
	var result models.ResetToken
	err := r.db.Collection("reset_tokens").FindOne(ctx, bson.M{"token": code}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, resettoken_errors.ErrTokenNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &result, nil
}

func (r *resetTokenRepository) CreateResetToken(ctx context.Context, req *requests.CreateResetTokenRequest) (*models.ResetToken, error) {
	expiryDate, err := time.Parse("2006-01-02 15:04:05", req.ExpiredAt)
	if err != nil {
		return nil, err
	}
	resetToken := &models.ResetToken{
		UserID:     int64(req.UserID),
		Token:      req.ResetToken,
		ExpiryDate: expiryDate,
	}
	_, err = r.db.Collection("reset_tokens").InsertOne(ctx, resetToken)
	if err != nil {
		return nil, resettoken_errors.ErrCreateResetToken.WithInternal(err)
	}
	return resetToken, nil
}

func (r *resetTokenRepository) CreateResetTokenInTx(ctx context.Context, tx *mongo.Database, req *requests.CreateResetTokenRequest) (*models.ResetToken, error) {
	expiryDate, err := time.Parse("2006-01-02 15:04:05", req.ExpiredAt)
	if err != nil {
		return nil, err
	}
	resetToken := &models.ResetToken{
		UserID:     int64(req.UserID),
		Token:      req.ResetToken,
		ExpiryDate: expiryDate,
	}
	_, err = tx.Collection("reset_tokens").InsertOne(ctx, resetToken)
	if err != nil {
		return nil, resettoken_errors.ErrCreateResetToken.WithInternal(err)
	}
	return resetToken, nil
}

func (r *resetTokenRepository) DeleteResetToken(ctx context.Context, userID int) error {
	_, err := r.db.Collection("reset_tokens").DeleteMany(ctx, bson.M{"user_id": userID})
	if err != nil {
		return resettoken_errors.ErrDeleteByUserID.WithInternal(err)
	}
	return nil
}
