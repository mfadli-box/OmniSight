package SM01

import "ict_rest/mechanic"

type UserListItem struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Fullname  string `json:"fullname"`
	Phone     string `json:"phone"`
	CompanyID string `json:"company_id"`
	Role      string `json:"role"`
	Job       string `json:"job"`
	IsAdmin   bool   `json:"is_admin"`
	IsHris    bool   `json:"is_hris"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
}

type UserCreateRequest struct {
	Username  string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Fullname  string `json:"fullname" binding:"required"`
	Phone     string `json:"phone"`
	CompanyID string `json:"company_id"`
	Role      string `json:"role"`
	Job       string `json:"job"`
	IsAdmin   bool   `json:"is_admin"`
	IsHris    bool   `json:"is_hris"`
	IsActive  bool   `json:"is_active"`
}

type UserUpdateRequest struct {
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
	Job      string `json:"job"`
	IsAdmin  bool   `json:"is_admin"`
	IsHris   bool   `json:"is_hris"`
	IsActive bool   `json:"is_active"`
}

type HrisCompanyItem struct {
	ID       string `json:"id"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	HrisLink string `json:"hris_link"`
}

type CompanySelectItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ModuleSelectItem struct {
	ID       string `json:"id"`
	ParentID string `json:"parent_id"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	IsPage   bool   `json:"is_page"`
	IsActive bool   `json:"is_active"`
}

type UserCompanyListItem struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	CompanyID   string `json:"company_id"`
	CompanyName string `json:"company_name"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at"`
}

type UserCompanyRequest struct {
	CompanyID string `json:"company_id"`
	IsActive  *bool  `json:"is_active"`
}

type UserPrivilegeListItem struct {
	ID            string `json:"id"`
	UserCompanyID string `json:"user_company_id"`
	ModuleID      string `json:"module_id"`
	CompanyName   string `json:"company_name"`
	ModuleCode    string `json:"module_code"`
	ModuleName    string `json:"module_name"`
	Level         string `json:"level"`
	CreatedAt     string `json:"created_at"`
}

type UserPrivilegeRequest struct {
	UserCompanyID string `json:"user_company_id" binding:"required"`
	ModuleID      string `json:"module_id" binding:"required"`
	Level         string `json:"level"`
}

type UserPrivilegeUpdateRequest struct {
	Level string `json:"level" binding:"required"`
}

type UserLocationListItem struct {
	ID             string `json:"id"`
	UserID         string `json:"user_id"`
	LocationTypeID string `json:"location_type_id"`
	CompanyName    string `json:"company_name"`
	TypeCode       string `json:"type_code"`
	TypeName       string `json:"type_name"`
	IsActive       bool   `json:"is_active"`
	CreatedAt      string `json:"created_at"`
}

type UserLocationCreateRequest struct {
	LocationTypeID string `json:"location_type_id" binding:"required"`
	IsActive       *bool  `json:"is_active"`
}

type UserLocationUpdateRequest struct {
	IsActive bool `json:"is_active"`
}

type LocationTypeSelectItem struct {
	ID   string `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type Repository interface {
	CountUser(companyID, search string) (int, error)
	ListUser(companyID, search string, page, size int, sortBy, sortOrder string) ([]UserListItem, error)
	UpdateUserPassword(id, hashedPassword string) error
	CreateUser(req UserCreateRequest, hashedPassword string) (string, error)
	UpdateUser(id string, req UserUpdateRequest) error
	ListHrisCompany() ([]HrisCompanyItem, error)
	ListAllCompany() ([]CompanySelectItem, error)
	ListAllModule() ([]ModuleSelectItem, error)
	AssignCompany(userID, companyID string) error
	CountUserCompany(userID, search string) (int, error)
	ListUserCompany(userID, search string, page, size int, sortBy, sortOrder string) ([]UserCompanyListItem, error)
	CreateUserCompany(userID string, req UserCompanyRequest) error
	UpdateUserCompany(userID, companyID string, req UserCompanyRequest) error
	DeleteUserCompany(userID, companyID string) error
	CountUserPrivilege(userID, search string) (int, error)
	ListUserPrivilege(userID, search string, page, size int, sortBy, sortOrder string) ([]UserPrivilegeListItem, error)
	CreateUserPrivilege(userID string, req UserPrivilegeRequest) error
	UpdateUserPrivilege(id string, req UserPrivilegeUpdateRequest) error
	DeleteUserPrivilege(id string) error
	CountUserLocation(userID, search string) (int, error)
	ListUserLocation(userID, search string, page, size int, sortBy, sortOrder string) ([]UserLocationListItem, error)
	CreateUserLocation(userID string, req UserLocationCreateRequest) error
	UpdateUserLocation(id string, req UserLocationUpdateRequest) error
	DeleteUserLocation(id string) error
	ListLocationTypeByCompany(companyID string) ([]LocationTypeSelectItem, error)
}

type UseCase interface {
	ListUser(companyID string, meta mechanic.ActionMeta) ([]UserListItem, mechanic.GridMeta, error)
	CreateUser(req UserCreateRequest) error
	UpdateUser(id string, req UserUpdateRequest) error
	ListHrisCompany() ([]HrisCompanyItem, error)
	ListAllCompany() ([]CompanySelectItem, error)
	ListAllModule() ([]ModuleSelectItem, error)
	AssignCompany(userID, companyID string) error
	ListUserCompany(userID string, meta mechanic.ActionMeta) ([]UserCompanyListItem, mechanic.GridMeta, error)
	CreateUserCompany(userID string, req UserCompanyRequest) error
	UpdateUserCompany(userID, companyID string, req UserCompanyRequest) error
	DeleteUserCompany(userID, companyID string) error
	ListUserPrivilege(userID string, meta mechanic.ActionMeta) ([]UserPrivilegeListItem, mechanic.GridMeta, error)
	CreateUserPrivilege(userID string, req UserPrivilegeRequest) error
	UpdateUserPrivilege(id string, req UserPrivilegeUpdateRequest) error
	DeleteUserPrivilege(id string) error
	ListUserLocation(userID string, meta mechanic.ActionMeta) ([]UserLocationListItem, mechanic.GridMeta, error)
	CreateUserLocation(userID string, req UserLocationCreateRequest) error
	UpdateUserLocation(id string, req UserLocationUpdateRequest) error
	DeleteUserLocation(id string) error
	ListLocationTypeByCompany(companyID string) ([]LocationTypeSelectItem, error)
}
