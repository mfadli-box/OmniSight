package SM03

import (
	"ict_rest/mechanic"
	"strings"
)

type useCase struct {
	repo Repository
}

func NCase(r Repository) UseCase {
	return &useCase{repo: r}
}

func (u *useCase) ListCompany(meta mechanic.ActionMeta) (
	[]CompanyListItem, mechanic.GridMeta, error) {
	page, size, err := mechanic.CheckMeta(meta.Page, meta.Size)
	if err != nil {
		return nil, mechanic.GridMeta{}, err
	}
	total, err := u.repo.CountCompany(meta.Search)
	if err != nil {
		return nil, mechanic.GridMeta{}, mechanic.InternalError("Failed to count companies", err)
	}
	list, err := u.repo.ListCompany(meta.Search, page, size, meta.SortBy, meta.SortOrder)
	if err != nil {
		return nil, mechanic.GridMeta{}, mechanic.InternalError("Failed to list companies", err)
	}
	return list, mechanic.BuildMeta(page, size, total), nil
}

func (u *useCase) ListModule() ([]ModuleItem, error) {
	list, err := u.repo.ListModule()
	if err != nil {
		return nil, mechanic.InternalError("Failed to list modules", err)
	}
	return list, nil
}

func (u *useCase) CreateCompany(req CompanyCreateRequest) error {
	if strings.TrimSpace(req.Code) == "" {
		return mechanic.ValidationError("Company code is required")
	}
	if strings.TrimSpace(req.Name) == "" {
		return mechanic.ValidationError("Company name is required")
	}
	if req.Valuta == "" {
		req.Valuta = "IDR"
	}
	if err := u.repo.CreateCompany(req); err != nil {
		return mechanic.InternalError("Failed to create company", err)
	}
	return nil
}

func (u *useCase) UpdateCompany(id string, req CompanyUpdateRequest) error {
	if strings.TrimSpace(id) == "" {
		return mechanic.ValidationError("Company ID is required")
	}
	if err := u.repo.UpdateCompany(id, req); err != nil {
		return mechanic.InternalError("Failed to update company", err)
	}
	return nil
}

func (u *useCase) ListCompanyModule(companyID string, meta mechanic.ActionMeta) (
	[]CompanyModuleItem, mechanic.GridMeta, error) {
	if strings.TrimSpace(companyID) == "" {
		return nil, mechanic.GridMeta{}, mechanic.ValidationError("Company ID is required")
	}
	page, size, err := mechanic.CheckMeta(meta.Page, meta.Size)
	if err != nil {
		return nil, mechanic.GridMeta{}, err
	}
	total, err := u.repo.CountCompanyModule(companyID, meta.Search)
	if err != nil {
		return nil, mechanic.GridMeta{}, mechanic.InternalError("Failed to count company modules", err)
	}
	list, err := u.repo.ListCompanyModule(companyID, meta.Search, page, size, meta.SortBy, meta.SortOrder)
	if err != nil {
		return nil, mechanic.GridMeta{}, mechanic.InternalError("Failed to list company modules", err)
	}
	return list, mechanic.BuildMeta(page, size, total), nil
}

func (u *useCase) AssignCompanyModule(companyID, moduleID string) error {
	if strings.TrimSpace(companyID) == "" {
		return mechanic.ValidationError("Company ID is required")
	}
	if strings.TrimSpace(moduleID) == "" {
		return mechanic.ValidationError("Module ID is required")
	}
	if err := u.repo.AssignCompanyModule(companyID, moduleID); err != nil {
		return mechanic.InternalError("Failed to assign company module", err)
	}
	return nil
}

func (u *useCase) UpdateCompanyModule(companyID, moduleID string, isActive bool) error {
	if strings.TrimSpace(companyID) == "" {
		return mechanic.ValidationError("Company ID is required")
	}
	if strings.TrimSpace(moduleID) == "" {
		return mechanic.ValidationError("Module ID is required")
	}
	if err := u.repo.UpdateCompanyModule(companyID, moduleID, isActive); err != nil {
		return mechanic.InternalError("Failed to update company module", err)
	}
	return nil
}

func (u *useCase) ListLocationType(companyID string, meta mechanic.ActionMeta) ([]LocationTypeItem, mechanic.GridMeta, error) {
	if strings.TrimSpace(companyID) == "" {
		return nil, mechanic.GridMeta{}, mechanic.ValidationError("Company ID is required")
	}
	page, size, err := mechanic.CheckMeta(meta.Page, meta.Size)
	if err != nil {
		return nil, mechanic.GridMeta{}, err
	}
	total, err := u.repo.CountLocationType(companyID, meta.Search)
	if err != nil {
		return nil, mechanic.GridMeta{}, mechanic.InternalError("Failed to count location types", err)
	}
	list, err := u.repo.ListLocationType(companyID, meta.Search, page, size, meta.SortBy, meta.SortOrder)
	if err != nil {
		return nil, mechanic.GridMeta{}, mechanic.InternalError("Failed to list location types", err)
	}
	return list, mechanic.BuildMeta(page, size, total), nil
}

func (u *useCase) CreateLocationType(companyID string, req LocationTypeCreateRequest) error {
	if strings.TrimSpace(companyID) == "" {
		return mechanic.ValidationError("Company ID is required")
	}
	if strings.TrimSpace(req.Code) == "" {
		return mechanic.ValidationError("Type code is required")
	}
	if strings.TrimSpace(req.Name) == "" {
		return mechanic.ValidationError("Type name is required")
	}
	if err := u.repo.CreateLocationType(companyID, req); err != nil {
		return mechanic.InternalError("Failed to create location type", err)
	}
	return nil
}

func (u *useCase) UpdateLocationType(id string, req LocationTypeUpdateRequest) error {
	if strings.TrimSpace(id) == "" {
		return mechanic.ValidationError("Type ID is required")
	}
	if err := u.repo.UpdateLocationType(id, req); err != nil {
		return mechanic.InternalError("Failed to update location type", err)
	}
	return nil
}

func (u *useCase) DeleteLocationType(id string) error {
	if strings.TrimSpace(id) == "" {
		return mechanic.ValidationError("Type ID is required")
	}
	if err := u.repo.DeleteLocationType(id); err != nil {
		return mechanic.InternalError("Failed to delete location type", err)
	}
	return nil
}
