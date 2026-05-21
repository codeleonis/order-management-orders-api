package update

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/your-org/your-service/internal/domain/product"
)

// ErrNoPatchFields is returned when a PATCH request provides no fields to update.
var ErrNoPatchFields = errors.New("at least one field must be provided")

// Command carries the partial update fields; nil means "do not change".
type Command struct {
	ID          uuid.UUID
	Name        *string
	Description *string
	Price       *float64
	SKU         *string
}

// Result holds the updated product on success.
type Result struct {
	Product *product.Product
}

// UseCase orchestrates partial product updates.
type UseCase struct {
	repo product.Repository
}

// New constructs a UseCase with the given repository.
func New(repo product.Repository) *UseCase {
	return &UseCase{repo: repo}
}

// Execute applies the patch fields to the product and persists the changes.
func (uc *UseCase) Execute(ctx context.Context, cmd Command) (Result, error) {
	if cmd.Name == nil && cmd.Description == nil && cmd.Price == nil && cmd.SKU == nil {
		return Result{}, ErrNoPatchFields
	}
	p, err := uc.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return Result{}, err
	}
	if err := uc.applyPatch(ctx, p, cmd); err != nil {
		return Result{}, err
	}
	if err := uc.repo.Update(ctx, p); err != nil {
		return Result{}, fmt.Errorf("update product: %w", err)
	}
	return Result{Product: p}, nil
}

func (uc *UseCase) applyPatch(ctx context.Context, p *product.Product, cmd Command) error {
	if cmd.Name != nil {
		name, err := product.NewName(*cmd.Name)
		if err != nil {
			return err
		}
		p.Rename(name)
	}
	if cmd.Description != nil {
		desc, err := product.NewDescription(*cmd.Description)
		if err != nil {
			return err
		}
		p.Redescribe(desc)
	}
	if cmd.Price != nil {
		price, err := product.NewPrice(*cmd.Price)
		if err != nil {
			return err
		}
		p.Reprice(price)
	}
	if cmd.SKU != nil {
		sku, err := product.NewSKU(*cmd.SKU)
		if err != nil {
			return err
		}
		exists, err := uc.repo.ExistsBySKU(ctx, sku)
		if err != nil {
			return fmt.Errorf("check sku uniqueness: %w", err)
		}
		if exists && sku.String() != p.SKU().String() {
			return product.ErrDuplicateSKU
		}
		p.UpdateSKU(sku)
	}
	return nil
}
