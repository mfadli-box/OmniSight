package SP02

type Repository interface {
	FindUserPasswordHash(id string) (string, error)
	UpdateUserPassword(id, hashedPassword string) error
}

type UseCase interface {
	ChangePassword(userID, currentPassword, newPassword string) error
}
