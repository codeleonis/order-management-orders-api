package findbyid_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/your-org/your-service/internal/domain/product/producttest"
	findbyid "github.com/your-org/your-service/internal/product/find_by_id"
)

func init() { gin.SetMode(gin.TestMode) }

func newRouter(repo *producttest.MemoryRepository) *gin.Engine {
	r := gin.New()
	uc := findbyid.New(repo)
	r.GET("/products/:id", findbyid.NewHandler(uc).Handle)
	return r
}

func TestHandler_Handle(t *testing.T) {
	t.Run("200 with product body on found", func(t *testing.T) {
		repo := producttest.NewMemoryRepository()
		p := seedProduct(t, repo)

		w := serve(t, newRouter(repo), "/products/"+p.ID().String())

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d — body: %s", w.Code, w.Body)
		}
		var resp map[string]any
		json.NewDecoder(w.Body).Decode(&resp)
		if resp["id"] != p.ID().String() {
			t.Errorf("expected id %s, got %v", p.ID(), resp["id"])
		}
		if resp["name"] != "Widget" {
			t.Errorf("expected name Widget, got %v", resp["name"])
		}
		if resp["sku"] != "WGT-001" {
			t.Errorf("expected sku WGT-001, got %v", resp["sku"])
		}
	})

	t.Run("400 on invalid UUID", func(t *testing.T) {
		w := serve(t, newRouter(producttest.NewMemoryRepository()), "/products/not-a-uuid")

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("404 on non-existent product", func(t *testing.T) {
		w := serve(t, newRouter(producttest.NewMemoryRepository()),
			"/products/00000000-0000-0000-0000-000000000001")

		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})
}
