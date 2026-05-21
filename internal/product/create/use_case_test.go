package create_test

import (
	"context"
	"errors"
	"testing"

	"github.com/your-org/your-service/internal/domain/product"
	"github.com/your-org/your-service/internal/domain/product/producttest"
	"github.com/your-org/your-service/internal/product/create"
)

func TestUseCase_Execute(t *testing.T) {
	t.Run("creates product successfully", func(t *testing.T) {
		repo := producttest.NewMemoryRepository()
		uc := create.New(repo)

		result, err := uc.Execute(context.Background(), create.Command{
			Name:        "Widget",
			Description: "A fine widget",
			Price:       9.99,
			SKU:         "WGT-001",
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Product == nil {
			t.Fatal("expected product, got nil")
		}
		if result.Product.Name().String() != "Widget" {
			t.Errorf("expected name Widget, got %s", result.Product.Name().String())
		}
	})

	t.Run("returns ErrInvalidName for empty name", func(t *testing.T) {
		repo := producttest.NewMemoryRepository()
		uc := create.New(repo)

		_, err := uc.Execute(context.Background(), create.Command{
			Name:  "",
			Price: 9.99,
			SKU:   "WGT-001",
		})

		if !errors.Is(err, product.ErrInvalidName) {
			t.Errorf("expected ErrInvalidName, got %v", err)
		}
	})

	t.Run("returns ErrInvalidPrice for zero price", func(t *testing.T) {
		repo := producttest.NewMemoryRepository()
		uc := create.New(repo)

		_, err := uc.Execute(context.Background(), create.Command{
			Name:  "Widget",
			Price: 0,
			SKU:   "WGT-001",
		})

		if !errors.Is(err, product.ErrInvalidPrice) {
			t.Errorf("expected ErrInvalidPrice, got %v", err)
		}
	})

	t.Run("returns ErrDuplicateSKU when SKU already exists", func(t *testing.T) {
		repo := producttest.NewMemoryRepository()
		uc := create.New(repo)

		cmd := create.Command{Name: "Widget", Price: 9.99, SKU: "WGT-001"}
		if _, err := uc.Execute(context.Background(), cmd); err != nil {
			t.Fatalf("first create failed: %v", err)
		}

		_, err := uc.Execute(context.Background(), cmd)
		if !errors.Is(err, product.ErrDuplicateSKU) {
			t.Errorf("expected ErrDuplicateSKU, got %v", err)
		}
	})

	t.Run("propagates repository error on save", func(t *testing.T) {
		repo := producttest.NewMemoryRepository()
		repo.SaveErr = errors.New("db down")
		uc := create.New(repo)

		_, err := uc.Execute(context.Background(), create.Command{
			Name:  "Widget",
			Price: 9.99,
			SKU:   "WGT-001",
		})

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
