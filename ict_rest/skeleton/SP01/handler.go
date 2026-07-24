package SP01

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

func (h *Handler) ListUserCompany(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		mechanic.Error(c, mechanic.Unauthorized("User not found in session"))
		return
	}
	list, err := h.usecase.ListUserCompany(userID.(string))
	if err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User companies retrieved",
		"data":    list,
	})
}

func (h *Handler) ListUserModule(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		mechanic.Error(c, mechanic.Unauthorized("User not found in session"))
		return
	}
	companyID := c.Query("company_id")
	if companyID == "" {
		c.JSON(http.StatusOK, gin.H{
			"message": "User modules retrieved",
			"data":    []ModuleTreeNode{},
		})
		return
	}
	isAdmin, _ := c.Get("isAdmin")
	admin, _ := isAdmin.(bool)

	var list []ModuleTreeNode
	var err error
	if admin {
		list, err = h.usecase.ListAllModuleTree()
	} else {
		list, err = h.usecase.ListUserModule(userID.(string), companyID)
	}
	if err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User modules retrieved",
		"data":    list,
	})
}
