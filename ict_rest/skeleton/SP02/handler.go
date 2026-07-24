package SP02

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

func (h *Handler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		mechanic.Error(c, mechanic.Unauthorized("Unauthorized"))
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
		userID.(string),
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
