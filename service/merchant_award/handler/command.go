package handler

import (
	"context"
	pbmerchant "github.com/MamangRust/microservice-ecommerce-grpc-pb/merchant"
	pbmerchant_award "github.com/MamangRust/microservice-ecommerce-grpc-pb/merchant_award"

	"github.com/MamangRust/microservice-ecommerce-grpc-merchant_award/service"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/errors"
	merchantaward_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/merchant_award"
	"google.golang.org/protobuf/types/known/emptypb"
)

type merchantAwardCommandHandler struct {
	pbmerchant_award.UnimplementedMerchantAwardCommandServiceServer
	merchantAwardCommand service.MerchantAwardCommandService
	logger               logger.LoggerInterface
}

func NewMerchantAwardCommandHandler(svc service.MerchantAwardCommandService, logger logger.LoggerInterface) MerchantAwardCommandHandler {
	return &merchantAwardCommandHandler{
		merchantAwardCommand: svc,
		logger:               logger,
	}
}

func (s *merchantAwardCommandHandler) Create(ctx context.Context, request *pbmerchant_award.CreateMerchantAwardRequest) (*pbmerchant_award.ApiResponseMerchantAward, error) {
	req := &requests.CreateMerchantCertificationOrAwardRequest{
		MerchantID:     int(request.GetMerchantId()),
		Title:          request.GetTitle(),
		Description:    request.GetDescription(),
		IssuedBy:       request.GetIssuedBy(),
		CertificateUrl: request.GetCertificateUrl(),
		IssueDate:      request.GetIssueDate(),
		ExpiryDate:     request.GetExpiryDate(),
	}

	if err := req.Validate(); err != nil {
		return nil, merchantaward_errors.ErrGrpcValidateCreateMerchantAward
	}

	merchant, err := s.merchantAwardCommand.Create(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchant_award.ApiResponseMerchantAward{
		Status:  "success",
		Message: "Successfully created merchant award",
		Data:    mapToProtoMerchantAwardResponseFromModel(merchant),
	}, nil
}

func (s *merchantAwardCommandHandler) Update(ctx context.Context, request *pbmerchant_award.UpdateMerchantAwardRequest) (*pbmerchant_award.ApiResponseMerchantAward, error) {
	id := int(request.GetMerchantCertificationId())
	req := &requests.UpdateMerchantCertificationOrAwardRequest{
		MerchantCertificationID: &id,
		Title:                   request.GetTitle(),
		Description:             request.GetDescription(),
		IssuedBy:                request.GetIssuedBy(),
		CertificateUrl:          request.GetCertificateUrl(),
		IssueDate:               request.GetIssueDate(),
		ExpiryDate:              request.GetExpiryDate(),
	}

	if err := req.Validate(); err != nil {
		return nil, merchantaward_errors.ErrGrpcValidateUpdateMerchantAward
	}

	merchant, err := s.merchantAwardCommand.Update(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchant_award.ApiResponseMerchantAward{
		Status:  "success",
		Message: "Successfully updated merchant award",
		Data:    mapToProtoMerchantAwardResponseFromModel(merchant),
	}, nil
}

func (s *merchantAwardCommandHandler) TrashedMerchantAward(ctx context.Context, request *pbmerchant_award.FindByIdMerchantAwardRequest) (*pbmerchant_award.ApiResponseMerchantAwardDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchantaward_errors.ErrGrpcMerchantInvalidId
	}

	merchant, err := s.merchantAwardCommand.Trash(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchant_award.ApiResponseMerchantAwardDeleteAt{
		Status:  "success",
		Message: "Successfully trashed merchant award",
		Data:    mapToProtoMerchantAwardResponseDeleteAtFromModel(merchant),
	}, nil
}

func (s *merchantAwardCommandHandler) RestoreMerchantAward(ctx context.Context, request *pbmerchant_award.FindByIdMerchantAwardRequest) (*pbmerchant_award.ApiResponseMerchantAwardDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchantaward_errors.ErrGrpcMerchantInvalidId
	}

	merchant, err := s.merchantAwardCommand.Restore(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchant_award.ApiResponseMerchantAwardDeleteAt{
		Status:  "success",
		Message: "Successfully restored merchant award",
		Data:    mapToProtoMerchantAwardResponseDeleteAtFromModel(merchant),
	}, nil
}

func (s *merchantAwardCommandHandler) DeleteMerchantAwardPermanent(ctx context.Context, request *pbmerchant_award.FindByIdMerchantAwardRequest) (*pbmerchant.ApiResponseMerchantDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchantaward_errors.ErrGrpcMerchantInvalidId
	}

	_, err := s.merchantAwardCommand.DeletePermanent(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchant.ApiResponseMerchantDelete{
		Status:  "success",
		Message: "Successfully deleted merchant award permanently",
	}, nil
}

func (s *merchantAwardCommandHandler) RestoreAllMerchantAward(ctx context.Context, _ *emptypb.Empty) (*pbmerchant.ApiResponseMerchantAll, error) {
	_, err := s.merchantAwardCommand.RestoreAll(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchant.ApiResponseMerchantAll{
		Status:  "success",
		Message: "Successfully restored all trashed merchant awards",
	}, nil
}

func (s *merchantAwardCommandHandler) DeleteAllMerchantAwardPermanent(ctx context.Context, _ *emptypb.Empty) (*pbmerchant.ApiResponseMerchantAll, error) {
	_, err := s.merchantAwardCommand.DeleteAll(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchant.ApiResponseMerchantAll{
		Status:  "success",
		Message: "Successfully deleted all merchant awards permanently",
	}, nil
}
