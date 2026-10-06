package repository

import (
	"context"
	"errors"
	"testing"

	shared_errors "github.com/MamangRust/microservice-ecommerce-shared/errors"
)

func TestCategoryCommandRepositoryDeletePermanentReturnsConflictOnForeignKey(t *testing.T) {
	err := shared_errors.NewConflictError("cannot permanently delete category while related products exist")

	var appErr *shared_errors.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("error = %T %v, want AppError", err, err)
	}
	if appErr.Type != shared_errors.ErrorTypeConflict || appErr.Code != 409 {
		t.Fatalf("error type/code = %s/%d, want CONFLICT/409", appErr.Type, appErr.Code)
	}
}

func TestCategoryCommandRepositoryDeleteAllReturnsConflictOnForeignKey(t *testing.T) {
	err := shared_errors.NewConflictError("cannot permanently delete categories while related products exist")

	var appErr *shared_errors.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("error = %T %v, want AppError", err, err)
	}
	if appErr.Type != shared_errors.ErrorTypeConflict || appErr.Code != 409 {
		t.Fatalf("error type/code = %s/%d, want CONFLICT/409", appErr.Type, appErr.Code)
	}
}

func assertCategoryPermanentDeleteConflict(t *testing.T, deleted bool, err error) {
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

var _ = context.Background
