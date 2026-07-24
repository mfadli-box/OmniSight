package SM06

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

func (h *Handler) ListLocationType(c *gin.Context) {
	companyID := c.GetString("companyId")
	var meta mechanic.ActionMeta
	if err := c.ShouldBindQuery(&meta); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid query parameters"))
		return
	}
	list, gridMeta, err := h.usecase.ListLocationType(companyID, meta)
	if err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Location types retrieved",
		"data":    list,
		"meta":    gridMeta,
	})
}

func (h *Handler) CreateLocationType(c *gin.Context) {
	companyID := c.GetString("companyId")
	var req LocationTypeCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid request body"))
		return
	}
	if err := h.usecase.CreateLocationType(companyID, req); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Location type created successfully",
	})
}

func (h *Handler) UpdateLocationType(c *gin.Context) {
	id := c.Param("id")
	var req LocationTypeUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid request body"))
		return
	}
	if err := h.usecase.UpdateLocationType(id, req); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Location type updated successfully",
	})
}

func (h *Handler) DeleteLocationType(c *gin.Context) {
	id := c.Param("id")
	if err := h.usecase.DeleteLocationType(id); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Location type deleted successfully",
	})
}

func (h *Handler) ListLocationTypeSelect(c *gin.Context) {
	companyID := c.GetString("companyId")
	list, err := h.usecase.ListLocationTypeSelect(companyID)
	if err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Location types retrieved",
		"data":    list,
	})
}

func (h *Handler) ListLocation(c *gin.Context) {
	companyID := c.GetString("companyId")
	var meta mechanic.ActionMeta
	if err := c.ShouldBindQuery(&meta); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid query parameters"))
		return
	}
	list, gridMeta, err := h.usecase.ListLocation(companyID, meta)
	if err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Locations retrieved",
		"data":    list,
		"meta":    gridMeta,
	})
}

func (h *Handler) CreateLocation(c *gin.Context) {
	companyID := c.GetString("companyId")
	var req LocationCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid request body"))
		return
	}
	if err := h.usecase.CreateLocation(companyID, req); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Location created successfully",
	})
}

func (h *Handler) UpdateLocation(c *gin.Context) {
	id := c.Param("id")
	var req LocationUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid request body"))
		return
	}
	if err := h.usecase.UpdateLocation(id, req); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Location updated successfully",
	})
}

func (h *Handler) DeleteLocation(c *gin.Context) {
	id := c.Param("id")
	if err := h.usecase.DeleteLocation(id); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Location deleted successfully",
	})
}

func (h *Handler) ListLocationSelect(c *gin.Context) {
	companyID := c.GetString("companyId")
	list, err := h.usecase.ListLocationSelect(companyID)
	if err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Locations retrieved",
		"data":    list,
	})
}
