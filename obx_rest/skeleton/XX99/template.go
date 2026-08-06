package XX99

import "obx_rest/mechanic"

type XX99Item struct {
	ID        string `json:"id"`
	CompanyID string `json:"company_id"`
	Code      string `json:"code"`
	Name      string `json:"name"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type XX99CreateRequest struct {
	CompanyID string `json:"company_id" binding:"required"`
	Code      string `json:"code" binding:"required"`
	Name      string `json:"name" binding:"required"`
	IsActive  bool   `json:"is_active"`
}

type XX99UpdateRequest struct {
	CompanyID string `json:"company_id"`
	Code      string `json:"code"`
	Name      string `json:"name"`
	IsActive  *bool  `json:"is_active"`
}

type Repository interface {
	List(meta mechanic.ActionMeta) ([]XX99Item, mechanic.GridMeta, error)
	GetByID(id string) (XX99Item, error)
	Create(req XX99CreateRequest) error
	Update(id string, req XX99UpdateRequest) error
	Delete(id string) error
}

type UseCase interface {
	List(meta mechanic.ActionMeta) ([]XX99Item, mechanic.GridMeta, error)
	GetByID(id string) (XX99Item, error)
	Create(req XX99CreateRequest) error
	Update(id string, req XX99UpdateRequest) error
	Delete(id string) error
}
