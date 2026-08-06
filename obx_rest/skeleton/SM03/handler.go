package SM03

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

func (h *Handler) ListCompany(c *gin.Context) {
	var meta mechanic.ActionMeta
	if err := c.ShouldBindQuery(&meta); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid query parameters"))
		return
	}
	list, gridMeta, err := h.usecase.ListCompany(meta)
	if err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Company retrieved",
		"data":    list,
		"meta":    gridMeta,
	})
}

func (h *Handler) ListModule(c *gin.Context) {
	list, err := h.usecase.ListModule()
	if err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Modules retrieved",
		"data":    list,
	})
}

func (h *Handler) CreateCompany(c *gin.Context) {
	var req CompanyCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid request body"))
		return
	}
	if err := h.usecase.CreateCompany(req); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Company created successfully",
	})
}

func (h *Handler) UpdateCompany(c *gin.Context) {
	id := c.Param("id")
	var req CompanyUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid request body"))
		return
	}
	if err := h.usecase.UpdateCompany(id, req); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Company updated successfully",
	})
}

func (h *Handler) ListCompanyModule(c *gin.Context) {
	id := c.Param("id")
	var meta mechanic.ActionMeta
	if err := c.ShouldBindQuery(&meta); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid query parameters"))
		return
	}
	list, gridMeta, err := h.usecase.ListCompanyModule(id, meta)
	if err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Company modules retrieved",
		"data":    list,
		"meta":    gridMeta,
	})
}

func (h *Handler) ListCompanyModuleSelect(c *gin.Context) {
	id := c.Param("id")
	list, err := h.usecase.ListCompanyModuleSelect(id)
	if err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Company module options retrieved",
		"data":    list,
	})
}

func (h *Handler) AssignCompanyModule(c *gin.Context) {
	id := c.Param("id")
	var body CompanyModuleAssignRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		mechanic.Error(c, mechanic.ValidationError("module_id is required"))
		return
	}
	if err := h.usecase.AssignCompanyModule(id, body.ModuleID); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Module assigned successfully",
	})
}

func (h *Handler) UpdateCompanyModule(c *gin.Context) {
	id := c.Param("id")
	moduleID := c.Param("moduleId")
	var body CompanyModuleUpdateRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid request body"))
		return
	}
	if err := h.usecase.UpdateCompanyModule(id, moduleID, body.IsActive); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Module updated successfully",
	})
}

func (h *Handler) ListArea(c *gin.Context) {
	id := c.Param("id")
	var meta mechanic.ActionMeta
	if err := c.ShouldBindQuery(&meta); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid query parameters"))
		return
	}
	list, gridMeta, err := h.usecase.ListArea(id, meta)
	if err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Areas retrieved",
		"data":    list,
		"meta":    gridMeta,
	})
}

func (h *Handler) CreateArea(c *gin.Context) {
	id := c.Param("id")
	var req AreaCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid request body"))
		return
	}
	if err := h.usecase.CreateArea(id, req); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Area created successfully",
	})
}

func (h *Handler) UpdateArea(c *gin.Context) {
	areaId := c.Param("areaId")
	var req AreaUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid request body"))
		return
	}
	if err := h.usecase.UpdateArea(areaId, req); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Area updated successfully",
	})
}

func (h *Handler) DeleteArea(c *gin.Context) {
	areaId := c.Param("areaId")
	if err := h.usecase.DeleteArea(areaId); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Area deleted successfully",
	})
}
