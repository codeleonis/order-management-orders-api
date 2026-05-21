package delete_test

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/your-org/your-service/internal/domain/product/producttest"
	deleteproduct "github.com/your-org/your-service/internal/product/delete"
)

func init() { gin.SetMode(gin.TestMode) }

func newRouter(repo *producttest.MemoryRepository) *gin.Engine {
	r := gin.New()
	uc := deleteproduct.New(repo)
	r.DELETE("/products/:id", deleteproduct.NewHandler(uc).Handle)
	return r
}

func TestHandler_Handle(t *testing.T) {
	t.Run("204 on successful delete", func(t *testing.T) {
		repo := producttest.NewMemoryRepository()
		p := seedProduct(t, repo)

		w := serve(t, newRouter(repo), "DELETE", "/products/"+p.ID().String())

		if w.Code != http.StatusNoContent {
			t.Errorf("expected 204, got %d — body: %s", w.Code, w.Body)
		}
	})

	t.Run("400 on invalid UUID", func(t *testing.T) {
		w := serve(t, newRouter(producttest.NewMemoryRepository()), "DELETE", "/products/not-a-uuid")

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("404 on non-existent product", func(t *testing.T) {
		w := serve(t, newRouter(producttest.NewMemoryRepository()),
			"DELETE", "/products/00000000-0000-0000-0000-000000000001")

		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})
}
