package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/errors/category_errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type categoryCommandRepository struct {
	db *mongo.Database
}

func NewCategoryCommandRepository(db *mongo.Database) *categoryCommandRepository {
	return &categoryCommandRepository{db: db}
}

func (r *categoryCommandRepository) Create(ctx context.Context, request *requests.CreateCategoryRequest) (*models.Category, error) {
	id, err := database.GetNextID(ctx, r.db, "categories")
	if err != nil {
		return nil, category_errors.ErrCreateCategory.WithInternal(err)
	}
	now := time.Now()
	category := &models.Category{
		CategoryID:    id,
		Name:          request.Name,
		Description:   &request.Description,
		SlugCategory:  request.SlugCategory,
		ImageCategory: &request.ImageCategory,
		CreatedAt:     &now,
		UpdatedAt:     &now,
	}
	_, err = r.db.Collection("categories").InsertOne(ctx, category)
	if err != nil {
		return nil, category_errors.ErrCreateCategory.WithInternal(err)
	}
	return category, nil
}

func (r *categoryCommandRepository) Update(ctx context.Context, request *requests.UpdateCategoryRequest) (*models.Category, error) {
	coll := r.db.Collection("categories")
	var category models.Category
	err := coll.FindOne(ctx, bson.M{"_id": *request.CategoryID, "deleted_at": bson.M{"$eq": nil}}).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, category_errors.ErrCategoryNotFound
		}
		return nil, category_errors.ErrUpdateCategory.WithInternal(err)
	}
	category.Name = request.Name
	category.Description = &request.Description
	category.SlugCategory = request.SlugCategory
	category.ImageCategory = &request.ImageCategory
	now := time.Now()
	category.UpdatedAt = &now
	_, err = coll.ReplaceOne(ctx, bson.M{"_id": *request.CategoryID}, category)
	if err != nil {
		return nil, category_errors.ErrUpdateCategory.WithInternal(err)
	}
	return &category, nil
}

func (r *categoryCommandRepository) Trash(ctx context.Context, categoryID int) (*models.Category, error) {
	coll := r.db.Collection("categories")
	now := time.Now()
	res, err := coll.UpdateOne(ctx, bson.M{"_id": categoryID, "deleted_at": bson.M{"$eq": nil}}, bson.M{"$set": bson.M{"deleted_at": now}})
	if err != nil {
		return nil, category_errors.ErrTrashedCategory.WithInternal(err)
	}
	if res.MatchedCount == 0 {
		return nil, category_errors.ErrCategoryNotFound
	}
	var category models.Category
	if err := coll.FindOne(ctx, bson.M{"_id": categoryID}).Decode(&category); err != nil {
		return nil, category_errors.ErrTrashedCategory.WithInternal(err)
	}
	return &category, nil
}

func (r *categoryCommandRepository) Restore(ctx context.Context, categoryID int) (*models.Category, error) {
	coll := r.db.Collection("categories")
	res, err := coll.UpdateOne(ctx, bson.M{"_id": categoryID, "deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	if err != nil {
		return nil, category_errors.ErrRestoreCategory.WithInternal(err)
	}
	if res.MatchedCount == 0 {
		return nil, category_errors.ErrCategoryNotFound
	}
	var category models.Category
	if err := coll.FindOne(ctx, bson.M{"_id": categoryID}).Decode(&category); err != nil {
		return nil, category_errors.ErrRestoreCategory.WithInternal(err)
	}
	return &category, nil
}

func (r *categoryCommandRepository) DeletePermanent(ctx context.Context, categoryID int) (bool, error) {
	res, err := r.db.Collection("categories").DeleteOne(ctx, bson.M{"_id": categoryID})
	if err != nil {
		return false, category_errors.ErrDeleteCategoryPermanently.WithInternal(err)
	}
	if res.DeletedCount == 0 {
		return false, category_errors.ErrCategoryNotFound
	}
	return true, nil
}

func (r *categoryCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	_, err := r.db.Collection("categories").UpdateMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	if err != nil {
		return false, category_errors.ErrRestoreAllCategories.WithInternal(err)
	}
	return true, nil
}

func (r *categoryCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	_, err := r.db.Collection("categories").DeleteMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, category_errors.ErrDeleteAllPermanentCategories.WithInternal(err)
	}
	return true, nil
}
