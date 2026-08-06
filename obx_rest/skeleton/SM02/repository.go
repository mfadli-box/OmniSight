package SM02

import (
	"context"
	"database/sql"
	"fmt"
	"obx_rest/mechanic"
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

func (r *repository) ListModuleSelect() ([]ModuleSelectItem, error) {
	query := `
		SELECT id, COALESCE(parent_id,''), code, name, is_active
		FROM   "dat_module"
		ORDER BY code ASC
	`
	rows, err := r.db.QueryContext(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []ModuleSelectItem
	for rows.Next() {
		var item ModuleSelectItem
		if err := rows.Scan(&item.ID, &item.ParentID, &item.Code, &item.Name, &item.IsActive); err != nil {
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
		RETURNING id
	`
	var moduleID string
	err := r.db.QueryRowContext(context.Background(), query,
		mechanic.NullableString(req.ParentID), req.Code, req.Name, req.Path,
		req.IsPage, req.IsActive,
	).Scan(&moduleID)
	if err != nil {
		return err
	}
	// Auto-assign the new module to all companies as inactive
	assignQuery := `
		INSERT INTO "dat_company_module" (
			id, company_id, module_id, is_active, created_at, updated_at
		)
		SELECT gen_random_uuid()::text, c.id, $1, false, NOW(), NOW()
		FROM   "dat_company" c
		ON CONFLICT (company_id, module_id) DO NOTHING
	`
	_, err = r.db.ExecContext(context.Background(), assignQuery, moduleID)
	if err != nil {
		return err
	}
	queryAuto := `
		INSERT INTO "dat_company_module" (
			id, company_id, module_id, is_active, created_at, updated_at
		)
		SELECT gen_random_uuid()::text, c.id, m.id, false, NOW(), NOW() FROM (
			SELECT 1 as x, id FROM "dat_module" WHERE parent_id IS NOT NULL
		) m LEFT OUTER JOIN (
			SELECT 1 as x, id FROM "dat_company"
		) c ON m.x = c.x LEFT OUTER JOIN (
			SELECT id, company_id, module_id FROM "dat_company_module"
		) cm ON m.id = cm.module_id AND c.id = cm.company_id
		WHERE cm.id IS NULL;
		INSERT INTO "dat_user_privilege" (
			id, user_company_id, module_id, level, created_at, updated_at
		)
		SELECT gen_random_uuid()::text, uc.id, cm.module_id, 'HIDE'::action_type, NOW(), NOW()
		FROM (
			SELECT id, company_id FROM "dat_user_company"
		) uc LEFT OUTER JOIN (
			SELECT company_id, module_id FROM "dat_company_module"
		) cm ON uc.company_id = cm.company_id LEFT OUTER JOIN (
			SELECT id, user_company_id, module_id FROM "dat_user_privilege"
		) up ON up.user_company_id = uc.id AND up.module_id = cm.module_id
		WHERE up.id IS NULL;
	`
	_, errAuto := r.db.ExecContext(context.Background(), queryAuto)
	return errAuto
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
