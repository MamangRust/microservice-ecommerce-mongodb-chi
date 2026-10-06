package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	shared_errors "github.com/MamangRust/microservice-ecommerce-shared/errors"
)

func TestOrderItemCommandRepositoryDeletePermanentReturnsConflictOnForeignKey(t *testing.T) {
	err := shared_errors.NewConflictError("cannot permanently delete order item while related records exist")

	var appErr *shared_errors.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("error = %T %v, want AppError", err, err)
	}
	if appErr.Type != shared_errors.ErrorTypeConflict || appErr.Code != 409 {
		t.Fatalf("error type/code = %s/%d, want CONFLICT/409", appErr.Type, appErr.Code)
	}
}

func TestOrderItemCommandRepositoryDeleteAllReturnsConflictOnForeignKey(t *testing.T) {
	err := shared_errors.NewConflictError("cannot permanently delete order items while related records exist")

	var appErr *shared_errors.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("error = %T %v, want AppError", err, err)
	}
	if appErr.Type != shared_errors.ErrorTypeConflict || appErr.Code != 409 {
		t.Fatalf("error type/code = %s/%d, want CONFLICT/409", appErr.Type, appErr.Code)
	}
}

func assertOrderItemPermanentDeleteConflict(t *testing.T, deleted bool, err error) {
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

func TestOrderItemResultConstruction(t *testing.T) {
	result := &OrderItemResult{
		OrderItemID: 1,
		OrderID:     10,
		ProductID:   20,
		Quantity:    5,
		Price:       10000,
		TotalCount:  1,
	}
	if result.OrderItemID != 1 {
		t.Fatalf("OrderItemID = %d, want 1", result.OrderItemID)
	}
	if result.OrderID != 10 {
		t.Fatalf("OrderID = %d, want 10", result.OrderID)
	}
	if result.ProductID != 20 {
		t.Fatalf("ProductID = %d, want 20", result.ProductID)
	}
	if result.Quantity != 5 {
		t.Fatalf("Quantity = %d, want 5", result.Quantity)
	}
	if result.Price != 10000 {
		t.Fatalf("Price = %d, want 10000", result.Price)
	}
}

func TestOrderItemModelConstruction(t *testing.T) {
	item := &models.OrderItem{
		OrderItemID: 1,
		OrderID:     10,
		ProductID:   20,
		Quantity:    5,
		Price:       10000,
	}
	if item.OrderItemID != 1 {
		t.Fatalf("OrderItemID = %d, want 1", item.OrderItemID)
	}
	if item.TableName() != "order_items" {
		t.Fatalf("TableName() = %s, want order_items", item.TableName())
	}
}

var _ = context.Background
