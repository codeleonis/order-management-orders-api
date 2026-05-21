package product

import (
	"context"

	"github.com/google/uuid"
)

// ListFilter specifies pagination constraints for repository queries.
type ListFilter struct {
	Limit  int32
	Offset int32
}

// Page holds a paginated slice of products together with the total count.
type Page struct {
	Items []*Product
	Total int64
}

// Repository is the persistence contract owned by the domain.
// Infrastructure provides the implementation; slices depend on this interface.
type Repository interface {
	Save(ctx context.Context, p *Product) error
	FindByID(ctx context.Context, id uuid.UUID) (*Product, error)
	List(ctx context.Context, f ListFilter) (Page, error)
	Update(ctx context.Context, p *Product) error
	Delete(ctx context.Context, id uuid.UUID) error
	ExistsBySKU(ctx context.Context, sku SKU) (bool, error)
}
