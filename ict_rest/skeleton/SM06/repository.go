package SM06

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

func (r *repository) CountLocationType(companyID, search string) (int, error) {
	query := `SELECT COUNT(*) FROM "dat_location_type" WHERE 1=1`
	args := []any{}
	argIdx := 1
	if companyID != "" {
		query += fmt.Sprintf(" AND company_id = $%d", argIdx)
		args = append(args, companyID)
		argIdx++
	}
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
		SELECT id, code, name, COALESCE(description,''), COALESCE(icon,''), COALESCE(color,''),
		       is_active, TO_CHAR(created_at, 'YYYY-MM-DD HH24:MI:SS')
		FROM   "dat_location_type" WHERE 1=1
	`
	args := []any{}
	argIdx := 1
	if companyID != "" {
		query += fmt.Sprintf(" AND company_id = $%d", argIdx)
		args = append(args, companyID)
		argIdx++
	}
	if search != "" {
		query += fmt.Sprintf(" AND (code ILIKE $%d OR name ILIKE $%d)", argIdx, argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}
	sortExpr := map[string]string{
		"code":       "code",
		"name":       "name",
		"created_at": "created_at",
	}
	expr, ok := sortExpr[sortBy]
	if !ok {
		expr = "created_at"
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

	var list []LocationTypeItem
	for rows.Next() {
		var item LocationTypeItem
		if err := rows.Scan(
			&item.ID, &item.Code, &item.Name, &item.Description,
			&item.Icon, &item.Color, &item.IsActive, &item.CreatedAt,
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

func (r *repository) CreateLocationType(companyID string, req LocationTypeCreateRequest) error {
	query := `
		INSERT INTO "dat_location_type" (id, company_id, code, name, description, icon, color, is_active, created_at, updated_at)
		VALUES (gen_random_uuid()::text, $1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
	`
	description := mechanic.NullableString(req.Description)
	icon := mechanic.NullableString(req.Icon)
	color := mechanic.NullableString(req.Color)
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	_, err := r.db.ExecContext(context.Background(), query,
		companyID, req.Code, req.Name, description, icon, color, isActive,
	)
	return err
}

func (r *repository) UpdateLocationType(id string, req LocationTypeUpdateRequest) error {
	query := `
		UPDATE "dat_location_type" SET
			code = $2, name = $3, description = $4,
			icon = $5, color = $6, is_active = $7
		WHERE id = $1
	`
	description := mechanic.NullableString(req.Description)
	icon := mechanic.NullableString(req.Icon)
	color := mechanic.NullableString(req.Color)
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	_, err := r.db.ExecContext(context.Background(), query,
		id, req.Code, req.Name, description, icon, color, isActive,
	)
	return err
}

func (r *repository) DeleteLocationType(id string) (int64, error) {
	result, err := r.db.ExecContext(context.Background(),
		`DELETE FROM "dat_location_type" WHERE id = $1`, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (r *repository) ListLocationTypeSelect(companyID string) ([]LocationTypeSelectItem, error) {
	query := `
		SELECT id, code, name
		FROM   "dat_location_type"
		WHERE  1=1
	`
	args := []any{}
	argIdx := 1
	if companyID != "" {
		query += fmt.Sprintf(" AND company_id = $%d", argIdx)
		args = append(args, companyID)
		argIdx++
	}
	query += ` AND is_active = true ORDER BY code`
	rows, err := r.db.QueryContext(context.Background(), query, args...)
	if err != nil {
		return nil, err
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
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) CountLocation(companyID, search string) (int, error) {
	query := `SELECT COUNT(*) FROM "dat_location" WHERE 1=1`
	args := []any{}
	argIdx := 1
	if companyID != "" {
		query += fmt.Sprintf(" AND company_id = $%d", argIdx)
		args = append(args, companyID)
		argIdx++
	}
	if search != "" {
		query += fmt.Sprintf(" AND (code ILIKE $%d OR name ILIKE $%d)", argIdx, argIdx)
		args = append(args, "%"+search+"%")
	}
	var total int
	err := r.db.QueryRowContext(context.Background(), query, args...).Scan(&total)
	return total, err
}

func (r *repository) ListLocation(companyID, search string, page, size int, sortBy, sortOrder string) ([]LocationItem, error) {
	query := `
		SELECT id, location_type_id, COALESCE(parent_id,''),
		       code, name, COALESCE(description,''), COALESCE(address,''),
		       COALESCE(city,''), COALESCE(province,''), country, COALESCE(postal_code,''),
		       COALESCE(latitude::text,''), COALESCE(longitude::text,''),
		       COALESCE(phone,''), COALESCE(email,''), COALESCE(timezone,''),
		       status, is_active, TO_CHAR(created_at, 'YYYY-MM-DD HH24:MI:SS')
		FROM   "dat_location" WHERE 1=1
	`
	args := []any{}
	argIdx := 1
	if companyID != "" {
		query += fmt.Sprintf(" AND company_id = $%d", argIdx)
		args = append(args, companyID)
		argIdx++
	}
	if search != "" {
		query += fmt.Sprintf(" AND (code ILIKE $%d OR name ILIKE $%d)", argIdx, argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}
	sortExpr := map[string]string{
		"code":       "code",
		"name":       "name",
		"created_at": "created_at",
		"status":     "status",
	}
	expr, ok := sortExpr[sortBy]
	if !ok {
		expr = "created_at"
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

	var list []LocationItem
	for rows.Next() {
		var item LocationItem
		if err := rows.Scan(
			&item.ID, &item.LocationTypeID, &item.ParentID,
			&item.Code, &item.Name, &item.Description, &item.Address,
			&item.City, &item.Province, &item.Country, &item.PostalCode,
			&item.Latitude, &item.Longitude,
			&item.Phone, &item.Email, &item.Timezone,
			&item.Status, &item.IsActive, &item.CreatedAt,
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

func (r *repository) CreateLocation(companyID string, req LocationCreateRequest) error {
	query := `
		INSERT INTO "dat_location" (id, company_id, location_type_id, parent_id,
			code, name, description, address, city, province, country, postal_code,
			latitude, longitude, phone, email, timezone, status, is_active, created_at, updated_at)
		VALUES (gen_random_uuid()::text, $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17::location_status, $18, NOW(), NOW())
	`
	parentID := mechanic.NullableString(req.ParentID)
	description := mechanic.NullableString(req.Description)
	address := mechanic.NullableString(req.Address)
	city := mechanic.NullableString(req.City)
	province := mechanic.NullableString(req.Province)
	postalCode := mechanic.NullableString(req.PostalCode)
	phone := mechanic.NullableString(req.Phone)
	email := mechanic.NullableString(req.Email)
	timezone := mechanic.NullableString(req.Timezone)
	country := req.Country
	if country == "" {
		country = "Indonesia"
	}
	status := req.Status
	if status == "" {
		status = "active"
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	latitude := mechanic.NullableFloat(req.Latitude)
	longitude := mechanic.NullableFloat(req.Longitude)
	_, err := r.db.ExecContext(context.Background(), query,
		companyID, req.LocationTypeID, parentID,
		req.Code, req.Name, description, address, city, province, country, postalCode,
		latitude, longitude, phone, email, timezone, status, isActive,
	)
	return err
}

func (r *repository) UpdateLocation(id string, req LocationUpdateRequest) error {
	query := `
		UPDATE "dat_location" SET
			location_type_id = $2, parent_id = $3,
			code = $4, name = $5, description = $6, address = $7,
			city = $8, province = $9, country = $10, postal_code = $11,
			latitude = $12, longitude = $13, phone = $14, email = $15,
			timezone = $16, status = $17::location_status, is_active = $18, updated_at = NOW()
		WHERE id = $1
	`
	parentID := mechanic.NullableString(req.ParentID)
	description := mechanic.NullableString(req.Description)
	address := mechanic.NullableString(req.Address)
	city := mechanic.NullableString(req.City)
	province := mechanic.NullableString(req.Province)
	postalCode := mechanic.NullableString(req.PostalCode)
	phone := mechanic.NullableString(req.Phone)
	email := mechanic.NullableString(req.Email)
	timezone := mechanic.NullableString(req.Timezone)
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	latitude := mechanic.NullableFloat(req.Latitude)
	longitude := mechanic.NullableFloat(req.Longitude)
	_, err := r.db.ExecContext(context.Background(), query,
		id, req.LocationTypeID, parentID,
		req.Code, req.Name, description, address, city, province, req.Country, postalCode,
		latitude, longitude, phone, email, timezone, req.Status, isActive,
	)
	return err
}

func (r *repository) DeleteLocation(id string) (int64, error) {
	result, err := r.db.ExecContext(context.Background(),
		`DELETE FROM "dat_location" WHERE id = $1`, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (r *repository) ListLocationSelect(companyID string) ([]LocationSelectItem, error) {
	query := `
		SELECT id, COALESCE(parent_id,''), code, name
		FROM   "dat_location"
		WHERE  1=1
	`
	args := []any{}
	argIdx := 1
	if companyID != "" {
		query += fmt.Sprintf(" AND company_id = $%d", argIdx)
		args = append(args, companyID)
		argIdx++
	}
	query += ` AND is_active = true ORDER BY code`
	rows, err := r.db.QueryContext(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []LocationSelectItem
	for rows.Next() {
		var item LocationSelectItem
		if err := rows.Scan(&item.ID, &item.ParentID, &item.Code, &item.Name); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}
