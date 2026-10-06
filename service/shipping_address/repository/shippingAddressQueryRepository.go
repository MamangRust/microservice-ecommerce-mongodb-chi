package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/errors"
	shippingaddress_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/shipping_address_errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type shippingAddressQueryRepository struct {
	db *mongo.Database
}

func NewShippingAddressQueryRepository(db *mongo.Database) *shippingAddressQueryRepository {
	return &shippingAddressQueryRepository{db: db}
}

func (r *shippingAddressQueryRepository) FindAll(ctx context.Context, req *requests.FindAllShippingAddress) ([]*ShippingAddressResult, error) {
	return r.runQuery(ctx, bson.M{"deleted_at": bson.M{"$eq": nil}}, req, shippingaddress_errors.ErrFindAllShippingAddress)
}

func (r *shippingAddressQueryRepository) FindActive(ctx context.Context, req *requests.FindAllShippingAddress) ([]*ShippingAddressResult, error) {
	return r.runQuery(ctx, bson.M{"deleted_at": bson.M{"$eq": nil}}, req, shippingaddress_errors.ErrFindActiveShippingAddress)
}

func (r *shippingAddressQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllShippingAddress) ([]*ShippingAddressResult, error) {
	return r.runQuery(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, req, shippingaddress_errors.ErrFindTrashedShippingAddress)
}

func (r *shippingAddressQueryRepository) FindByID(ctx context.Context, shippingID int) (*ShippingAddressResult, error) {
	var d bson.M
	err := r.db.Collection("shipping_addresses").FindOne(ctx, bson.M{"_id": shippingID, "deleted_at": bson.M{"$eq": nil}}).Decode(&d)
	if err != nil {
		return nil, shippingaddress_errors.ErrFindShippingAddressByID
	}
	return bsonMToShippingAddressResult(d), nil
}

func (r *shippingAddressQueryRepository) FindByOrder(ctx context.Context, orderID int) (*ShippingAddressResult, error) {
	var d bson.M
	err := r.db.Collection("shipping_addresses").FindOne(ctx, bson.M{"order_id": orderID}).Decode(&d)
	if err != nil {
		return nil, shippingaddress_errors.ErrFindShippingAddressByOrder
	}
	return bsonMToShippingAddressResult(d), nil
}

func (r *shippingAddressQueryRepository) runQuery(ctx context.Context, filter bson.M, req *requests.FindAllShippingAddress, errFactory *errors.AppError) ([]*ShippingAddressResult, error) {
	offset := (req.Page - 1) * req.PageSize
	if req.Search != "" {
		filter["$or"] = []bson.M{
			{"alamat": bson.M{"$regex": req.Search, "$options": "i"}},
			{"provinsi": bson.M{"$regex": req.Search, "$options": "i"}},
		}
	}
	pipeline := []bson.M{
		{"$match": filter},
		{"$sort": bson.M{"shipping_address_id": -1}},
		{"$facet": bson.M{
			"metadata": []bson.M{{"$count": "total_count"}},
			"data":     []bson.M{{"$skip": offset}, {"$limit": req.PageSize}},
		}},
	}
	cursor, err := r.db.Collection("shipping_addresses").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, errFactory.WithInternal(err)
	}
	defer cursor.Close(ctx)
	var facetResults []struct {
		Metadata []struct {
			TotalCount int64 `bson:"total_count"`
		} `bson:"metadata"`
		Data []bson.M `bson:"data"`
	}
	if err := cursor.All(ctx, &facetResults); err != nil {
		return nil, errFactory.WithInternal(err)
	}
	if len(facetResults) == 0 {
		return []*ShippingAddressResult{}, nil
	}
	var results []*ShippingAddressResult
	for _, d := range facetResults[0].Data {
		item := bsonMToShippingAddressResult(d)
		if len(facetResults[0].Metadata) > 0 {
			item.TotalCount = facetResults[0].Metadata[0].TotalCount
		}
		results = append(results, item)
	}
	return results, nil
}

func bsonMToShippingAddressResult(m bson.M) *ShippingAddressResult {
	return &ShippingAddressResult{
		ShippingAddressID: bsonMToInt32SA(m, "_id"),
		OrderID:           bsonMToInt32SA(m, "order_id"),
		Alamat:            bsonMToStrSA(m, "alamat"),
		Provinsi:          bsonMToStrSA(m, "provinsi"),
		Kota:              bsonMToStrSA(m, "kota"),
		Negara:            bsonMToStrSA(m, "negara"),
		Courier:           bsonMToStrSA(m, "courier"),
		ShippingMethod:    bsonMToStrSA(m, "shipping_method"),
		ShippingCost:      bsonMToFloat64SA(m, "shipping_cost"),
		CreatedAt:         bsonMToTimeStrPtrSA(m, "created_at"),
		UpdatedAt:         bsonMToTimeStrPtrSA(m, "updated_at"),
		DeletedAt:         bsonMToTimeStrPtrSA(m, "deleted_at"),
	}
}

func bsonMToInt32SA(m bson.M, key string) int32 {
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

func bsonMToStrSA(m bson.M, key string) string {
	v, ok := m[key]
	if !ok {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func bsonMToFloat64SA(m bson.M, key string) float64 {
	v, ok := m[key]
	if !ok {
		return 0
	}
	switch t := v.(type) {
	case float64:
		return t
	case int32:
		return float64(t)
	case int64:
		return float64(t)
	default:
		return 0
	}
}

func bsonMToTimeStrPtrSA(m bson.M, key string) *string {
	v, ok := m[key]
	if !ok || v == nil {
		return nil
	}
	switch t := v.(type) {
	case time.Time:
		s := t.Format("2006-01-02 15:04:05.000")
		return &s
	case bson.DateTime:
		s := t.Time().Format("2006-01-02 15:04:05.000")
		return &s
	case string:
		return &t
	default:
		return nil
	}
}
