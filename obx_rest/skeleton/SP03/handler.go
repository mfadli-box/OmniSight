package SP03

import (
	"net/http"
	"obx_rest/mechanic"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	usecase UseCase
}

func NHand(u UseCase) *Handler {
	return &Handler{usecase: u}
}

func requireUserID(c *gin.Context) (string, bool) {
	value, exists := c.Get("userId")
	if !exists {
		mechanic.Error(c, mechanic.Unauthorized("User not found in session"))
		return "", false
	}
	userID, ok := value.(string)
	if !ok || userID == "" {
		mechanic.Error(c, mechanic.Unauthorized("User not found in session"))
		return "", false
	}
	return userID, true
}

func (h *Handler) ListActions(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	var meta mechanic.ActionMeta
	if err := c.ShouldBindQuery(&meta); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid query parameters"))
		return
	}
	list, gridMeta, err := h.usecase.ListActions(
		userID,
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
