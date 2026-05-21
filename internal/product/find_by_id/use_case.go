package findbyid

import (
	"context"

	"github.com/google/uuid"
	"github.com/your-org/your-service/internal/domain/product"
)

// Command carries the ID of the product to retrieve.
type Command struct {
	ID uuid.UUID
}

// Result holds the found product on success.
type Result struct {
	Product *product.Product
}

// UseCase retrieves a single product by its identifier.
type UseCase struct {
	repo product.Repository
}

// New constructs a UseCase with the given repository.
func New(repo product.Repository) *UseCase {
	return &UseCase{repo: repo}
}

// Execute fetches the product from the repository.
func (uc *UseCase) Execute(ctx context.Context, cmd Command) (Result, error) {
	p, err := uc.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return Result{}, err
	}
	return Result{Product: p}, nil
}
