package delete

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/your-org/your-service/internal/domain/product"
)

// Command carries the ID of the product to delete.
type Command struct {
	ID uuid.UUID
}

// UseCase orchestrates product deletion.
type UseCase struct {
	repo product.Repository
}

// New constructs a UseCase with the given repository.
func New(repo product.Repository) *UseCase {
	return &UseCase{repo: repo}
}

// Execute verifies the product exists and then removes it.
func (uc *UseCase) Execute(ctx context.Context, cmd Command) error {
	if _, err := uc.repo.FindByID(ctx, cmd.ID); err != nil {
		return err
	}
	if err := uc.repo.Delete(ctx, cmd.ID); err != nil {
		return fmt.Errorf("delete product: %w", err)
	}
	return nil
}
