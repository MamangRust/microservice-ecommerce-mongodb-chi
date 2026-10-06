package repository

import (
	"testing"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
)

func TestProductResultConstruction(t *testing.T) {
	result := &ProductResult{
		ProductID:    1,
		MerchantID:   10,
		CategoryID:   20,
		Name:         "Test Product",
		Price:        50000,
		CountInStock: 100,
	}

	if result.ProductID != 1 {
		t.Errorf("expected ProductID=1, got %d", result.ProductID)
	}
	if result.Name != "Test Product" {
		t.Errorf("expected Name='Test Product', got %s", result.Name)
	}
}

func TestProductCommandRepositoryCreateArgs(t *testing.T) {
	product := &models.Product{
		MerchantID:   1,
		CategoryID:   1,
		Name:         "Test",
		Price:        50000,
		CountInStock: 100,
	}

	if product.Name != "Test" {
		t.Errorf("expected Name='Test', got %s", product.Name)
	}
	if product.Price != 50000 {
		t.Errorf("expected Price=50000, got %d", product.Price)
	}
}
