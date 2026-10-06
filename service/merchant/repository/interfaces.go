package repository

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MerchantResult struct {
	MerchantID   int32
	UserID       int32
	Name         string
	Description  *string
	Address      *string
	ContactEmail *string
	ContactPhone *string
	Status       string
	CreatedAt    *string
	UpdatedAt    *string
	DeletedAt    *string
	TotalCount   int64
}

type MerchantDocumentResult struct {
	DocumentID   int32
	MerchantID   int32
	DocumentType string
	DocumentUrl  string
	Status       string
	Note         *string
	UploadedAt   *string
	CreatedAt    *string
	UpdatedAt    *string
	DeletedAt    *string
	TotalCount   int64
}

type MerchantDocumentQueryRepository interface {
	FindAll(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*MerchantDocumentResult, *int, error)
	FindByID(ctx context.Context, id int) (*MerchantDocumentResult, error)
	FindActive(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*MerchantDocumentResult, *int, error)
	FindTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*MerchantDocumentResult, *int, error)
}

type MerchantDocumentCommandRepository interface {
	Create(ctx context.Context, request *requests.CreateMerchantDocumentRequest) (*models.MerchantDocument, error)
	CreateInTx(ctx context.Context, tx *mongo.Database, request *requests.CreateMerchantDocumentRequest) (*models.MerchantDocument, error)
	Update(ctx context.Context, request *requests.UpdateMerchantDocumentRequest) (*models.MerchantDocument, error)
	UpdateStatus(ctx context.Context, request *requests.UpdateMerchantDocumentStatusRequest) (*models.MerchantDocument, error)
	UpdateStatusInTx(ctx context.Context, tx *mongo.Database, request *requests.UpdateMerchantDocumentStatusRequest) (*models.MerchantDocument, error)
	Trash(ctx context.Context, merchant_document_id int) (*models.MerchantDocument, error)
	Restore(ctx context.Context, merchant_document_id int) (*models.MerchantDocument, error)
	DeletePermanent(ctx context.Context, merchant_document_id int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}

type MerchantQueryRepository interface {
	FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantResult, error)
	FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantResult, error)
	FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantResult, error)
	FindByID(ctx context.Context, user_id int) (*MerchantResult, error)
}

type MerchantCommandRepository interface {
	Create(ctx context.Context, request *requests.CreateMerchantRequest) (*models.Merchant, error)
	CreateInTx(ctx context.Context, tx *mongo.Database, request *requests.CreateMerchantRequest) (*models.Merchant, error)
	Update(ctx context.Context, request *requests.UpdateMerchantRequest) (*models.Merchant, error)
	Trash(ctx context.Context, merchant_id int) (*models.Merchant, error)
	Restore(ctx context.Context, merchant_id int) (*models.Merchant, error)
	DeletePermanent(ctx context.Context, merchant_id int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
	UpdateStatus(ctx context.Context, request *requests.UpdateMerchantStatusRequest) (*models.Merchant, error)
	UpdateStatusInTx(ctx context.Context, tx *mongo.Database, request *requests.UpdateMerchantStatusRequest) (*models.Merchant, error)
}
