package SM04

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

func (r *repository) ListSignatureType() ([]SignatureTypeItem, error) {
	query := `
		SELECT id, code, name, TO_CHAR(created_at, 'YYYY-MM-DD HH24:MI:SS')
		FROM   "dat_signature_type" ORDER BY code`
	rows, err := r.db.QueryContext(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []SignatureTypeItem
	for rows.Next() {
		var item SignatureTypeItem
		if err := rows.Scan(&item.ID, &item.Code, &item.Name, &item.CreatedAt); err != nil {
			return nil, err
		}
		steps, err := r.ListApprovalStepsByType(item.ID)
		if err == nil {
			item.Steps = steps
		}
		list = append(list, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) ListApprovalStepsByType(typeID string) ([]ApprovalStepI, error) {
	query := `SELECT id, step, condition FROM "dat_approval_step" WHERE type_id = $1 ORDER BY step`
	rows, err := r.db.QueryContext(context.Background(), query, typeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []ApprovalStepI
	for rows.Next() {
		var s ApprovalStepI
		if err := rows.Scan(&s.ID, &s.Step, &s.Condition); err != nil {
			return nil, err
		}
		signs, err := r.ListApprovalSignsByStep(s.ID)
		if err == nil {
			s.Signers = signs
		}
		list = append(list, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) ListApprovalSignsByStep(stepID string) ([]ApprovalSignI, error) {
	query := `SELECT user_id FROM "dat_approval_sign" WHERE step_id = $1`
	rows, err := r.db.QueryContext(context.Background(), query, stepID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []ApprovalSignI
	for rows.Next() {
		var s ApprovalSignI
		if err := rows.Scan(&s.UserID); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) CreateSignatureType(code, name string) (string, error) {
	var id string
	query := `
		INSERT INTO "dat_signature_type" (
			id, code, name, created_at, updated_at
		) VALUES (gen_random_uuid()::text,
			$1, $2, NOW(), NOW()
		) RETURNING id`
	err := r.db.QueryRowContext(context.Background(), query, code, name).Scan(&id)
	return id, err
}

func (r *repository) InsertApprovalStep(typeID string, step int, condition string) (string, error) {
	var id string
	query := `
		INSERT INTO "dat_approval_step" (
			id, type_id, step, condition
		) VALUES (gen_random_uuid()::text,
			$1, $2, $3::approval_flag
		) RETURNING id`
	err := r.db.QueryRowContext(context.Background(), query, typeID, step, condition).Scan(&id)
	return id, err
}

func (r *repository) InsertApprovalSign(stepID, userID string) error {
	query := `
		INSERT INTO "dat_approval_sign" (
			id, step_id, user_id
		) VALUES (gen_random_uuid()::text,
			$1, $2
		) ON CONFLICT DO NOTHING`
	_, err := r.db.ExecContext(context.Background(), query, stepID, userID)
	return err
}

func (r *repository) UpdateSignatureType(id, code, name string) error {
	query := `UPDATE "dat_signature_type" SET code = $2, name = $3, updated_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(context.Background(), query, id, code, name)
	return err
}

func (r *repository) DeleteApprovalStepsByType(typeID string) error {
	query := `DELETE FROM "dat_approval_step" WHERE type_id = $1`
	_, err := r.db.ExecContext(context.Background(), query, typeID)
	return err
}

func (r *repository) DeleteApprovalSignsByStep(stepID string) error {
	query := `DELETE FROM "dat_approval_sign" WHERE step_id = $1`
	_, err := r.db.ExecContext(context.Background(), query, stepID)
	return err
}

func (r *repository) ListUser(companyID, search string) ([]UserListItem, error) {
	query := `
		SELECT id, username, fullname, COALESCE(company_id, '')
		FROM   "dat_user"
		WHERE  is_active = TRUE`
	args := []any{}
	argIdx := 1
	if companyID != "" {
		query += fmt.Sprintf(` AND company_id = $%d`, argIdx)
		args = append(args, companyID)
		argIdx++
	}
	if search != "" {
		query += fmt.Sprintf(` AND (username ILIKE $%d OR fullname ILIKE $%d)`, argIdx, argIdx)
		args = append(args, "%"+search+"%")
	}
	query += ` ORDER BY fullname`
	rows, err := r.db.QueryContext(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []UserListItem
	for rows.Next() {
		var item UserListItem
		if err := rows.Scan(
			&item.ID, &item.Username, &item.Fullname,
			&item.CompanyID,
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
