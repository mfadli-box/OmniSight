package SM03

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

func (r *repository) CountCompany(search string) (int, error) {
	query := `SELECT COUNT(*) FROM "dat_company" WHERE 1=1`
	args := []any{}
	argIdx := 1
	if search != "" {
		query += fmt.Sprintf(" AND (code ILIKE $%d OR name ILIKE $%d)", argIdx, argIdx)
		args = append(args, "%"+search+"%")
	}
	var total int
	err := r.db.QueryRowContext(context.Background(), query, args...).Scan(&total)
	return total, err
}

func (r *repository) ListCompany(search string, page, size int, sortBy, sortOrder string) (
	[]CompanyListItem, error) {
	query := `
		SELECT id, code, name, COALESCE(vat_id,''), COALESCE(reg_no,''), COALESCE(address,''),
		       valuta, COALESCE(hris_link,''), is_active,
		       TO_CHAR(created_at, 'YYYY-MM-DD HH24:MI:SS')
		FROM   "dat_company" WHERE 1=1
	`
	args := []any{}
	argIdx := 1
	if search != "" {
		query += fmt.Sprintf(" AND (code ILIKE $%d OR name ILIKE $%d)", argIdx, argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}
	validSort := map[string]bool{"code": true, "name": true, "created_at": true}
	if !validSort[sortBy] {
		sortBy = "created_at"
	}
	if sortOrder != "asc" {
		sortOrder = "desc"
	}
	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)
	offset := (page - 1) * size
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, size, offset)

	rows, err := r.db.QueryContext(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	defer rows.Close()

	var list []CompanyListItem
	for rows.Next() {
		var item CompanyListItem
		if err := rows.Scan(
			&item.ID, &item.Code, &item.Name, &item.VatID, &item.RegNo,
			&item.Address, &item.Valuta, &item.HrisLink, &item.IsActive,
			&item.CreatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, nil
}

func (r *repository) ListModule() ([]ModuleItem, error) {
	query := `
		SELECT id, code, name, parent_id, is_active
		FROM   "dat_module"
		WHERE  is_active = true
		  AND  parent_id IS NOT NULL
		ORDER BY code ASC
	`
	rows, err := r.db.QueryContext(context.Background(), query)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	defer rows.Close()

	var list []ModuleItem
	for rows.Next() {
		var item ModuleItem
		if err := rows.Scan(
			&item.ID, &item.Code, &item.Name, &item.ParentID, &item.IsActive,
		); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, nil
}

func (r *repository) CreateCompany(req CompanyCreateRequest) error {
	query := `
		INSERT INTO "dat_company" (
			id, code, name, vat_id, reg_no, address, valuta, hris_link,
			is_active, created_at, updated_at
		) VALUES (gen_random_uuid()::text,
			$1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW()
		)
	`
	_, err := r.db.ExecContext(context.Background(), query,
		req.Code, req.Name,
		mechanic.NullableString(req.VatID), mechanic.NullableString(req.RegNo),
		mechanic.NullableString(req.Address), mechanic.NullableString(req.Valuta),
		mechanic.NullableString(req.HrisLink), req.IsActive,
	)
	return err
}

func (r *repository) UpdateCompany(id string, req CompanyUpdateRequest) error {
	query := `
		UPDATE "dat_company"
		SET    name = $2, vat_id = $3, reg_no = $4, address = $5, valuta = $6,
		       hris_link = $7, is_active = $8, updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.ExecContext(context.Background(), query,
		id, req.Name,
		mechanic.NullableString(req.VatID), mechanic.NullableString(req.RegNo),
		mechanic.NullableString(req.Address), mechanic.NullableString(req.Valuta),
		mechanic.NullableString(req.HrisLink), req.IsActive,
	)
	return err
}

func (r *repository) CountCompanyModule(companyID, search string) (int, error) {
	query := `SELECT COUNT(*) FROM "dat_company_module" cm JOIN "dat_module" m ON cm.module_id = m.id WHERE cm.company_id = $1`
	args := []any{companyID}
	argIdx := 2
	if search != "" {
		query += fmt.Sprintf(" AND (m.code ILIKE $%d OR m.name ILIKE $%d)", argIdx, argIdx)
		args = append(args, "%"+search+"%")
	}
	var total int
	err := r.db.QueryRowContext(context.Background(), query, args...).Scan(&total)
	return total, err
}

func (r *repository) ListCompanyModule(companyID, search string, page, size int, sortBy, sortOrder string) ([]CompanyModuleItem, error) {
	query := `
		SELECT cm.id, cm.module_id, m.code, m.name, m.path, cm.is_active
		FROM   "dat_company_module" cm
		JOIN   "dat_module" m ON cm.module_id = m.id
		WHERE  cm.company_id = $1
	`
	args := []any{companyID}
	argIdx := 2
	if search != "" {
		query += fmt.Sprintf(" AND (m.code ILIKE $%d OR m.name ILIKE $%d)", argIdx, argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}
	validSort := map[string]bool{"code": true, "name": true, "created_at": true}
	if !validSort[sortBy] {
		sortBy = "code"
	}
	if sortOrder != "asc" {
		sortOrder = "desc"
	}
	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)
	offset := (page - 1) * size
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, size, offset)

	rows, err := r.db.QueryContext(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	defer rows.Close()

	var list []CompanyModuleItem
	for rows.Next() {
		var item CompanyModuleItem
		if err := rows.Scan(
			&item.ID, &item.ModuleID, &item.Code, &item.Name, &item.Path, &item.IsActive,
		); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, nil
}

func (r *repository) AssignCompanyModule(companyID, moduleID string) error {
	query := `
		INSERT INTO "dat_company_module" (
			id, company_id, module_id, is_active, created_at, updated_at
		) VALUES (gen_random_uuid()::text,
			$1, $2, true, NOW(), NOW()
		) ON CONFLICT (company_id, module_id)
		  DO UPDATE SET is_active = true, updated_at = NOW()
	`
	_, err := r.db.ExecContext(context.Background(), query, companyID, moduleID)
	return err
}

func (r *repository) UpdateCompanyModule(companyID, moduleID string, isActive bool) error {
	query := `
		UPDATE "dat_company_module"
		SET    is_active = $3, updated_at = NOW()
		WHERE  company_id = $1 AND module_id = $2
	`
	_, err := r.db.ExecContext(context.Background(), query, companyID, moduleID, isActive)
	return err
}

func (r *repository) CountLocationType(companyID, search string) (int, error) {
	query := `SELECT COUNT(*) FROM "dat_location_type" WHERE company_id = $1`
	args := []any{companyID}
	argIdx := 2
	if search != "" {
		query += fmt.Sprintf(" AND (code ILIKE $%d OR name ILIKE $%d)", argIdx, argIdx)
		args = append(args, "%"+search+"%")
	}
	var total int
	err := r.db.QueryRowContext(context.Background(), query, args...).Scan(&total)
	return total, err
}

func (r *repository) ListLocationType(companyID, search string, page, size int, sortBy, sortOrder string) ([]LocationTypeItem, error) {
	query := `
		SELECT id, code, name, COALESCE(description,''), is_active,
		       TO_CHAR(created_at, 'YYYY-MM-DD HH24:MI:SS')
		FROM   "dat_location_type"
		WHERE  company_id = $1
	`
	args := []any{companyID}
	argIdx := 2
	if search != "" {
		query += fmt.Sprintf(" AND (code ILIKE $%d OR name ILIKE $%d)", argIdx, argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}
	validSort := map[string]bool{"code": true, "name": true, "created_at": true}
	if !validSort[sortBy] {
		sortBy = "code"
	}
	if sortOrder != "asc" {
		sortOrder = "desc"
	}
	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)
	offset := (page - 1) * size
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, size, offset)

	rows, err := r.db.QueryContext(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	defer rows.Close()
	var list []LocationTypeItem
	for rows.Next() {
		var item LocationTypeItem
		if err := rows.Scan(
			&item.ID, &item.Code, &item.Name, &item.Description, &item.IsActive, &item.CreatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, nil
}

func (r *repository) CreateLocationType(companyID string, req LocationTypeCreateRequest) error {
	query := `
		INSERT INTO "dat_location_type" (
			id, company_id, code, name, description, is_active, created_at, updated_at
		) VALUES (gen_random_uuid()::text,
			$1, $2, $3, $4, $5, NOW(), NOW()
		)
	`
	_, err := r.db.ExecContext(context.Background(), query,
		companyID, req.Code, req.Name,
		mechanic.NullableString(req.Description), req.IsActive,
	)
	return err
}

func (r *repository) UpdateLocationType(id string, req LocationTypeUpdateRequest) error {
	query := `
		UPDATE "dat_location_type"
		SET    code = $2, name = $3, description = $4, is_active = $5, updated_at = NOW()
		WHERE  id = $1
	`
	_, err := r.db.ExecContext(context.Background(), query,
		id, req.Code, req.Name,
		mechanic.NullableString(req.Description), req.IsActive,
	)
	return err
}

func (r *repository) DeleteLocationType(id string) error {
	query := `DELETE FROM "dat_location_type" WHERE id = $1`
	_, err := r.db.ExecContext(context.Background(), query, id)
	return err
}
