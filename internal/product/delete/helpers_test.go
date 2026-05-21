package delete_test

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/your-org/your-service/internal/domain/product"
	"github.com/your-org/your-service/internal/domain/product/producttest"
)

func seedProduct(t *testing.T, repo *producttest.MemoryRepository) *product.Product {
	t.Helper()
	name, _ := product.NewName("Widget")
	desc, _ := product.NewDescription("")
	price, _ := product.NewPrice(9.99)
	sku, _ := product.NewSKU("WGT-001")
	p := product.New(name, desc, price, sku)
	if err := repo.Save(context.Background(), p); err != nil {
		t.Fatalf("seed: %v", err)
	}
	return p
}

func serve(t *testing.T, r *gin.Engine, method, path string) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, strings.NewReader(""))
	r.ServeHTTP(w, req)
	return w
}
