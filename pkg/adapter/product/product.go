package product

import (
	"context"
	"time"

	pbproduct "github.com/MamangRust/microservice-ecommerce-grpc-pb/product"
	"github.com/MamangRust/microservice-ecommerce-pkg/adapter"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-pkg/resilience"
	product_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/product_errors"
)

// QueryRepository is the contract consumers depend on for product reads.
type QueryRepository interface {
	FindByID(ctx context.Context, productID int) (*models.Product, error)
}

// CommandRepository is the contract consumers depend on for product writes
// (stock adjustments). CreateProduct/UpdateProduct live entirely inside the
// owning service and are not exposed here.
type CommandRepository interface {
	UpdateProductCountStock(ctx context.Context, productID int, stock int) (*models.Product, error)
	AdjustProductStock(ctx context.Context, productID int, delta int, operationID string) (*models.Product, error)
}

// Repository implements both QueryRepository and CommandRepository on top of the
// generated product query/command clients, guarding every call.
type Repository struct {
	query   pbproduct.ProductQueryServiceClient
	command pbproduct.ProductCommandServiceClient
	guard   *resilience.DependencyGuard
}

// New builds the product adapter from its query and command clients.
func New(query pbproduct.ProductQueryServiceClient, command pbproduct.ProductCommandServiceClient, opts ...adapter.GuardOption) *Repository {
	r := &Repository{query: query, command: command}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// NewAdapter wraps the generated product clients.
func NewAdapter(query pbproduct.ProductQueryServiceClient, command pbproduct.ProductCommandServiceClient, opts ...adapter.GuardOption) *Repository {
	return New(query, command, opts...)
}

// NewQueryAdapter returns the adapter restricted to the QueryRepository surface
// for consumers that only read (cart, review).
func NewQueryAdapter(query pbproduct.ProductQueryServiceClient, opts ...adapter.GuardOption) *Repository {
	return New(query, nil, opts...)
}

// NewCommandAdapter returns the adapter restricted to the CommandRepository surface
// for consumers that only write (currently none — included for symmetry).
func NewCommandAdapter(command pbproduct.ProductCommandServiceClient, opts ...adapter.GuardOption) *Repository {
	return New(nil, command, opts...)
}

// SetGuard attaches the dependency guard used to wrap every outbound call.
func (r *Repository) SetGuard(g *resilience.DependencyGuard) {
	r.guard = g
}

func (r *Repository) FindByID(ctx context.Context, productID int) (*models.Product, error) {
	res, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbproduct.ApiResponseProduct, error) {
		return r.query.FindById(callCtx, &pbproduct.FindByIdProductRequest{Id: int32(productID)})
	})
	if err != nil {
		return nil, product_errors.ErrProductNotFound.WithInternal(err)
	}
	return toModel(res.Data), nil
}

func (r *Repository) UpdateProductCountStock(ctx context.Context, productID int, stock int) (*models.Product, error) {
	res, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbproduct.ApiResponseProduct, error) {
		return r.command.UpdateProductCountStock(callCtx, &pbproduct.UpdateProductCountStockRequest{
			ProductId: int32(productID),
			Stock:     int32(stock),
		})
	})
	if err != nil {
		return nil, product_errors.ErrUpdateProductCountStock.WithInternal(err)
	}
	return toModel(res.Data), nil
}

func (r *Repository) AdjustProductStock(ctx context.Context, productID int, delta int, operationID string) (*models.Product, error) {
	res, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbproduct.ApiResponseProduct, error) {
		return r.command.AdjustProductStock(callCtx, &pbproduct.AdjustProductStockRequest{
			ProductId:   int32(productID),
			Delta:       int32(delta),
			OperationId: operationID,
		})
	})
	if err != nil {
		return nil, product_errors.ErrUpdateProductCountStock.WithInternal(err)
	}
	return toModel(res.Data), nil
}

func toModel(data *pbproduct.ProductResponse) *models.Product {
	if data == nil {
		return nil
	}
	description := data.Description
	weight := data.Weight
	image := data.ImageProduct
	return &models.Product{
		ProductID:    data.Id,
		MerchantID:   data.MerchantId,
		CategoryID:   data.CategoryId,
		Name:         data.Name,
		Description:  &description,
		Price:        data.Price,
		CountInStock: data.CountInStock,
		Weight:       &weight,
		ImageProduct: &image,
		CreatedAt:    parseTime(data.CreatedAt),
		UpdatedAt:    parseTime(data.UpdatedAt),
	}
}

// parseTime accepts the loose RFC3339-ish timestamp strings the product service
// emits and returns nil for empty strings. Anything that fails to parse is treated
// as nil rather than an error so the caller keeps a usable *models.Product with
// just the textual fields populated.
func parseTime(s string) *time.Time {
	if s == "" {
		return nil
	}
	for _, layout := range []string{time.RFC3339Nano, time.RFC3339, "2006-01-02 15:04:05"} {
		if t, err := time.Parse(layout, s); err == nil {
			return &t
		}
	}
	return nil
}

var (
	_ QueryRepository   = (*Repository)(nil)
	_ CommandRepository = (*Repository)(nil)
)
