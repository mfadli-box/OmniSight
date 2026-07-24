package SM01

import (
	"ict_rest/mechanic"

	"golang.org/x/crypto/bcrypt"
)

type useCase struct {
	repo Repository
}

func NCase(r Repository) UseCase {
	return &useCase{repo: r}
}

func (u *useCase) ListUser(companyID string, meta mechanic.ActionMeta) (
	[]UserListItem, mechanic.GridMeta, error) {
	page, size, err := mechanic.CheckMeta(meta.Page, meta.Size)
	if err != nil {
		return nil, mechanic.GridMeta{}, err
	}
	total, err := u.repo.CountUser(companyID, meta.Search)
	if err != nil {
		return nil, mechanic.GridMeta{}, mechanic.InternalError("Failed to count users", err)
	}
	list, err := u.repo.ListUser(companyID, meta.Search, page, size, meta.SortBy, meta.SortOrder)
	if err != nil {
		return nil, mechanic.GridMeta{}, mechanic.InternalError("Failed to list users", err)
	}
	return list, mechanic.BuildMeta(page, size, total), nil
}

func (u *useCase) CreateUser(req UserCreateRequest) error {
	if req.Username == "" {
		return mechanic.ValidationError("Username is required")
	}
	if req.Email == "" {
		return mechanic.ValidationError("Email is required")
	}
	if req.Password == "" {
		return mechanic.ValidationError("Password is required")
	}
	if req.Fullname == "" {
		return mechanic.ValidationError("Fullname is required")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return mechanic.InternalError("Failed to hash password", err)
	}
	if req.Role == "" {
		req.Role = "staff"
	}
	userID, err := u.repo.CreateUser(req, string(hashed))
	if err != nil {
		return mechanic.InternalError("Failed to create user", err)
	}
	if req.IsHris && req.CompanyID != "" {
		_ = u.repo.CreateUserCompany(userID, UserCompanyRequest{CompanyID: req.CompanyID})
	}
	return nil
}

func (u *useCase) UpdateUser(id string, req UserUpdateRequest) error {
	if id == "" {
		return mechanic.ValidationError("User ID is required")
	}
	if req.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return mechanic.InternalError("Failed to hash password", err)
		}
		if err := u.repo.UpdateUserPassword(id, string(hashed)); err != nil {
			return mechanic.InternalError("Failed to update user password", err)
		}
	}
	if err := u.repo.UpdateUser(id, req); err != nil {
		return mechanic.InternalError("Failed to update user", err)
	}
	return nil
}

func (u *useCase) ListHrisCompany() ([]HrisCompanyItem, error) {
	list, err := u.repo.ListHrisCompany()
	if err != nil {
		return nil, mechanic.InternalError("Failed to list HRIS companies", err)
	}
	return list, nil
}

func (u *useCase) ListAllCompany() ([]CompanySelectItem, error) {
	list, err := u.repo.ListAllCompany()
	if err != nil {
		return nil, mechanic.InternalError("Failed to list companies", err)
	}
	return list, nil
}

func (u *useCase) ListAllModule() ([]ModuleSelectItem, error) {
	list, err := u.repo.ListAllModule()
	if err != nil {
		return nil, mechanic.InternalError("Failed to list modules", err)
	}
	return list, nil
}

func (u *useCase) AssignCompany(userID, companyID string) error {
	if userID == "" || companyID == "" {
		return mechanic.ValidationError("User ID and Company ID are required")
	}
	if err := u.repo.AssignCompany(userID, companyID); err != nil {
		return mechanic.InternalError("Failed to assign company", err)
	}
	return nil
}

func (u *useCase) ListUserCompany(userID string, meta mechanic.ActionMeta) ([]UserCompanyListItem, mechanic.GridMeta, error) {
	page, size, err := mechanic.CheckMeta(meta.Page, meta.Size)
	if err != nil {
		return nil, mechanic.GridMeta{}, err
	}
	total, err := u.repo.CountUserCompany(userID, meta.Search)
	if err != nil {
		return nil, mechanic.GridMeta{}, mechanic.InternalError("Failed to count user companies", err)
	}
	list, err := u.repo.ListUserCompany(userID, meta.Search, page, size, meta.SortBy, meta.SortOrder)
	if err != nil {
		return nil, mechanic.GridMeta{}, mechanic.InternalError("Failed to list user companies", err)
	}
	return list, mechanic.BuildMeta(page, size, total), nil
}

func (u *useCase) CreateUserCompany(userID string, req UserCompanyRequest) error {
	if userID == "" || req.CompanyID == "" {
		return mechanic.ValidationError("User ID and Company ID are required")
	}
	if err := u.repo.CreateUserCompany(userID, req); err != nil {
		return mechanic.InternalError("Failed to create user company", err)
	}
	return nil
}

func (u *useCase) UpdateUserCompany(userID, companyID string, req UserCompanyRequest) error {
	if userID == "" || companyID == "" {
		return mechanic.ValidationError("User ID and Company ID are required")
	}
	if err := u.repo.UpdateUserCompany(userID, companyID, req); err != nil {
		return mechanic.InternalError("Failed to update user company", err)
	}
	return nil
}

func (u *useCase) DeleteUserCompany(userID, companyID string) error {
	if userID == "" || companyID == "" {
		return mechanic.ValidationError("User ID and Company ID are required")
	}
	if err := u.repo.DeleteUserCompany(userID, companyID); err != nil {
		return mechanic.InternalError("Failed to delete user company", err)
	}
	return nil
}

func (u *useCase) ListUserPrivilege(userID string, meta mechanic.ActionMeta) ([]UserPrivilegeListItem, mechanic.GridMeta, error) {
	page, size, err := mechanic.CheckMeta(meta.Page, meta.Size)
	if err != nil {
		return nil, mechanic.GridMeta{}, err
	}
	total, err := u.repo.CountUserPrivilege(userID, meta.Search)
	if err != nil {
		return nil, mechanic.GridMeta{}, mechanic.InternalError("Failed to count user privileges", err)
	}
	list, err := u.repo.ListUserPrivilege(userID, meta.Search, page, size, meta.SortBy, meta.SortOrder)
	if err != nil {
		return nil, mechanic.GridMeta{}, mechanic.InternalError("Failed to list user privileges", err)
	}
	return list, mechanic.BuildMeta(page, size, total), nil
}

func (u *useCase) CreateUserPrivilege(userID string, req UserPrivilegeRequest) error {
	if userID == "" || req.UserCompanyID == "" || req.ModuleID == "" {
		return mechanic.ValidationError("User Company ID and Module ID are required")
	}
	if err := u.repo.CreateUserPrivilege(userID, req); err != nil {
		return mechanic.InternalError("Failed to create user privilege", err)
	}
	return nil
}

func (u *useCase) UpdateUserPrivilege(id string, req UserPrivilegeUpdateRequest) error {
	if id == "" {
		return mechanic.ValidationError("Privilege ID is required")
	}
	if err := u.repo.UpdateUserPrivilege(id, req); err != nil {
		return mechanic.InternalError("Failed to update user privilege", err)
	}
	return nil
}

func (u *useCase) DeleteUserPrivilege(id string) error {
	if id == "" {
		return mechanic.ValidationError("Privilege ID is required")
	}
	if err := u.repo.DeleteUserPrivilege(id); err != nil {
		return mechanic.InternalError("Failed to delete user privilege", err)
	}
	return nil
}

func (u *useCase) ListUserLocation(userID string, meta mechanic.ActionMeta) ([]UserLocationListItem, mechanic.GridMeta, error) {
	page, size, err := mechanic.CheckMeta(meta.Page, meta.Size)
	if err != nil {
		return nil, mechanic.GridMeta{}, err
	}
	total, err := u.repo.CountUserLocation(userID, meta.Search)
	if err != nil {
		return nil, mechanic.GridMeta{}, mechanic.InternalError("Failed to count user locations", err)
	}
	list, err := u.repo.ListUserLocation(userID, meta.Search, page, size, meta.SortBy, meta.SortOrder)
	if err != nil {
		return nil, mechanic.GridMeta{}, mechanic.InternalError("Failed to list user locations", err)
	}
	return list, mechanic.BuildMeta(page, size, total), nil
}

func (u *useCase) CreateUserLocation(userID string, req UserLocationCreateRequest) error {
	if userID == "" {
		return mechanic.ValidationError("User ID is required")
	}
	if req.LocationTypeID == "" {
		return mechanic.ValidationError("Location type is required")
	}
	if err := u.repo.CreateUserLocation(userID, req); err != nil {
		return mechanic.InternalError("Failed to create user location", err)
	}
	return nil
}

func (u *useCase) UpdateUserLocation(id string, req UserLocationUpdateRequest) error {
	if id == "" {
		return mechanic.ValidationError("Location ID is required")
	}
	if err := u.repo.UpdateUserLocation(id, req); err != nil {
		return mechanic.InternalError("Failed to update user location", err)
	}
	return nil
}

func (u *useCase) DeleteUserLocation(id string) error {
	if id == "" {
		return mechanic.ValidationError("Location ID is required")
	}
	if err := u.repo.DeleteUserLocation(id); err != nil {
		return mechanic.InternalError("Failed to delete user location", err)
	}
	return nil
}

func (u *useCase) ListLocationTypeByCompany(companyID string) ([]LocationTypeSelectItem, error) {
	if companyID == "" {
		return nil, mechanic.ValidationError("Company ID is required")
	}
	list, err := u.repo.ListLocationTypeByCompany(companyID)
	if err != nil {
		return nil, mechanic.InternalError("Failed to list location types by company", err)
	}
	return list, nil
}
