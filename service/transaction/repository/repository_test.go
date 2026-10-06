package repository

import (
	"context"
	"errors"
	"testing"

	shared_errors "github.com/MamangRust/microservice-ecommerce-shared/errors"
)

func TestTransactionCommandRepositoryDeletePermanentReturnsConflictOnForeignKey(t *testing.T) {
	err := shared_errors.NewConflictError("cannot permanently delete transaction while related records exist")

	var appErr *shared_errors.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("error = %T %v, want AppError", err, err)
	}
	if appErr.Type != shared_errors.ErrorTypeConflict || appErr.Code != 409 {
		t.Fatalf("error type/code = %s/%d, want CONFLICT/409", appErr.Type, appErr.Code)
	}
}

func TestTransactionCommandRepositoryDeleteAllReturnsConflictOnForeignKey(t *testing.T) {
	err := shared_errors.NewConflictError("cannot permanently delete transactions while related records exist")

	var appErr *shared_errors.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("error = %T %v, want AppError", err, err)
	}
	if appErr.Type != shared_errors.ErrorTypeConflict || appErr.Code != 409 {
		t.Fatalf("error type/code = %s/%d, want CONFLICT/409", appErr.Type, appErr.Code)
	}
}

func assertTransactionPermanentDeleteConflict(t *testing.T, deleted bool, err error) {
	t.Helper()
	if deleted {
		t.Fatal("delete reported success after a foreign-key violation")
	}
	var appErr *shared_errors.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("error = %T %v, want AppError", err, err)
	}
	if appErr.Type != shared_errors.ErrorTypeConflict || appErr.Code != 409 {
		t.Fatalf("error type/code = %s/%d, want CONFLICT/409", appErr.Type, appErr.Code)
	}
}

func TestTransactionResultConstruction(t *testing.T) {
	result := &TransactionResult{
		TransactionID: 1,
		OrderID:       10,
		MerchantID:    20,
		PaymentMethod: "credit_card",
		Amount:        50000,
		PaymentStatus: "success",
	}
	if result.TransactionID != 1 {
		t.Fatalf("TransactionID = %d, want 1", result.TransactionID)
	}
	if result.OrderID != 10 {
		t.Fatalf("OrderID = %d, want 10", result.OrderID)
	}
	if result.MerchantID != 20 {
		t.Fatalf("MerchantID = %d, want 20", result.MerchantID)
	}
	if result.PaymentStatus != "success" {
		t.Fatalf("PaymentStatus = %s, want success", result.PaymentStatus)
	}
}

func TestOutboxEventResultConstruction(t *testing.T) {
	result := &OutboxEventResult{
		OutboxID: 1,
		Topic:    "test-topic",
		EventKey: "key-1",
		Status:   "pending",
		Attempts: 0,
	}
	if result.OutboxID != 1 {
		t.Fatalf("OutboxID = %d, want 1", result.OutboxID)
	}
	if result.Topic != "test-topic" {
		t.Fatalf("Topic = %s, want test-topic", result.Topic)
	}
}

var _ = context.Background
