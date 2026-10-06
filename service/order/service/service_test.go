package service

import (
	"context"
	"testing"
	"time"

	"github.com/MamangRust/microservice-ecommerce-grpc-order/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/observability"
	"go.uber.org/zap"
)

// --- Query Repo Stub ---

type stubOrderQueryRepo struct{}

func (s *stubOrderQueryRepo) FindAll(_ context.Context, _ *requests.FindAllOrder) ([]*repository.OrderResult, error) {
	now := time.Now()
	return []*repository.OrderResult{
		{OrderID: 1, UserID: 10, MerchantID: 20, TotalPrice: 50000, TotalCount: 1, CreatedAt: &now, UpdatedAt: &now},
	}, nil
}
func (s *stubOrderQueryRepo) FindActive(_ context.Context, _ *requests.FindAllOrder) ([]*repository.OrderResult, error) {
	now := time.Now()
	return []*repository.OrderResult{
		{OrderID: 1, UserID: 10, MerchantID: 20, TotalPrice: 50000, TotalCount: 1, CreatedAt: &now, UpdatedAt: &now},
	}, nil
}
func (s *stubOrderQueryRepo) FindTrashed(_ context.Context, _ *requests.FindAllOrder) ([]*repository.OrderResult, error) {
	return []*repository.OrderResult{}, nil
}
func (s *stubOrderQueryRepo) FindByMerchant(_ context.Context, _ *requests.FindAllOrderByMerchant) ([]*repository.OrderResult, error) {
	now := time.Now()
	return []*repository.OrderResult{
		{OrderID: 1, UserID: 10, MerchantID: 20, TotalPrice: 50000, TotalCount: 1, CreatedAt: &now, UpdatedAt: &now},
	}, nil
}
func (s *stubOrderQueryRepo) FindByID(_ context.Context, id int) (*models.Order, error) {
	now := time.Now()
	return &models.Order{OrderID: int32(id), UserID: 10, MerchantID: 20, TotalPrice: 50000, CreatedAt: &now, UpdatedAt: &now}, nil
}

// --- Query Cache Stub ---

type stubOrderQueryCache struct {
	order *models.Order
}

func (s *stubOrderQueryCache) GetOrderAllCache(_ context.Context, _ *requests.FindAllOrder) ([]*repository.OrderResult, *int, bool) {
	return nil, nil, false
}
func (s *stubOrderQueryCache) SetOrderAllCache(_ context.Context, _ *requests.FindAllOrder, _ []*repository.OrderResult, _ *int) {
}
func (s *stubOrderQueryCache) GetOrderActiveCache(_ context.Context, _ *requests.FindAllOrder) ([]*repository.OrderResult, *int, bool) {
	return nil, nil, false
}
func (s *stubOrderQueryCache) SetOrderActiveCache(_ context.Context, _ *requests.FindAllOrder, _ []*repository.OrderResult, _ *int) {
}
func (s *stubOrderQueryCache) GetOrderTrashedCache(_ context.Context, _ *requests.FindAllOrder) ([]*repository.OrderResult, *int, bool) {
	return nil, nil, false
}
func (s *stubOrderQueryCache) SetOrderTrashedCache(_ context.Context, _ *requests.FindAllOrder, _ []*repository.OrderResult, _ *int) {
}
func (s *stubOrderQueryCache) GetOrderByMerchantCache(_ context.Context, _ *requests.FindAllOrderByMerchant) ([]*repository.OrderResult, *int, bool) {
	return nil, nil, false
}
func (s *stubOrderQueryCache) SetOrderByMerchantCache(_ context.Context, _ *requests.FindAllOrderByMerchant, _ []*repository.OrderResult, _ *int) {
}
func (s *stubOrderQueryCache) GetCachedOrderCache(_ context.Context, _ int) (*models.Order, bool) {
	if s.order == nil {
		return nil, false
	}
	return s.order, true
}
func (s *stubOrderQueryCache) SetCachedOrderCache(_ context.Context, _ *models.Order) {}
func (s *stubOrderQueryCache) InvalidateOrderCache(_ context.Context)                 {}

// --- Query Service Tests ---

func TestOrderQueryServiceFindAll(t *testing.T) {
	log := &logger.Logger{Log: zap.NewNop()}
	obs, _ := observability.NewObservability("test", log)
	svc := &orderQueryService{observability: obs, cache: &stubOrderQueryCache{}, orderRepository: &stubOrderQueryRepo{}, logger: log}

	data, total, err := svc.FindAll(context.Background(), &requests.FindAllOrder{Page: 1, PageSize: 10, Search: ""})
	if err != nil {
		t.Fatalf("FindAll() error = %v", err)
	}
	if len(data) != 1 {
		t.Fatalf("len(data) = %d, want 1", len(data))
	}
	if total == nil || *total != 1 {
		t.Fatalf("total = %v, want 1", total)
	}
}

func TestOrderQueryServiceFindActive(t *testing.T) {
	log := &logger.Logger{Log: zap.NewNop()}
	obs, _ := observability.NewObservability("test", log)
	svc := &orderQueryService{observability: obs, cache: &stubOrderQueryCache{}, orderRepository: &stubOrderQueryRepo{}, logger: log}

	data, total, err := svc.FindActive(context.Background(), &requests.FindAllOrder{Page: 1, PageSize: 10, Search: ""})
	if err != nil {
		t.Fatalf("FindActive() error = %v", err)
	}
	if len(data) != 1 || total == nil || *total != 1 {
		t.Fatalf("FindActive() data=%v total=%v", data, total)
	}
}

func TestOrderQueryServiceFindTrashed(t *testing.T) {
	log := &logger.Logger{Log: zap.NewNop()}
	obs, _ := observability.NewObservability("test", log)
	svc := &orderQueryService{observability: obs, cache: &stubOrderQueryCache{}, orderRepository: &stubOrderQueryRepo{}, logger: log}

	data, total, err := svc.FindTrashed(context.Background(), &requests.FindAllOrder{Page: 1, PageSize: 10, Search: ""})
	if err != nil {
		t.Fatalf("FindTrashed() error = %v", err)
	}
	if len(data) != 0 || total == nil || *total != 0 {
		t.Fatalf("FindTrashed() data=%v total=%v", data, total)
	}
}

func TestOrderQueryServiceFindByID(t *testing.T) {
	log := &logger.Logger{Log: zap.NewNop()}
	obs, _ := observability.NewObservability("test", log)
	svc := &orderQueryService{observability: obs, cache: &stubOrderQueryCache{}, orderRepository: &stubOrderQueryRepo{}, logger: log}

	result, err := svc.FindByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}
	if result.OrderID != 1 {
		t.Fatalf("OrderID = %d, want 1", result.OrderID)
	}
}

func TestOrderQueryServiceFindByIDFromCache(t *testing.T) {
	log := &logger.Logger{Log: zap.NewNop()}
	obs, _ := observability.NewObservability("test", log)
	now := time.Now()
	svc := &orderQueryService{
		observability:   obs,
		cache:           &stubOrderQueryCache{order: &models.Order{OrderID: 42, CreatedAt: &now, UpdatedAt: &now}},
		orderRepository: &stubOrderQueryRepo{},
		logger:          log,
	}

	result, err := svc.FindByID(context.Background(), 42)
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}
	if result.OrderID != 42 {
		t.Fatalf("OrderID = %d, want 42", result.OrderID)
	}
}

func TestOrderQueryServiceFindByMerchant(t *testing.T) {
	log := &logger.Logger{Log: zap.NewNop()}
	obs, _ := observability.NewObservability("test", log)
	svc := &orderQueryService{observability: obs, cache: &stubOrderQueryCache{}, orderRepository: &stubOrderQueryRepo{}, logger: log}

	data, total, err := svc.FindByMerchant(context.Background(), &requests.FindAllOrderByMerchant{MerchantID: 20, Page: 1, PageSize: 10, Search: ""})
	if err != nil {
		t.Fatalf("FindByMerchant() error = %v", err)
	}
	if len(data) != 1 || total == nil || *total != 1 {
		t.Fatalf("FindByMerchant() data=%v total=%v", data, total)
	}
}

func TestPointerInt32ToInt(t *testing.T) {
	result := pointerInt32ToInt(42)
	if result == nil || *result != 42 {
		t.Fatalf("pointerInt32ToInt(42) = %v, want 42", result)
	}

	result = pointerInt32ToInt(0)
	if result == nil || *result != 0 {
		t.Fatalf("pointerInt32ToInt(0) = %v, want 0", result)
	}
}

func TestReservationOperationID(t *testing.T) {
	now := time.Now()
	reservation := &models.OrderStockReservation{
		OrderID:       1,
		ProductID:     2,
		ReservationID: 3,
		UpdatedAt:     &now,
	}
	opID := reservationOperationID("prefix", reservation)
	if opID == "" {
		t.Fatal("reservationOperationID() returned empty string")
	}
}
