package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	shared_errors "github.com/MamangRust/microservice-ecommerce-shared/errors"
	"github.com/MamangRust/microservice-ecommerce-shared/errors/user_errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type userCommandRepository struct {
	db *mongo.Database
}

func NewUserCommandRepository(db *mongo.Database) *userCommandRepository {
	return &userCommandRepository{db: db}
}

func (r *userCommandRepository) Create(ctx context.Context, request *requests.CreateUserRequest) (*models.User, error) {
	id, err := database.GetNextID(ctx, r.db, "users")
	if err != nil {
		return nil, user_errors.ErrCreateUser.WithInternal(err)
	}
	now := time.Now()
	user := &models.User{
		UserID:    id,
		Firstname: request.FirstName,
		Lastname:  request.LastName,
		Email:     request.Email,
		Password:  request.Password,
		CreatedAt: &now,
		UpdatedAt: &now,
	}
	_, err = r.db.Collection("users").InsertOne(ctx, user)
	if err != nil {
		return nil, user_errors.ErrCreateUser.WithInternal(err)
	}
	return user, nil
}

func (r *userCommandRepository) Update(ctx context.Context, request *requests.UpdateUserRequest) (*models.User, error) {
	coll := r.db.Collection("users")
	var user models.User
	err := coll.FindOne(ctx, bson.M{"_id": *request.UserID, "deleted_at": bson.M{"$eq": nil}}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, user_errors.ErrUserNotFound
		}
		return nil, user_errors.ErrUpdateUser.WithInternal(err)
	}

	user.Firstname = request.FirstName
	user.Lastname = request.LastName
	user.Email = request.Email
	if request.Password != "" {
		user.Password = request.Password
	}
	now := time.Now()
	user.UpdatedAt = &now

	_, err = coll.ReplaceOne(ctx, bson.M{"_id": *request.UserID}, user)
	if err != nil {
		return nil, user_errors.ErrUpdateUser.WithInternal(err)
	}
	return &user, nil
}

func (r *userCommandRepository) Trash(ctx context.Context, userID int) (*models.User, error) {
	coll := r.db.Collection("users")
	now := time.Now()
	res, err := coll.UpdateOne(ctx, bson.M{"_id": userID, "deleted_at": bson.M{"$eq": nil}}, bson.M{"$set": bson.M{"deleted_at": now}})
	if err != nil {
		return nil, user_errors.ErrTrashedUser.WithInternal(err)
	}
	if res.MatchedCount == 0 {
		return nil, user_errors.ErrUserNotFound
	}
	var user models.User
	if err := coll.FindOne(ctx, bson.M{"_id": userID}).Decode(&user); err != nil {
		return nil, user_errors.ErrTrashedUser.WithInternal(err)
	}
	return &user, nil
}

func (r *userCommandRepository) Restore(ctx context.Context, userID int) (*models.User, error) {
	coll := r.db.Collection("users")
	res, err := coll.UpdateOne(ctx, bson.M{"_id": userID, "deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	if err != nil {
		return nil, user_errors.ErrRestoreUser.WithInternal(err)
	}
	if res.MatchedCount == 0 {
		return nil, user_errors.ErrUserNotFound
	}
	var user models.User
	if err := coll.FindOne(ctx, bson.M{"_id": userID}).Decode(&user); err != nil {
		return nil, user_errors.ErrRestoreUser.WithInternal(err)
	}
	return &user, nil
}

func (r *userCommandRepository) DeletePermanent(ctx context.Context, userID int) (bool, error) {
	res, err := r.db.Collection("users").DeleteOne(ctx, bson.M{"_id": userID})
	if err != nil {
		return false, user_errors.ErrDeleteUserPermanent.WithInternal(err)
	}
	if res.DeletedCount == 0 {
		return false, user_errors.ErrUserNotFound
	}
	return true, nil
}

func (r *userCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	_, err := r.db.Collection("users").UpdateMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	if err != nil {
		return false, user_errors.ErrRestoreAllUsers.WithInternal(err)
	}
	return true, nil
}

func (r *userCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	_, err := r.db.Collection("users").DeleteMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, user_errors.ErrDeleteAllUsers.WithInternal(err)
	}
	return true, nil
}

func (r *userCommandRepository) UpdateIsVerified(ctx context.Context, userID int, isVerified bool) (*models.User, error) {
	coll := r.db.Collection("users")
	var user models.User
	err := coll.FindOne(ctx, bson.M{"_id": userID, "deleted_at": bson.M{"$eq": nil}}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, user_errors.ErrUserNotFound
		}
		return nil, user_errors.ErrUpdateUser.WithInternal(err)
	}
	user.IsVerified = &isVerified
	now := time.Now()
	user.UpdatedAt = &now
	_, err = coll.ReplaceOne(ctx, bson.M{"_id": userID}, user)
	if err != nil {
		return nil, user_errors.ErrUpdateUser.WithInternal(err)
	}
	return &user, nil
}

func (r *userCommandRepository) UpdatePassword(ctx context.Context, userID int, password string) (*models.User, error) {
	coll := r.db.Collection("users")
	var user models.User
	err := coll.FindOne(ctx, bson.M{"_id": userID, "deleted_at": bson.M{"$eq": nil}}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, user_errors.ErrUserNotFound
		}
		return nil, user_errors.ErrUpdateUser.WithInternal(err)
	}
	user.Password = password
	now := time.Now()
	user.UpdatedAt = &now
	_, err = coll.ReplaceOne(ctx, bson.M{"_id": userID}, user)
	if err != nil {
		return nil, user_errors.ErrUpdateUser.WithInternal(err)
	}
	return &user, nil
}

// Ensure shared_errors is used
var _ = shared_errors.ErrInternal
