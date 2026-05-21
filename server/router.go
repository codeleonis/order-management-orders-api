package server

import (
	"github.com/gin-gonic/gin"
	createproduct "github.com/your-org/your-service/internal/product/create"
	deleteproduct "github.com/your-org/your-service/internal/product/delete"
	findbyid "github.com/your-org/your-service/internal/product/find_by_id"
	listproducts "github.com/your-org/your-service/internal/product/list"
	updateproduct "github.com/your-org/your-service/internal/product/update"
)

// Handlers groups the HTTP handler for each use case.
type Handlers struct {
	Create   *createproduct.Handler
	Delete   *deleteproduct.Handler
	FindByID *findbyid.Handler
	List     *listproducts.Handler
	Update   *updateproduct.Handler
}

// NewRouter builds the gin engine and registers all routes.
func NewRouter(h Handlers) *gin.Engine {
	r := gin.New()
	r.Use(RequestID())
	r.Use(RequestLogger())
	r.Use(gin.Recovery())

	r.GET("/health", healthCheck)

	v1 := r.Group("/api/v1")
	{
		products := v1.Group("/products")
		products.POST("", h.Create.Handle)
		products.GET("", h.List.Handle)
		products.GET("/:id", h.FindByID.Handle)
		products.PATCH("/:id", h.Update.Handle)
		products.DELETE("/:id", h.Delete.Handle)
	}

	return r
}
