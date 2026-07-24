package SP03

import (
	"context"
	"database/sql"
	"fmt"
	"ict_rest/mechanic"
)

type repository struct {
	db *sql.DB
}

func NRepo(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) ListActions(userID, search string, page, size int, sortBy, sortOrder string) (
	[]UserActionItem, mechanic.GridMeta, error) {
	where := `WHERE 1=1`
	args := []any{}
	argIdx := 1
	if userID != "" {
		where += fmt.Sprintf(` AND user_id = $%d`, argIdx)
		args = append(args, userID)
		argIdx++
	}
	if search != "" {
		where += fmt.Sprintf(`
			AND (module_code ILIKE $%d
			 OR  action      ILIKE $%d
			 OR  path        ILIKE $%d
			 OR  ip_address  ILIKE $%d)
		`, argIdx, argIdx, argIdx, argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}
	countQuery := `
		SELECT COUNT(*)
		FROM   "dat_user_action"
	` + where
	var total int
	if err := r.db.QueryRowContext(context.Background(),
		countQuery,
		args...,
	).Scan(&total); err != nil {
		return nil, mechanic.GridMeta{}, err
	}
	sortExpr := map[string]string{
		"created_at":  "created_at",
		"action":      "action",
		"module_code": "module_code",
		"ip_address":  "ip_address",
	}
	orderCol, ok := sortExpr[sortBy]
	if !ok {
		orderCol = "created_at"
	}
	if sortOrder != "asc" {
		sortOrder = "desc"
	}
	query := `
		SELECT COALESCE(id,         ''), COALESCE(user_id,     ''),
		       COALESCE(company_id, ''), COALESCE(module_code, ''),
		       COALESCE(action,     ''), COALESCE(path,        ''),
			   COALESCE(ip_address, ''), COALESCE(user_agent,  ''),
		       TO_CHAR(created_at, 'YYYY-MM-DD HH24:MI:SS')
		FROM   "dat_user_action"
	` + where + fmt.Sprintf(` ORDER BY %s %s`, orderCol, sortOrder)
	offset := (page - 1) * size
	query += fmt.Sprintf(` LIMIT $%d OFFSET $%d`, argIdx, argIdx+1)
	args = append(args, size, offset)
	rows, err := r.db.QueryContext(context.Background(), query, args...)
	if err != nil {
		return nil, mechanic.GridMeta{}, err
	}
	defer rows.Close()
	var list []UserActionItem
	for rows.Next() {
		var item UserActionItem
		if err := rows.Scan(
			&item.ID, &item.UserID, &item.CompanyID, &item.ModuleCode,
			&item.Action, &item.Path, &item.IPAddress, &item.UserAgent,
			&item.CreatedAt,
		); err != nil {
			return nil, mechanic.GridMeta{}, err
		}
		list = append(list, item)
	}
	if err := rows.Err(); err != nil {
		return nil, mechanic.GridMeta{}, err
	}
	meta := mechanic.GridMeta{
		Page:       page,
		Size:       size,
		Total:      total,
		TotalPages: (total + size - 1) / size,
	}
	return list, meta, nil
}
