package SP03

import "ict_rest/mechanic"

type useCase struct {
	repo Repository
}

func NCase(r Repository) UseCase {
	return &useCase{repo: r}
}

func (u *useCase) ListActions(userID, search string, page, size int, sortBy, sortOrder string) (
	[]UserActionItem, mechanic.GridMeta, error) {
	page, size, err := mechanic.CheckMeta(page, size)
	if err != nil {
		return nil, mechanic.GridMeta{}, err
	}
	return u.repo.ListActions(
		userID,
		search,
		page,
		size,
		sortBy,
		sortOrder,
	)
}
