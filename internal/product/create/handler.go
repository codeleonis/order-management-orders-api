package create

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/your-org/your-service/internal/domain/product"
)

// Handler handles HTTP requests for creating a product.
type Handler struct {
	useCase *UseCase
}

// NewHandler constructs a Handler with the given use case.
func NewHandler(useCase *UseCase) *Handler {
	return &Handler{useCase: useCase}
}

// Handle parses the request, delegates to the use case, and writes the HTTP response.
func (h *Handler) Handle(c *gin.Context) {
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	result, err := h.useCase.Execute(c.Request.Context(), Command(req))
	if err != nil {
		c.JSON(errToStatus(err), gin.H{"error": err.Error()})
		return
	}

	p := result.Product
	c.JSON(http.StatusCreated, Response{
		ID:          p.ID().String(),
		Name:        p.Name().String(),
		Description: p.Description().String(),
		Price:       p.Price().Value(),
		SKU:         p.SKU().String(),
		CreatedAt:   p.CreatedAt(),
	})
}

func errToStatus(err error) int {
	switch {
	case errors.Is(err, product.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, product.ErrDuplicateSKU):
		return http.StatusConflict
	case errors.Is(err, product.ErrInvalidName),
		errors.Is(err, product.ErrInvalidDesc),
		errors.Is(err, product.ErrInvalidPrice),
		errors.Is(err, product.ErrInvalidSKU):
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}
