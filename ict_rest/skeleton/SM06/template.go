package SM06

import "ict_rest/mechanic"

type LocationTypeItem struct {
	ID          string `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at"`
}

type LocationTypeCreateRequest struct {
	Code        string `json:"code" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
	IsActive    *bool  `json:"is_active"`
}

type LocationTypeUpdateRequest struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
	IsActive    *bool  `json:"is_active"`
}

type LocationItem struct {
	ID             string `json:"id"`
	LocationTypeID string `json:"location_type_id"`
	ParentID       string `json:"parent_id"`
	Code           string `json:"code"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Address        string `json:"address"`
	City           string `json:"city"`
	Province       string `json:"province"`
	Country        string `json:"country"`
	PostalCode     string `json:"postal_code"`
	Latitude       string `json:"latitude"`
	Longitude      string `json:"longitude"`
	Phone          string `json:"phone"`
	Email          string `json:"email"`
	Timezone       string `json:"timezone"`
	Status         string `json:"status"`
	IsActive       bool   `json:"is_active"`
	CreatedAt      string `json:"created_at"`
}

type LocationCreateRequest struct {
	LocationTypeID string  `json:"location_type_id" binding:"required"`
	ParentID       string  `json:"parent_id"`
	Code           string  `json:"code" binding:"required"`
	Name           string  `json:"name" binding:"required"`
	Description    string  `json:"description"`
	Address        string  `json:"address"`
	City           string  `json:"city"`
	Province       string  `json:"province"`
	Country        string  `json:"country"`
	PostalCode     string  `json:"postal_code"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	Phone          string  `json:"phone"`
	Email          string  `json:"email"`
	Timezone       string  `json:"timezone"`
	Status         string  `json:"status"`
	IsActive       *bool   `json:"is_active"`
}

type LocationUpdateRequest struct {
	LocationTypeID string  `json:"location_type_id"`
	ParentID       string  `json:"parent_id"`
	Code           string  `json:"code"`
	Name           string  `json:"name"`
	Description    string  `json:"description"`
	Address        string  `json:"address"`
	City           string  `json:"city"`
	Province       string  `json:"province"`
	Country        string  `json:"country"`
	PostalCode     string  `json:"postal_code"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	Phone          string  `json:"phone"`
	Email          string  `json:"email"`
	Timezone       string  `json:"timezone"`
	Status         string  `json:"status"`
	IsActive       *bool   `json:"is_active"`
}

type LocationTypeSelectItem struct {
	ID   string `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type LocationSelectItem struct {
	ID       string `json:"id"`
	ParentID string `json:"parent_id"`
	Code     string `json:"code"`
	Name     string `json:"name"`
}

type Repository interface {
	CountLocationType(companyID, search string) (int, error)
	ListLocationType(companyID, search string, page, size int, sortBy, sortOrder string) ([]LocationTypeItem, error)
	CreateLocationType(companyID string, req LocationTypeCreateRequest) error
	UpdateLocationType(id string, req LocationTypeUpdateRequest) error
	DeleteLocationType(id string) (int64, error)
	ListLocationTypeSelect(companyID string) ([]LocationTypeSelectItem, error)

	CountLocation(companyID, search string) (int, error)
	ListLocation(companyID, search string, page, size int, sortBy, sortOrder string) ([]LocationItem, error)
	CreateLocation(companyID string, req LocationCreateRequest) error
	UpdateLocation(id string, req LocationUpdateRequest) error
	DeleteLocation(id string) (int64, error)
	ListLocationSelect(companyID string) ([]LocationSelectItem, error)

}

type UseCase interface {
	ListLocationType(companyID string, meta mechanic.ActionMeta) ([]LocationTypeItem, mechanic.GridMeta, error)
	CreateLocationType(companyID string, req LocationTypeCreateRequest) error
	UpdateLocationType(id string, req LocationTypeUpdateRequest) error
	DeleteLocationType(id string) error
	ListLocationTypeSelect(companyID string) ([]LocationTypeSelectItem, error)

	ListLocation(companyID string, meta mechanic.ActionMeta) ([]LocationItem, mechanic.GridMeta, error)
	CreateLocation(companyID string, req LocationCreateRequest) error
	UpdateLocation(id string, req LocationUpdateRequest) error
	DeleteLocation(id string) error
	ListLocationSelect(companyID string) ([]LocationSelectItem, error)
}
