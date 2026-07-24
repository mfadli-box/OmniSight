package SM05

import "ict_rest/mechanic"

type SessionListItem struct {
	ID           string `json:"id"`
	UserID       string `json:"user_id"`
	Username     string `json:"username"`
	Fullname     string `json:"fullname"`
	CompanyName  string `json:"company_name"`
	IPAddress    string `json:"ip_address"`
	UserAgent    string `json:"user_agent"`
	TokenPreview string `json:"token_preview"`
	CreatedAt    string `json:"created_at"`
	ExpiresAt    string `json:"expires_at"`
	IsActive     bool   `json:"is_active"`
}

type SessionDetailItem struct {
	ID            string `json:"id"`
	UserID        string `json:"user_id"`
	Username      string `json:"username"`
	Fullname      string `json:"fullname"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	Role          string `json:"role"`
	IsAdmin       bool   `json:"is_admin"`
	IsHris        bool   `json:"is_hris"`
	IsActive      bool   `json:"is_active"`
	CompanyID     string `json:"company_id"`
	CompanyName   string `json:"company_name"`
	Token         string `json:"token"`
	IPAddress     string `json:"ip_address"`
	UserAgent     string `json:"user_agent"`
	CreatedAt     string `json:"created_at"`
	ExpiresAt     string `json:"expires_at"`
	SessionActive bool   `json:"session_active"`
}

type Repository interface {
	CountSession(search string) (int, error)
	ListSession(search string, page, size int, sortBy, sortOrder string) ([]SessionListItem, error)
	GetSessionDetail(id string) (*SessionDetailItem, error)
	RevokeSession(id string) (int64, error)
}

type UseCase interface {
	ListSession(meta mechanic.ActionMeta) ([]SessionListItem, mechanic.GridMeta, error)
	GetSessionDetail(id string) (*SessionDetailItem, error)
	RevokeSession(id string) error
}
