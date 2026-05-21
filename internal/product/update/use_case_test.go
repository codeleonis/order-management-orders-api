package update_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/your-org/your-service/internal/domain/product"
	"github.com/your-org/your-service/internal/domain/product/producttest"
	updateproduct "github.com/your-org/your-service/internal/product/update"
)

func TestUseCase_Execute(t *testing.T) {
	t.Run("patches name only", func(t *testing.T) {
		repo := producttest.NewMemoryRepository()
		p := seedProductWithSKU(t, repo, "WGT-001")

		uc := updateproduct.New(repo)
		result, err := uc.Execute(context.Background(), updateproduct.Command{
			ID:   p.ID(),
			Name: ptr("Updated Name"),
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Product.Name().String() != "Updated Name" {
			t.Errorf("expected name Updated Name, got %s", result.Product.Name().String())
		}
		if result.Product.Price().Value() != p.Price().Value() {
			t.Error("price should not have changed")
		}
	})

	t.Run("patches multiple fields", func(t *testing.T) {
		repo := producttest.NewMemoryRepository()
		p := seedProductWithSKU(t, repo, "WGT-001")

		uc := updateproduct.New(repo)
		result, err := uc.Execute(context.Background(), updateproduct.Command{
			ID:    p.ID(),
			Name:  ptr("New Name"),
			Price: ptr(49.99),
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Product.Price().Value() != 49.99 {
			t.Errorf("expected price 49.99, got %v", result.Product.Price().Value())
		}
	})

	t.Run("returns ErrNoPatchFields when no fields provided", func(t *testing.T) {
		repo := producttest.NewMemoryRepository()
		p := seedProductWithSKU(t, repo, "WGT-001")

		uc := updateproduct.New(repo)
		_, err := uc.Execute(context.Background(), updateproduct.Command{ID: p.ID()})

		if !errors.Is(err, updateproduct.ErrNoPatchFields) {
			t.Errorf("expected ErrNoPatchFields, got %v", err)
		}
	})

	t.Run("returns ErrNotFound for unknown id", func(t *testing.T) {
		repo := producttest.NewMemoryRepository()
		uc := updateproduct.New(repo)

		_, err := uc.Execute(context.Background(), updateproduct.Command{
			ID:   uuid.New(),
			Name: ptr("X"),
		})
		if !errors.Is(err, product.ErrNotFound) {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("returns ErrInvalidPrice for zero price", func(t *testing.T) {
		repo := producttest.NewMemoryRepository()
		p := seedProductWithSKU(t, repo, "WGT-001")

		uc := updateproduct.New(repo)
		_, err := uc.Execute(context.Background(), updateproduct.Command{
			ID:    p.ID(),
			Price: ptr(0.0),
		})
		if !errors.Is(err, product.ErrInvalidPrice) {
			t.Errorf("expected ErrInvalidPrice, got %v", err)
		}
	})

	t.Run("returns ErrDuplicateSKU when SKU belongs to another product", func(t *testing.T) {
		repo := producttest.NewMemoryRepository()
		seedProductWithSKU(t, repo, "WGT-001")
		p2 := seedProductWithSKU(t, repo, "WGT-002")

		uc := updateproduct.New(repo)
		_, err := uc.Execute(context.Background(), updateproduct.Command{
			ID:  p2.ID(),
			SKU: ptr("WGT-001"),
		})
		if !errors.Is(err, product.ErrDuplicateSKU) {
			t.Errorf("expected ErrDuplicateSKU, got %v", err)
		}
	})
}
