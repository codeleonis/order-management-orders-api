package update_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/your-org/your-service/internal/domain/product/producttest"
	updateproduct "github.com/your-org/your-service/internal/product/update"
)

func init() { gin.SetMode(gin.TestMode) }

func newRouter(repo *producttest.MemoryRepository) *gin.Engine {
	r := gin.New()
	uc := updateproduct.New(repo)
	r.PATCH("/products/:id", updateproduct.NewHandler(uc).Handle)
	return r
}

func TestHandler_Handle(t *testing.T) {
	t.Run("200 on partial update (name only)", func(t *testing.T) {
		repo := producttest.NewMemoryRepository()
		p := seedProductWithSKU(t, repo, "WGT-001")

		w := serve(t, newRouter(repo), "PATCH", "/products/"+p.ID().String(),
			`{"name":"Updated Widget"}`)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d — body: %s", w.Code, w.Body)
		}
		var resp map[string]any
		json.NewDecoder(w.Body).Decode(&resp)
		if resp["name"] != "Updated Widget" {
			t.Errorf("expected name Updated Widget, got %v", resp["name"])
		}
	})

	t.Run("200 on multi-field patch", func(t *testing.T) {
		repo := producttest.NewMemoryRepository()
		p := seedProductWithSKU(t, repo, "WGT-002")

		w := serve(t, newRouter(repo), "PATCH", "/products/"+p.ID().String(),
			`{"name":"New Name","price":49.99}`)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
		var resp map[string]any
		json.NewDecoder(w.Body).Decode(&resp)
		if resp["price"] != 49.99 {
			t.Errorf("expected price 49.99, got %v", resp["price"])
		}
	})

	t.Run("400 on invalid UUID", func(t *testing.T) {
		w := serve(t, newRouter(producttest.NewMemoryRepository()),
			"PATCH", "/products/not-a-uuid", `{"name":"x"}`)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("400 on malformed JSON body", func(t *testing.T) {
		repo := producttest.NewMemoryRepository()
		p := seedProductWithSKU(t, repo, "WGT-003")

		w := serve(t, newRouter(repo), "PATCH", "/products/"+p.ID().String(), `{bad`)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("422 on no fields provided", func(t *testing.T) {
		repo := producttest.NewMemoryRepository()
		p := seedProductWithSKU(t, repo, "WGT-004")

		w := serve(t, newRouter(repo), "PATCH", "/products/"+p.ID().String(), `{}`)

		if w.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected 422, got %d", w.Code)
		}
	})

	t.Run("422 on invalid price", func(t *testing.T) {
		repo := producttest.NewMemoryRepository()
		p := seedProductWithSKU(t, repo, "WGT-005")

		w := serve(t, newRouter(repo), "PATCH", "/products/"+p.ID().String(), `{"price":0}`)

		if w.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected 422, got %d", w.Code)
		}
	})

	t.Run("404 on non-existent product", func(t *testing.T) {
		w := serve(t, newRouter(producttest.NewMemoryRepository()),
			"PATCH", "/products/00000000-0000-0000-0000-000000000001", `{"name":"x"}`)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})

	t.Run("409 on duplicate SKU", func(t *testing.T) {
		repo := producttest.NewMemoryRepository()
		seedProductWithSKU(t, repo, "WGT-001")
		p2 := seedProductWithSKU(t, repo, "WGT-006")

		w := serve(t, newRouter(repo), "PATCH", "/products/"+p2.ID().String(),
			`{"sku":"WGT-001"}`)

		if w.Code != http.StatusConflict {
			t.Errorf("expected 409, got %d", w.Code)
		}
	})
}
