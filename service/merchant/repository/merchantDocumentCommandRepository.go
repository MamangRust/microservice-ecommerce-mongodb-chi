package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	merchant_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/merchant"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type merchantDocumentCommandRepository struct {
	db *mongo.Database
}

func NewMerchantDocumentCommandRepository(db *mongo.Database) *merchantDocumentCommandRepository {
	return &merchantDocumentCommandRepository{db: db}
}

func (r *merchantDocumentCommandRepository) Create(ctx context.Context, request *requests.CreateMerchantDocumentRequest) (*models.MerchantDocument, error) {
	return r.create(ctx, r.db, request)
}

func (r *merchantDocumentCommandRepository) CreateInTx(ctx context.Context, tx *mongo.Database, request *requests.CreateMerchantDocumentRequest) (*models.MerchantDocument, error) {
	return r.create(ctx, tx, request)
}

func (r *merchantDocumentCommandRepository) create(ctx context.Context, db *mongo.Database, request *requests.CreateMerchantDocumentRequest) (*models.MerchantDocument, error) {
	id, err := database.GetNextID(ctx, db, "merchant_documents")
	if err != nil {
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}
	now := time.Now()
	doc := &models.MerchantDocument{
		DocumentID:   id,
		MerchantID:   int32(request.MerchantID),
		DocumentType: request.DocumentType,
		DocumentUrl:  request.DocumentUrl,
		Status:       "pending",
		Note:         stringPtr(""),
		UploadedAt:   &now,
		CreatedAt:    &now,
		UpdatedAt:    &now,
	}
	if _, err := db.Collection("merchant_documents").InsertOne(ctx, doc); err != nil {
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}
	return doc, nil
}

func (r *merchantDocumentCommandRepository) Update(ctx context.Context, request *requests.UpdateMerchantDocumentRequest) (*models.MerchantDocument, error) {
	coll := r.db.Collection("merchant_documents")
	var doc models.MerchantDocument
	if err := coll.FindOne(ctx, bson.M{"_id": *request.DocumentID}).Decode(&doc); err != nil {
		return nil, merchant_errors.ErrMerchantNotFound
	}
	now := time.Now()
	update := bson.M{
		"document_type": request.DocumentType,
		"document_url":  request.DocumentUrl,
		"status":        request.Status,
		"note":          stringPtr(request.Note),
		"updated_at":    now,
	}
	if _, err := coll.UpdateOne(ctx, bson.M{"_id": *request.DocumentID}, bson.M{"$set": update}); err != nil {
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}
	if err := coll.FindOne(ctx, bson.M{"_id": *request.DocumentID}).Decode(&doc); err != nil {
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}
	return &doc, nil
}

func (r *merchantDocumentCommandRepository) UpdateStatus(ctx context.Context, request *requests.UpdateMerchantDocumentStatusRequest) (*models.MerchantDocument, error) {
	coll := r.db.Collection("merchant_documents")
	var doc models.MerchantDocument
	if err := coll.FindOne(ctx, bson.M{"_id": *request.DocumentID}).Decode(&doc); err != nil {
		return nil, merchant_errors.ErrMerchantNotFound
	}
	now := time.Now()
	update := bson.M{
		"status":     request.Status,
		"note":       stringPtr(request.Note),
		"updated_at": now,
	}
	if _, err := coll.UpdateOne(ctx, bson.M{"_id": *request.DocumentID}, bson.M{"$set": update}); err != nil {
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}
	if err := coll.FindOne(ctx, bson.M{"_id": *request.DocumentID}).Decode(&doc); err != nil {
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}
	return &doc, nil
}

func (r *merchantDocumentCommandRepository) UpdateStatusInTx(ctx context.Context, tx *mongo.Database, request *requests.UpdateMerchantDocumentStatusRequest) (*models.MerchantDocument, error) {
	coll := tx.Collection("merchant_documents")
	var doc models.MerchantDocument
	if err := coll.FindOne(ctx, bson.M{"_id": *request.DocumentID}).Decode(&doc); err != nil {
		return nil, merchant_errors.ErrMerchantNotFound
	}
	now := time.Now()
	update := bson.M{
		"status":     request.Status,
		"note":       stringPtr(request.Note),
		"updated_at": now,
	}
	if _, err := coll.UpdateOne(ctx, bson.M{"_id": *request.DocumentID}, bson.M{"$set": update}); err != nil {
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}
	if err := coll.FindOne(ctx, bson.M{"_id": *request.DocumentID}).Decode(&doc); err != nil {
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}
	return &doc, nil
}

func (r *merchantDocumentCommandRepository) Trash(ctx context.Context, documentID int) (*models.MerchantDocument, error) {
	coll := r.db.Collection("merchant_documents")
	var doc models.MerchantDocument
	if err := coll.FindOne(ctx, bson.M{"_id": documentID}).Decode(&doc); err != nil {
		return nil, merchant_errors.ErrMerchantNotFound
	}
	now := time.Now()
	if _, err := coll.UpdateOne(ctx, bson.M{"_id": documentID}, bson.M{"$set": bson.M{"deleted_at": now}}); err != nil {
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}
	doc.DeletedAt = &now
	return &doc, nil
}

func (r *merchantDocumentCommandRepository) Restore(ctx context.Context, documentID int) (*models.MerchantDocument, error) {
	coll := r.db.Collection("merchant_documents")
	res, err := coll.UpdateOne(ctx, bson.M{"_id": documentID, "deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	if err != nil {
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}
	if res.MatchedCount == 0 {
		return nil, merchant_errors.ErrMerchantNotFound
	}
	var doc models.MerchantDocument
	if err := coll.FindOne(ctx, bson.M{"_id": documentID}).Decode(&doc); err != nil {
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}
	return &doc, nil
}

func (r *merchantDocumentCommandRepository) DeletePermanent(ctx context.Context, documentID int) (bool, error) {
	res, err := r.db.Collection("merchant_documents").DeleteOne(ctx, bson.M{"_id": documentID})
	if err != nil {
		return false, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}
	if res.DeletedCount == 0 {
		return false, merchant_errors.ErrMerchantNotFound
	}
	return true, nil
}

func (r *merchantDocumentCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	_, err := r.db.Collection("merchant_documents").UpdateMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, bson.M{"$set": bson.M{"deleted_at": nil}})
	if err != nil {
		return false, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}
	return true, nil
}

func (r *merchantDocumentCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	_, err := r.db.Collection("merchant_documents").DeleteMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}
	return true, nil
}
