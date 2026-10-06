package transactionapimapper

import (
	pbtxstats "github.com/MamangRust/microservice-ecommerce-grpc-pb/transaction"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/response"
)

type transactionStatsResponseMapper struct{}

func NewTransactionStatsResponseMapper() TransactionStatsResponseMapper {
	return &transactionStatsResponseMapper{}
}

func (m *transactionStatsResponseMapper) ToTransactionMonthAmountSuccess(row *pbtxstats.TransactionMonthlyAmountSuccess) *response.TransactionMonthlyAmountSuccessResponse {
	return &response.TransactionMonthlyAmountSuccessResponse{
		Year:         row.Year,
		Month:        row.Month,
		TotalSuccess: int(row.TotalSuccess),
		TotalAmount:  int(row.TotalAmount),
	}
}

func (m *transactionStatsResponseMapper) ToTransactionMonthlyAmountSuccess(rows []*pbtxstats.TransactionMonthlyAmountSuccess) []*response.TransactionMonthlyAmountSuccessResponse {
	var mapped []*response.TransactionMonthlyAmountSuccessResponse
	for _, row := range rows {
		mapped = append(mapped, m.ToTransactionMonthAmountSuccess(row))
	}
	return mapped
}

func (m *transactionStatsResponseMapper) ToTransactionYearAmountSuccess(row *pbtxstats.TransactionYearlyAmountSuccess) *response.TransactionYearlyAmountSuccessResponse {
	return &response.TransactionYearlyAmountSuccessResponse{
		Year:         row.Year,
		TotalSuccess: int(row.TotalSuccess),
		TotalAmount:  int(row.TotalAmount),
	}
}

func (m *transactionStatsResponseMapper) ToTransactionYearlyAmountSuccess(rows []*pbtxstats.TransactionYearlyAmountSuccess) []*response.TransactionYearlyAmountSuccessResponse {
	var mapped []*response.TransactionYearlyAmountSuccessResponse
	for _, row := range rows {
		mapped = append(mapped, m.ToTransactionYearAmountSuccess(row))
	}
	return mapped
}

func (m *transactionStatsResponseMapper) ToTransactionMonthAmountFailed(row *pbtxstats.TransactionMonthlyAmountFailed) *response.TransactionMonthlyAmountFailedResponse {
	return &response.TransactionMonthlyAmountFailedResponse{
		Year:        row.Year,
		Month:       row.Month,
		TotalFailed: int(row.TotalFailed),
		TotalAmount: int(row.TotalAmount),
	}
}

func (m *transactionStatsResponseMapper) ToTransactionMonthlyAmountFailed(rows []*pbtxstats.TransactionMonthlyAmountFailed) []*response.TransactionMonthlyAmountFailedResponse {
	var mapped []*response.TransactionMonthlyAmountFailedResponse
	for _, row := range rows {
		mapped = append(mapped, m.ToTransactionMonthAmountFailed(row))
	}
	return mapped
}

func (m *transactionStatsResponseMapper) ToTransactionYearAmountFailed(row *pbtxstats.TransactionYearlyAmountFailed) *response.TransactionYearlyAmountFailedResponse {
	return &response.TransactionYearlyAmountFailedResponse{
		Year:        row.Year,
		TotalFailed: int(row.TotalFailed),
		TotalAmount: int(row.TotalAmount),
	}
}

func (m *transactionStatsResponseMapper) ToTransactionYearlyAmountFailed(rows []*pbtxstats.TransactionYearlyAmountFailed) []*response.TransactionYearlyAmountFailedResponse {
	var mapped []*response.TransactionYearlyAmountFailedResponse
	for _, row := range rows {
		mapped = append(mapped, m.ToTransactionYearAmountFailed(row))
	}
	return mapped
}

func (m *transactionStatsResponseMapper) ToTransactionMonthMethod(row *pbtxstats.TransactionMonthlyMethod) *response.TransactionMonthlyMethodResponse {
	return &response.TransactionMonthlyMethodResponse{
		Month:             row.Month,
		PaymentMethod:     row.PaymentMethod,
		TotalTransactions: int(row.TotalTransactions),
		TotalAmount:       int(row.TotalAmount),
	}
}

func (m *transactionStatsResponseMapper) ToTransactionMonthlyMethod(rows []*pbtxstats.TransactionMonthlyMethod) []*response.TransactionMonthlyMethodResponse {
	var mapped []*response.TransactionMonthlyMethodResponse
	for _, row := range rows {
		mapped = append(mapped, m.ToTransactionMonthMethod(row))
	}
	return mapped
}

func (m *transactionStatsResponseMapper) ToTransactionYearMethod(row *pbtxstats.TransactionYearlyMethod) *response.TransactionYearlyMethodResponse {
	return &response.TransactionYearlyMethodResponse{
		Year:              row.Year,
		PaymentMethod:     row.PaymentMethod,
		TotalTransactions: int(row.TotalTransactions),
		TotalAmount:       int(row.TotalAmount),
	}
}

func (m *transactionStatsResponseMapper) ToTransactionYearlyMethod(rows []*pbtxstats.TransactionYearlyMethod) []*response.TransactionYearlyMethodResponse {
	var mapped []*response.TransactionYearlyMethodResponse
	for _, row := range rows {
		mapped = append(mapped, m.ToTransactionYearMethod(row))
	}
	return mapped
}

func (m *transactionStatsResponseMapper) ToApiResponseTransactionMonthAmountSuccess(pbResponse *pbtxstats.ApiResponseTransactionMonthAmountSuccess) *response.ApiResponsesTransactionMonthSuccess {
	return &response.ApiResponsesTransactionMonthSuccess{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToTransactionMonthlyAmountSuccess(pbResponse.Data),
	}
}

func (m *transactionStatsResponseMapper) ToApiResponseTransactionMonthAmountFailed(pbResponse *pbtxstats.ApiResponseTransactionMonthAmountFailed) *response.ApiResponsesTransactionMonthFailed {
	return &response.ApiResponsesTransactionMonthFailed{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToTransactionMonthlyAmountFailed(pbResponse.Data),
	}
}

func (m *transactionStatsResponseMapper) ToApiResponseTransactionYearAmountSuccess(pbResponse *pbtxstats.ApiResponseTransactionYearAmountSuccess) *response.ApiResponsesTransactionYearSuccess {
	return &response.ApiResponsesTransactionYearSuccess{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToTransactionYearlyAmountSuccess(pbResponse.Data),
	}
}

func (m *transactionStatsResponseMapper) ToApiResponseTransactionYearAmountFailed(pbResponse *pbtxstats.ApiResponseTransactionYearAmountFailed) *response.ApiResponsesTransactionYearFailed {
	return &response.ApiResponsesTransactionYearFailed{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToTransactionYearlyAmountFailed(pbResponse.Data),
	}
}

func (m *transactionStatsResponseMapper) ToApiResponseTransactionMonthMethod(pbResponse *pbtxstats.ApiResponseTransactionMonthPaymentMethod) *response.ApiResponsesTransactionMonthMethod {
	return &response.ApiResponsesTransactionMonthMethod{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToTransactionMonthlyMethod(pbResponse.Data),
	}
}

func (m *transactionStatsResponseMapper) ToApiResponseTransactionYearMethod(pbResponse *pbtxstats.ApiResponseTransactionYearPaymentmethod) *response.ApiResponsesTransactionYearMethod {
	return &response.ApiResponsesTransactionYearMethod{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToTransactionYearlyMethod(pbResponse.Data),
	}
}
