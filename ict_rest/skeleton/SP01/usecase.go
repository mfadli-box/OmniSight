package SP01

import "ict_rest/mechanic"

type useCase struct {
	repo Repository
}

func NCase(r Repository) UseCase {
	return &useCase{repo: r}
}

func (u *useCase) ListUserCompany(userID string) ([]UserCompanyItem, error) {
	list, err := u.repo.ListUserCompany(userID)
	if err != nil {
		return nil, mechanic.InternalError("Failed to list user companies", err)
	}
	return list, nil
}

func (u *useCase) ListAllModuleTree() ([]ModuleTreeNode, error) {
	list, err := u.repo.ListAllModuleTree()
	if err != nil {
		return nil, mechanic.InternalError("Failed to list module tree", err)
	}
	return list, nil
}

func (u *useCase) ListUserModule(userID, companyID string) ([]ModuleTreeNode, error) {
	if companyID == "" {
		return []ModuleTreeNode{}, nil
	}
	list, err := u.repo.ListUserModule(userID, companyID)
	if err != nil {
		return nil, mechanic.InternalError("Failed to list user modules", err)
	}
	return list, nil
}
