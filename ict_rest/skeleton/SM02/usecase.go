package SM02

import (
	"ict_rest/mechanic"
)

type useCase struct {
	repo Repository
}

func NCase(r Repository) UseCase {
	return &useCase{repo: r}
}

func (u *useCase) ListModule(meta mechanic.ActionMeta) (
	[]ModuleListItem, mechanic.GridMeta, error) {
	page, size, err := mechanic.CheckMeta(meta.Page, meta.Size)
	if err != nil {
		return nil, mechanic.GridMeta{}, err
	}
	total, err := u.repo.CountModule(meta.Search)
	if err != nil {
		return nil, mechanic.GridMeta{}, mechanic.InternalError("Failed to count modules", err)
	}
	list, err := u.repo.ListModule(
		meta.Search,
		page,
		size,
		meta.SortBy,
		meta.SortOrder,
	)
	if err != nil {
		return nil, mechanic.GridMeta{}, mechanic.InternalError("Failed to list modules", err)
	}
	return list, mechanic.BuildMeta(page, size, total), nil
}

func (u *useCase) CreateModule(req ModuleCreateRequest) error {
	if req.Code == "" {
		return mechanic.ValidationError("Code is required")
	}
	if req.Name == "" {
		return mechanic.ValidationError("Name is required")
	}
	if req.Path == "" {
		return mechanic.ValidationError("Path is required")
	}
	if err := u.repo.CreateModule(req); err != nil {
		return mechanic.InternalError("Failed to create module", err)
	}
	return nil
}

func (u *useCase) UpdateModule(id string, req ModuleUpdateRequest) error {
	if id == "" {
		return mechanic.ValidationError("Module ID is required")
	}
	if err := u.repo.UpdateModule(id, req); err != nil {
		return mechanic.InternalError("Failed to update module", err)
	}
	return nil
}
