package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	shared_errors "github.com/MamangRust/microservice-ecommerce-shared/errors"
	"github.com/MamangRust/microservice-ecommerce-shared/errors/product_errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type productCommandRepository struct {
	db *mongo.Database
}

func NewProductCommandRepository(db *mongo.Database) *productCommandRepository {
	return &productCommandRepository{db: db}
}

func productTimePtr(t time.Time) *time.Time { return &t }

func productInt32Ptr(v int) *int32 {
	i := int32(v)
	return &i
}

func (r *productCommandRepository) Create(ctx context.Context, request *requests.CreateProductRequest) (*models.Product, error) {
	id, err := database.GetNextID(ctx, r.db, "products")
	if err != nil {
		return nil, product_errors.ErrCreateProduct.WithInternal(err)
	}
	now := time.Now()
	product := &models.Product{
		ProductID:    id,
		MerchantID:   int32(request.MerchantID),
		CategoryID:   int32(request.CategoryID),
		Name:         request.Name,
		Description:  &request.Description,
		Price:        int32(request.Price),
		CountInStock: int32(request.CountInStock),
		Brand:        &request.Brand,
		Weight:       productInt32Ptr(request.Weight),
		SlugProduct:  request.SlugProduct,
		ImageProduct: &request.ImageProduct,
		CreatedAt:    &now,
		UpdatedAt:    &now,
	}
	if _, err := r.db.Collection("products").InsertOne(ctx, product); err != nil {
		return nil, product_errors.ErrCreateProduct.WithInternal(err)
	}
	return product, nil
}

func (r *productCommandRepository) Update(ctx context.Context, request *requests.UpdateProductRequest) (*models.Product, error) {
	coll := r.db.Collection("products")
	var product models.Product
	if err := coll.FindOne(ctx, bson.M{"_id": *request.ProductID}).Decode(&product); err != nil {
		return nil, product_errors.ErrProductNotFound
	}
	desc := request.Description
	brand := request.Brand
	img := request.ImageProduct
	now := time.Now()
	update := bson.M{
		"category_id":    int32(request.CategoryID),
		"name":           request.Name,
		"description":    &desc,
		"price":          int32(request.Price),
		"count_in_stock": int32(request.CountInStock),
		"brand":          &brand,
		"weight":         productInt32Ptr(request.Weight),
		"slug_product":   request.SlugProduct,
		"image_product":  &img,
		"updated_at":     now,
	}
	if _, err := coll.UpdateOne(ctx, bson.M{"_id": *request.ProductID}, bson.M{"$set": update}); err != nil {
		return nil, product_errors.ErrUpdateProduct.WithInternal(err)
	}
	if err := coll.FindOne(ctx, bson.M{"_id": *request.ProductID}).Decode(&product); err != nil {
		return nil, product_errors.ErrUpdateProduct.WithInternal(err)
	}
	return &product, nil
}

func (r *productCommandRepository) UpdateProductCountStock(ctx context.Context, productID int, stock int) (*models.Product, error) {
	coll := r.db.Collection("products")
	var product models.Product
	if err := coll.FindOne(ctx, bson.M{"_id": productID}).Decode(&product); err != nil {
		return nil, product_errors.ErrProductNotFound
	}
	now := time.Now()
	if _, err := coll.UpdateOne(ctx, bson.M{"_id": productID}, bson.M{"$set": bson.M{"count_in_stock": int32(stock), "updated_at": now}}); err != nil {
		return nil, product_errors.ErrProductInternal.WithInternal(err)
	}
	product.CountInStock = int32(stock)
	product.UpdatedAt = &now
	return &product, nil
}

func (r *productCommandRepository) AdjustProductStock(ctx context.Context, productID int, delta int, operationID string) (*models.Product, error) {
	coll := r.db.Collection("products")
	var product models.Product
	if err := coll.FindOne(ctx, bson.M{"_id": productID}).Decode(&product); err != nil {
		return nil, product_errors.ErrProductNotFound
	}
	newStock := int(product.CountInStock) + delta
	if newStock < 0 {
		return nil, product_errors.ErrUpdateProductCountStock
	}
	now := time.Now()
	if _, err := coll.UpdateOne(ctx, bson.M{"_id": productID}, bson.M{"$set": bson.M{"count_in_stock": int32(newStock), "updated_at": now}}); err != nil {
		return nil, product_errors.ErrUpdateProductCountStock.WithInternal(err)
	}
	if operationID != "" {
		_, _ = r.db.Collection("product_stock_adjustments").InsertOne(ctx, models.ProductStockAdjustment{
			OperationID: operationID,
			ProductID:   int32(productID),
			Delta:       int32(delta),
			CreatedAt:   &now,
		})
	}
	product.CountInStock = int32(newStock)
	product.UpdatedAt = &now
	return &product, nil
}

func (r *productCommandRepository) Trash(ctx context.Context, productID int) (*models.Product, error) {
	coll := r.db.Collection("products")
	var product models.Product
	if err := coll.FindOne(ctx, bson.M{"_id": productID}).Decode(&product); err != nil {
		return nil, product_errors.ErrProductNotFound
	}
	now := time.Now()
	if _, err := coll.UpdateOne(ctx, bson.M{"_id": productID}, bson.M{"$set": bson.M{"deleted_at": now}}); err != nil {
		return nil, product_errors.ErrTrashedProduct.WithInternal(err)
	}
	product.DeletedAt = &now
	return &product, nil
}

func (r *productCommandRepository) Restore(ctx context.Context, productID int) (*models.Product, error) {
	coll := r.db.Collection("products")
	res, err := coll.UpdateOne(ctx, bson.M{"_id": productID, "deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	if err != nil {
		return nil, product_errors.ErrRestoreProduct.WithInternal(err)
	}
	if res.MatchedCount == 0 {
		return nil, product_errors.ErrProductNotFound
	}
	var product models.Product
	if err := coll.FindOne(ctx, bson.M{"_id": productID}).Decode(&product); err != nil {
		return nil, product_errors.ErrRestoreProduct.WithInternal(err)
	}
	return &product, nil
}

func (r *productCommandRepository) DeletePermanent(ctx context.Context, productID int) (bool, error) {
	res, err := r.db.Collection("products").DeleteOne(ctx, bson.M{"_id": productID})
	if err != nil {
		return false, shared_errors.NewConflictError("cannot permanently delete product while related records exist").WithInternal(err)
	}
	if res.DeletedCount == 0 {
		return false, product_errors.ErrProductNotFound
	}
	return true, nil
}

func (r *productCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	_, err := r.db.Collection("products").UpdateMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	if err != nil {
		return false, product_errors.ErrRestoreAllProducts.WithInternal(err)
	}
	return true, nil
}

func (r *productCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	_, err := r.db.Collection("products").DeleteMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, shared_errors.NewConflictError("cannot permanently delete products while related records exist").WithInternal(err)
	}
	return true, nil
}

func timePtr(t time.Time) *time.Time { return &t }
