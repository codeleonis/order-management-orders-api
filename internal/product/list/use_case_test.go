package list_test

import (
	"context"
	"testing"

	"github.com/your-org/your-service/internal/domain/product/producttest"
	listproducts "github.com/your-org/your-service/internal/product/list"
)

func TestUseCase_Execute(t *testing.T) {
	t.Run("returns paginated results", func(t *testing.T) {
		repo := producttest.NewMemoryRepository()
		seedN(t, repo, 5)

		uc := listproducts.New(repo)
		result, err := uc.Execute(context.Background(), listproducts.Command{Page: 1, PageSize: 3})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Total != 5 {
			t.Errorf("expected total 5, got %d", result.Total)
		}
		if len(result.Items) != 3 {
			t.Errorf("expected 3 items on page 1, got %d", len(result.Items))
		}
		if result.TotalPages != 2 {
			t.Errorf("expected 2 total pages, got %d", result.TotalPages)
		}
	})

	t.Run("returns empty page when no products", func(t *testing.T) {
		repo := producttest.NewMemoryRepository()
		uc := listproducts.New(repo)

		result, err := uc.Execute(context.Background(), listproducts.Command{Page: 1, PageSize: 20})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Total != 0 {
			t.Errorf("expected total 0, got %d", result.Total)
		}
	})
}
