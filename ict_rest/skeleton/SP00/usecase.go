package SP00

import (
	"crypto/rand"
	"encoding/hex"
	"ict_rest/mechanic"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type useCase struct {
	repo Repository
}

func NCase(r Repository) UseCase {
	return &useCase{repo: r}
}

func (u *useCase) ListHrisCompany() ([]HrisCompanyItem, error) {
	list, err := u.repo.ListHrisCompany()
	if err != nil {
		return nil, mechanic.InternalError("Failed to list HRIS companies", err)
	}
	return list, nil
}

func (u *useCase) GetCompanyHrisLink(companyID string) (string, error) {
	link, err := u.repo.GetCompanyHrisLink(companyID)
	if err != nil {
		return "", mechanic.InternalError("Failed to get company HRIS link", err)
	}
	return link, nil
}

func GenerateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (u *useCase) Login(username, password, ipAddress, userAgent string) (
	LoginResponse, error) {
	id, passwordHash, fullname, email, companyID, isAdmin, isHris, isActive,
		err := u.repo.FindUserByUsername(username)
	if err != nil {
		return LoginResponse{}, mechanic.Unauthorized("invalid credentials")
	}
	if strings.TrimSpace(companyID) != "" {
		return LoginResponse{}, mechanic.Unauthorized("this account requires HRIS login")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return LoginResponse{}, mechanic.Unauthorized("invalid credentials")
	}
	if !isActive {
		return LoginResponse{}, mechanic.Unauthorized("account is deactivated")
	}

	token, err := GenerateToken()
	if err != nil {
		return LoginResponse{}, mechanic.InternalError("Failed to generate token", err)
	}
	expiresAt := time.Now().Add(24 * time.Hour)
	if err := u.repo.CreateUserSession(id, token, ipAddress, userAgent, expiresAt); err != nil {
		return LoginResponse{}, mechanic.InternalError("Failed to create session", err)
	}

	profile := ProfileData{
		ID:        id,
		Username:  username,
		Email:     email,
		Fullname:  fullname,
		CompanyID: companyID,
		IsAdmin:   isAdmin,
		IsHris:    isHris,
		IsActive:  isActive,
	}
	return LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		Profile:   profile,
	}, nil
}

func (u *useCase) LoginHris(username, password, companyID, ipAddress, userAgent string) (
	LoginResponse, error) {
	id, fullname, email, isAdmin, isHris, isActive,
		err := u.repo.FindUserByUsernameAndCompany(username, companyID)
	if err != nil {
		return LoginResponse{}, mechanic.Unauthorized("HRIS user not found in system")
	}
	if !isActive {
		return LoginResponse{}, mechanic.Unauthorized("account is deactivated")
	}

	encryptedKey, err := mechanic.Encrypt(password)
	if err == nil {
		_ = u.repo.UpdateUserKey(id, encryptedKey)
	}

	token, err := GenerateToken()
	if err != nil {
		return LoginResponse{}, mechanic.InternalError("Failed to generate token", err)
	}
	expiresAt := time.Now().Add(24 * time.Hour)
	if err := u.repo.CreateUserSession(
		id,
		token,
		ipAddress,
		userAgent,
		expiresAt,
	); err != nil {
		return LoginResponse{}, mechanic.InternalError("Failed to create session", err)
	}

	profile := ProfileData{
		ID:        id,
		Username:  username,
		Email:     email,
		Fullname:  fullname,
		CompanyID: companyID,
		IsAdmin:   isAdmin,
		IsHris:    isHris,
		IsActive:  isActive,
	}
	return LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		Profile:   profile,
	}, nil
}

func (u *useCase) Logout(token string) error {
	if err := u.repo.DeleteSession(token); err != nil {
		return mechanic.InternalError("Failed to delete session", err)
	}
	return nil
}
