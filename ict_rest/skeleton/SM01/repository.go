package SM01

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

func (r *repository) CountUser(companyID, search string) (int, error) {
	query := `
		SELECT COUNT(*) FROM "dat_user" WHERE 1=1
	`
	args := []any{}
	argIdx := 1
	if companyID != "" {
		query += fmt.Sprintf(" AND company_id = $%d", argIdx)
		args = append(args, companyID)
		argIdx++
	}
	if search != "" {
		query += fmt.Sprintf(" AND (username ILIKE $%d OR fullname ILIKE $%d OR email ILIKE $%d)", argIdx, argIdx, argIdx)
		args = append(args, "%"+search+"%")
	}
	var total int
	err := r.db.QueryRowContext(context.Background(), query, args...).Scan(&total)
	return total, err
}

func (r *repository) ListUser(companyID, search string, page, size int, sortBy, sortOrder string) (
	[]UserListItem, error) {
	query := `
		SELECT id, username, email, fullname, COALESCE(phone,''), company_id,
		       role, job, is_admin, is_hris, is_active,
		       TO_CHAR(created_at, 'YYYY-MM-DD HH24:MI:SS')
		FROM   "dat_user"
		WHERE  1=1
	`
	args := []any{}
	argIdx := 1
	if companyID != "" {
		query += fmt.Sprintf(" AND company_id = $%d", argIdx)
		args = append(args, companyID)
		argIdx++
	}
	if search != "" {
		query += fmt.Sprintf(" AND (username ILIKE $%d OR fullname ILIKE $%d OR email ILIKE $%d)", argIdx, argIdx, argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}
	validSort := map[string]bool{"username": true, "fullname": true, "email": true, "created_at": true}
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

	var list []UserListItem
	for rows.Next() {
		var item UserListItem
		if err := rows.Scan(
			&item.ID, &item.Username, &item.Email, &item.Fullname, &item.Phone,
			&item.CompanyID, &item.Role, &item.Job, &item.IsAdmin, &item.IsHris,
			&item.IsActive, &item.CreatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, nil
}

func (r *repository) UpdateUserPassword(id, hashedPassword string) error {
	query := `UPDATE "dat_user" SET password = $2, updated_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(context.Background(), query, id, hashedPassword)
	return err
}

func (r *repository) CreateUser(req UserCreateRequest, hashedPassword string) (string, error) {
	var exists bool
	err := r.db.QueryRowContext(context.Background(),
		`SELECT EXISTS(SELECT 1 FROM "dat_user" WHERE username = $1 AND company_id = $2)`,
		req.Username, req.CompanyID,
	).Scan(&exists)
	if err == nil && exists {
		return "", mechanic.Conflict("Username already exists in this company")
	}
	query := `
		INSERT INTO "dat_user" (
			id, username, email, password, fullname, phone, company_id, role,
			job, is_admin, is_hris, is_active, created_at, updated_at
		) VALUES (gen_random_uuid()::text,
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW()
		)
		RETURNING id
	`
	var id string
	err = r.db.QueryRowContext(context.Background(), query,
		req.Username, req.Email, hashedPassword, req.Fullname,
		mechanic.NullableString(req.Phone), req.CompanyID,
		req.Role, req.Job,
		req.IsAdmin, req.IsHris, req.IsActive,
	).Scan(&id)
	return id, err
}

func (r *repository) UpdateUser(id string, req UserUpdateRequest) error {
	query := `
		UPDATE "dat_user" SET
			email = $2, fullname = $3, phone = $4, role = $5, job = $6,
			is_admin = $7, is_hris = $8, is_active = $9, updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.ExecContext(context.Background(), query,
		id, req.Email, req.Fullname, mechanic.NullableString(req.Phone),
		req.Role, req.Job,
		req.IsAdmin, req.IsHris, req.IsActive,
	)
	return err
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
	if rows.Err() != nil {
		return nil, rows.Err()
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
	return list, nil
}

func (r *repository) ListAllCompany() ([]CompanySelectItem, error) {
	query := `SELECT id, name FROM "dat_company" WHERE is_active = true ORDER BY name ASC`
	rows, err := r.db.QueryContext(context.Background(), query)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	defer rows.Close()
	var list []CompanySelectItem
	for rows.Next() {
		var item CompanySelectItem
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, nil
}

func (r *repository) ListAllModule() ([]ModuleSelectItem, error) {
	query := `
		SELECT id, COALESCE(parent_id,''), code, name, path, is_page, is_active
		FROM   "dat_module"
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
	var list []ModuleSelectItem
	for rows.Next() {
		var item ModuleSelectItem
		if err := rows.Scan(
			&item.ID,
			&item.ParentID,
			&item.Code,
			&item.Name,
			&item.Path,
			&item.IsPage,
			&item.IsActive,
		); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, nil
}

func (r *repository) AssignCompany(userID, companyID string) error {
	query := `
		INSERT INTO "dat_user_company" (id, user_id, company_id, is_active, created_at, updated_at)
		VALUES (gen_random_uuid()::text, $1, $2, true, NOW(), NOW())
		ON CONFLICT (user_id, company_id) DO UPDATE SET is_active = true, updated_at = NOW()
	`
	_, err := r.db.ExecContext(context.Background(), query, userID, companyID)
	return err
}

func (r *repository) CountUserCompany(userID, search string) (int, error) {
	query := `SELECT COUNT(*) FROM "dat_user_company" uc JOIN "dat_company" c ON c.id = uc.company_id WHERE uc.user_id = $1`
	args := []any{userID}
	argIdx := 2
	if search != "" {
		query += fmt.Sprintf(" AND c.name ILIKE $%d", argIdx)
		args = append(args, "%"+search+"%")
	}
	var total int
	err := r.db.QueryRowContext(context.Background(), query, args...).Scan(&total)
	return total, err
}

func (r *repository) ListUserCompany(userID, search string, page, size int, sortBy, sortOrder string) ([]UserCompanyListItem, error) {
	query := `
		SELECT uc.id, uc.user_id, uc.company_id, c.name AS company_name, uc.is_active,
			TO_CHAR(uc.created_at, 'YYYY-MM-DD HH24:MI:SS')
		FROM "dat_user_company" uc
		JOIN "dat_company" c ON c.id = uc.company_id
		WHERE uc.user_id = $1
	`
	args := []any{userID}
	argIdx := 2
	if search != "" {
		query += fmt.Sprintf(" AND c.name ILIKE $%d", argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}
	sortExpr := map[string]string{
		"company_name": "company_name",
		"created_at":   "uc.created_at",
	}
	expr, ok := sortExpr[sortBy]
	if !ok {
		expr = "uc.created_at"
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
	var list []UserCompanyListItem
	for rows.Next() {
		var item UserCompanyListItem
		if err := rows.Scan(
			&item.ID, &item.UserID, &item.CompanyID, &item.CompanyName,
			&item.IsActive, &item.CreatedAt,
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

func (r *repository) CreateUserCompany(userID string, req UserCompanyRequest) error {
	query := `
		INSERT INTO "dat_user_company" (id, user_id, company_id, is_active, created_at, updated_at)
		VALUES (gen_random_uuid()::text, $1, $2, $3, NOW(), NOW())
		ON CONFLICT (user_id, company_id) DO UPDATE SET is_active = $3, updated_at = NOW()
	`
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	_, err := r.db.ExecContext(context.Background(), query, userID, req.CompanyID, isActive)
	return err
}

func (r *repository) UpdateUserCompany(userID, companyID string, req UserCompanyRequest) error {
	query := `UPDATE "dat_user_company" SET is_active = $3, updated_at = NOW() WHERE user_id = $1 AND company_id = $2`
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	_, err := r.db.ExecContext(context.Background(), query, userID, companyID, isActive)
	return err
}

func (r *repository) DeleteUserCompany(userID, companyID string) error {
	query := `DELETE FROM "dat_user_company" WHERE user_id = $1 AND company_id = $2`
	_, err := r.db.ExecContext(context.Background(), query, userID, companyID)
	return err
}

func (r *repository) CountUserPrivilege(userID, search string) (int, error) {
	query := `
		SELECT COUNT(*) FROM "dat_user_privilege" up
		JOIN "dat_user_company" uc ON uc.id = up.user_company_id
		JOIN "dat_company" c ON c.id = uc.company_id
		JOIN "dat_module" m ON m.id = up.module_id
		WHERE uc.user_id = $1
	`
	args := []any{userID}
	argIdx := 2
	if search != "" {
		query += fmt.Sprintf(" AND (c.name ILIKE $%d OR m.code ILIKE $%d OR m.name ILIKE $%d)", argIdx, argIdx, argIdx)
		args = append(args, "%"+search+"%")
	}
	var total int
	err := r.db.QueryRowContext(context.Background(), query, args...).Scan(&total)
	return total, err
}

func (r *repository) ListUserPrivilege(userID, search string, page, size int, sortBy, sortOrder string) ([]UserPrivilegeListItem, error) {
	query := `
		SELECT up.id, up.user_company_id, up.module_id, c.name AS company_name,
			m.code AS module_code, m.name AS module_name, up.level::text,
			TO_CHAR(up.created_at, 'YYYY-MM-DD HH24:MI:SS')
		FROM "dat_user_privilege" up
		JOIN "dat_user_company" uc ON uc.id = up.user_company_id
		JOIN "dat_company" c ON c.id = uc.company_id
		JOIN "dat_module" m ON m.id = up.module_id
		WHERE uc.user_id = $1
	`
	args := []any{userID}
	argIdx := 2
	if search != "" {
		query += fmt.Sprintf(" AND (c.name ILIKE $%d OR m.code ILIKE $%d OR m.name ILIKE $%d)", argIdx, argIdx, argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}
	sortExpr := map[string]string{
		"company_name": "company_name",
		"module_code":  "module_code",
		"module_name":  "module_name",
		"level":        "up.level",
		"created_at":   "up.created_at",
	}
	expr, ok := sortExpr[sortBy]
	if !ok {
		expr = "up.created_at"
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
	var list []UserPrivilegeListItem
	for rows.Next() {
		var item UserPrivilegeListItem
		if err := rows.Scan(
			&item.ID, &item.UserCompanyID, &item.ModuleID, &item.CompanyName,
			&item.ModuleCode, &item.ModuleName, &item.Level, &item.CreatedAt,
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

func (r *repository) CreateUserPrivilege(userID string, req UserPrivilegeRequest) error {
	level := "hide"
	if req.Level != "" {
		level = req.Level
	}
	query := `
		INSERT INTO "dat_user_privilege" (
			id, user_company_id, module_id, level, created_at, updated_at
		) VALUES (gen_random_uuid()::text,
			$1, $2, $3::action_type, NOW(), NOW()
		) ON CONFLICT (user_company_id, module_id)
		  DO UPDATE SET level = $3::action_type, updated_at = NOW()
	`
	_, err := r.db.ExecContext(context.Background(), query, req.UserCompanyID, req.ModuleID, level)
	return err
}

func (r *repository) UpdateUserPrivilege(id string, req UserPrivilegeUpdateRequest) error {
	query := `UPDATE "dat_user_privilege" SET level = $2::action_type, updated_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(context.Background(), query, id, req.Level)
	return err
}

func (r *repository) DeleteUserPrivilege(id string) error {
	query := `DELETE FROM "dat_user_privilege" WHERE id = $1`
	_, err := r.db.ExecContext(context.Background(), query, id)
	return err
}

func (r *repository) CountUserLocation(userID, search string) (int, error) {
	query := `
		SELECT COUNT(*) FROM "dat_user_location" ul
		JOIN "dat_location_type" lt ON lt.id = ul.location_type_id
		JOIN "dat_company" co ON co.id = lt.company_id
		WHERE ul.user_id = $1
	`
	args := []any{userID}
	argIdx := 2
	if search != "" {
		query += fmt.Sprintf(" AND (co.name ILIKE $%d OR lt.code ILIKE $%d OR lt.name ILIKE $%d)", argIdx, argIdx, argIdx)
		args = append(args, "%"+search+"%")
	}
	var total int
	err := r.db.QueryRowContext(context.Background(), query, args...).Scan(&total)
	return total, err
}

func (r *repository) ListUserLocation(userID, search string, page, size int, sortBy, sortOrder string) ([]UserLocationListItem, error) {
	query := `
		SELECT ul.id, ul.user_id, ul.location_type_id, co.name AS company_name,
			lt.code AS type_code, lt.name AS type_name, ul.is_active,
			TO_CHAR(ul.created_at, 'YYYY-MM-DD HH24:MI:SS')
		FROM "dat_user_location" ul
		JOIN "dat_location_type" lt ON lt.id = ul.location_type_id
		JOIN "dat_company" co ON co.id = lt.company_id
		WHERE ul.user_id = $1
	`
	args := []any{userID}
	argIdx := 2
	if search != "" {
		query += fmt.Sprintf(" AND (co.name ILIKE $%d OR lt.code ILIKE $%d OR lt.name ILIKE $%d)", argIdx, argIdx, argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}
	sortExpr := map[string]string{
		"company_name": "company_name",
		"type_code":    "type_code",
		"type_name":    "type_name",
		"created_at":   "ul.created_at",
	}
	expr, ok := sortExpr[sortBy]
	if !ok {
		expr = "ul.created_at"
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
	var list []UserLocationListItem
	for rows.Next() {
		var item UserLocationListItem
		if err := rows.Scan(
			&item.ID, &item.UserID, &item.LocationTypeID, &item.CompanyName,
			&item.TypeCode, &item.TypeName, &item.IsActive, &item.CreatedAt,
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

func (r *repository) CreateUserLocation(userID string, req UserLocationCreateRequest) error {
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	query := `
		INSERT INTO "dat_user_location" (id, user_id, location_type_id, is_active, created_at, updated_at)
		VALUES (gen_random_uuid()::text, $1, $2, $3, NOW(), NOW())
		ON CONFLICT (user_id, location_type_id) DO UPDATE SET is_active = $3, updated_at = NOW()
	`
	_, err := r.db.ExecContext(context.Background(), query, userID, req.LocationTypeID, isActive)
	return err
}

func (r *repository) UpdateUserLocation(id string, req UserLocationUpdateRequest) error {
	query := `UPDATE "dat_user_location" SET is_active = $2, updated_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(context.Background(), query, id, req.IsActive)
	return err
}

func (r *repository) DeleteUserLocation(id string) error {
	query := `DELETE FROM "dat_user_location" WHERE id = $1`
	_, err := r.db.ExecContext(context.Background(), query, id)
	return err
}

func (r *repository) ListLocationTypeByCompany(companyID string) ([]LocationTypeSelectItem, error) {
	query := `SELECT id, code, name FROM "dat_location_type" WHERE is_active = true AND company_id = $1 ORDER BY code ASC`
	rows, err := r.db.QueryContext(context.Background(), query, companyID)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	defer rows.Close()
	var list []LocationTypeSelectItem
	for rows.Next() {
		var item LocationTypeSelectItem
		if err := rows.Scan(&item.ID, &item.Code, &item.Name); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, nil
}
