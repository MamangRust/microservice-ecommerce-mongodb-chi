package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	sharedErrors "github.com/MamangRust/microservice-ecommerce-shared/errors"
	refreshtoken_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/refresh_token_errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type refreshTokenRepository struct {
	db *mongo.Database
}

func NewRefreshTokenRepository(db *mongo.Database) *refreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) FindByToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	var result models.RefreshToken
	err := r.db.Collection("refresh_tokens").FindOne(ctx, bson.M{"token": token}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, refreshtoken_errors.ErrTokenNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &result, nil
}

func (r *refreshTokenRepository) FindByUserId(ctx context.Context, userID int) (*models.RefreshToken, error) {
	var result models.RefreshToken
	err := r.db.Collection("refresh_tokens").FindOne(ctx, bson.M{"user_id": userID}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, refreshtoken_errors.ErrTokenNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &result, nil
}

func (r *refreshTokenRepository) CreateRefreshToken(ctx context.Context, req *requests.CreateRefreshToken) (*models.RefreshToken, error) {
	layout := "2006-01-02 15:04:05"
	expirationTime, err := time.Parse(layout, req.ExpiresAt)
	if err != nil {
		return nil, refreshtoken_errors.ErrParseDate
	}
	now := time.Now()
	refreshToken := &models.RefreshToken{
		UserID:     int32(req.UserId),
		Token:      req.Token,
		Expiration: expirationTime,
		CreatedAt:  &now,
		UpdatedAt:  &now,
	}
	_, err = r.db.Collection("refresh_tokens").InsertOne(ctx, refreshToken)
	if err != nil {
		return nil, refreshtoken_errors.ErrCreateRefreshToken.WithInternal(err)
	}
	return refreshToken, nil
}

func (r *refreshTokenRepository) UpdateRefreshToken(ctx context.Context, req *requests.UpdateRefreshToken) (*models.RefreshToken, error) {
	layout := "2006-01-02 15:04:05"
	expirationTime, err := time.Parse(layout, req.ExpiresAt)
	if err != nil {
		return nil, refreshtoken_errors.ErrParseDate
	}
	coll := r.db.Collection("refresh_tokens")
	var result models.RefreshToken
	err = coll.FindOne(ctx, bson.M{"user_id": req.UserId}).Decode(&result)
	if err != nil {
		return nil, refreshtoken_errors.ErrUpdateRefreshToken.WithInternal(err)
	}
	result.Token = req.Token
	result.Expiration = expirationTime
	now := time.Now()
	result.UpdatedAt = &now
	_, err = coll.ReplaceOne(ctx, bson.M{"user_id": req.UserId}, result)
	if err != nil {
		return nil, refreshtoken_errors.ErrUpdateRefreshToken.WithInternal(err)
	}
	return &result, nil
}

func (r *refreshTokenRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	_, err := r.db.Collection("refresh_tokens").DeleteOne(ctx, bson.M{"token": token})
	if err != nil {
		return refreshtoken_errors.ErrDeleteRefreshToken.WithInternal(err)
	}
	return nil
}

func (r *refreshTokenRepository) DeleteRefreshTokenByUserId(ctx context.Context, userID int) error {
	_, err := r.db.Collection("refresh_tokens").DeleteMany(ctx, bson.M{"user_id": userID})
	if err != nil {
		return refreshtoken_errors.ErrDeleteByUserID.WithInternal(err)
	}
	return nil
}
