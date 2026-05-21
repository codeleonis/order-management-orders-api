package findbyid

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/your-org/your-service/internal/domain/product"
)

// Handler handles HTTP requests for retrieving a product by ID.
type Handler struct {
	useCase *UseCase
}

// NewHandler constructs a Handler with the given use case.
func NewHandler(useCase *UseCase) *Handler {
	return &Handler{useCase: useCase}
}

// Handle parses the request, delegates to the use case, and writes the HTTP response.
func (h *Handler) Handle(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	result, err := h.useCase.Execute(c.Request.Context(), Command{ID: id})
	if err != nil {
		c.JSON(errToStatus(err), gin.H{"error": err.Error()})
		return
	}

	p := result.Product
	c.JSON(http.StatusOK, Response{
		ID:          p.ID().String(),
		Name:        p.Name().String(),
		Description: p.Description().String(),
		Price:       p.Price().Value(),
		SKU:         p.SKU().String(),
		CreatedAt:   p.CreatedAt(),
		UpdatedAt:   p.UpdatedAt(),
	})
}

func errToStatus(err error) int {
	switch {
	case errors.Is(err, product.ErrNotFound):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
