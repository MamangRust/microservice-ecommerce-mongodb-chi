package handler

import (
	pbcommon "github.com/MamangRust/microservice-ecommerce-grpc-pb/common"
	pbmerchant_policy "github.com/MamangRust/microservice-ecommerce-grpc-pb/merchant_policy"

	"github.com/MamangRust/microservice-ecommerce-grpc-merchant_policy/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/convert"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func mapToSingleResponseFromModel(data *models.MerchantPolicy) *pbmerchant_policy.ApiResponseMerchantPolicies {
	if data == nil {
		return nil
	}
	return &pbmerchant_policy.ApiResponseMerchantPolicies{
		Status:  "success",
		Message: "Successfully fetched merchant policy",
		Data: &pbmerchant_policy.MerchantPoliciesResponse{
			Id:          data.MerchantPolicyID,
			MerchantId:  data.MerchantID,
			PolicyType:  data.PolicyType,
			Title:       data.Title,
			Description: data.Description,
			CreatedAt:   convert.FormatTimePtr(data.CreatedAt),
			UpdatedAt:   convert.FormatTimePtr(data.UpdatedAt),
		},
	}
}

func mapToSingleResponse(data *repository.MerchantPolicyResult) *pbmerchant_policy.ApiResponseMerchantPolicies {
	return &pbmerchant_policy.ApiResponseMerchantPolicies{
		Status:  "success",
		Message: "Successfully fetched merchant policy",
		Data:    mapToMerchantPolicyResponseFromResult(data),
	}
}

func mapToPaginationResponse(data []*repository.MerchantPolicyResult, total *int) *pbmerchant_policy.ApiResponsePaginationMerchantPolicies {
	var policies []*pbmerchant_policy.MerchantPoliciesResponse
	for _, v := range data {
		policies = append(policies, mapToMerchantPolicyResponseFromResult(v))
	}
	return &pbmerchant_policy.ApiResponsePaginationMerchantPolicies{
		Status:  "success",
		Message: "Successfully fetched merchant policies",
		Data:    policies,
		Pagination: &pbcommon.PaginationMeta{
			TotalRecords: int32(*total),
		},
	}
}

func mapToPaginationDeleteAtResponse(data []*repository.MerchantPolicyResult, total *int) *pbmerchant_policy.ApiResponsePaginationMerchantPoliciesDeleteAt {
	var policies []*pbmerchant_policy.MerchantPoliciesResponseDeleteAt
	for _, item := range data {
		policies = append(policies, mapToMerchantPolicyResponseDeleteAtFromResult(item))
	}
	return &pbmerchant_policy.ApiResponsePaginationMerchantPoliciesDeleteAt{
		Status:  "success",
		Message: "Successfully fetched merchant policies",
		Data:    policies,
		Pagination: &pbcommon.PaginationMeta{
			TotalRecords: int32(*total),
		},
	}
}

func mapToSingleDeleteAtResponse(data *models.MerchantPolicy) *pbmerchant_policy.ApiResponseMerchantPoliciesDeleteAt {
	return &pbmerchant_policy.ApiResponseMerchantPoliciesDeleteAt{
		Status:  "success",
		Message: "Successfully processed merchant policy",
		Data:    mapToMerchantPolicyResponseDeleteAtFromModel(data),
	}
}

func mapToMerchantPolicyResponseFromResult(v *repository.MerchantPolicyResult) *pbmerchant_policy.MerchantPoliciesResponse {
	if v == nil {
		return nil
	}
	return &pbmerchant_policy.MerchantPoliciesResponse{
		Id:           v.MerchantPolicyID,
		MerchantId:   v.MerchantID,
		PolicyType:   v.PolicyType,
		Title:        v.Title,
		Description:  v.Description,
		MerchantName: v.MerchantName,
		CreatedAt:    convert.FormatTimePtr(v.CreatedAt),
		UpdatedAt:    convert.FormatTimePtr(v.UpdatedAt),
	}
}

func mapToMerchantPolicyResponseDeleteAtFromResult(v *repository.MerchantPolicyResult) *pbmerchant_policy.MerchantPoliciesResponseDeleteAt {
	if v == nil {
		return nil
	}
	res := &pbmerchant_policy.MerchantPoliciesResponseDeleteAt{
		Id:           v.MerchantPolicyID,
		MerchantId:   v.MerchantID,
		PolicyType:   v.PolicyType,
		Title:        v.Title,
		Description:  v.Description,
		MerchantName: v.MerchantName,
		CreatedAt:    convert.FormatTimePtr(v.CreatedAt),
		UpdatedAt:    convert.FormatTimePtr(v.UpdatedAt),
	}
	if v.DeletedAt != nil {
		res.DeletedAt = &wrapperspb.StringValue{Value: convert.FormatTimePtr(v.DeletedAt)}
	}
	return res
}

func mapToMerchantPolicyResponseDeleteAtFromModel(v *models.MerchantPolicy) *pbmerchant_policy.MerchantPoliciesResponseDeleteAt {
	if v == nil {
		return nil
	}
	res := &pbmerchant_policy.MerchantPoliciesResponseDeleteAt{
		Id:          v.MerchantPolicyID,
		MerchantId:  v.MerchantID,
		PolicyType:  v.PolicyType,
		Title:       v.Title,
		Description: v.Description,
		CreatedAt:   convert.FormatTimePtr(v.CreatedAt),
		UpdatedAt:   convert.FormatTimePtr(v.UpdatedAt),
	}
	if v.DeletedAt != nil {
		res.DeletedAt = &wrapperspb.StringValue{Value: convert.FormatTimePtr(v.DeletedAt)}
	}
	return res
}
