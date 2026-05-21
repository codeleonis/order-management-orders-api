package delete_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/your-org/your-service/internal/domain/product"
	"github.com/your-org/your-service/internal/domain/product/producttest"
	deleteproduct "github.com/your-org/your-service/internal/product/delete"
)

func TestUseCase_Execute(t *testing.T) {
	t.Run("deletes existing product", func(t *testing.T) {
		repo := producttest.NewMemoryRepository()
		p := seedProduct(t, repo)

		uc := deleteproduct.New(repo)
		if err := uc.Execute(context.Background(), deleteproduct.Command{ID: p.ID()}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if _, err := repo.FindByID(context.Background(), p.ID()); !errors.Is(err, product.ErrNotFound) {
			t.Error("expected product to be deleted")
		}
	})

	t.Run("returns ErrNotFound for non-existent product", func(t *testing.T) {
		repo := producttest.NewMemoryRepository()
		uc := deleteproduct.New(repo)

		err := uc.Execute(context.Background(), deleteproduct.Command{ID: uuid.New()})
		if !errors.Is(err, product.ErrNotFound) {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}
