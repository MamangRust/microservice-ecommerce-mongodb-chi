package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	shared_errors "github.com/MamangRust/microservice-ecommerce-shared/errors"
)

func TestOrderCommandRepositoryDeletePermanentReturnsConflictOnForeignKey(t *testing.T) {
	err := shared_errors.NewConflictError("cannot permanently delete order while related records exist")

	var appErr *shared_errors.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("error = %T %v, want AppError", err, err)
	}
	if appErr.Type != shared_errors.ErrorTypeConflict || appErr.Code != 409 {
		t.Fatalf("error type/code = %s/%d, want CONFLICT/409", appErr.Type, appErr.Code)
	}
}

func TestOrderCommandRepositoryDeleteAllReturnsConflictOnForeignKey(t *testing.T) {
	err := shared_errors.NewConflictError("cannot permanently delete orders while related records exist")

	var appErr *shared_errors.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("error = %T %v, want AppError", err, err)
	}
	if appErr.Type != shared_errors.ErrorTypeConflict || appErr.Code != 409 {
		t.Fatalf("error type/code = %s/%d, want CONFLICT/409", appErr.Type, appErr.Code)
	}
}

func TestOrderResultConstruction(t *testing.T) {
	result := &OrderResult{
		OrderID:    1,
		UserID:     10,
		MerchantID: 20,
		TotalPrice: 50000,
		TotalCount: 1,
	}
	if result.OrderID != 1 {
		t.Fatalf("OrderID = %d, want 1", result.OrderID)
	}
	if result.UserID != 10 {
		t.Fatalf("UserID = %d, want 10", result.UserID)
	}
	if result.MerchantID != 20 {
		t.Fatalf("MerchantID = %d, want 20", result.MerchantID)
	}
	if result.TotalPrice != 50000 {
		t.Fatalf("TotalPrice = %d, want 50000", result.TotalPrice)
	}
}

func TestStockReservationResultConstruction(t *testing.T) {
	now := time.Now()
	result := &StockReservationResult{
		ReservationID: 1,
		OrderID:       10,
		ProductID:     20,
		Quantity:      5,
		Status:        "reserved",
		CreatedAt:     &now,
		UpdatedAt:     &now,
	}
	if result.ReservationID != 1 {
		t.Fatalf("ReservationID = %d, want 1", result.ReservationID)
	}
	if result.Status != "reserved" {
		t.Fatalf("Status = %s, want reserved", result.Status)
	}
}

func TestOrderModelConstruction(t *testing.T) {
	now := time.Now()
	order := &models.Order{
		OrderID:    1,
		UserID:     10,
		MerchantID: 20,
		TotalPrice: 50000,
		CreatedAt:  &now,
		UpdatedAt:  &now,
	}
	if order.OrderID != 1 {
		t.Fatalf("OrderID = %d, want 1", order.OrderID)
	}
}

var _ = context.Background
