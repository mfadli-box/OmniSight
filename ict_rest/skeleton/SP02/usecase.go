package SP02

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

func (u *useCase) ChangePassword(userID, currentPassword, newPassword string) error {
	if userID == "" {
		return mechanic.ValidationError("User ID is required")
	}
	if currentPassword == "" {
		return mechanic.ValidationError("Current password is required")
	}
	if newPassword == "" {
		return mechanic.ValidationError("New password is required")
	}
	hash, err := u.repo.FindUserPasswordHash(userID)
	if err != nil {
		return mechanic.NotFound("User not found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(currentPassword)); err != nil {
		return mechanic.ValidationError("Current password is incorrect")
	}
	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return mechanic.InternalError("Failed to hash password", err)
	}
	if err := u.repo.UpdateUserPassword(userID, string(newHash)); err != nil {
		return mechanic.InternalError("Failed to update user password", err)
	}
	return nil
}
