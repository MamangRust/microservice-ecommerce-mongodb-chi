package merchant

import (
	"context"
	"time"

	pbmerchant "github.com/MamangRust/microservice-ecommerce-grpc-pb/merchant"
	"github.com/MamangRust/microservice-ecommerce-pkg/adapter"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-pkg/resilience"
	merchant_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/merchant"
)

// QueryRepository is the contract consumers depend on for merchant reads.
// It is implemented by Repository and can be satisfied by a test double.
type QueryRepository interface {
	FindByID(ctx context.Context, merchantID int) (*models.Merchant, error)
}

// Repository implements QueryRepository on top of the generated merchant query
// client, guarding every call.
type Repository struct {
	client pbmerchant.MerchantQueryServiceClient
	guard  *resilience.DependencyGuard
}

// New builds the merchant adapter from its query client.
func New(client pbmerchant.MerchantQueryServiceClient, opts ...adapter.GuardOption) *Repository {
	r := &Repository{client: client}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// NewQueryAdapter wraps a generated merchant query client into a QueryRepository.
func NewQueryAdapter(client pbmerchant.MerchantQueryServiceClient, opts ...adapter.GuardOption) QueryRepository {
	return New(client, opts...)
}

// SetGuard attaches the dependency guard used to wrap every outbound call.
func (r *Repository) SetGuard(g *resilience.DependencyGuard) {
	r.guard = g
}

func (r *Repository) FindByID(ctx context.Context, merchantID int) (*models.Merchant, error) {
	res, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbmerchant.ApiResponseMerchant, error) {
		return r.client.FindById(callCtx, &pbmerchant.FindByIdMerchantRequest{Id: int32(merchantID)})
	})
	if err != nil {
		return nil, merchant_errors.ErrMerchantNotFound.WithInternal(err)
	}
	return toModel(res.Data), nil
}

func toModel(m *pbmerchant.MerchantResponse) *models.Merchant {
	if m == nil {
		return nil
	}
	return &models.Merchant{
		MerchantID:   m.Id,
		UserID:       m.UserId,
		Name:         m.Name,
		Description:  &m.Description,
		Address:      &m.Address,
		ContactEmail: &m.ContactEmail,
		ContactPhone: &m.ContactPhone,
		Status:       m.Status,
		CreatedAt:    parseTimeString(m.CreatedAt),
		UpdatedAt:    parseTimeString(m.UpdatedAt),
	}
}

func parseTimeString(s string) *time.Time {
	if s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02 15:04:05.000", s)
	if err != nil {
		t, err = time.Parse("2006-01-02T15:04:05Z", s)
		if err != nil {
			return nil
		}
	}
	return &t
}

var _ QueryRepository = (*Repository)(nil)
