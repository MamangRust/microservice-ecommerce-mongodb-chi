package handler

import (
	"context"
	"fmt"
	pbmerchant_document "github.com/MamangRust/microservice-ecommerce-grpc-pb/merchant_document"
	"testing"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"go.uber.org/zap"
)

type merchantDocumentCommandServiceStub struct {
	update       func(context.Context, *requests.UpdateMerchantDocumentRequest) (*models.MerchantDocument, error)
	updateStatus func(context.Context, *requests.UpdateMerchantDocumentStatusRequest) (*models.MerchantDocument, error)
}

func (s merchantDocumentCommandServiceStub) Create(context.Context, *requests.CreateMerchantDocumentRequest) (*models.MerchantDocument, error) {
	return nil, fmt.Errorf("unexpected Create call")
}

func (s merchantDocumentCommandServiceStub) Update(ctx context.Context, req *requests.UpdateMerchantDocumentRequest) (*models.MerchantDocument, error) {
	if s.update == nil {
		return nil, fmt.Errorf("unexpected Update call")
	}
	return s.update(ctx, req)
}

func (s merchantDocumentCommandServiceStub) UpdateStatus(ctx context.Context, req *requests.UpdateMerchantDocumentStatusRequest) (*models.MerchantDocument, error) {
	if s.updateStatus == nil {
		return nil, fmt.Errorf("unexpected UpdateStatus call")
	}
	return s.updateStatus(ctx, req)
}

func (s merchantDocumentCommandServiceStub) Trash(context.Context, int) (*models.MerchantDocument, error) {
	return nil, fmt.Errorf("unexpected Trash call")
}

func (s merchantDocumentCommandServiceStub) Restore(context.Context, int) (*models.MerchantDocument, error) {
	return nil, fmt.Errorf("unexpected Restore call")
}

func (s merchantDocumentCommandServiceStub) DeletePermanent(context.Context, int) (bool, error) {
	return false, fmt.Errorf("unexpected DeletePermanent call")
}

func (s merchantDocumentCommandServiceStub) RestoreAll(context.Context) (bool, error) {
	return false, fmt.Errorf("unexpected RestoreAll call")
}

func (s merchantDocumentCommandServiceStub) DeleteAll(context.Context) (bool, error) {
	return false, fmt.Errorf("unexpected DeleteAll call")
}

func newMerchantDocumentCommandHandlerForTest(stub merchantDocumentCommandServiceStub) pbmerchant_document.MerchantDocumentCommandServiceServer {
	return NewMerchantDocumentCommandHandler(stub, &logger.Logger{Log: zap.NewNop()})
}

func TestMerchantDocumentCommandHandlerUpdateUsesDocumentID(t *testing.T) {
	const (
		documentID = 42
		merchantID = 7
	)

	var gotRequest *requests.UpdateMerchantDocumentRequest
	handler := newMerchantDocumentCommandHandlerForTest(merchantDocumentCommandServiceStub{
		update: func(_ context.Context, req *requests.UpdateMerchantDocumentRequest) (*models.MerchantDocument, error) {
			gotRequest = req
			note := req.Note
			now := time.Now()
			return &models.MerchantDocument{
				DocumentID:   int32(*req.DocumentID),
				MerchantID:   int32(req.MerchantID),
				DocumentType: req.DocumentType,
				DocumentUrl:  req.DocumentUrl,
				Status:       req.Status,
				Note:         &note,
				CreatedAt:    &now,
				UpdatedAt:    &now,
			}, nil
		},
	})

	got, err := handler.Update(context.Background(), &pbmerchant_document.UpdateMerchantDocumentRequest{
		DocumentId:   documentID,
		MerchantId:   merchantID,
		DocumentType: "business_license",
		DocumentUrl:  "https://example.com/license.pdf",
		Status:       "pending",
		Note:         "replace before review",
	})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if gotRequest == nil {
		t.Fatal("Update() did not call the command service")
	}
	if gotRequest.DocumentID == nil || *gotRequest.DocumentID != documentID {
		t.Fatalf("service received DocumentID %v, want %d", gotRequest.DocumentID, documentID)
	}
	if gotRequest.MerchantID != merchantID {
		t.Fatalf("service received MerchantID %d, want %d", gotRequest.MerchantID, merchantID)
	}
	if got.GetStatus() != "success" {
		t.Fatalf("response status = %q, want success", got.GetStatus())
	}
}

func TestMerchantDocumentCommandHandlerUpdateStatusUsesDocumentID(t *testing.T) {
	const (
		documentID = 84
		merchantID = 13
	)

	var gotRequest *requests.UpdateMerchantDocumentStatusRequest
	handler := newMerchantDocumentCommandHandlerForTest(merchantDocumentCommandServiceStub{
		updateStatus: func(_ context.Context, req *requests.UpdateMerchantDocumentStatusRequest) (*models.MerchantDocument, error) {
			gotRequest = req
			note := req.Note
			now := time.Now()
			return &models.MerchantDocument{
				DocumentID: int32(*req.DocumentID),
				MerchantID: int32(req.MerchantID),
				Status:     req.Status,
				Note:       &note,
				CreatedAt:  &now,
				UpdatedAt:  &now,
			}, nil
		},
	})

	got, err := handler.UpdateStatus(context.Background(), &pbmerchant_document.UpdateMerchantDocumentStatusRequest{
		DocumentId: documentID,
		MerchantId: merchantID,
		Status:     "approved",
		Note:       "verified",
	})
	if err != nil {
		t.Fatalf("UpdateStatus() error = %v", err)
	}
	if gotRequest == nil {
		t.Fatal("UpdateStatus() did not call the command service")
	}
	if gotRequest.DocumentID == nil || *gotRequest.DocumentID != documentID {
		t.Fatalf("service received DocumentID %v, want %d", gotRequest.DocumentID, documentID)
	}
	if gotRequest.MerchantID != merchantID {
		t.Fatalf("service received MerchantID %d, want %d", gotRequest.MerchantID, merchantID)
	}
	if got.GetStatus() != "success" {
		t.Fatalf("response status = %q, want success", got.GetStatus())
	}
}
