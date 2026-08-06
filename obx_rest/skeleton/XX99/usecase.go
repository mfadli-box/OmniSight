package XX99

import "obx_rest/mechanic"

type useCase struct {
	repo Repository
}

func NCase(r Repository) UseCase {
	return &useCase{repo: r}
}

func (u *useCase) List(meta mechanic.ActionMeta) ([]XX99Item, mechanic.GridMeta, error) {
	page, size, err := mechanic.CheckMeta(meta.Page, meta.Size)
	if err != nil {
		return nil, mechanic.GridMeta{}, err
	}
	meta.Page = page
	meta.Size = size
	return u.repo.List(meta)
}

func (u *useCase) GetByID(id string) (XX99Item, error) {
	if id == "" {
		return XX99Item{}, mechanic.ValidationError("ID is required")
	}
	return u.repo.GetByID(id)
}

func (u *useCase) Create(req XX99CreateRequest) error {
	if req.CompanyID == "" || req.Code == "" || req.Name == "" {
		return mechanic.ValidationError("company_id, code, and name are required")
	}
	return u.repo.Create(req)
}

func (u *useCase) Update(id string, req XX99UpdateRequest) error {
	if id == "" {
		return mechanic.ValidationError("ID is required")
	}
	return u.repo.Update(id, req)
}

func (u *useCase) Delete(id string) error {
	if id == "" {
		return mechanic.ValidationError("ID is required")
	}
	return u.repo.Delete(id)
}
