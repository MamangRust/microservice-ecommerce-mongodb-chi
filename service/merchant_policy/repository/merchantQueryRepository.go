package repository

import (
	"context"
	pbmerchant "github.com/MamangRust/microservice-ecommerce-grpc-pb/merchant"

	merchant_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/merchant"
)

type merchantQueryRepository struct {
	client pbmerchant.MerchantQueryServiceClient
}

func NewMerchantQueryRepository(client pbmerchant.MerchantQueryServiceClient) MerchantQueryRepository {
	return &merchantQueryRepository{
		client: client,
	}
}

func (r *merchantQueryRepository) FindByID(ctx context.Context, merchant_id int) (string, error) {
	res, err := r.client.FindById(ctx, &pbmerchant.FindByIdMerchantRequest{Id: int32(merchant_id)})
	if err != nil {
		return "", merchant_errors.ErrMerchantNotFound.WithInternal(err)
	}
	return res.Data.Name, nil
}
