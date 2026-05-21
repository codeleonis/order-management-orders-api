package list_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/your-org/your-service/internal/domain/product/producttest"
	listproducts "github.com/your-org/your-service/internal/product/list"
)

func init() { gin.SetMode(gin.TestMode) }

func newRouter(repo *producttest.MemoryRepository) *gin.Engine {
	r := gin.New()
	uc := listproducts.New(repo)
	r.GET("/products", listproducts.NewHandler(uc).Handle)
	return r
}

func TestHandler_Handle(t *testing.T) {
	t.Run("200 with paginated results", func(t *testing.T) {
		repo := producttest.NewMemoryRepository()
		seedN(t, repo, 5)

		w := serve(t, newRouter(repo), "/products?page=1&page_size=3")

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d — body: %s", w.Code, w.Body)
		}
		var resp map[string]any
		json.NewDecoder(w.Body).Decode(&resp)
		if resp["total"] != float64(5) {
			t.Errorf("expected total 5, got %v", resp["total"])
		}
		items := resp["items"].([]any)
		if len(items) != 3 {
			t.Errorf("expected 3 items, got %d", len(items))
		}
		if resp["total_pages"] != float64(2) {
			t.Errorf("expected total_pages 2, got %v", resp["total_pages"])
		}
	})

	t.Run("200 with defaults when no query params", func(t *testing.T) {
		repo := producttest.NewMemoryRepository()
		seedN(t, repo, 2)

		w := serve(t, newRouter(repo), "/products")

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
		var resp map[string]any
		json.NewDecoder(w.Body).Decode(&resp)
		if resp["page"] != float64(1) {
			t.Errorf("expected page 1 (default), got %v", resp["page"])
		}
		if resp["page_size"] != float64(20) {
			t.Errorf("expected page_size 20 (default), got %v", resp["page_size"])
		}
	})

	t.Run("200 with empty list when no products", func(t *testing.T) {
		w := serve(t, newRouter(producttest.NewMemoryRepository()), "/products")

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
		var resp map[string]any
		json.NewDecoder(w.Body).Decode(&resp)
		if resp["total"] != float64(0) {
			t.Errorf("expected total 0, got %v", resp["total"])
		}
	})

	t.Run("enforces max page_size of 100", func(t *testing.T) {
		repo := producttest.NewMemoryRepository()
		seedN(t, repo, 2)

		w := serve(t, newRouter(repo), "/products?page_size=999")

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
		var resp map[string]any
		json.NewDecoder(w.Body).Decode(&resp)
		if resp["page_size"] != float64(100) {
			t.Errorf("expected page_size capped at 100, got %v", resp["page_size"])
		}
	})
}
