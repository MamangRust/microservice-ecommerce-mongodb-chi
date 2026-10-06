package service

import (
	"context"
	"testing"

	"github.com/MamangRust/microservice-ecommerce-grpc-order-item/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/observability"
	"go.uber.org/zap"
)

// --- Stubs ---

type stubOrderItemCmdCache struct{}

func (s *stubOrderItemCmdCache) InvalidateOrderItemCache(_ context.Context) error { return nil }

type stubOrderItemCmdRepo struct{}

func (s *stubOrderItemCmdRepo) Create(_ context.Context, req *requests.CreateOrderItemRecordRequest) (*models.OrderItem, error) {
	return &models.OrderItem{
		OrderItemID: 99, OrderID: int32(req.OrderID), ProductID: int32(req.ProductID),
		Quantity: int32(req.Quantity), Price: int32(req.Price),
	}, nil
}
func (s *stubOrderItemCmdRepo) Update(_ context.Context, req *requests.UpdateOrderItemRecordRequest) (*models.OrderItem, error) {
	return &models.OrderItem{
		OrderItemID: int32(req.OrderItemID), OrderID: 10, ProductID: 20,
		Quantity: int32(req.Quantity), Price: int32(req.Price),
	}, nil
}
func (s *stubOrderItemCmdRepo) Trash(_ context.Context, id int) (*models.OrderItem, error) {
	return &models.OrderItem{OrderItemID: int32(id), OrderID: 10, ProductID: 20}, nil
}
func (s *stubOrderItemCmdRepo) Restore(_ context.Context, id int) (*models.OrderItem, error) {
	return &models.OrderItem{OrderItemID: int32(id), OrderID: 10, ProductID: 20}, nil
}
func (s *stubOrderItemCmdRepo) DeletePermanent(_ context.Context, _ int) (bool, error) {
	return true, nil
}
func (s *stubOrderItemCmdRepo) DeleteOrderItemByOrderPermanent(_ context.Context, _ int) (bool, error) {
	return true, nil
}
func (s *stubOrderItemCmdRepo) RestoreAll(_ context.Context) (bool, error) { return true, nil }
func (s *stubOrderItemCmdRepo) DeleteAll(_ context.Context) (bool, error)  { return true, nil }
func (s *stubOrderItemCmdRepo) CalculateTotalPrice(_ context.Context, _ int) (int, error) {
	return 25000, nil
}

type stubOrderItemQueryRepo struct{}

func (s *stubOrderItemQueryRepo) FindAll(_ context.Context, _ *requests.FindAllOrderItems) ([]*repository.OrderItemResult, error) {
	return []*repository.OrderItemResult{
		{OrderItemID: 1, OrderID: 10, ProductID: 20, Quantity: 2, Price: 10000, TotalCount: 1},
	}, nil
}
func (s *stubOrderItemQueryRepo) FindActive(_ context.Context, _ *requests.FindAllOrderItems) ([]*repository.OrderItemResult, error) {
	return []*repository.OrderItemResult{
		{OrderItemID: 1, OrderID: 10, ProductID: 20, Quantity: 2, Price: 10000, TotalCount: 1},
	}, nil
}
func (s *stubOrderItemQueryRepo) FindTrashed(_ context.Context, _ *requests.FindAllOrderItems) ([]*repository.OrderItemResult, error) {
	return []*repository.OrderItemResult{}, nil
}
func (s *stubOrderItemQueryRepo) FindOrderItemByOrder(_ context.Context, orderID int) ([]*repository.OrderItemResult, error) {
	return []*repository.OrderItemResult{
		{OrderItemID: 1, OrderID: int32(orderID), ProductID: 20, Quantity: 2, Price: 10000},
		{OrderItemID: 2, OrderID: int32(orderID), ProductID: 30, Quantity: 1, Price: 5000},
	}, nil
}

type stubOrderItemQueryCache struct {
	data  []*repository.OrderItemResult
	total int
}

func (s *stubOrderItemQueryCache) GetCachedOrderItemsAll(_ context.Context, _ *requests.FindAllOrderItems) ([]*repository.OrderItemResult, *int, bool) {
	if s.data == nil {
		return nil, nil, false
	}
	t := s.total
	return s.data, &t, true
}
func (s *stubOrderItemQueryCache) SetCachedOrderItemsAll(_ context.Context, _ *requests.FindAllOrderItems, _ []*repository.OrderItemResult, _ *int) {
}
func (s *stubOrderItemQueryCache) GetCachedOrderItemActive(_ context.Context, _ *requests.FindAllOrderItems) ([]*repository.OrderItemResult, *int, bool) {
	return nil, nil, false
}
func (s *stubOrderItemQueryCache) SetCachedOrderItemActive(_ context.Context, _ *requests.FindAllOrderItems, _ []*repository.OrderItemResult, _ *int) {
}
func (s *stubOrderItemQueryCache) GetCachedOrderItemTrashed(_ context.Context, _ *requests.FindAllOrderItems) ([]*repository.OrderItemResult, *int, bool) {
	return nil, nil, false
}
func (s *stubOrderItemQueryCache) SetCachedOrderItemTrashed(_ context.Context, _ *requests.FindAllOrderItems, _ []*repository.OrderItemResult, _ *int) {
}
func (s *stubOrderItemQueryCache) GetCachedOrderItems(_ context.Context, _ int) ([]*repository.OrderItemResult, bool) {
	return nil, false
}
func (s *stubOrderItemQueryCache) SetCachedOrderItems(_ context.Context, _ []*repository.OrderItemResult) {
}

// --- Command Service Tests ---

func TestOrderItemCommandServiceCreate(t *testing.T) {
	log := &logger.Logger{Log: zap.NewNop()}
	obs, _ := observability.NewObservability("test", log)

	svc := &orderItemCommandService{
		observability:       obs,
		cache:               &stubOrderItemCmdCache{},
		orderItemRepository: &stubOrderItemCmdRepo{},
		logger:              log,
	}

	result, err := svc.Create(context.Background(), &requests.CreateOrderItemRecordRequest{
		OrderID:   10,
		ProductID: 20,
		Quantity:  3,
		Price:     10000,
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if result.OrderItemID != 99 {
		t.Fatalf("OrderItemID = %d, want 99", result.OrderItemID)
	}
	if result.Quantity != 3 {
		t.Fatalf("Quantity = %d, want 3", result.Quantity)
	}
}

func TestOrderItemCommandServiceUpdate(t *testing.T) {
	log := &logger.Logger{Log: zap.NewNop()}
	obs, _ := observability.NewObservability("test", log)

	svc := &orderItemCommandService{
		observability:       obs,
		cache:               &stubOrderItemCmdCache{},
		orderItemRepository: &stubOrderItemCmdRepo{},
		logger:              log,
	}

	result, err := svc.Update(context.Background(), &requests.UpdateOrderItemRecordRequest{
		OrderItemID: 5,
		Quantity:    10,
		Price:       20000,
	})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if result.OrderItemID != 5 {
		t.Fatalf("OrderItemID = %d, want 5", result.OrderItemID)
	}
	if result.Quantity != 10 {
		t.Fatalf("Quantity = %d, want 10", result.Quantity)
	}
}

func TestOrderItemCommandServiceTrash(t *testing.T) {
	log := &logger.Logger{Log: zap.NewNop()}
	obs, _ := observability.NewObservability("test", log)

	svc := &orderItemCommandService{
		observability:       obs,
		cache:               &stubOrderItemCmdCache{},
		orderItemRepository: &stubOrderItemCmdRepo{},
		logger:              log,
	}

	result, err := svc.Trash(context.Background(), 3)
	if err != nil {
		t.Fatalf("Trash() error = %v", err)
	}
	if result.OrderItemID != 3 {
		t.Fatalf("OrderItemID = %d, want 3", result.OrderItemID)
	}
}

func TestOrderItemCommandServiceRestore(t *testing.T) {
	log := &logger.Logger{Log: zap.NewNop()}
	obs, _ := observability.NewObservability("test", log)

	svc := &orderItemCommandService{
		observability:       obs,
		cache:               &stubOrderItemCmdCache{},
		orderItemRepository: &stubOrderItemCmdRepo{},
		logger:              log,
	}

	result, err := svc.Restore(context.Background(), 3)
	if err != nil {
		t.Fatalf("Restore() error = %v", err)
	}
	if result.OrderItemID != 3 {
		t.Fatalf("OrderItemID = %d, want 3", result.OrderItemID)
	}
}

func TestOrderItemCommandServiceDeletePermanent(t *testing.T) {
	log := &logger.Logger{Log: zap.NewNop()}
	obs, _ := observability.NewObservability("test", log)

	svc := &orderItemCommandService{
		observability:       obs,
		cache:               &stubOrderItemCmdCache{},
		orderItemRepository: &stubOrderItemCmdRepo{},
		logger:              log,
	}

	ok, err := svc.DeletePermanent(context.Background(), 1)
	if err != nil {
		t.Fatalf("DeletePermanent() error = %v", err)
	}
	if !ok {
		t.Fatal("DeletePermanent() = false, want true")
	}
}

func TestOrderItemCommandServiceDeleteByOrderPermanent(t *testing.T) {
	log := &logger.Logger{Log: zap.NewNop()}
	obs, _ := observability.NewObservability("test", log)

	svc := &orderItemCommandService{
		observability:       obs,
		cache:               &stubOrderItemCmdCache{},
		orderItemRepository: &stubOrderItemCmdRepo{},
		logger:              log,
	}

	ok, err := svc.DeleteByOrderPermanent(context.Background(), 10)
	if err != nil {
		t.Fatalf("DeleteByOrderPermanent() error = %v", err)
	}
	if !ok {
		t.Fatal("DeleteByOrderPermanent() = false, want true")
	}
}

func TestOrderItemCommandServiceRestoreAll(t *testing.T) {
	log := &logger.Logger{Log: zap.NewNop()}
	obs, _ := observability.NewObservability("test", log)

	svc := &orderItemCommandService{
		observability:       obs,
		cache:               &stubOrderItemCmdCache{},
		orderItemRepository: &stubOrderItemCmdRepo{},
		logger:              log,
	}

	ok, err := svc.RestoreAll(context.Background())
	if err != nil {
		t.Fatalf("RestoreAll() error = %v", err)
	}
	if !ok {
		t.Fatal("RestoreAll() = false, want true")
	}
}

func TestOrderItemCommandServiceDeleteAll(t *testing.T) {
	log := &logger.Logger{Log: zap.NewNop()}
	obs, _ := observability.NewObservability("test", log)

	svc := &orderItemCommandService{
		observability:       obs,
		cache:               &stubOrderItemCmdCache{},
		orderItemRepository: &stubOrderItemCmdRepo{},
		logger:              log,
	}

	ok, err := svc.DeleteAll(context.Background())
	if err != nil {
		t.Fatalf("DeleteAll() error = %v", err)
	}
	if !ok {
		t.Fatal("DeleteAll() = false, want true")
	}
}

func TestOrderItemCommandServiceCalculateTotalPrice(t *testing.T) {
	log := &logger.Logger{Log: zap.NewNop()}
	obs, _ := observability.NewObservability("test", log)

	svc := &orderItemCommandService{
		observability:       obs,
		cache:               &stubOrderItemCmdCache{},
		orderItemRepository: &stubOrderItemCmdRepo{},
		logger:              log,
	}

	total, err := svc.CalculateTotalPrice(context.Background(), 10)
	if err != nil {
		t.Fatalf("CalculateTotalPrice() error = %v", err)
	}
	if total != 25000 {
		t.Fatalf("total = %d, want 25000", total)
	}
}

// --- Query Service Tests ---

func TestOrderItemQueryServiceFindAll(t *testing.T) {
	log := &logger.Logger{Log: zap.NewNop()}
	obs, _ := observability.NewObservability("test", log)

	svc := &orderItemQueryService{
		observability:       obs,
		cache:               &stubOrderItemQueryCache{},
		orderItemRepository: &stubOrderItemQueryRepo{},
		logger:              log,
	}

	data, total, err := svc.FindAll(context.Background(), &requests.FindAllOrderItems{
		Page: 1, PageSize: 10, Search: "",
	})
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

func TestOrderItemQueryServiceFindAllFromCache(t *testing.T) {
	log := &logger.Logger{Log: zap.NewNop()}
	obs, _ := observability.NewObservability("test", log)

	cached := []*repository.OrderItemResult{
		{OrderItemID: 42, TotalCount: 1},
	}
	svc := &orderItemQueryService{
		observability:       obs,
		cache:               &stubOrderItemQueryCache{data: cached, total: 1},
		orderItemRepository: &stubOrderItemQueryRepo{},
		logger:              log,
	}

	data, total, err := svc.FindAll(context.Background(), &requests.FindAllOrderItems{
		Page: 1, PageSize: 10, Search: "",
	})
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

func TestOrderItemQueryServiceFindActive(t *testing.T) {
	log := &logger.Logger{Log: zap.NewNop()}
	obs, _ := observability.NewObservability("test", log)

	svc := &orderItemQueryService{
		observability:       obs,
		cache:               &stubOrderItemQueryCache{},
		orderItemRepository: &stubOrderItemQueryRepo{},
		logger:              log,
	}

	data, total, err := svc.FindActive(context.Background(), &requests.FindAllOrderItems{
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

func TestOrderItemQueryServiceFindTrashed(t *testing.T) {
	log := &logger.Logger{Log: zap.NewNop()}
	obs, _ := observability.NewObservability("test", log)

	svc := &orderItemQueryService{
		observability:       obs,
		cache:               &stubOrderItemQueryCache{},
		orderItemRepository: &stubOrderItemQueryRepo{},
		logger:              log,
	}

	data, total, err := svc.FindTrashed(context.Background(), &requests.FindAllOrderItems{
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

func TestOrderItemQueryServiceFindByOrder(t *testing.T) {
	log := &logger.Logger{Log: zap.NewNop()}
	obs, _ := observability.NewObservability("test", log)

	svc := &orderItemQueryService{
		observability:       obs,
		cache:               &stubOrderItemQueryCache{},
		orderItemRepository: &stubOrderItemQueryRepo{},
		logger:              log,
	}

	data, err := svc.FindByOrder(context.Background(), 10)
	if err != nil {
		t.Fatalf("FindByOrder() error = %v", err)
	}
	if len(data) != 2 {
		t.Fatalf("len(data) = %d, want 2", len(data))
	}
	if data[0].OrderID != 10 {
		t.Fatalf("OrderID = %d, want 10", data[0].OrderID)
	}
}
