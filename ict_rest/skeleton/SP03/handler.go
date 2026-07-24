package SP03

import (
	"ict_rest/mechanic"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	usecase UseCase
}

func NHand(u UseCase) *Handler {
	return &Handler{usecase: u}
}

func (h *Handler) ListActions(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		mechanic.Error(c, mechanic.Unauthorized("User not found in session"))
		return
	}
	var meta mechanic.ActionMeta
	if err := c.ShouldBindQuery(&meta); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid query parameters"))
		return
	}
	list, gridMeta, err := h.usecase.ListActions(
		userID.(string),
		meta.Search,
		meta.Page,
		meta.Size,
		meta.SortBy,
		meta.SortOrder,
	)
	if err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Actions retrieved",
		"data":    list,
		"meta":    gridMeta,
	})
}
