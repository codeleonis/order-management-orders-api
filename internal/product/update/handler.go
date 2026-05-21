package update

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/your-org/your-service/internal/domain/product"
)

const errKey = "error"

// Handler handles HTTP requests for updating a product.
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
		c.JSON(http.StatusBadRequest, gin.H{errKey: "invalid product id"})
		return
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{errKey: "invalid request body"})
		return
	}

	result, err := h.useCase.Execute(c.Request.Context(), Command{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		SKU:         req.SKU,
	})
	if err != nil {
		c.JSON(errToStatus(err), gin.H{errKey: err.Error()})
		return
	}

	p := result.Product
	c.JSON(http.StatusOK, Response{
		ID:          p.ID().String(),
		Name:        p.Name().String(),
		Description: p.Description().String(),
		Price:       p.Price().Value(),
		SKU:         p.SKU().String(),
		UpdatedAt:   p.UpdatedAt(),
	})
}

func errToStatus(err error) int {
	switch {
	case errors.Is(err, product.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, product.ErrDuplicateSKU):
		return http.StatusConflict
	case errors.Is(err, ErrNoPatchFields),
		errors.Is(err, product.ErrInvalidName),
		errors.Is(err, product.ErrInvalidDesc),
		errors.Is(err, product.ErrInvalidPrice),
		errors.Is(err, product.ErrInvalidSKU):
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}
