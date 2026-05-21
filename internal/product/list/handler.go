package list

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for listing products.
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
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameters"})
		return
	}
	req.normalize()

	result, err := h.useCase.Execute(c.Request.Context(), Command(req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	items := make([]ProductItem, 0, len(result.Items))
	for _, p := range result.Items {
		items = append(items, ProductItem{
			ID:          p.ID().String(),
			Name:        p.Name().String(),
			Description: p.Description().String(),
			Price:       p.Price().Value(),
			SKU:         p.SKU().String(),
			CreatedAt:   p.CreatedAt(),
			UpdatedAt:   p.UpdatedAt(),
		})
	}

	c.JSON(http.StatusOK, Response{
		Items:      items,
		Total:      result.Total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: result.TotalPages,
	})
}
