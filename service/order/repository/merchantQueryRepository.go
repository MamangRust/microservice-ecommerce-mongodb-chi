package repository

import (
	"context"
	pbmerchant "github.com/MamangRust/microservice-ecommerce-grpc-pb/merchant"

	dto "github.com/MamangRust/microservice-ecommerce-grpc-order/dto"
	merchant_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/merchant"
)

type merchantQueryRepository struct {
	client pbmerchant.MerchantQueryServiceClient
}

func NewMerchantQueryRepository(client pbmerchant.MerchantQueryServiceClient) *merchantQueryRepository {
	return &merchantQueryRepository{
		client: client,
	}
}

func (r *merchantQueryRepository) FindByID(ctx context.Context, merchant_id int) (*dto.GetMerchantByIDRow, error) {
	res, err := r.client.FindById(ctx, &pbmerchant.FindByIdMerchantRequest{Id: int32(merchant_id)})
	if err != nil {
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}

	return &dto.GetMerchantByIDRow{
		MerchantID:   res.Data.Id,
		UserID:       res.Data.UserId,
		Name:         res.Data.Name,
		Description:  &res.Data.Description,
		Address:      &res.Data.Address,
		ContactEmail: &res.Data.ContactEmail,
		ContactPhone: &res.Data.ContactPhone,
		Status:       res.Data.Status,
	}, nil
}
