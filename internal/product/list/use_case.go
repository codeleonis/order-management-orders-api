package list

import (
	"context"
	"math"

	"github.com/your-org/your-service/internal/domain/product"
)

// Command carries the pagination parameters for listing products.
type Command struct {
	Page     int
	PageSize int
}

// Result holds the paginated products and metadata.
type Result struct {
	Items      []*product.Product
	Total      int64
	TotalPages int
}

// UseCase retrieves a paginated list of products.
type UseCase struct {
	repo product.Repository
}

// New constructs a UseCase with the given repository.
func New(repo product.Repository) *UseCase {
	return &UseCase{repo: repo}
}

// Execute queries the repository and computes pagination metadata.
func (uc *UseCase) Execute(ctx context.Context, cmd Command) (Result, error) {
	//nolint:gosec // PageSize and Page are bounded by request normalizer (max 100)
	page, err := uc.repo.List(ctx, product.ListFilter{
		Limit:  int32(cmd.PageSize),
		Offset: int32((cmd.Page - 1) * cmd.PageSize),
	})
	if err != nil {
		return Result{}, err
	}

	totalPages := int(math.Ceil(float64(page.Total) / float64(cmd.PageSize)))
	return Result{
		Items:      page.Items,
		Total:      page.Total,
		TotalPages: totalPages,
	}, nil
}
