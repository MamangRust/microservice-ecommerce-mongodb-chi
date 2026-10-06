package repository

import (
	"context"
	pbmerchant "github.com/MamangRust/microservice-ecommerce-grpc-pb/merchant"

	dto "github.com/MamangRust/microservice-ecommerce-grpc-transaction/dto"
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
		// pertahankan status gRPC dari dependency service (NotFound -> 404, dst)
		return nil, err
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
