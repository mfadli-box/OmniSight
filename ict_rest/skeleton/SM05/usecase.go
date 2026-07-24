package SM05

import (
	"ict_rest/mechanic"
)

type useCase struct {
	repo Repository
}

func NCase(r Repository) UseCase {
	return &useCase{repo: r}
}

func (u *useCase) ListSession(meta mechanic.ActionMeta) ([]SessionListItem, mechanic.GridMeta, error) {
	page, size, err := mechanic.CheckMeta(meta.Page, meta.Size)
	if err != nil {
		return nil, mechanic.GridMeta{}, err
	}
	total, err := u.repo.CountSession(meta.Search)
	if err != nil {
		return nil, mechanic.GridMeta{}, mechanic.InternalError("Failed to count sessions", err)
	}
	list, err := u.repo.ListSession(meta.Search, page, size, meta.SortBy, meta.SortOrder)
	if err != nil {
		return nil, mechanic.GridMeta{}, mechanic.InternalError("Failed to load sessions", err)
	}
	return list, mechanic.BuildMeta(page, size, total), nil
}

func (u *useCase) GetSessionDetail(id string) (*SessionDetailItem, error) {
	if id == "" {
		return nil, mechanic.ValidationError("Session ID is required")
	}
	detail, err := u.repo.GetSessionDetail(id)
	if err != nil {
		return nil, mechanic.InternalError("Failed to load session detail", err)
	}
	return detail, nil
}

func (u *useCase) RevokeSession(id string) error {
	if id == "" {
		return mechanic.ValidationError("Session ID is required")
	}
	n, err := u.repo.RevokeSession(id)
	if err != nil {
		return mechanic.InternalError("Failed to revoke session", err)
	}
	if n == 0 {
		return mechanic.NotFound("Session not found")
	}
	return nil
}
