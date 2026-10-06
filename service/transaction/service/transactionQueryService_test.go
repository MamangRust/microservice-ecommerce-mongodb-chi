package service

import (
	"context"
	"testing"
	"time"

	"github.com/MamangRust/microservice-ecommerce-grpc-transaction/cache"
	"github.com/MamangRust/microservice-ecommerce-grpc-transaction/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/observability"
	"go.uber.org/zap"
)

// --- Cache stub ---

type stubTxQueryCache struct {
	data   []*repository.TransactionResult
	total  int
	single *repository.TransactionResult
}

func (s *stubTxQueryCache) GetCachedTransactionsCache(_ context.Context, _ *requests.FindAllTransaction) ([]*repository.TransactionResult, *int, bool) {
	if s.data == nil {
		return nil, nil, false
	}
	t := s.total
	return s.data, &t, true
}
func (s *stubTxQueryCache) SetCachedTransactionsCache(_ context.Context, _ *requests.FindAllTransaction, _ []*repository.TransactionResult, _ *int) {
}
func (s *stubTxQueryCache) GetCachedTransactionByMerchant(_ context.Context, _ *requests.FindAllTransactionByMerchant) ([]*repository.TransactionResult, *int, bool) {
	return nil, nil, false
}
func (s *stubTxQueryCache) SetCachedTransactionByMerchant(_ context.Context, _ *requests.FindAllTransactionByMerchant, _ []*repository.TransactionResult, _ *int) {
}
func (s *stubTxQueryCache) GetCachedTransactionActiveCache(_ context.Context, _ *requests.FindAllTransaction) ([]*repository.TransactionResult, *int, bool) {
	return nil, nil, false
}
func (s *stubTxQueryCache) SetCachedTransactionActiveCache(_ context.Context, _ *requests.FindAllTransaction, _ []*repository.TransactionResult, _ *int) {
}
func (s *stubTxQueryCache) GetCachedTransactionTrashedCache(_ context.Context, _ *requests.FindAllTransaction) ([]*repository.TransactionResult, *int, bool) {
	return nil, nil, false
}
func (s *stubTxQueryCache) SetCachedTransactionTrashedCache(_ context.Context, _ *requests.FindAllTransaction, _ []*repository.TransactionResult, _ *int) {
}
func (s *stubTxQueryCache) GetCachedTransactionCache(_ context.Context, _ int) (*repository.TransactionResult, bool) {
	if s.single == nil {
		return nil, false
	}
	return s.single, true
}
func (s *stubTxQueryCache) SetCachedTransactionCache(_ context.Context, _ *repository.TransactionResult) {
}
func (s *stubTxQueryCache) GetCachedTransactionByOrderId(_ context.Context, _ int) (*repository.TransactionResult, bool) {
	return nil, false
}
func (s *stubTxQueryCache) SetCachedTransactionByOrderId(_ context.Context, _ int, _ *repository.TransactionResult) {
}

var _ cache.TransactionQueryCache = (*stubTxQueryCache)(nil)

// --- Repo stub ---

type stubTxQueryRepoAll struct{}

func (s *stubTxQueryRepoAll) FindAll(_ context.Context, _ *requests.FindAllTransaction) ([]*repository.TransactionResult, error) {
	now := time.Now()
	return []*repository.TransactionResult{
		{TransactionID: 1, OrderID: 10, MerchantID: 20, PaymentMethod: "credit_card", Amount: 50000, PaymentStatus: "success", TotalCount: 1, CreatedAt: &now, UpdatedAt: &now},
		{TransactionID: 2, OrderID: 11, MerchantID: 20, PaymentMethod: "bank_transfer", Amount: 30000, PaymentStatus: "pending", TotalCount: 2, CreatedAt: &now, UpdatedAt: &now},
	}, nil
}
func (s *stubTxQueryRepoAll) FindActive(_ context.Context, _ *requests.FindAllTransaction) ([]*repository.TransactionResult, error) {
	now := time.Now()
	return []*repository.TransactionResult{
		{TransactionID: 1, OrderID: 10, MerchantID: 20, PaymentMethod: "credit_card", Amount: 50000, PaymentStatus: "success", TotalCount: 1, CreatedAt: &now, UpdatedAt: &now},
	}, nil
}
func (s *stubTxQueryRepoAll) FindTrashed(_ context.Context, _ *requests.FindAllTransaction) ([]*repository.TransactionResult, error) {
	return []*repository.TransactionResult{}, nil
}
func (s *stubTxQueryRepoAll) FindByMerchant(_ context.Context, _ *requests.FindAllTransactionByMerchant) ([]*repository.TransactionResult, error) {
	now := time.Now()
	return []*repository.TransactionResult{
		{TransactionID: 1, OrderID: 10, MerchantID: 20, PaymentMethod: "credit_card", Amount: 50000, PaymentStatus: "success", TotalCount: 1, CreatedAt: &now, UpdatedAt: &now},
	}, nil
}
func (s *stubTxQueryRepoAll) FindByID(_ context.Context, id int) (*repository.TransactionResult, error) {
	now := time.Now()
	return &repository.TransactionResult{
		TransactionID: int32(id), OrderID: 10, MerchantID: 20,
		PaymentMethod: "credit_card", Amount: 50000,
		PaymentStatus: "success", CreatedAt: &now, UpdatedAt: &now,
	}, nil
}
func (s *stubTxQueryRepoAll) FindByOrderID(_ context.Context, _ int) (*repository.TransactionResult, error) {
	now := time.Now()
	return &repository.TransactionResult{
		TransactionID: 1, OrderID: 10, MerchantID: 20,
		PaymentMethod: "credit_card", Amount: 50000,
		PaymentStatus: "success", CreatedAt: &now, UpdatedAt: &now,
	}, nil
}

// --- Tests ---

func TestTransactionQueryServiceFindAll(t *testing.T) {
	log := &logger.Logger{Log: zap.NewNop()}
	obs, _ := observability.NewObservability("test", log)

	svc := &transactionQueryService{
		observability: obs,
		cache:         &stubTxQueryCache{},
		repository:    &stubTxQueryRepoAll{},
		logger:        log,
	}

	data, total, err := svc.FindAll(context.Background(), &requests.FindAllTransaction{
		Page: 1, PageSize: 10, Search: "",
	})
	if err != nil {
		t.Fatalf("FindAll() error = %v", err)
	}
	if len(data) != 2 {
		t.Fatalf("len(data) = %d, want 2", len(data))
	}
	if total == nil || *total != 1 {
		t.Fatalf("total = %v, want 1", total)
	}
}

func TestTransactionQueryServiceFindAllFromCache(t *testing.T) {
	log := &logger.Logger{Log: zap.NewNop()}
	obs, _ := observability.NewObservability("test", log)

	now := time.Now()
	cached := []*repository.TransactionResult{
		{TransactionID: 1, TotalCount: 1, CreatedAt: &now, UpdatedAt: &now},
	}
	svc := &transactionQueryService{
		observability: obs,
		cache:         &stubTxQueryCache{data: cached, total: 1},
		repository:    &stubTxQueryRepoAll{},
		logger:        log,
	}

	data, total, err := svc.FindAll(context.Background(), &requests.FindAllTransaction{
		Page: 1, PageSize: 10, Search: "",
	})
	if err != nil {
		t.Fatalf("FindAll() error = %v", err)
	}
	if len(data) != 1 {
		t.Fatalf("len(data) = %d, want 1 (from cache)", len(data))
	}
	if total == nil || *total != 1 {
		t.Fatalf("total = %v, want 1", total)
	}
}

func TestTransactionQueryServiceFindActive(t *testing.T) {
	log := &logger.Logger{Log: zap.NewNop()}
	obs, _ := observability.NewObservability("test", log)

	svc := &transactionQueryService{
		observability: obs,
		cache:         &stubTxQueryCache{},
		repository:    &stubTxQueryRepoAll{},
		logger:        log,
	}

	data, total, err := svc.FindActive(context.Background(), &requests.FindAllTransaction{
		Page: 1, PageSize: 10, Search: "",
	})
	if err != nil {
		t.Fatalf("FindActive() error = %v", err)
	}
	if len(data) != 1 {
		t.Fatalf("len(data) = %d, want 1", len(data))
	}
	if total == nil || *total != 1 {
		t.Fatalf("total = %v, want 1", total)
	}
}

func TestTransactionQueryServiceFindTrashed(t *testing.T) {
	log := &logger.Logger{Log: zap.NewNop()}
	obs, _ := observability.NewObservability("test", log)

	svc := &transactionQueryService{
		observability: obs,
		cache:         &stubTxQueryCache{},
		repository:    &stubTxQueryRepoAll{},
		logger:        log,
	}

	data, total, err := svc.FindTrashed(context.Background(), &requests.FindAllTransaction{
		Page: 1, PageSize: 10, Search: "",
	})
	if err != nil {
		t.Fatalf("FindTrashed() error = %v", err)
	}
	if len(data) != 0 {
		t.Fatalf("len(data) = %d, want 0", len(data))
	}
	if total == nil || *total != 0 {
		t.Fatalf("total = %v, want 0", total)
	}
}

func TestTransactionQueryServiceFindByID(t *testing.T) {
	log := &logger.Logger{Log: zap.NewNop()}
	obs, _ := observability.NewObservability("test", log)

	svc := &transactionQueryService{
		observability: obs,
		cache:         &stubTxQueryCache{},
		repository:    &stubTxQueryRepoAll{},
		logger:        log,
	}

	result, err := svc.FindByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}
	if result.TransactionID != 1 {
		t.Fatalf("TransactionID = %d, want 1", result.TransactionID)
	}
}

func TestTransactionQueryServiceFindByIDFromCache(t *testing.T) {
	log := &logger.Logger{Log: zap.NewNop()}
	obs, _ := observability.NewObservability("test", log)

	now := time.Now()
	svc := &transactionQueryService{
		observability: obs,
		cache:         &stubTxQueryCache{single: &repository.TransactionResult{TransactionID: 42, CreatedAt: &now, UpdatedAt: &now}},
		repository:    &stubTxQueryRepoAll{},
		logger:        log,
	}

	result, err := svc.FindByID(context.Background(), 42)
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}
	if result.TransactionID != 42 {
		t.Fatalf("TransactionID = %d, want 42", result.TransactionID)
	}
}

func TestTransactionQueryServiceFindByOrderID(t *testing.T) {
	log := &logger.Logger{Log: zap.NewNop()}
	obs, _ := observability.NewObservability("test", log)

	svc := &transactionQueryService{
		observability: obs,
		cache:         &stubTxQueryCache{},
		repository:    &stubTxQueryRepoAll{},
		logger:        log,
	}

	result, err := svc.FindByOrderID(context.Background(), 10)
	if err != nil {
		t.Fatalf("FindByOrderID() error = %v", err)
	}
	if result.OrderID != 10 {
		t.Fatalf("OrderID = %d, want 10", result.OrderID)
	}
}

func TestTransactionQueryServiceFindByMerchant(t *testing.T) {
	log := &logger.Logger{Log: zap.NewNop()}
	obs, _ := observability.NewObservability("test", log)

	svc := &transactionQueryService{
		observability: obs,
		cache:         &stubTxQueryCache{},
		repository:    &stubTxQueryRepoAll{},
		logger:        log,
	}

	data, total, err := svc.FindByMerchant(context.Background(), &requests.FindAllTransactionByMerchant{
		MerchantID: 20, Page: 1, PageSize: 10, Search: "",
	})
	if err != nil {
		t.Fatalf("FindByMerchant() error = %v", err)
	}
	if len(data) != 1 {
		t.Fatalf("len(data) = %d, want 1", len(data))
	}
	if total == nil || *total != 1 {
		t.Fatalf("total = %v, want 1", total)
	}
}
