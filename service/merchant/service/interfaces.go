package service

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-grpc-merchant/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

type MerchantQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*repository.MerchantResult, *int, error)
	FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*repository.MerchantResult, *int, error)
	FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*repository.MerchantResult, *int, error)
	FindByID(ctx context.Context, merchantID int) (*repository.MerchantResult, error)
}

type MerchantCommandService interface {
	Create(ctx context.Context, request *requests.CreateMerchantRequest) (*models.Merchant, error)
	Update(ctx context.Context, request *requests.UpdateMerchantRequest) (*models.Merchant, error)
	Trash(ctx context.Context, merchantID int) (*models.Merchant, error)
	Restore(ctx context.Context, merchantID int) (*models.Merchant, error)
	DeletePermanent(ctx context.Context, merchantID int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
	UpdateMerchantStatus(ctx context.Context, request *requests.UpdateMerchantStatusRequest) (*models.Merchant, error)
}

type MerchantDocumentQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*repository.MerchantDocumentResult, *int, error)
	FindActive(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*repository.MerchantDocumentResult, *int, error)
	FindTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*repository.MerchantDocumentResult, *int, error)
	FindByID(ctx context.Context, documentID int) (*repository.MerchantDocumentResult, error)
}

type MerchantDocumentCommandService interface {
	Create(ctx context.Context, request *requests.CreateMerchantDocumentRequest) (*models.MerchantDocument, error)
	Update(ctx context.Context, request *requests.UpdateMerchantDocumentRequest) (*models.MerchantDocument, error)
	UpdateStatus(ctx context.Context, request *requests.UpdateMerchantDocumentStatusRequest) (*models.MerchantDocument, error)
	Trash(ctx context.Context, documentID int) (*models.MerchantDocument, error)
	Restore(ctx context.Context, documentID int) (*models.MerchantDocument, error)
	DeletePermanent(ctx context.Context, documentID int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
