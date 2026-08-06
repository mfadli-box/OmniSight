package SM01

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

func requireParam(c *gin.Context, key, label string) (string, bool) {
	value := c.Param(key)
	if value == "" {
		mechanic.Error(c, mechanic.ValidationError(label+" is required"))
		return "", false
	}
	return value, true
}

func (h *Handler) ListUser(c *gin.Context) {
	var meta mechanic.ActionMeta
	if err := c.ShouldBindQuery(&meta); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid query parameters"))
		return
	}
	companyID := c.GetString("companyId")
	list, gridMeta, err := h.usecase.ListUser(companyID, meta)
	if err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Users retrieved",
		"data":    list,
		"meta":    gridMeta,
	})
}

func (h *Handler) CreateUser(c *gin.Context) {
	var req UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid request body"))
		return
	}
	if err := h.usecase.CreateUser(req); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
	})
}

func (h *Handler) UpdateUser(c *gin.Context) {
	id, ok := requireParam(c, "id", "User ID")
	if !ok {
		return
	}
	var req UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid request body"))
		return
	}
	if err := h.usecase.UpdateUser(id, req); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
	})
}

func (h *Handler) ListHRISCompany(c *gin.Context) {
	list, err := h.usecase.ListHrisCompany()
	if err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "HRIS companies retrieved",
		"data":    list,
	})
}

func (h *Handler) ListAllCompany(c *gin.Context) {
	list, err := h.usecase.ListAllCompany()
	if err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Companies retrieved",
		"data":    list,
	})
}

func (h *Handler) ListAllModule(c *gin.Context) {
	list, err := h.usecase.ListAllModule()
	if err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Modules retrieved",
		"data":    list,
	})
}

func (h *Handler) AssignCompany(c *gin.Context) {
	id, ok := requireParam(c, "id", "User ID")
	if !ok {
		return
	}
	var body struct {
		CompanyID string `json:"company_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		mechanic.Error(c, mechanic.ValidationError("company_id is required"))
		return
	}
	if err := h.usecase.AssignCompany(id, body.CompanyID); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Company assigned successfully",
	})
}

func (h *Handler) ListUserCompany(c *gin.Context) {
	id, ok := requireParam(c, "id", "User ID")
	if !ok {
		return
	}
	var meta mechanic.ActionMeta
	if err := c.ShouldBindQuery(&meta); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid query parameters"))
		return
	}
	list, gridMeta, err := h.usecase.ListUserCompany(id, meta)
	if err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User companies retrieved",
		"data":    list,
		"meta":    gridMeta,
	})
}

func (h *Handler) ListUserCompanySelect(c *gin.Context) {
	id, ok := requireParam(c, "id", "User ID")
	if !ok {
		return
	}
	list, err := h.usecase.ListUserCompanySelect(id)
	if err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User company options retrieved",
		"data":    list,
	})
}

func (h *Handler) CreateUserCompany(c *gin.Context) {
	id, ok := requireParam(c, "id", "User ID")
	if !ok {
		return
	}
	var req UserCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mechanic.Error(c, mechanic.ValidationError("company_id is required"))
		return
	}
	if err := h.usecase.CreateUserCompany(id, req); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Company assigned successfully",
	})
}

func (h *Handler) UpdateUserCompany(c *gin.Context) {
	id, ok := requireParam(c, "id", "User ID")
	if !ok {
		return
	}
	companyID, ok := requireParam(c, "companyId", "Company ID")
	if !ok {
		return
	}
	var req UserCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid request body"))
		return
	}
	if err := h.usecase.UpdateUserCompany(id, companyID, req); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Company updated successfully",
	})
}

func (h *Handler) DeleteUserCompany(c *gin.Context) {
	id, ok := requireParam(c, "id", "User ID")
	if !ok {
		return
	}
	companyID, ok := requireParam(c, "companyId", "Company ID")
	if !ok {
		return
	}
	if err := h.usecase.DeleteUserCompany(id, companyID); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Company removed successfully",
	})
}

func (h *Handler) ListUserPrivilege(c *gin.Context) {
	id, ok := requireParam(c, "id", "User ID")
	if !ok {
		return
	}
	var meta mechanic.ActionMeta
	if err := c.ShouldBindQuery(&meta); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid query parameters"))
		return
	}
	list, gridMeta, err := h.usecase.ListUserPrivilege(id, meta)
	if err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User privileges retrieved",
		"data":    list,
		"meta":    gridMeta,
	})
}

func (h *Handler) CreateUserPrivilege(c *gin.Context) {
	id, ok := requireParam(c, "id", "User ID")
	if !ok {
		return
	}
	var req UserPrivilegeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mechanic.Error(c, mechanic.ValidationError("user_company_id and module_id are required"))
		return
	}
	if err := h.usecase.CreateUserPrivilege(id, req); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Privilege assigned successfully",
	})
}

func (h *Handler) UpdateUserPrivilege(c *gin.Context) {
	id, ok := requireParam(c, "privilegeId", "Privilege ID")
	if !ok {
		return
	}
	var req UserPrivilegeUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mechanic.Error(c, mechanic.ValidationError("level is required"))
		return
	}
	if err := h.usecase.UpdateUserPrivilege(id, req); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Privilege updated successfully",
	})
}

func (h *Handler) DeleteUserPrivilege(c *gin.Context) {
	id, ok := requireParam(c, "privilegeId", "Privilege ID")
	if !ok {
		return
	}
	if err := h.usecase.DeleteUserPrivilege(id); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Privilege removed successfully",
	})
}

func (h *Handler) ListUserArea(c *gin.Context) {
	id, ok := requireParam(c, "id", "User ID")
	if !ok {
		return
	}
	var meta mechanic.ActionMeta
	if err := c.ShouldBindQuery(&meta); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid query parameters"))
		return
	}
	list, gridMeta, err := h.usecase.ListUserArea(id, meta)
	if err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User areas retrieved",
		"data":    list,
		"meta":    gridMeta,
	})
}

func (h *Handler) CreateUserArea(c *gin.Context) {
	id, ok := requireParam(c, "id", "User ID")
	if !ok {
		return
	}
	var req UserAreaCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mechanic.Error(c, mechanic.ValidationError("area_id is required"))
		return
	}
	if err := h.usecase.CreateUserArea(id, req); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Company area assigned successfully",
	})
}

func (h *Handler) UpdateUserArea(c *gin.Context) {
	areaId, ok := requireParam(c, "areaId", "Area ID")
	if !ok {
		return
	}
	var req UserAreaUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid request body"))
		return
	}
	if err := h.usecase.UpdateUserArea(areaId, req); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Company area updated successfully",
	})
}

func (h *Handler) DeleteUserArea(c *gin.Context) {
	areaId, ok := requireParam(c, "areaId", "Area ID")
	if !ok {
		return
	}
	if err := h.usecase.DeleteUserArea(areaId); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Company area removed successfully",
	})
}

func (h *Handler) ListAreaByCompany(c *gin.Context) {
	companyID := c.Query("company_id")
	if companyID == "" {
		mechanic.Error(c, mechanic.ValidationError("company_id is required"))
		return
	}
	list, err := h.usecase.ListAreaByCompany(companyID)
	if err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Company areas retrieved",
		"data":    list,
	})
}
