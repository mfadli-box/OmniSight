package SP02

import (
	"context"
	"database/sql"
)

type repository struct {
	db *sql.DB
}

func NRepo(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindUserPasswordHash(id string) (string, error) {
	var hash string
	err := r.db.QueryRowContext(context.Background(), `
		SELECT password
		FROM   "dat_user"
		WHERE  id = $1
	`, id).Scan(&hash)
	return hash, err
}

func (r *repository) UpdateUserPassword(id, hashedPassword string) error {
	query := `
		UPDATE "dat_user"
		SET    password = $2, updated_at = NOW()
		WHERE  id = $1
	`
	_, err := r.db.ExecContext(context.Background(), query, id, hashedPassword)
	return err
}
