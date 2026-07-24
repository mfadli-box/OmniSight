package SP01

type UserCompanyItem struct {
	ID        string `json:"id"`
	CompanyID string `json:"company_id"`
	Name      string `json:"name"`
}

type ModuleTreeNode struct {
	ID       string           `json:"id"`
	Code     string           `json:"code"`
	Name     string           `json:"name"`
	Path     string           `json:"path"`
	IsPage   bool             `json:"is_page"`
	Children []ModuleTreeNode `json:"children,omitempty"`
}

type Repository interface {
	ListUserCompany(userID string) ([]UserCompanyItem, error)
	ListAllModuleTree() ([]ModuleTreeNode, error)
	ListUserModule(userID, companyID string) ([]ModuleTreeNode, error)
}

type UseCase interface {
	ListUserCompany(userID string) ([]UserCompanyItem, error)
	ListAllModuleTree() ([]ModuleTreeNode, error)
	ListUserModule(userID, companyID string) ([]ModuleTreeNode, error)
}
