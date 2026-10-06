package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	shippingaddress_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/shipping_address_errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type shippingAddressCommandRepository struct {
	db *mongo.Database
}

func NewShippingAddressCommandRepository(db *mongo.Database) *shippingAddressCommandRepository {
	return &shippingAddressCommandRepository{db: db}
}

func (r *shippingAddressCommandRepository) Create(ctx context.Context, request *requests.CreateShippingAddressRequest) (*models.ShippingAddress, error) {
	var orderID int32
	if request.OrderID != nil {
		orderID = int32(*request.OrderID)
	}
	id, err := database.GetNextID(ctx, r.db, "shipping_addresses")
	if err != nil {
		return nil, shippingaddress_errors.ErrCreateShippingAddress
	}
	now := time.Now()
	addr := &models.ShippingAddress{
		ShippingAddressID: id,
		OrderID:           orderID,
		Alamat:            request.Alamat,
		Provinsi:          request.Provinsi,
		Kota:              request.Kota,
		Negara:            request.Negara,
		Courier:           request.Courier,
		ShippingMethod:    request.ShippingMethod,
		ShippingCost:      float64(request.ShippingCost),
		CreatedAt:         &now,
		UpdatedAt:         &now,
	}
	if _, err := r.db.Collection("shipping_addresses").InsertOne(ctx, addr); err != nil {
		return nil, shippingaddress_errors.ErrCreateShippingAddress
	}
	return addr, nil
}

func (r *shippingAddressCommandRepository) Update(ctx context.Context, request *requests.UpdateShippingAddressRequest) (*models.ShippingAddress, error) {
	coll := r.db.Collection("shipping_addresses")
	var addr models.ShippingAddress
	if err := coll.FindOne(ctx, bson.M{"_id": *request.ShippingID}).Decode(&addr); err != nil {
		return nil, shippingaddress_errors.ErrShippingAddressNotFound
	}
	now := time.Now()
	update := bson.M{
		"alamat":          request.Alamat,
		"provinsi":        request.Provinsi,
		"kota":            request.Kota,
		"negara":          request.Negara,
		"courier":         request.Courier,
		"shipping_method": request.ShippingMethod,
		"shipping_cost":   float64(request.ShippingCost),
		"updated_at":      now,
	}
	if _, err := coll.UpdateOne(ctx, bson.M{"_id": *request.ShippingID}, bson.M{"$set": update}); err != nil {
		return nil, shippingaddress_errors.ErrUpdateShippingAddress
	}
	if err := coll.FindOne(ctx, bson.M{"_id": *request.ShippingID}).Decode(&addr); err != nil {
		return nil, shippingaddress_errors.ErrUpdateShippingAddress
	}
	return &addr, nil
}

func (r *shippingAddressCommandRepository) Trash(ctx context.Context, shippingID int) (*models.ShippingAddress, error) {
	coll := r.db.Collection("shipping_addresses")
	var addr models.ShippingAddress
	if err := coll.FindOne(ctx, bson.M{"_id": shippingID}).Decode(&addr); err != nil {
		return nil, shippingaddress_errors.ErrShippingAddressNotFound
	}
	now := time.Now()
	if _, err := coll.UpdateOne(ctx, bson.M{"_id": shippingID}, bson.M{"$set": bson.M{"deleted_at": now}}); err != nil {
		return nil, shippingaddress_errors.ErrTrashShippingAddress
	}
	addr.DeletedAt = &now
	return &addr, nil
}

func (r *shippingAddressCommandRepository) Restore(ctx context.Context, shippingID int) (*models.ShippingAddress, error) {
	coll := r.db.Collection("shipping_addresses")
	res, err := coll.UpdateOne(ctx, bson.M{"_id": shippingID, "deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	if err != nil {
		return nil, shippingaddress_errors.ErrRestoreShippingAddress
	}
	if res.MatchedCount == 0 {
		return nil, shippingaddress_errors.ErrShippingAddressNotFound
	}
	var addr models.ShippingAddress
	if err := coll.FindOne(ctx, bson.M{"_id": shippingID}).Decode(&addr); err != nil {
		return nil, shippingaddress_errors.ErrRestoreShippingAddress
	}
	return &addr, nil
}

func (r *shippingAddressCommandRepository) DeletePermanent(ctx context.Context, shippingID int) (bool, error) {
	res, err := r.db.Collection("shipping_addresses").DeleteOne(ctx, bson.M{"_id": shippingID})
	if err != nil {
		return false, shippingaddress_errors.ErrDeleteShippingAddressPermanent
	}
	if res.DeletedCount == 0 {
		return false, shippingaddress_errors.ErrShippingAddressNotFound
	}
	return true, nil
}

func (r *shippingAddressCommandRepository) DeleteByOrderIDPermanent(ctx context.Context, orderID int) (bool, error) {
	_, err := r.db.Collection("shipping_addresses").DeleteMany(ctx, bson.M{"order_id": orderID})
	if err != nil {
		return false, shippingaddress_errors.ErrDeleteShippingAddressPermanent
	}
	return true, nil
}

func (r *shippingAddressCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	_, err := r.db.Collection("shipping_addresses").UpdateMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	if err != nil {
		return false, shippingaddress_errors.ErrRestoreAllShippingAddresses
	}
	return true, nil
}

func (r *shippingAddressCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	_, err := r.db.Collection("shipping_addresses").DeleteMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, shippingaddress_errors.ErrDeleteAllPermanentShippingAddress
	}
	return true, nil
}
