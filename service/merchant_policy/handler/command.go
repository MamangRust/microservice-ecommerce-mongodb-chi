package handler

import (
	"context"
	pbmerchant "github.com/MamangRust/microservice-ecommerce-grpc-pb/merchant"
	pbmerchant_policy "github.com/MamangRust/microservice-ecommerce-grpc-pb/merchant_policy"

	"github.com/MamangRust/microservice-ecommerce-grpc-merchant_policy/service"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

type merchantPolicyCommandHandler struct {
	pbmerchant_policy.UnimplementedMerchantPolicyCommandServiceServer
	merchantPolicyService service.MerchantPoliciesCommandService
	logger                logger.LoggerInterface
}

func NewMerchantPolicyCommandHandler(
	merchantPolicyService service.MerchantPoliciesCommandService,
	logger logger.LoggerInterface,
) pbmerchant_policy.MerchantPolicyCommandServiceServer {
	return &merchantPolicyCommandHandler{
		merchantPolicyService: merchantPolicyService,
		logger:                logger,
	}
}

func (h *merchantPolicyCommandHandler) Create(ctx context.Context, req *pbmerchant_policy.CreateMerchantPoliciesRequest) (*pbmerchant_policy.ApiResponseMerchantPolicies, error) {
	policy, err := h.merchantPolicyService.Create(ctx, &requests.CreateMerchantPolicyRequest{
		MerchantID:  int(req.GetMerchantId()),
		PolicyType:  req.GetPolicyType(),
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
	})

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return mapToSingleResponseFromModel(policy), nil
}

func (h *merchantPolicyCommandHandler) Update(ctx context.Context, req *pbmerchant_policy.UpdateMerchantPoliciesRequest) (*pbmerchant_policy.ApiResponseMerchantPolicies, error) {
	id := int(req.GetMerchantPolicyId())
	policy, err := h.merchantPolicyService.Update(ctx, &requests.UpdateMerchantPolicyRequest{
		MerchantPolicyID: &id,
		PolicyType:       req.GetPolicyType(),
		Title:            req.GetTitle(),
		Description:      req.GetDescription(),
	})

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return mapToSingleResponseFromModel(policy), nil
}

func (h *merchantPolicyCommandHandler) TrashedMerchantPolicies(ctx context.Context, req *pbmerchant_policy.FindByIdMerchantPoliciesRequest) (*pbmerchant_policy.ApiResponseMerchantPoliciesDeleteAt, error) {
	policy, err := h.merchantPolicyService.Trash(ctx, int(req.GetId()))

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return mapToSingleDeleteAtResponse(policy), nil
}

func (h *merchantPolicyCommandHandler) RestoreMerchantPolicies(ctx context.Context, req *pbmerchant_policy.FindByIdMerchantPoliciesRequest) (*pbmerchant_policy.ApiResponseMerchantPoliciesDeleteAt, error) {
	policy, err := h.merchantPolicyService.Restore(ctx, int(req.GetId()))

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return mapToSingleDeleteAtResponse(policy), nil
}

func (h *merchantPolicyCommandHandler) DeleteMerchantPoliciesPermanent(ctx context.Context, req *pbmerchant_policy.FindByIdMerchantPoliciesRequest) (*pbmerchant.ApiResponseMerchantDelete, error) {
	_, err := h.merchantPolicyService.DeletePermanent(ctx, int(req.GetId()))

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchant.ApiResponseMerchantDelete{
		Status:  "success",
		Message: "Successfully deleted merchant policy permanently",
	}, nil
}

func (h *merchantPolicyCommandHandler) RestoreAllMerchantPolicies(ctx context.Context, req *emptypb.Empty) (*pbmerchant.ApiResponseMerchantAll, error) {
	_, err := h.merchantPolicyService.RestoreAll(ctx)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchant.ApiResponseMerchantAll{
		Status:  "success",
		Message: "Successfully restored all merchant policies",
	}, nil
}

func (h *merchantPolicyCommandHandler) DeleteAllMerchantPoliciesPermanent(ctx context.Context, req *emptypb.Empty) (*pbmerchant.ApiResponseMerchantAll, error) {
	_, err := h.merchantPolicyService.DeleteAll(ctx)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchant.ApiResponseMerchantAll{
		Status:  "success",
		Message: "Successfully deleted all merchant policies permanently",
	}, nil
}
