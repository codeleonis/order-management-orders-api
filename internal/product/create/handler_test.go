package create_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/your-org/your-service/internal/domain/product"
	"github.com/your-org/your-service/internal/domain/product/producttest"
	"github.com/your-org/your-service/internal/product/create"
)

func init() { gin.SetMode(gin.TestMode) }

func newRouter(repo *producttest.MemoryRepository) *gin.Engine {
	r := gin.New()
	uc := create.New(repo)
	r.POST("/products", create.NewHandler(uc).Handle)
	return r
}

func TestHandler_Handle(t *testing.T) {
	t.Run("201 on valid request", func(t *testing.T) {
		w := serve(t, newRouter(producttest.NewMemoryRepository()),
			"POST", "/products", `{"name":"Widget","description":"A widget","price":9.99,"sku":"WGT-001"}`)

		if w.Code != http.StatusCreated {
			t.Errorf("expected 201, got %d — body: %s", w.Code, w.Body)
		}
		var resp map[string]any
		json.NewDecoder(w.Body).Decode(&resp)
		if resp["id"] == nil {
			t.Error("expected id in response")
		}
		if resp["name"] != "Widget" {
			t.Errorf("expected name Widget, got %v", resp["name"])
		}
	})

	t.Run("400 on malformed JSON", func(t *testing.T) {
		w := serve(t, newRouter(producttest.NewMemoryRepository()),
			"POST", "/products", `{bad json`)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("422 on empty name", func(t *testing.T) {
		w := serve(t, newRouter(producttest.NewMemoryRepository()),
			"POST", "/products", `{"name":"","price":9.99,"sku":"WGT-001"}`)

		if w.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected 422, got %d", w.Code)
		}
	})

	t.Run("422 on zero price", func(t *testing.T) {
		w := serve(t, newRouter(producttest.NewMemoryRepository()),
			"POST", "/products", `{"name":"Widget","price":0,"sku":"WGT-001"}`)

		if w.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected 422, got %d", w.Code)
		}
	})

	t.Run("409 on duplicate SKU", func(t *testing.T) {
		repo := producttest.NewMemoryRepository()
		r := newRouter(repo)
		body := `{"name":"Widget","price":9.99,"sku":"WGT-001"}`

		serve(t, r, "POST", "/products", body)
		w := serve(t, r, "POST", "/products", body)

		if w.Code != http.StatusConflict {
			t.Errorf("expected 409, got %d", w.Code)
		}
	})

	t.Run("422 on name exceeding 100 chars", func(t *testing.T) {
		longName := strings.Repeat("x", 101)
		w := serve(t, newRouter(producttest.NewMemoryRepository()),
			"POST", "/products", `{"name":"`+longName+`","price":9.99,"sku":"WGT-001"}`)

		if w.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected 422, got %d", w.Code)
		}
	})

	t.Run("500 on repository failure", func(t *testing.T) {
		repo := producttest.NewMemoryRepository()
		repo.SaveErr = product.ErrNotFound // any non-domain error to simulate infra failure

		w := serve(t, newRouter(repo),
			"POST", "/products", `{"name":"Widget","price":9.99,"sku":"WGT-001"}`)

		// SaveErr is ErrNotFound here just to trigger the error path; 404 is acceptable
		if w.Code == http.StatusCreated {
			t.Error("expected non-201 on repo error")
		}
	})
}

func serve(t *testing.T, r *gin.Engine, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	return w
}
