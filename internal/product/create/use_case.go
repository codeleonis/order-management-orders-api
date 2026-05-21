package create

import (
	"context"
	"fmt"

	"github.com/your-org/your-service/internal/domain/product"
)

// Command carries the input data for creating a product.
type Command struct {
	Name        string
	Description string
	Price       float64
	SKU         string
}

// Result holds the newly created product on success.
type Result struct {
	Product *product.Product
}

// UseCase orchestrates product creation.
type UseCase struct {
	repo product.Repository
}

// New constructs a UseCase with the given repository.
func New(repo product.Repository) *UseCase {
	return &UseCase{repo: repo}
}

// Execute validates the command, enforces SKU uniqueness, and persists the product.
func (uc *UseCase) Execute(ctx context.Context, cmd Command) (Result, error) {
	name, err := product.NewName(cmd.Name)
	if err != nil {
		return Result{}, err
	}
	desc, err := product.NewDescription(cmd.Description)
	if err != nil {
		return Result{}, err
	}
	price, err := product.NewPrice(cmd.Price)
	if err != nil {
		return Result{}, err
	}
	sku, err := product.NewSKU(cmd.SKU)
	if err != nil {
		return Result{}, err
	}

	exists, err := uc.repo.ExistsBySKU(ctx, sku)
	if err != nil {
		return Result{}, fmt.Errorf("check sku uniqueness: %w", err)
	}
	if exists {
		return Result{}, product.ErrDuplicateSKU
	}

	p := product.New(name, desc, price, sku)
	if err := uc.repo.Save(ctx, p); err != nil {
		return Result{}, fmt.Errorf("save product: %w", err)
	}
	return Result{Product: p}, nil
}
