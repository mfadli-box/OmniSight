package SM06

import "ict_rest/mechanic"

type useCase struct {
	repo Repository
}

func NCase(r Repository) UseCase {
	return &useCase{repo: r}
}

func (u *useCase) ListLocationType(companyID string, meta mechanic.ActionMeta) ([]LocationTypeItem, mechanic.GridMeta, error) {
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
		return nil, mechanic.GridMeta{}, mechanic.InternalError("Failed to load location types", err)
	}
	return list, mechanic.BuildMeta(page, size, total), nil
}

func (u *useCase) CreateLocationType(companyID string, req LocationTypeCreateRequest) error {
	if req.Code == "" {
		return mechanic.ValidationError("Code is required")
	}
	if req.Name == "" {
		return mechanic.ValidationError("Name is required")
	}
	if err := u.repo.CreateLocationType(companyID, req); err != nil {
		return mechanic.InternalError("Failed to create location type", err)
	}
	return nil
}

func (u *useCase) UpdateLocationType(id string, req LocationTypeUpdateRequest) error {
	if id == "" {
		return mechanic.ValidationError("ID is required")
	}
	if err := u.repo.UpdateLocationType(id, req); err != nil {
		return mechanic.InternalError("Failed to update location type", err)
	}
	return nil
}

func (u *useCase) DeleteLocationType(id string) error {
	if id == "" {
		return mechanic.ValidationError("ID is required")
	}
	n, err := u.repo.DeleteLocationType(id)
	if err != nil {
		return mechanic.InternalError("Failed to delete location type", err)
	}
	if n == 0 {
		return mechanic.NotFound("Location type not found")
	}
	return nil
}

func (u *useCase) ListLocationTypeSelect(companyID string) ([]LocationTypeSelectItem, error) {
	list, err := u.repo.ListLocationTypeSelect(companyID)
	if err != nil {
		return nil, mechanic.InternalError("Failed to load location types", err)
	}
	return list, nil
}

func (u *useCase) ListLocation(companyID string, meta mechanic.ActionMeta) ([]LocationItem, mechanic.GridMeta, error) {
	page, size, err := mechanic.CheckMeta(meta.Page, meta.Size)
	if err != nil {
		return nil, mechanic.GridMeta{}, err
	}
	total, err := u.repo.CountLocation(companyID, meta.Search)
	if err != nil {
		return nil, mechanic.GridMeta{}, mechanic.InternalError("Failed to count locations", err)
	}
	list, err := u.repo.ListLocation(companyID, meta.Search, page, size, meta.SortBy, meta.SortOrder)
	if err != nil {
		return nil, mechanic.GridMeta{}, mechanic.InternalError("Failed to load locations", err)
	}
	return list, mechanic.BuildMeta(page, size, total), nil
}

func (u *useCase) CreateLocation(companyID string, req LocationCreateRequest) error {
	if req.Code == "" {
		return mechanic.ValidationError("Code is required")
	}
	if req.Name == "" {
		return mechanic.ValidationError("Name is required")
	}
	if req.LocationTypeID == "" {
		return mechanic.ValidationError("Location type is required")
	}
	if err := u.repo.CreateLocation(companyID, req); err != nil {
		return mechanic.InternalError("Failed to create location", err)
	}
	return nil
}

func (u *useCase) UpdateLocation(id string, req LocationUpdateRequest) error {
	if id == "" {
		return mechanic.ValidationError("ID is required")
	}
	if err := u.repo.UpdateLocation(id, req); err != nil {
		return mechanic.InternalError("Failed to update location", err)
	}
	return nil
}

func (u *useCase) DeleteLocation(id string) error {
	if id == "" {
		return mechanic.ValidationError("ID is required")
	}
	n, err := u.repo.DeleteLocation(id)
	if err != nil {
		return mechanic.InternalError("Failed to delete location", err)
	}
	if n == 0 {
		return mechanic.NotFound("Location not found")
	}
	return nil
}

func (u *useCase) ListLocationSelect(companyID string) ([]LocationSelectItem, error) {
	list, err := u.repo.ListLocationSelect(companyID)
	if err != nil {
		return nil, mechanic.InternalError("Failed to load locations", err)
	}
	return list, nil
}
