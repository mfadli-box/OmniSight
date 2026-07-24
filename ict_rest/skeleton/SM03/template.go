package SM03

import "ict_rest/mechanic"

type CompanyListItem struct {
	ID        string `json:"id"`
	Code      string `json:"code"`
	Name      string `json:"name"`
	VatID     string `json:"vat_id"`
	RegNo     string `json:"reg_no"`
	Address   string `json:"address"`
	Valuta    string `json:"valuta"`
	HrisLink  string `json:"hris_link"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
}

type ModuleItem struct {
	ID       string `json:"id"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	ParentID string `json:"parent_id"`
	IsActive bool   `json:"is_active"`
}

type CompanyCreateRequest struct {
	Code     string `json:"code" binding:"required"`
	Name     string `json:"name" binding:"required"`
	VatID    string `json:"vat_id"`
	RegNo    string `json:"reg_no"`
	Address  string `json:"address"`
	Valuta   string `json:"valuta"`
	HrisLink string `json:"hris_link"`
	IsActive bool   `json:"is_active"`
}

type CompanyUpdateRequest struct {
	Name     string `json:"name"`
	VatID    string `json:"vat_id"`
	RegNo    string `json:"reg_no"`
	Address  string `json:"address"`
	Valuta   string `json:"valuta"`
	HrisLink string `json:"hris_link"`
	IsActive bool   `json:"is_active"`
}

type CompanyModuleItem struct {
	ID       string `json:"id"`
	ModuleID string `json:"module_id"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	IsActive bool   `json:"is_active"`
}

type CompanyModuleAssignRequest struct {
	ModuleID string `json:"module_id" binding:"required"`
}

type CompanyModuleUpdateRequest struct {
	IsActive bool `json:"is_active"`
}

type LocationTypeItem struct {
	ID          string `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at"`
}

type LocationTypeCreateRequest struct {
	Code        string `json:"code" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

type LocationTypeUpdateRequest struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

type Repository interface {
	CountCompany(search string) (int, error)
	ListCompany(search string, page, size int, sortBy, sortOrder string) ([]CompanyListItem, error)
	ListModule() ([]ModuleItem, error)
	CreateCompany(req CompanyCreateRequest) error
	UpdateCompany(id string, req CompanyUpdateRequest) error
	CountCompanyModule(companyID, search string) (int, error)
	ListCompanyModule(companyID, search string, page, size int, sortBy, sortOrder string) ([]CompanyModuleItem, error)
	AssignCompanyModule(companyID, moduleID string) error
	UpdateCompanyModule(companyID, moduleID string, isActive bool) error
	CountLocationType(companyID, search string) (int, error)
	ListLocationType(companyID, search string, page, size int, sortBy, sortOrder string) ([]LocationTypeItem, error)
	CreateLocationType(companyID string, req LocationTypeCreateRequest) error
	UpdateLocationType(id string, req LocationTypeUpdateRequest) error
	DeleteLocationType(id string) error
}

type UseCase interface {
	ListCompany(meta mechanic.ActionMeta) ([]CompanyListItem, mechanic.GridMeta, error)
	ListModule() ([]ModuleItem, error)
	CreateCompany(req CompanyCreateRequest) error
	UpdateCompany(id string, req CompanyUpdateRequest) error
	ListCompanyModule(companyID string, meta mechanic.ActionMeta) ([]CompanyModuleItem, mechanic.GridMeta, error)
	AssignCompanyModule(companyID, moduleID string) error
	UpdateCompanyModule(companyID, moduleID string, isActive bool) error
	ListLocationType(companyID string, meta mechanic.ActionMeta) ([]LocationTypeItem, mechanic.GridMeta, error)
	CreateLocationType(companyID string, req LocationTypeCreateRequest) error
	UpdateLocationType(id string, req LocationTypeUpdateRequest) error
	DeleteLocationType(id string) error
}
