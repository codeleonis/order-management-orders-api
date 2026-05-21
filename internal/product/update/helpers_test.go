package update_test

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/your-org/your-service/internal/domain/product"
	"github.com/your-org/your-service/internal/domain/product/producttest"
)

func ptr[T any](v T) *T { return &v }

func seedProductWithSKU(t *testing.T, repo *producttest.MemoryRepository, skuVal string) *product.Product {
	t.Helper()
	name, _ := product.NewName("Widget")
	desc, _ := product.NewDescription("")
	price, _ := product.NewPrice(9.99)
	sku, _ := product.NewSKU(skuVal)
	p := product.New(name, desc, price, sku)
	if err := repo.Save(context.Background(), p); err != nil {
		t.Fatalf("seed product: %v", err)
	}
	return p
}

func serve(t *testing.T, r *gin.Engine, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	return w
}
