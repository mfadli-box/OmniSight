package SP02

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
		mechanic.Error(c, mechanic.Unauthorized("Unauthorized"))
		return "", false
	}
	userID, ok := value.(string)
	if !ok || userID == "" {
		mechanic.Error(c, mechanic.Unauthorized("Unauthorized"))
		return "", false
	}
	return userID, true
}

func (h *Handler) ChangePassword(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	var req struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Current and new password are required"))
		return
	}
	if err := h.usecase.ChangePassword(
		userID,
		req.CurrentPassword,
		req.NewPassword,
	); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Password changed successfully",
	})
}
