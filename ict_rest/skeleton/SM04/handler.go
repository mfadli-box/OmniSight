package SM04

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

func (h *Handler) ListSignatureType(c *gin.Context) {
	list, err := h.usecase.ListSignatureType()
	if err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Signature types retrieved",
		"data":    list,
	})
}

func (h *Handler) CreateSignatureType(c *gin.Context) {
	var req TypeCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid request body"))
		return
	}
	if err := h.usecase.CreateSignatureType(req); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Signature type created",
	})
}

func (h *Handler) UpdateSignatureType(c *gin.Context) {
	id := c.Param("id")
	var req TypeCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid request body"))
		return
	}
	if err := h.usecase.UpdateSignatureType(id, req); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Signature type updated",
	})
}

func (h *Handler) ListUser(c *gin.Context) {
	companyID := c.GetString("companyId")
	isAdmin, _ := c.Get("isAdmin")
	if isAdmin.(bool) {
		companyID = ""
	}
	search := c.Query("search")
	list, err := h.usecase.ListUser(companyID, search)
	if err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Users retrieved",
		"data":    list,
	})
}
