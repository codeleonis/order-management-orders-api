package list_test

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/your-org/your-service/internal/domain/product"
	"github.com/your-org/your-service/internal/domain/product/producttest"
)

func seedN(t *testing.T, repo *producttest.MemoryRepository, n int) {
	t.Helper()
	for i := range n {
		name, _ := product.NewName("Product")
		desc, _ := product.NewDescription("")
		price, _ := product.NewPrice(float64(i+1) * 1.99)
		sku, _ := product.NewSKU("SKU-" + string(rune('A'+i)))
		p := product.New(name, desc, price, sku)
		if err := repo.Save(context.Background(), p); err != nil {
			t.Fatalf("seed: %v", err)
		}
	}
}

func serve(t *testing.T, r *gin.Engine, path string) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", path, strings.NewReader(""))
	r.ServeHTTP(w, req)
	return w
}
