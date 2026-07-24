package SM02

import "ict_rest/mechanic"

type ModuleListItem struct {
	ID        string `json:"id"`
	ParentID  string `json:"parent_id"`
	Code      string `json:"code"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	IsPage    bool   `json:"is_page"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
}

type ModuleCreateRequest struct {
	ParentID string `json:"parent_id"`
	Code     string `json:"code" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Path     string `json:"path" binding:"required"`
	IsPage   bool   `json:"is_page"`
	IsActive bool   `json:"is_active"`
}

type ModuleUpdateRequest struct {
	ParentID string `json:"parent_id"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	IsPage   *bool  `json:"is_page"`
	IsActive *bool  `json:"is_active"`
}

type Repository interface {
	CountModule(search string) (int, error)
	ListModule(search string, page, size int, sortBy, sortOrder string) ([]ModuleListItem, error)
	CreateModule(req ModuleCreateRequest) error
	UpdateModule(id string, req ModuleUpdateRequest) error
}

type UseCase interface {
	ListModule(meta mechanic.ActionMeta) ([]ModuleListItem, mechanic.GridMeta, error)
	CreateModule(req ModuleCreateRequest) error
	UpdateModule(id string, req ModuleUpdateRequest) error
}
