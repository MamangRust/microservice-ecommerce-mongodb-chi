package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

type MerchantBusinessResult struct {
	MerchantBusinessInfoID int32
	MerchantID             int32
	BusinessType           *string
	TaxID                  *string
	EstablishedYear        *int32
	NumberOfEmployees      *int32
	WebsiteUrl             *string
	CreatedAt              *time.Time
	UpdatedAt              *time.Time
	DeletedAt              *time.Time
	TotalCount             int64
}

type MerchantQueryRepository interface {
	FindByID(ctx context.Context, id int) (string, error)
}

type MerchantBusinessQueryRepository interface {
	FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantBusinessResult, error)
	FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantBusinessResult, error)
	FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantBusinessResult, error)
	FindByID(ctx context.Context, id int) (*MerchantBusinessResult, error)
}

type MerchantBusinessCommandRepository interface {
	Create(ctx context.Context, request *requests.CreateMerchantBusinessInformationRequest) (*models.MerchantBusinessInformation, error)
	Update(ctx context.Context, request *requests.UpdateMerchantBusinessInformationRequest) (*models.MerchantBusinessInformation, error)
	Trash(ctx context.Context, id int) (*models.MerchantBusinessInformation, error)
	Restore(ctx context.Context, id int) (*models.MerchantBusinessInformation, error)
	DeletePermanent(ctx context.Context, id int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
