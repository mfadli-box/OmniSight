package SP03

import "ict_rest/mechanic"

type UserActionItem struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	CompanyID  string `json:"company_id"`
	ModuleCode string `json:"module_code"`
	Action     string `json:"action"`
	Path       string `json:"path"`
	IPAddress  string `json:"ip_address"`
	UserAgent  string `json:"user_agent"`
	CreatedAt  string `json:"created_at"`
}

type Repository interface {
	ListActions(userID, search string, page, size int, sortBy, sortOrder string) (
		[]UserActionItem, mechanic.GridMeta, error)
}

type UseCase interface {
	ListActions(userID, search string, page, size int, sortBy, sortOrder string) (
		[]UserActionItem, mechanic.GridMeta, error)
}
