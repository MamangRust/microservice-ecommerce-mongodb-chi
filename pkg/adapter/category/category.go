package category

import (
	"context"
	"time"

	pbcategory "github.com/MamangRust/microservice-ecommerce-grpc-pb/category"
	"github.com/MamangRust/microservice-ecommerce-pkg/adapter"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-pkg/resilience"
	"github.com/MamangRust/microservice-ecommerce-shared/errors/category_errors"
)

// QueryRepository is the contract consumers depend on for category reads.
// It is implemented by Repository and can be satisfied by a test double.
type QueryRepository interface {
	FindByID(ctx context.Context, categoryID int) (*models.Category, error)
	FindAll(ctx context.Context, page, pageSize int, search string) ([]*models.Category, error)
}

// Repository implements QueryRepository on top of the generated category query
// client, guarding every call.
type Repository struct {
	client pbcategory.CategoryQueryServiceClient
	guard  *resilience.DependencyGuard
}

// New builds the category adapter from its query client.
func New(client pbcategory.CategoryQueryServiceClient, opts ...adapter.GuardOption) *Repository {
	r := &Repository{client: client}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// NewQueryAdapter wraps a generated category query client into a QueryRepository.
func NewQueryAdapter(client pbcategory.CategoryQueryServiceClient, opts ...adapter.GuardOption) QueryRepository {
	return New(client, opts...)
}

// SetGuard attaches the dependency guard used to wrap every outbound call.
func (r *Repository) SetGuard(g *resilience.DependencyGuard) {
	r.guard = g
}

func (r *Repository) FindByID(ctx context.Context, categoryID int) (*models.Category, error) {
	res, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbcategory.ApiResponseCategory, error) {
		return r.client.FindById(callCtx, &pbcategory.FindByIdCategoryRequest{Id: int32(categoryID)})
	})
	if err != nil {
		return nil, category_errors.ErrCategoryNotFound.WithInternal(err)
	}
	return toModel(res.Data), nil
}

func (r *Repository) FindAll(ctx context.Context, page, pageSize int, search string) ([]*models.Category, error) {
	res, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbcategory.ApiResponsePaginationCategory, error) {
		return r.client.FindAll(callCtx, &pbcategory.FindAllCategoryRequest{
			Page:     int32(page),
			PageSize: int32(pageSize),
			Search:   search,
		})
	})
	if err != nil {
		return nil, category_errors.ErrFindAllCategory.WithInternal(err)
	}
	return toModels(res.Data), nil
}

func toModel(c *pbcategory.CategoryResponse) *models.Category {
	if c == nil {
		return nil
	}
	return &models.Category{
		CategoryID:    c.Id,
		Name:          c.Name,
		Description:   &c.Description,
		SlugCategory:  &c.SlugCategory,
		ImageCategory: &c.ImageCategory,
		CreatedAt:     parseTimeString(c.CreatedAt),
		UpdatedAt:     parseTimeString(c.UpdatedAt),
	}
}

func toModels(cs []*pbcategory.CategoryResponse) []*models.Category {
	out := make([]*models.Category, 0, len(cs))
	for _, c := range cs {
		out = append(out, toModel(c))
	}
	return out
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
