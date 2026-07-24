package SM05

import (
	"context"
	"database/sql"
	"fmt"
)

type repository struct {
	db *sql.DB
}

func NRepo(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CountSession(search string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM   "dat_user_session" s
		JOIN   "dat_user" u ON u.id = s.user_id
		WHERE  1=1
	`
	args := []any{}
	argIdx := 1
	if search != "" {
		query += fmt.Sprintf(" AND (u.username ILIKE $%d OR u.fullname ILIKE $%d OR COALESCE(s.ip_address,'') ILIKE $%d)", argIdx, argIdx, argIdx)
		args = append(args, "%"+search+"%")
	}
	var total int
	err := r.db.QueryRowContext(context.Background(), query, args...).Scan(&total)
	return total, err
}

func (r *repository) ListSession(search string, page, size int, sortBy, sortOrder string) (
	[]SessionListItem, error) {
	query := `
		SELECT s.id, s.user_id, u.username, u.fullname,
		       COALESCE(
		           (SELECT c.name FROM "dat_user_company" uc
		            JOIN "dat_company" c ON c.id = uc.company_id
		            WHERE uc.user_id = u.id AND uc.is_active = true
		            ORDER BY uc.created_at ASC LIMIT 1), ''
		       ),
		       COALESCE(s.ip_address, ''),
		       COALESCE(s.user_agent, ''),
		       LEFT(s.token, 8) || '...',
		       TO_CHAR(s.created_at, 'YYYY-MM-DD HH24:MI:SS'),
		       TO_CHAR(s.expires_at, 'YYYY-MM-DD HH24:MI:SS'),
		       s.expires_at > NOW()
		FROM   "dat_user_session" s
		JOIN   "dat_user" u ON u.id = s.user_id
		WHERE  1=1
	`
	args := []any{}
	argIdx := 1
	if search != "" {
		query += fmt.Sprintf(" AND (u.username ILIKE $%d OR u.fullname ILIKE $%d OR COALESCE(s.ip_address,'') ILIKE $%d)", argIdx, argIdx, argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}
	sortExpr := map[string]string{
		"username":   "u.username",
		"created_at": "s.created_at",
		"expires_at": "s.expires_at",
		"ip_address": "s.ip_address",
	}
	expr, ok := sortExpr[sortBy]
	if !ok {
		expr = "s.created_at"
	}
	if sortOrder != "asc" {
		sortOrder = "desc"
	}
	query += fmt.Sprintf(" ORDER BY %s %s", expr, sortOrder)
	offset := (page - 1) * size
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, size, offset)

	rows, err := r.db.QueryContext(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []SessionListItem
	for rows.Next() {
		var item SessionListItem
		if err := rows.Scan(
			&item.ID, &item.UserID, &item.Username, &item.Fullname,
			&item.CompanyName, &item.IPAddress, &item.UserAgent,
			&item.TokenPreview, &item.CreatedAt, &item.ExpiresAt, &item.IsActive,
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

func (r *repository) GetSessionDetail(id string) (*SessionDetailItem, error) {
	query := `
		SELECT s.id, s.user_id, u.username, u.fullname, u.email,
		       COALESCE(u.phone, ''), u.role, u.is_admin, u.is_hris, u.is_active,
		       COALESCE(u.company_id, ''),
		       COALESCE(
		           (SELECT c.name FROM "dat_user_company" uc
		            JOIN "dat_company" c ON c.id = uc.company_id
		            WHERE uc.user_id = u.id AND uc.is_active = true
		            ORDER BY uc.created_at ASC LIMIT 1), ''
		       ),
		       s.token,
		       COALESCE(s.ip_address, ''),
		       COALESCE(s.user_agent, ''),
		       TO_CHAR(s.created_at, 'YYYY-MM-DD HH24:MI:SS'),
		       TO_CHAR(s.expires_at, 'YYYY-MM-DD HH24:MI:SS'),
		       s.expires_at > NOW()
		FROM   "dat_user_session" s
		JOIN   "dat_user" u ON u.id = s.user_id
		WHERE  s.id = $1
	`
	item := &SessionDetailItem{}
	err := r.db.QueryRowContext(context.Background(), query, id).Scan(
		&item.ID, &item.UserID, &item.Username, &item.Fullname, &item.Email,
		&item.Phone, &item.Role, &item.IsAdmin, &item.IsHris, &item.IsActive,
		&item.CompanyID, &item.CompanyName,
		&item.Token, &item.IPAddress, &item.UserAgent,
		&item.CreatedAt, &item.ExpiresAt, &item.SessionActive,
	)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (r *repository) RevokeSession(id string) (int64, error) {
	result, err := r.db.ExecContext(context.Background(),
		`DELETE FROM "dat_user_session" WHERE id = $1`, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
