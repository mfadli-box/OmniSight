package SM05

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

func (h *Handler) ListSession(c *gin.Context) {
	var meta mechanic.ActionMeta
	if err := c.ShouldBindQuery(&meta); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid query parameters"))
		return
	}
	list, gridMeta, err := h.usecase.ListSession(meta)
	if err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Sessions retrieved",
		"data":    list,
		"meta":    gridMeta,
	})
}

func (h *Handler) GetSessionDetail(c *gin.Context) {
	id := c.Param("id")
	detail, err := h.usecase.GetSessionDetail(id)
	if err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Session detail retrieved",
		"data":    detail,
	})
}

func (h *Handler) RevokeSession(c *gin.Context) {
	id := c.Param("id")
	if err := h.usecase.RevokeSession(id); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Session revoked successfully",
	})
}
