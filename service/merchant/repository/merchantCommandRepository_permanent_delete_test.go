package repository

import (
	"testing"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
)

func TestMerchantResultConstruction(t *testing.T) {
	r := &MerchantResult{
		MerchantID:   1,
		UserID:       10,
		Name:         "Test Merchant",
		Description:  strPtr("desc"),
		Address:      strPtr("addr"),
		ContactEmail: strPtr("e@e.com"),
		ContactPhone: strPtr("123"),
		Status:       "active",
		TotalCount:   5,
	}
	if r.MerchantID != 1 || r.Name != "Test Merchant" {
		t.Fatal("MerchantResult fields mismatch")
	}
}

func TestMerchantDocumentResultConstruction(t *testing.T) {
	r := &MerchantDocumentResult{
		DocumentID:   1,
		MerchantID:   10,
		DocumentType: "license",
		DocumentUrl:  "https://example.com/license.pdf",
		Status:       "pending",
		TotalCount:   3,
	}
	if r.DocumentID != 1 || r.DocumentType != "license" {
		t.Fatal("MerchantDocumentResult fields mismatch")
	}
}

func TestMerchantModelConstruction(t *testing.T) {
	name := "Test"
	desc := "desc"
	m := &models.Merchant{
		MerchantID:  1,
		UserID:      10,
		Name:        name,
		Description: &desc,
		Status:      "active",
	}
	if m.MerchantID != 1 || m.Name != "Test" {
		t.Fatal("Merchant model fields mismatch")
	}
}

func TestMerchantDocumentModelConstruction(t *testing.T) {
	note := "test note"
	d := &models.MerchantDocument{
		DocumentID:   1,
		MerchantID:   10,
		DocumentType: "business_license",
		DocumentUrl:  "https://example.com/doc.pdf",
		Status:       "pending",
		Note:         &note,
	}
	if d.DocumentID != 1 || d.DocumentType != "business_license" {
		t.Fatal("MerchantDocument model fields mismatch")
	}
}

func TestStringPtr(t *testing.T) {
	if s := stringPtr("hello"); s == nil || *s != "hello" {
		t.Fatal("stringPtr(\"hello\") failed")
	}
	if s := stringPtr(""); s != nil {
		t.Fatal("stringPtr(\"\") should return nil")
	}
}

func strPtr(s string) *string { return &s }
