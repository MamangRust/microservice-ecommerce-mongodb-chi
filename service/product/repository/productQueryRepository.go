package repository

import (
	"context"
	"time"

	categoryadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/category"

	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/errors/product_errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type productQueryRepository struct {
	db            *mongo.Database
	categoryQuery categoryadapter.QueryRepository
}

func NewProductQueryRepository(db *mongo.Database, categoryQuery categoryadapter.QueryRepository) *productQueryRepository {
	return &productQueryRepository{db: db, categoryQuery: categoryQuery}
}

func (r *productQueryRepository) findPaginated(ctx context.Context, filter bson.M) ([]*ProductResult, error) {
	cursor, err := r.db.Collection("products").Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var docs []bson.M
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}
	var results []*ProductResult
	for _, d := range docs {
		results = append(results, bsonMToProductResult(d))
	}
	return results, nil
}

func (r *productQueryRepository) FindAll(ctx context.Context, req *requests.FindAllProduct) ([]*ProductResult, error) {
	offset := (req.Page - 1) * req.PageSize
	filter := bson.M{"deleted_at": bson.M{"$eq": nil}}
	if req.Search != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": req.Search, "$options": "i"}},
			{"description": bson.M{"$regex": req.Search, "$options": "i"}},
		}
	}
	pipeline := []bson.M{
		{"$match": filter},
		{"$sort": bson.M{"product_id": -1}},
		{"$facet": bson.M{
			"metadata": []bson.M{{"$count": "total_count"}},
			"data":     []bson.M{{"$skip": offset}, {"$limit": req.PageSize}},
		}},
	}
	results, total, err := runPaginatedProductQuery(ctx, r.db.Collection("products"), pipeline)
	if err != nil {
		return nil, product_errors.ErrFindAllProducts.WithInternal(err)
	}
	for _, item := range results {
		item.TotalCount = total
	}
	return results, nil
}

func (r *productQueryRepository) FindActive(ctx context.Context, req *requests.FindAllProduct) ([]*ProductResult, error) {
	return r.FindAll(ctx, req)
}

func (r *productQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllProduct) ([]*ProductResult, error) {
	offset := (req.Page - 1) * req.PageSize
	filter := bson.M{"deleted_at": bson.M{"$ne": nil}}
	if req.Search != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": req.Search, "$options": "i"}},
			{"description": bson.M{"$regex": req.Search, "$options": "i"}},
		}
	}
	pipeline := []bson.M{
		{"$match": filter},
		{"$sort": bson.M{"product_id": -1}},
		{"$facet": bson.M{
			"metadata": []bson.M{{"$count": "total_count"}},
			"data":     []bson.M{{"$skip": offset}, {"$limit": req.PageSize}},
		}},
	}
	results, total, err := runPaginatedProductQuery(ctx, r.db.Collection("products"), pipeline)
	if err != nil {
		return nil, product_errors.ErrFindTrashedProducts.WithInternal(err)
	}
	for _, item := range results {
		item.TotalCount = total
	}
	return results, nil
}

func (r *productQueryRepository) FindByMerchant(ctx context.Context, req *requests.FindAllProductByMerchant) ([]*ProductResult, error) {
	offset := (req.Page - 1) * req.PageSize
	filter := bson.M{"deleted_at": bson.M{"$eq": nil}, "merchant_id": req.MerchantID}
	if req.Search != "" {
		filter["name"] = bson.M{"$regex": req.Search, "$options": "i"}
	}
	if req.CategoryID != 0 {
		filter["category_id"] = req.CategoryID
	}
	if req.MinPrice != nil {
		filter["price"] = bson.M{"$gte": *req.MinPrice}
	}
	if req.MaxPrice != nil {
		filter["price"] = bson.M{"$lte": *req.MaxPrice}
	}
	pipeline := []bson.M{
		{"$match": filter},
		{"$sort": bson.M{"product_id": -1}},
		{"$facet": bson.M{
			"metadata": []bson.M{{"$count": "total_count"}},
			"data":     []bson.M{{"$skip": offset}, {"$limit": req.PageSize}},
		}},
	}
	results, total, err := runPaginatedProductQuery(ctx, r.db.Collection("products"), pipeline)
	if err != nil {
		return nil, product_errors.ErrFindProductsByMerchant.WithInternal(err)
	}
	for _, item := range results {
		item.TotalCount = total
	}
	return results, nil
}

func (r *productQueryRepository) FindByCategory(ctx context.Context, req *requests.FindAllProductByCategory) ([]*ProductResult, error) {
	offset := (req.Page - 1) * req.PageSize
	categoryID := 0
	if req.CategoryName != "" {
		cats, err := r.categoryQuery.FindAll(ctx, 1, 1, req.CategoryName)
		if err != nil {
			return nil, product_errors.ErrFindProductsByCategory.WithInternal(err)
		}
		if len(cats) > 0 {
			categoryID = int(cats[0].CategoryID)
		}
	}
	filter := bson.M{"deleted_at": bson.M{"$eq": nil}}
	if categoryID != 0 {
		filter["category_id"] = categoryID
	}
	if req.Search != "" {
		filter["name"] = bson.M{"$regex": req.Search, "$options": "i"}
	}
	if req.MinPrice != nil {
		filter["price"] = bson.M{"$gte": *req.MinPrice}
	}
	if req.MaxPrice != nil {
		filter["price"] = bson.M{"$lte": *req.MaxPrice}
	}
	pipeline := []bson.M{
		{"$match": filter},
		{"$sort": bson.M{"product_id": -1}},
		{"$facet": bson.M{
			"metadata": []bson.M{{"$count": "total_count"}},
			"data":     []bson.M{{"$skip": offset}, {"$limit": req.PageSize}},
		}},
	}
	results, total, err := runPaginatedProductQuery(ctx, r.db.Collection("products"), pipeline)
	if err != nil {
		return nil, product_errors.ErrFindProductsByCategory.WithInternal(err)
	}
	for _, item := range results {
		item.TotalCount = total
	}
	return results, nil
}

func (r *productQueryRepository) FindByID(ctx context.Context, productID int) (*ProductResult, error) {
	var d bson.M
	err := r.db.Collection("products").FindOne(ctx, bson.M{"_id": productID, "deleted_at": bson.M{"$eq": nil}}).Decode(&d)
	if err != nil {
		return nil, product_errors.ErrProductNotFound
	}
	return bsonMToProductResult(d), nil
}

func (r *productQueryRepository) FindByIDTrashed(ctx context.Context, productID int) (*ProductResult, error) {
	var d bson.M
	err := r.db.Collection("products").FindOne(ctx, bson.M{"_id": productID}).Decode(&d)
	if err != nil {
		return nil, product_errors.ErrProductNotFound
	}
	return bsonMToProductResult(d), nil
}

func runPaginatedProductQuery(ctx context.Context, coll *mongo.Collection, pipeline []bson.M) ([]*ProductResult, int64, error) {
	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	var facetResults []struct {
		Metadata []struct{ TotalCount int64 `bson:"total_count"` } `bson:"metadata"`
		Data     []bson.M `bson:"data"`
	}
	if err := cursor.All(ctx, &facetResults); err != nil {
		return nil, 0, err
	}
	if len(facetResults) == 0 {
		return []*ProductResult{}, 0, nil
	}
	var total int64
	if len(facetResults[0].Metadata) > 0 {
		total = facetResults[0].Metadata[0].TotalCount
	}
	var results []*ProductResult
	for _, d := range facetResults[0].Data {
		results = append(results, bsonMToProductResult(d))
	}
	return results, total, nil
}

func bsonMToProductResult(m bson.M) *ProductResult {
	return &ProductResult{
		ProductID:    bsonMToInt32P(m, "_id"),
		MerchantID:   bsonMToInt32P(m, "merchant_id"),
		CategoryID:   bsonMToInt32P(m, "category_id"),
		Name:         bsonMToStrP(m, "name"),
		Description:  bsonMToStrPtrP(m, "description"),
		Price:        bsonMToInt32P(m, "price"),
		CountInStock: bsonMToInt32P(m, "count_in_stock"),
		Brand:        bsonMToStrPtrP(m, "brand"),
		Weight:       bsonMToInt32PtrP(m, "weight"),
		Rating:       bsonMToFloat64PtrP(m, "rating"),
		SlugProduct:  bsonMToStrPtrP(m, "slug_product"),
		ImageProduct: bsonMToStrPtrP(m, "image_product"),
		CreatedAt:    bsonMToStrPtrP(m, "created_at"),
		UpdatedAt:    bsonMToStrPtrP(m, "updated_at"),
		DeletedAt:    bsonMToStrPtrP(m, "deleted_at"),
	}
}

func bsonMToInt32P(m bson.M, key string) int32 {
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

func bsonMToInt32PtrP(m bson.M, key string) *int32 {
	v, ok := m[key]
	if !ok || v == nil {
		return nil
	}
	switch t := v.(type) {
	case int32:
		return &t
	case int64:
		i := int32(t)
		return &i
	case float64:
		i := int32(t)
		return &i
	default:
		return nil
	}
}

func bsonMToFloat64PtrP(m bson.M, key string) *float64 {
	v, ok := m[key]
	if !ok || v == nil {
		return nil
	}
	switch t := v.(type) {
	case float64:
		return &t
	case int32:
		f := float64(t)
		return &f
	case int64:
		f := float64(t)
		return &f
	default:
		return nil
	}
}

func bsonMToStrP(m bson.M, key string) string {
	v, ok := m[key]
	if !ok {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func bsonMToStrPtrP(m bson.M, key string) *string {
	v, ok := m[key]
	if !ok || v == nil {
		return nil
	}
	switch t := v.(type) {
	case string:
		return &t
	case time.Time:
		s := t.Format("2006-01-02 15:04:05")
		return &s
	case bson.DateTime:
		s := t.Time().Format("2006-01-02 15:04:05")
		return &s
	default:
		return nil
	}
}
