package service

import (
	"context"
	"fmt"
	"testing"

	dto "github.com/MamangRust/microservice-ecommerce-grpc-transaction/dto"
	"github.com/MamangRust/microservice-ecommerce-grpc-transaction/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/observability"
	"go.uber.org/zap"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// --- Stubs for repositories ---

type stubTxQueryRepo struct{}

func (s *stubTxQueryRepo) FindAll(_ context.Context, _ *requests.FindAllTransaction) ([]*repository.TransactionResult, error) {
	return nil, fmt.Errorf("unexpected FindAll")
}
func (s *stubTxQueryRepo) FindActive(_ context.Context, _ *requests.FindAllTransaction) ([]*repository.TransactionResult, error) {
	return nil, fmt.Errorf("unexpected FindActive")
}
func (s *stubTxQueryRepo) FindTrashed(_ context.Context, _ *requests.FindAllTransaction) ([]*repository.TransactionResult, error) {
	return nil, fmt.Errorf("unexpected FindTrashed")
}
func (s *stubTxQueryRepo) FindByMerchant(_ context.Context, _ *requests.FindAllTransactionByMerchant) ([]*repository.TransactionResult, error) {
	return nil, fmt.Errorf("unexpected FindByMerchant")
}
func (s *stubTxQueryRepo) FindByID(_ context.Context, id int) (*repository.TransactionResult, error) {
	if id == 1 {
		return &repository.TransactionResult{
			TransactionID: 1, OrderID: 10, MerchantID: 20,
			PaymentMethod: "credit_card", Amount: 50000,
			PaymentStatus: "success",
		}, nil
	}
	return nil, fmt.Errorf("transaction not found")
}
func (s *stubTxQueryRepo) FindByOrderID(_ context.Context, _ int) (*repository.TransactionResult, error) {
	return nil, fmt.Errorf("unexpected FindByOrderID")
}

type stubTxCmdRepo struct{}

func (s *stubTxCmdRepo) Create(_ context.Context, req *requests.CreateTransactionRequest) (*repository.TransactionResult, error) {
	return &repository.TransactionResult{
		TransactionID: 99, OrderID: int32(req.OrderID), MerchantID: int32(req.MerchantID),
		PaymentMethod: req.PaymentMethod, Amount: int32(req.Amount),
		PaymentStatus: "success",
	}, nil
}
func (s *stubTxCmdRepo) CreateInTx(_ context.Context, _ *mongo.Database, req *requests.CreateTransactionRequest) (*repository.TransactionResult, error) {
	return s.Create(nil, req)
}
func (s *stubTxCmdRepo) Update(_ context.Context, req *requests.UpdateTransactionRequest) (*repository.TransactionResult, error) {
	return &repository.TransactionResult{
		TransactionID: int32(*req.TransactionID), OrderID: int32(req.OrderID),
		MerchantID: int32(req.MerchantID), PaymentMethod: req.PaymentMethod,
		Amount: int32(req.Amount), PaymentStatus: *req.PaymentStatus,
	}, nil
}
func (s *stubTxCmdRepo) Trash(_ context.Context, id int) (*repository.TransactionResult, error) {
	return &repository.TransactionResult{TransactionID: int32(id), PaymentStatus: "success"}, nil
}
func (s *stubTxCmdRepo) Restore(_ context.Context, id int) (*repository.TransactionResult, error) {
	return &repository.TransactionResult{TransactionID: int32(id), PaymentStatus: "success"}, nil
}
func (s *stubTxCmdRepo) DeletePermanent(_ context.Context, _ int) (bool, error) {
	return true, nil
}
func (s *stubTxCmdRepo) DeleteByOrderIDPermanent(_ context.Context, _ int) (bool, error) {
	return true, nil
}
func (s *stubTxCmdRepo) RestoreAll(_ context.Context) (bool, error) {
	return true, nil
}
func (s *stubTxCmdRepo) DeleteAll(_ context.Context) (bool, error) {
	return true, nil
}

type stubUserQuery struct{}

func (s *stubUserQuery) FindByID(_ context.Context, id int) (*models.User, error) {
	return &models.User{UserID: int32(id), Email: "test@example.com", Firstname: "Test", Lastname: "User"}, nil
}
func (s *stubUserQuery) FindByEmail(_ context.Context, _ string) (*models.User, error) {
	return nil, fmt.Errorf("not used in this test")
}
func (s *stubUserQuery) FindByEmailAndVerify(_ context.Context, _ string) (*models.User, error) {
	return nil, fmt.Errorf("not used in this test")
}
func (s *stubUserQuery) FindByVerificationCode(_ context.Context, _ string) (*models.User, error) {
	return nil, fmt.Errorf("not used in this test")
}

type stubMerchantQuery struct{}

func (s *stubMerchantQuery) FindByID(_ context.Context, id int) (*dto.GetMerchantByIDRow, error) {
	return &dto.GetMerchantByIDRow{MerchantID: int32(id), Name: "Merchant"}, nil
}

type stubOrderQuery struct{}

func (s *stubOrderQuery) FindByID(_ context.Context, id int) (*models.Order, error) {
	return &models.Order{OrderID: int32(id)}, nil
}

type stubOrderItem struct{}

func (s *stubOrderItem) FindOrderItemByOrder(_ context.Context, orderID int) ([]*models.OrderItem, error) {
	return []*models.OrderItem{
		{OrderItemID: 1, OrderID: int32(orderID), ProductID: 1, Quantity: 2, Price: 10000},
	}, nil
}

func (s *stubOrderItem) CalculateTotalPrice(_ context.Context, _ int) (*int32, error) {
	total := int32(20000)
	return &total, nil
}

type stubShippingQuery struct{}

func (s *stubShippingQuery) FindByID(_ context.Context, _ int) (*models.ShippingAddress, error) {
	return &models.ShippingAddress{ShippingCost: 5000}, nil
}

func (s *stubShippingQuery) FindByOrder(_ context.Context, _ int) (*models.ShippingAddress, error) {
	return &models.ShippingAddress{ShippingCost: 5000}, nil
}

type stubTxCmdCache struct{}

func (s *stubTxCmdCache) DeleteTransactionCache(_ context.Context, _ int) {}
func (s *stubTxCmdCache) InvalidateTransactionCache(_ context.Context)    {}

// --- Tests ---

func TestTransactionCommandServiceCreate(t *testing.T) {
	log := &logger.Logger{Log: zap.NewNop()}
	obs, _ := observability.NewObservability("test", log)

	svc := &transactionCommandService{
		observability:      obs,
		cache:              &stubTxCmdCache{},
		transactionQuery:   &stubTxQueryRepo{},
		transactionCommand: &stubTxCmdRepo{},
		userQuery:          &stubUserQuery{},
		merchantQuery:      &stubMerchantQuery{},
		orderQuery:         &stubOrderQuery{},
		orderItem:          &stubOrderItem{},
		shippingAddress:    &stubShippingQuery{},
		logger:             log,
	}

	status := "success"
	result, err := svc.Create(context.Background(), &requests.CreateTransactionRequest{
		OrderID:       10,
		MerchantID:    20,
		UserID:        30,
		PaymentMethod: "credit_card",
		Amount:        30000,
		PaymentStatus: &status,
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if result.TransactionID != 99 {
		t.Fatalf("TransactionID = %d, want 99", result.TransactionID)
	}
}

func TestTransactionCommandServiceTrash(t *testing.T) {
	log := &logger.Logger{Log: zap.NewNop()}
	obs, _ := observability.NewObservability("test", log)

	svc := &transactionCommandService{
		observability:      obs,
		cache:              &stubTxCmdCache{},
		transactionCommand: &stubTxCmdRepo{},
		logger:             log,
	}

	result, err := svc.Trash(context.Background(), 5)
	if err != nil {
		t.Fatalf("Trash() error = %v", err)
	}
	if result.TransactionID != 5 {
		t.Fatalf("TransactionID = %d, want 5", result.TransactionID)
	}
}

func TestTransactionCommandServiceRestore(t *testing.T) {
	log := &logger.Logger{Log: zap.NewNop()}
	obs, _ := observability.NewObservability("test", log)

	svc := &transactionCommandService{
		observability:      obs,
		cache:              &stubTxCmdCache{},
		transactionCommand: &stubTxCmdRepo{},
		logger:             log,
	}

	result, err := svc.Restore(context.Background(), 5)
	if err != nil {
		t.Fatalf("Restore() error = %v", err)
	}
	if result.TransactionID != 5 {
		t.Fatalf("TransactionID = %d, want 5", result.TransactionID)
	}
}

func TestTransactionCommandServiceDeletePermanent(t *testing.T) {
	log := &logger.Logger{Log: zap.NewNop()}
	obs, _ := observability.NewObservability("test", log)

	svc := &transactionCommandService{
		observability:      obs,
		cache:              &stubTxCmdCache{},
		transactionCommand: &stubTxCmdRepo{},
		logger:             log,
	}

	ok, err := svc.DeletePermanent(context.Background(), 1)
	if err != nil {
		t.Fatalf("DeletePermanent() error = %v", err)
	}
	if !ok {
		t.Fatal("DeletePermanent() = false, want true")
	}
}

func TestTransactionCommandServiceDeleteAll(t *testing.T) {
	log := &logger.Logger{Log: zap.NewNop()}
	obs, _ := observability.NewObservability("test", log)

	svc := &transactionCommandService{
		observability:      obs,
		cache:              &stubTxCmdCache{},
		transactionCommand: &stubTxCmdRepo{},
		logger:             log,
	}

	ok, err := svc.DeleteAll(context.Background())
	if err != nil {
		t.Fatalf("DeleteAll() error = %v", err)
	}
	if !ok {
		t.Fatal("DeleteAll() = false, want true")
	}
}

func TestTransactionCommandServiceRestoreAll(t *testing.T) {
	log := &logger.Logger{Log: zap.NewNop()}
	obs, _ := observability.NewObservability("test", log)

	svc := &transactionCommandService{
		observability:      obs,
		cache:              &stubTxCmdCache{},
		transactionCommand: &stubTxCmdRepo{},
		logger:             log,
	}

	ok, err := svc.RestoreAll(context.Background())
	if err != nil {
		t.Fatalf("RestoreAll() error = %v", err)
	}
	if !ok {
		t.Fatal("RestoreAll() = false, want true")
	}
}

func TestTransactionCommandServiceDeleteByOrderIDPermanent(t *testing.T) {
	log := &logger.Logger{Log: zap.NewNop()}
	obs, _ := observability.NewObservability("test", log)

	svc := &transactionCommandService{
		observability:      obs,
		cache:              &stubTxCmdCache{},
		transactionCommand: &stubTxCmdRepo{},
		logger:             log,
	}

	ok, err := svc.DeleteByOrderIDPermanent(context.Background(), 10)
	if err != nil {
		t.Fatalf("DeleteByOrderIDPermanent() error = %v", err)
	}
	if !ok {
		t.Fatal("DeleteByOrderIDPermanent() = false, want true")
	}
}

func TestTransactionCommandServiceUpdate(t *testing.T) {
	log := &logger.Logger{Log: zap.NewNop()}
	obs, _ := observability.NewObservability("test", log)

	svc := &transactionCommandService{
		observability:      obs,
		cache:              &stubTxCmdCache{},
		transactionQuery:   &stubTxQueryRepo{},
		transactionCommand: &stubTxCmdRepo{},
		merchantQuery:      &stubMerchantQuery{},
		orderQuery:         &stubOrderQuery{},
		orderItem:          &stubOrderItem{},
		shippingAddress:    &stubShippingQuery{},
		logger:             log,
	}

	txID := 1
	status := "success"
	result, err := svc.Update(context.Background(), &requests.UpdateTransactionRequest{
		TransactionID: &txID,
		MerchantID:    20,
		OrderID:       10,
		PaymentMethod: "credit_card",
		Amount:        60000,
		PaymentStatus: &status,
	})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if result.TransactionID != 1 {
		t.Fatalf("TransactionID = %d, want 1", result.TransactionID)
	}
}
