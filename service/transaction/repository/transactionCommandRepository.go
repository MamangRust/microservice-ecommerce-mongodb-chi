package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/errors/transaction_errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type transactionCommandRepository struct {
	db *mongo.Database
}

func NewTransactionCommandRepository(db *mongo.Database) *transactionCommandRepository {
	return &transactionCommandRepository{db: db}
}

func toTransactionResult(t *models.Transaction) *TransactionResult {
	if t == nil {
		return nil
	}
	return &TransactionResult{
		TransactionID: t.TransactionID,
		OrderID:       t.OrderID,
		MerchantID:    t.MerchantID,
		PaymentMethod: t.PaymentMethod,
		Amount:        t.Amount,
		PaymentStatus: t.PaymentStatus,
		CreatedAt:     t.CreatedAt,
		UpdatedAt:     t.UpdatedAt,
		DeletedAt:     t.DeletedAt,
	}
}

func (r *transactionCommandRepository) Create(ctx context.Context, request *requests.CreateTransactionRequest) (*TransactionResult, error) {
	return r.create(ctx, r.db, request)
}

func (r *transactionCommandRepository) CreateInTx(ctx context.Context, tx *mongo.Database, request *requests.CreateTransactionRequest) (*TransactionResult, error) {
	return r.create(ctx, tx, request)
}

func (r *transactionCommandRepository) create(ctx context.Context, db *mongo.Database, request *requests.CreateTransactionRequest) (*TransactionResult, error) {
	id, err := database.GetNextID(ctx, db, "transactions")
	if err != nil {
		return nil, transaction_errors.ErrCreateTransaction.WithInternal(err)
	}
	now := time.Now()
	txn := &models.Transaction{
		TransactionID: id,
		OrderID:       int32(request.OrderID),
		MerchantID:    int32(request.MerchantID),
		PaymentMethod: request.PaymentMethod,
		Amount:        int32(request.Amount),
		PaymentStatus: *request.PaymentStatus,
		CreatedAt:     &now,
		UpdatedAt:     &now,
	}
	if _, err := db.Collection("transactions").InsertOne(ctx, txn); err != nil {
		return nil, transaction_errors.ErrCreateTransaction.WithInternal(err)
	}
	return toTransactionResult(txn), nil
}

func (r *transactionCommandRepository) Update(ctx context.Context, request *requests.UpdateTransactionRequest) (*TransactionResult, error) {
	coll := r.db.Collection("transactions")
	var txn models.Transaction
	if err := coll.FindOne(ctx, bson.M{"_id": *request.TransactionID, "deleted_at": bson.M{"$eq": nil}}).Decode(&txn); err != nil {
		return nil, transaction_errors.ErrTransactionNotFound
	}
	now := time.Now()
	update := bson.M{
		"merchant_id":    int32(request.MerchantID),
		"payment_method": request.PaymentMethod,
		"amount":         int32(request.Amount),
		"order_id":       int32(request.OrderID),
		"payment_status": *request.PaymentStatus,
		"updated_at":     now,
	}
	if _, err := coll.UpdateOne(ctx, bson.M{"_id": *request.TransactionID}, bson.M{"$set": update}); err != nil {
		return nil, transaction_errors.ErrUpdateTransaction.WithInternal(err)
	}
	if err := coll.FindOne(ctx, bson.M{"_id": *request.TransactionID}).Decode(&txn); err != nil {
		return nil, transaction_errors.ErrUpdateTransaction.WithInternal(err)
	}
	return toTransactionResult(&txn), nil
}

func (r *transactionCommandRepository) Trash(ctx context.Context, transactionID int) (*TransactionResult, error) {
	coll := r.db.Collection("transactions")
	var txn models.Transaction
	if err := coll.FindOne(ctx, bson.M{"_id": transactionID, "deleted_at": bson.M{"$eq": nil}}).Decode(&txn); err != nil {
		return nil, transaction_errors.ErrTransactionNotFound
	}
	now := time.Now()
	if _, err := coll.UpdateOne(ctx, bson.M{"_id": transactionID}, bson.M{"$set": bson.M{"deleted_at": now}}); err != nil {
		return nil, transaction_errors.ErrTrashTransaction.WithInternal(err)
	}
	txn.DeletedAt = &now
	return toTransactionResult(&txn), nil
}

func (r *transactionCommandRepository) Restore(ctx context.Context, transactionID int) (*TransactionResult, error) {
	coll := r.db.Collection("transactions")
	res, err := coll.UpdateOne(ctx, bson.M{"_id": transactionID, "deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	if err != nil {
		return nil, transaction_errors.ErrRestoreTransaction.WithInternal(err)
	}
	if res.MatchedCount == 0 {
		return nil, transaction_errors.ErrTransactionNotFound
	}
	var txn models.Transaction
	if err := coll.FindOne(ctx, bson.M{"_id": transactionID}).Decode(&txn); err != nil {
		return nil, transaction_errors.ErrRestoreTransaction.WithInternal(err)
	}
	return toTransactionResult(&txn), nil
}

func (r *transactionCommandRepository) DeletePermanent(ctx context.Context, transactionID int) (bool, error) {
	res, err := r.db.Collection("transactions").DeleteOne(ctx, bson.M{"_id": transactionID})
	if err != nil {
		return false, transaction_errors.ErrDeleteTransactionPermanently.WithInternal(err)
	}
	if res.DeletedCount == 0 {
		return false, transaction_errors.ErrTransactionNotFound
	}
	return true, nil
}

func (r *transactionCommandRepository) DeleteByOrderIDPermanent(ctx context.Context, orderID int) (bool, error) {
	_, err := r.db.Collection("transactions").DeleteMany(ctx, bson.M{"order_id": orderID})
	if err != nil {
		return false, transaction_errors.ErrDeleteTransactionPermanently.WithInternal(err)
	}
	return true, nil
}

func (r *transactionCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	_, err := r.db.Collection("transactions").UpdateMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	if err != nil {
		return false, transaction_errors.ErrRestoreAllTransactions.WithInternal(err)
	}
	return true, nil
}

func (r *transactionCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	_, err := r.db.Collection("transactions").DeleteMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, transaction_errors.ErrDeleteAllTransactionPermanent.WithInternal(err)
	}
	return true, nil
}
