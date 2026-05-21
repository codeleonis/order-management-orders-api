package findbyid_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/your-org/your-service/internal/domain/product"
	"github.com/your-org/your-service/internal/domain/product/producttest"
	findbyid "github.com/your-org/your-service/internal/product/find_by_id"
)

func TestUseCase_Execute(t *testing.T) {
	t.Run("returns product by id", func(t *testing.T) {
		repo := producttest.NewMemoryRepository()
		p := seedProduct(t, repo)

		uc := findbyid.New(repo)
		result, err := uc.Execute(context.Background(), findbyid.Command{ID: p.ID()})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Product.ID() != p.ID() {
			t.Errorf("expected id %s, got %s", p.ID(), result.Product.ID())
		}
	})

	t.Run("returns ErrNotFound for unknown id", func(t *testing.T) {
		repo := producttest.NewMemoryRepository()
		uc := findbyid.New(repo)

		_, err := uc.Execute(context.Background(), findbyid.Command{ID: uuid.New()})
		if !errors.Is(err, product.ErrNotFound) {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}
