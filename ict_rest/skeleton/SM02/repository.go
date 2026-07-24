package SM02

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

func (r *repository) CountModule(search string) (int, error) {
	query := `SELECT COUNT(*) FROM "dat_module" WHERE 1=1`
	args := []any{}
	argIdx := 1
	if search != "" {
		query += fmt.Sprintf(" AND (code ILIKE $%d OR name ILIKE $%d OR path ILIKE $%d)",
			argIdx, argIdx, argIdx)
		args = append(args, "%"+search+"%")
	}
	var total int
	err := r.db.QueryRowContext(context.Background(), query, args...).Scan(&total)
	return total, err
}

func (r *repository) ListModule(search string, page, size int, sortBy, sortOrder string) (
	[]ModuleListItem, error) {
	query := `
		SELECT id, COALESCE(parent_id,''), code, name, path, is_page, is_active,
		       TO_CHAR(created_at, 'YYYY-MM-DD HH24:MI:SS')
		FROM   "dat_module" WHERE 1=1
	`
	args := []any{}
	argIdx := 1
	if search != "" {
		query += fmt.Sprintf(" AND (code ILIKE $%d OR name ILIKE $%d OR path ILIKE $%d)",
			argIdx, argIdx, argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}
	sortExpr := map[string]string{
		"code":       "code",
		"name":       "name",
		"path":       "path",
		"created_at": "created_at",
	}
	orderCol, ok := sortExpr[sortBy]
	if !ok {
		orderCol = "created_at"
	}
	if sortOrder != "asc" {
		sortOrder = "desc"
	}
	query += fmt.Sprintf(" ORDER BY %s %s", orderCol, sortOrder)
	offset := (page - 1) * size
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, size, offset)

	rows, err := r.db.QueryContext(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []ModuleListItem
	for rows.Next() {
		var item ModuleListItem
		if err := rows.Scan(
			&item.ID, &item.ParentID, &item.Code, &item.Name, &item.Path,
			&item.IsPage, &item.IsActive, &item.CreatedAt,
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

func (r *repository) CreateModule(req ModuleCreateRequest) error {
	query := `
		INSERT INTO "dat_module" (id, parent_id, code, name, path, is_page, is_active, created_at, updated_at)
		VALUES (gen_random_uuid()::text, $1, $2, $3, $4, $5, $6, NOW(), NOW())
	`
	_, err := r.db.ExecContext(context.Background(), query,
		mechanic.NullableString(req.ParentID), req.Code, req.Name, req.Path,
		req.IsPage, req.IsActive,
	)
	return err
}

func (r *repository) UpdateModule(id string, req ModuleUpdateRequest) error {
	query := `
		UPDATE "dat_module" SET
			parent_id = $2, code = $3, name = $4, path = $5,
			is_page = $6, is_active = $7, updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.ExecContext(context.Background(), query,
		id, mechanic.NullableString(req.ParentID), req.Code, req.Name, req.Path,
		req.IsPage, req.IsActive,
	)
	return err
}
