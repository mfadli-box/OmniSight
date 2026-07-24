package SP00

import (
	"encoding/xml"
	"time"
)

type HrisCompanyItem struct {
	ID       string `json:"id"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	HrisLink string `json:"hris_link"`
}

type LoginRequest struct {
	CompanyID string `json:"company_id"`
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

type HrisSoapResponse struct {
	XMLName xml.Name `xml:"string"`
	Text    string   `xml:",chardata"`
}

type LoginResponse struct {
	Token     string      `json:"token"`
	ExpiresAt time.Time   `json:"expires_at"`
	Profile   ProfileData `json:"user_profile"`
}

type ProfileData struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Fullname    string `json:"fullname"`
	Phone       string `json:"phone"`
	CompanyID   string `json:"company_id"`
	CompanyName string `json:"company_name,omitempty"`
	Role        string `json:"role"`
	Job         string `json:"job"`
	IsAdmin     bool   `json:"is_admin"`
	IsHris      bool   `json:"is_hris"`
	IsActive    bool   `json:"is_active"`
}

type Repository interface {
	ListHrisCompany() ([]HrisCompanyItem, error)
	GetCompanyHrisLink(companyID string) (string, error)
	FindUserByUsername(username string) (
		id, passwordHash, fullname, email, companyID string, isAdmin, isHris, isActive bool, err error)
	FindUserByUsernameAndCompany(username, companyID string) (
		id, fullname, email string, isAdmin, isHris, isActive bool, err error)
	CreateUserSession(userID, token, ipAddress, userAgent string, expiresAt time.Time) error
	UpdateUserKey(id, key string) error
	DeleteSession(token string) error
}

type UseCase interface {
	ListHrisCompany() ([]HrisCompanyItem, error)
	GetCompanyHrisLink(companyID string) (string, error)
	Login(username, password, ipAddress, userAgent string) (LoginResponse, error)
	LoginHris(username, password, companyID, ipAddress, userAgent string) (LoginResponse, error)
	Logout(token string) error
}
