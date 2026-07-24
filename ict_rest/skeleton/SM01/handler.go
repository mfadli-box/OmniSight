package SM01

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
	id := c.Param("id")
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
	id := c.Param("id")
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
	id := c.Param("id")
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

func (h *Handler) CreateUserCompany(c *gin.Context) {
	id := c.Param("id")
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
	id := c.Param("id")
	companyID := c.Param("companyId")
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
	id := c.Param("id")
	companyID := c.Param("companyId")
	if err := h.usecase.DeleteUserCompany(id, companyID); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Company removed successfully",
	})
}

func (h *Handler) ListUserPrivilege(c *gin.Context) {
	id := c.Param("id")
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
	id := c.Param("id")
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
	id := c.Param("privilegeId")
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
	id := c.Param("privilegeId")
	if err := h.usecase.DeleteUserPrivilege(id); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Privilege removed successfully",
	})
}

func (h *Handler) ListUserLocation(c *gin.Context) {
	id := c.Param("id")
	var meta mechanic.ActionMeta
	if err := c.ShouldBindQuery(&meta); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid query parameters"))
		return
	}
	list, gridMeta, err := h.usecase.ListUserLocation(id, meta)
	if err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User locations retrieved",
		"data":    list,
		"meta":    gridMeta,
	})
}

func (h *Handler) CreateUserLocation(c *gin.Context) {
	id := c.Param("id")
	var req UserLocationCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mechanic.Error(c, mechanic.ValidationError("location_type_id is required"))
		return
	}
	if err := h.usecase.CreateUserLocation(id, req); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Location assigned successfully",
	})
}

func (h *Handler) UpdateUserLocation(c *gin.Context) {
	locationId := c.Param("locationId")
	var req UserLocationUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mechanic.Error(c, mechanic.ValidationError("Invalid request body"))
		return
	}
	if err := h.usecase.UpdateUserLocation(locationId, req); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Location updated successfully",
	})
}

func (h *Handler) DeleteUserLocation(c *gin.Context) {
	locationId := c.Param("locationId")
	if err := h.usecase.DeleteUserLocation(locationId); err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Location removed successfully",
	})
}

func (h *Handler) ListLocationTypeByCompany(c *gin.Context) {
	companyID := c.Query("company_id")
	if companyID == "" {
		mechanic.Error(c, mechanic.ValidationError("company_id is required"))
		return
	}
	list, err := h.usecase.ListLocationTypeByCompany(companyID)
	if err != nil {
		mechanic.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Location types retrieved",
		"data":    list,
	})
}
