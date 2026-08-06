package SP01

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

func requireContextString(c *gin.Context, key, label string) (string, bool) {
	value, exists := c.Get(key)
	if !exists {
		mechanic.Error(c, mechanic.Unauthorized(label+" not found in session"))
		return "", false
	}
	text, ok := value.(string)
	if !ok || text == "" {
		mechanic.Error(c, mechanic.Unauthorized(label+" not found in session"))
		return "", false
	}
	return text, true
}

func getContextBool(c *gin.Context, key string) bool {
	value, exists := c.Get(key)
	if !exists {
		return false
	}
	flag, ok := value.(bool)
	if !ok {
		return false
	}
	return flag
}

func (h *Handler) GetPrivilege(c *gin.Context) {
	userID, ok := requireContextString(c, "userId", "User")
	if !ok {
		return
	}
	companyID := c.GetString("companyId")
	moduleID := c.GetString("moduleId")
	isAdmin := getContextBool(c, "isAdmin")
	level, err := h.usecase.GetPrivilege(userID, companyID, moduleID, isAdmin)
	if err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{"level": level},
	})
}

func (h *Handler) ListUserCompany(c *gin.Context) {
	userID, ok := requireContextString(c, "userId", "User")
	if !ok {
		return
	}
	list, err := h.usecase.ListUserCompany(userID)
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
	userID, ok := requireContextString(c, "userId", "User")
	if !ok {
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
	admin := getContextBool(c, "isAdmin")

	var list []ModuleTreeNode
	var err error
	if admin {
		list, err = h.usecase.ListAllModuleTree()
	} else {
		list, err = h.usecase.ListUserModule(userID, companyID)
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
