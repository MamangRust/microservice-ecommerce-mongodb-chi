package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	sharedErrors "github.com/MamangRust/microservice-ecommerce-shared/errors"
	"github.com/MamangRust/microservice-ecommerce-shared/errors/user_errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type userQueryRepository struct {
	db *mongo.Database
}

func NewUserQueryRepository(db *mongo.Database) *userQueryRepository {
	return &userQueryRepository{db: db}
}

func userBuildSearchFilter(search string) bson.M {
	if search == "" {
		return nil
	}
	return bson.M{"$or": []bson.M{
		{"firstname": bson.M{"$regex": search, "$options": "i"}},
		{"lastname": bson.M{"$regex": search, "$options": "i"}},
		{"email": bson.M{"$regex": search, "$options": "i"}},
	}}
}

func userRunPaginatedQuery(ctx context.Context, coll *mongo.Collection, filter bson.M, sort bson.M, offset int, pageSize int) ([]*UserResult, error) {
	pipeline := []bson.M{
		{"$match": filter},
		{"$sort": sort},
		{"$facet": bson.M{
			"metadata": []bson.M{{"$count": "total_count"}},
			"data": []bson.M{
				{"$skip": offset},
				{"$limit": pageSize},
			},
		}},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var facetResults []struct {
		Metadata []struct{ TotalCount int64 `bson:"total_count"` } `bson:"metadata"`
		Data     []bson.M `bson:"data"`
	}
	if err := cursor.All(ctx, &facetResults); err != nil {
		return nil, err
	}
	if len(facetResults) == 0 {
		return []*UserResult{}, nil
	}

	var results []*UserResult
	for _, d := range facetResults[0].Data {
		item := &UserResult{
			UserID:    bsonMToInt32(d, "_id"),
			Firstname: bsonMToStr(d, "firstname"),
			Lastname:  bsonMToStr(d, "lastname"),
			Email:     bsonMToStr(d, "email"),
			Password:  bsonMToStr(d, "password"),
			CreatedAt: bsonMToTimePtr(d, "created_at"),
			UpdatedAt: bsonMToTimePtr(d, "updated_at"),
			DeletedAt: bsonMToTimePtr(d, "deleted_at"),
		}
		if len(facetResults[0].Metadata) > 0 {
			item.TotalCount = facetResults[0].Metadata[0].TotalCount
		}
		results = append(results, item)
	}
	return results, nil
}

func (r *userQueryRepository) FindAll(ctx context.Context, req *requests.FindAllUsers) ([]*UserResult, error) {
	offset := (req.Page - 1) * req.PageSize
	filter := bson.M{"deleted_at": bson.M{"$eq": nil}}
	if searchFilter := userBuildSearchFilter(req.Search); searchFilter != nil {
		filter["$or"] = searchFilter["$or"]
	}
	results, err := userRunPaginatedQuery(ctx, r.db.Collection("users"), filter, bson.M{"created_at": -1}, offset, req.PageSize)
	if err != nil {
		return nil, user_errors.ErrFindAllUsers.WithInternal(err)
	}
	return results, nil
}

func (r *userQueryRepository) FindByID(ctx context.Context, userID int) (*models.User, error) {
	var user models.User
	err := r.db.Collection("users").FindOne(ctx, bson.M{"_id": userID, "deleted_at": bson.M{"$eq": nil}}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, user_errors.ErrUserNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &user, nil
}

func (r *userQueryRepository) FindByIDWithPassword(ctx context.Context, userID int) (*models.User, error) {
	var user models.User
	err := r.db.Collection("users").FindOne(ctx, bson.M{"_id": userID, "deleted_at": bson.M{"$eq": nil}}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, user_errors.ErrUserNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &user, nil
}

func (r *userQueryRepository) FindActive(ctx context.Context, req *requests.FindAllUsers) ([]*UserResult, error) {
	offset := (req.Page - 1) * req.PageSize
	filter := bson.M{"deleted_at": bson.M{"$eq": nil}}
	if searchFilter := userBuildSearchFilter(req.Search); searchFilter != nil {
		filter["$or"] = searchFilter["$or"]
	}
	results, err := userRunPaginatedQuery(ctx, r.db.Collection("users"), filter, bson.M{"created_at": -1}, offset, req.PageSize)
	if err != nil {
		return nil, user_errors.ErrFindActiveUsers.WithInternal(err)
	}
	return results, nil
}

func (r *userQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllUsers) ([]*UserResult, error) {
	offset := (req.Page - 1) * req.PageSize
	filter := bson.M{"deleted_at": bson.M{"$ne": nil}}
	if searchFilter := userBuildSearchFilter(req.Search); searchFilter != nil {
		filter["$or"] = searchFilter["$or"]
	}
	results, err := userRunPaginatedQuery(ctx, r.db.Collection("users"), filter, bson.M{"created_at": -1}, offset, req.PageSize)
	if err != nil {
		return nil, user_errors.ErrFindTrashedUsers.WithInternal(err)
	}
	return results, nil
}

func (r *userQueryRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.Collection("users").FindOne(ctx, bson.M{"email": email, "deleted_at": bson.M{"$eq": nil}}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, user_errors.ErrUserNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &user, nil
}

func (r *userQueryRepository) FindByEmailWithPassword(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.Collection("users").FindOne(ctx, bson.M{"email": email, "deleted_at": bson.M{"$eq": nil}}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, user_errors.ErrUserNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &user, nil
}

func (r *userQueryRepository) FindByVerificationCode(ctx context.Context, code string) (*models.User, error) {
	var user models.User
	err := r.db.Collection("users").FindOne(ctx, bson.M{"verification_code": code}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, user_errors.ErrUserNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &user, nil
}

func bsonMToTimePtr(m bson.M, key string) *time.Time {
	v, ok := m[key]
	if !ok || v == nil {
		return nil
	}
	switch t := v.(type) {
	case time.Time:
		return &t
	case bson.DateTime:
		ts := t.Time()
		return &ts
	default:
		return nil
	}
}

func bsonMToInt32(m bson.M, key string) int32 {
	v, ok := m[key]
	if !ok {
		return 0
	}
	switch t := v.(type) {
	case int32:
		return t
	case int64:
		return int32(t)
	case float64:
		return int32(t)
	default:
		return 0
	}
}

func bsonMToStr(m bson.M, key string) string {
	v, ok := m[key]
	if !ok {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
