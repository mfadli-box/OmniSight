package SP00

import (
	"context"
	"database/sql"
	"ict_rest/mechanic"
	"time"
)

type repository struct {
	db *sql.DB
}

func NRepo(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) ListHrisCompany() ([]HrisCompanyItem, error) {
	query := `
		SELECT id, code, name, COALESCE(hris_link,'')
		FROM   "dat_company"
		WHERE  is_active = true AND (hris_link IS NOT NULL AND hris_link != '')
		ORDER BY name ASC
	`
	rows, err := r.db.QueryContext(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []HrisCompanyItem
	for rows.Next() {
		var item HrisCompanyItem
		if err := rows.Scan(
			&item.ID,
			&item.Code,
			&item.Name,
			&item.HrisLink,
		); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) GetCompanyHrisLink(companyID string) (string, error) {
	query := `
		SELECT COALESCE(hris_link,'')
		FROM   "dat_company"
		WHERE  id = $1 AND is_active = true
	`
	var link string
	err := r.db.QueryRowContext(context.Background(), query,
		companyID,
	).Scan(
		&link,
	)
	return link, err
}

func (r *repository) FindUserByUsername(username string) (
	id, passwordHash, fullname, email, companyID string, isAdmin, isHris, isActive bool, err error) {
	query := `
		SELECT id, password, fullname, email, company_id,
		       is_admin, is_hris, is_active
		FROM   "dat_user"
		WHERE  username = $1 AND is_active = true
	`
	err = r.db.QueryRowContext(context.Background(), query,
		username,
	).Scan(
		&id,
		&passwordHash,
		&fullname,
		&email,
		&companyID,
		&isAdmin,
		&isHris,
		&isActive,
	)
	return
}

func (r *repository) FindUserByUsernameAndCompany(username, companyID string) (
	id, fullname, email string, isAdmin, isHris, isActive bool, err error) {
	query := `
		SELECT id, fullname, email, is_admin, is_hris, is_active
		FROM   "dat_user"
		WHERE  username = $1 AND company_id = $2
		  AND  is_hris = true AND is_active = true
	`
	err = r.db.QueryRowContext(context.Background(), query,
		username,
		companyID,
	).Scan(
		&id,
		&fullname,
		&email,
		&isAdmin,
		&isHris,
		&isActive,
	)
	return
}

func (r *repository) CreateUserSession(userID, token, ipAddress, userAgent string, expiresAt time.Time) error {
	query := `
		INSERT INTO "dat_user_session" (
			id, user_id, token, ip_address, user_agent, created_at, expires_at
		) VALUES (
		 	gen_random_uuid()::text, $1, $2, $3, $4, NOW(), $5
		)
	`
	_, err := r.db.ExecContext(context.Background(), query,
		userID,
		token,
		mechanic.NullableString(ipAddress),
		mechanic.NullableString(userAgent),
		expiresAt,
	)
	return err
}

func (r *repository) UpdateUserKey(id, key string) error {
	query := `
		UPDATE "dat_user"
		SET    key = $2, updated_at = NOW()
		WHERE  id = $1`
	_, err := r.db.ExecContext(context.Background(), query,
		id,
		key,
	)
	return err
}

func (r *repository) DeleteSession(token string) error {
	query := `
		DELETE FROM "dat_user_session"
		WHERE  token = $1
	`
	_, err := r.db.ExecContext(context.Background(), query, token)
	return err
}
