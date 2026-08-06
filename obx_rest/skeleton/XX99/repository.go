package XX99

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

func (r *repository) List(meta mechanic.ActionMeta) ([]XX99Item, mechanic.GridMeta, error) {
	where := `WHERE 1=1`
	args := []any{}
	argIdx := 1

	if meta.Search != "" {
		where += fmt.Sprintf(` AND (company_id ILIKE $%d OR code ILIKE $%d OR name ILIKE $%d)`, argIdx, argIdx, argIdx)
		args = append(args, "%"+meta.Search+"%")
		argIdx++
	}

	countQuery := `SELECT COUNT(*) FROM "dat_xx99" ` + where
	var total int
	if err := r.db.QueryRowContext(context.Background(), countQuery, args...).Scan(&total); err != nil {
		return nil, mechanic.GridMeta{}, err
	}

	sortExpr := map[string]string{
		"company_id": "company_id",
		"code":       "code",
		"name":       "name",
		"created_at": "created_at",
		"updated_at": "updated_at",
	}
	orderCol, ok := sortExpr[meta.SortBy]
	if !ok {
		orderCol = "created_at"
	}
	if meta.SortOrder != "asc" {
		meta.SortOrder = "desc"
	}

	query := `
		SELECT id, company_id, code, name, is_active,
		       TO_CHAR(created_at, 'YYYY-MM-DD HH24:MI:SS'),
		       TO_CHAR(updated_at, 'YYYY-MM-DD HH24:MI:SS')
		FROM "dat_xx99"
	` + where + fmt.Sprintf(` ORDER BY %s %s`, orderCol, meta.SortOrder)

	offset := (meta.Page - 1) * meta.Size
	query += fmt.Sprintf(` LIMIT $%d OFFSET $%d`, argIdx, argIdx+1)
	args = append(args, meta.Size, offset)

	rows, err := r.db.QueryContext(context.Background(), query, args...)
	if err != nil {
		return nil, mechanic.GridMeta{}, err
	}
	defer rows.Close()

	list := []XX99Item{}
	for rows.Next() {
		var item XX99Item
		if err := rows.Scan(
			&item.ID,
			&item.CompanyID,
			&item.Code,
			&item.Name,
			&item.IsActive,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, mechanic.GridMeta{}, err
		}
		list = append(list, item)
	}
	if err := rows.Err(); err != nil {
		return nil, mechanic.GridMeta{}, err
	}

	return list, mechanic.BuildMeta(meta.Page, meta.Size, total), nil
}

func (r *repository) GetByID(id string) (XX99Item, error) {
	query := `
		SELECT id, company_id, code, name, is_active,
		       TO_CHAR(created_at, 'YYYY-MM-DD HH24:MI:SS'),
		       TO_CHAR(updated_at, 'YYYY-MM-DD HH24:MI:SS')
		FROM "dat_xx99"
		WHERE id = $1
	`

	var item XX99Item
	err := r.db.QueryRowContext(context.Background(), query, id).Scan(
		&item.ID,
		&item.CompanyID,
		&item.Code,
		&item.Name,
		&item.IsActive,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return XX99Item{}, mechanic.NotFound("XX99 data not found")
	}
	if err != nil {
		return XX99Item{}, err
	}
	return item, nil
}

func (r *repository) Create(req XX99CreateRequest) error {
	query := `
		INSERT INTO "dat_xx99" (
			id, company_id, code, name, is_active, created_at, updated_at
		)
		VALUES (
			gen_random_uuid()::text, $1, $2, $3, $4, NOW(), NOW()
		)
	`
	_, err := r.db.ExecContext(context.Background(), query,
		req.CompanyID,
		req.Code,
		req.Name,
		req.IsActive,
	)
	return err
}

func (r *repository) Update(id string, req XX99UpdateRequest) error {
	query := `
		UPDATE "dat_xx99"
		SET company_id = COALESCE($2, company_id),
		    code = COALESCE($3, code),
		    name = COALESCE($4, name),
		    is_active = COALESCE($5, is_active),
		    updated_at = NOW()
		WHERE id = $1
	`
	result, err := r.db.ExecContext(context.Background(), query,
		id,
		mechanic.NullableString(req.CompanyID),
		mechanic.NullableString(req.Code),
		mechanic.NullableString(req.Name),
		req.IsActive,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return mechanic.NotFound("XX99 data not found")
	}
	return nil
}

func (r *repository) Delete(id string) error {
	query := `DELETE FROM "dat_xx99" WHERE id = $1`
	result, err := r.db.ExecContext(context.Background(), query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return mechanic.NotFound("XX99 data not found")
	}
	return nil
}
