package handler

import (
	"context"
	pbcategory "github.com/MamangRust/microservice-ecommerce-grpc-pb/category"
	pbcommon "github.com/MamangRust/microservice-ecommerce-grpc-pb/common"
	"math"

	"github.com/MamangRust/microservice-ecommerce-grpc-category/service"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	category_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/category_errors"
)

type categoryQueryHandler struct {
	pbcategory.UnimplementedCategoryQueryServiceServer
	service service.CategoryQueryService
	logger  logger.LoggerInterface
}

func NewCategoryQueryHandler(service service.CategoryQueryService, logger logger.LoggerInterface) pbcategory.CategoryQueryServiceServer {
	return &categoryQueryHandler{service: service, logger: logger}
}

func (h *categoryQueryHandler) FindAll(ctx context.Context, request *pbcategory.FindAllCategoryRequest) (*pbcategory.ApiResponsePaginationCategory, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllCategory{Page: page, PageSize: pageSize, Search: search}
	categories, totalRecords, err := h.service.FindAll(ctx, &reqService)
	if err != nil {
		return nil, category_errors.ErrGrpcFindAllCategory
	}

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage: int32(page), PageSize: int32(pageSize),
		TotalPages: int32(math.Ceil(float64(*totalRecords) / float64(pageSize))), TotalRecords: int32(*totalRecords),
	}
	results := make([]*pbcategory.CategoryResponse, len(categories))
	for i, v := range categories {
		results[i] = (&Handler{}).mapToCategoryResponseFromResult(v)
	}
	return &pbcategory.ApiResponsePaginationCategory{Status: "success", Message: "Successfully fetched categories", Data: results, Pagination: paginationMeta}, nil
}

func (h *categoryQueryHandler) FindById(ctx context.Context, request *pbcategory.FindByIdCategoryRequest) (*pbcategory.ApiResponseCategory, error) {
	id := int(request.GetId())
	if id == 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidId
	}
	category, err := h.service.FindByID(ctx, id)
	if err != nil {
		return nil, category_errors.ErrGrpcCategoryNotFound
	}
	return &pbcategory.ApiResponseCategory{Status: "success", Message: "Successfully fetched category", Data: (&Handler{}).mapToCategoryResponseFromModel(category)}, nil
}

func (h *categoryQueryHandler) FindByActive(ctx context.Context, request *pbcategory.FindAllCategoryRequest) (*pbcategory.ApiResponsePaginationCategoryDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllCategory{Page: page, PageSize: pageSize, Search: search}
	categories, totalRecords, err := h.service.FindActive(ctx, &reqService)
	if err != nil {
		return nil, category_errors.ErrGrpcFindAllCategory
	}

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage: int32(page), PageSize: int32(pageSize),
		TotalPages: int32(math.Ceil(float64(*totalRecords) / float64(pageSize))), TotalRecords: int32(*totalRecords),
	}
	results := make([]*pbcategory.CategoryResponseDeleteAt, len(categories))
	for i, v := range categories {
		results[i] = (&Handler{}).mapToDeleteAtResponseFromResult(v)
	}
	return &pbcategory.ApiResponsePaginationCategoryDeleteAt{Status: "success", Message: "Successfully fetched active categories", Data: results, Pagination: paginationMeta}, nil
}

func (h *categoryQueryHandler) FindByTrashed(ctx context.Context, request *pbcategory.FindAllCategoryRequest) (*pbcategory.ApiResponsePaginationCategoryDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllCategory{Page: page, PageSize: pageSize, Search: search}
	categories, totalRecords, err := h.service.FindTrashed(ctx, &reqService)
	if err != nil {
		return nil, category_errors.ErrGrpcFindAllCategory
	}

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage: int32(page), PageSize: int32(pageSize),
		TotalPages: int32(math.Ceil(float64(*totalRecords) / float64(pageSize))), TotalRecords: int32(*totalRecords),
	}
	results := make([]*pbcategory.CategoryResponseDeleteAt, len(categories))
	for i, v := range categories {
		results[i] = (&Handler{}).mapToDeleteAtResponseFromResult(v)
	}
	return &pbcategory.ApiResponsePaginationCategoryDeleteAt{Status: "success", Message: "Successfully fetched trashed categories", Data: results, Pagination: paginationMeta}, nil
}
