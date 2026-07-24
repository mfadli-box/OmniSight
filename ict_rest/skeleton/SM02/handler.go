package SM02

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

func (h *Handler) ListModule(c *gin.Context) {
	var meta mechanic.ActionMeta
	if err := c.ShouldBindQuery(&meta); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid query parameters"))
		return
	}
	list, gridMeta, err := h.usecase.ListModule(meta)
	if err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Modules retrieved",
		"data":    list,
		"meta":    gridMeta,
	})
}

func (h *Handler) CreateModule(c *gin.Context) {
	var req ModuleCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid request body"))
		return
	}
	if err := h.usecase.CreateModule(req); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Module created successfully",
	})
}

func (h *Handler) UpdateModule(c *gin.Context) {
	id := c.Param("id")
	var req ModuleUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid request body"))
		return
	}
	if err := h.usecase.UpdateModule(id, req); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Module updated successfully",
	})
}
