package SM04

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

func (r *repository) CountRequest(companyID, search string) (int, error) {
	query := `SELECT COUNT(*) FROM "dat_request" WHERE 1=1`
	args := []any{}
	argIdx := 1
	if companyID != "" {
		query += fmt.Sprintf(" AND company_id = $%d", argIdx)
		args = append(args, companyID)
		argIdx++
	}
	if search != "" {
		query += fmt.Sprintf(
			" AND (code ILIKE $%d OR title ILIKE $%d OR status::text ILIKE $%d)",
			argIdx, argIdx, argIdx)
		args = append(args, "%"+search+"%")
	}
	var total int
	err := r.db.QueryRowContext(context.Background(), query, args...).Scan(&total)
	return total, err
}

func (r *repository) ListRequest(companyID, search string, page, size int, sortBy, sortOrder string) (
	[]RequestItem, error) {
	query := `
		SELECT r.id, r.company_id, c.name AS company_name,
		       r.type_id, t.name AS type_name,
		       r.requester_id, u.fullname AS requester_name,
		       r.code, r.title, COALESCE(r.description, ''),
		       r.priority, r.status::text, r.current_step,
		       COALESCE(r.completion_note, ''),
		       COALESCE(r.completed_by, ''),
		       COALESCE(TO_CHAR(r.completed_at, 'YYYY-MM-DD HH24:MI:SS'), ''),
		       TO_CHAR(r.created_at, 'YYYY-MM-DD HH24:MI:SS')
		FROM   "dat_request" r
		JOIN   "dat_company" c ON c.id = r.company_id
		JOIN   "dat_signature_type" t ON t.id = r.type_id
		JOIN   "dat_user" u ON u.id = r.requester_id
		WHERE  1=1
	`
	args := []any{}
	argIdx := 1
	if companyID != "" {
		query += fmt.Sprintf(" AND r.company_id = $%d", argIdx)
		args = append(args, companyID)
		argIdx++
	}
	if search != "" {
		query += fmt.Sprintf(
			" AND (r.code ILIKE $%d OR r.title ILIKE $%d OR r.status::text ILIKE $%d)",
			argIdx, argIdx, argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}
	sortExpr := map[string]string{
		"code":       "r.code",
		"title":      "r.title",
		"priority":   "r.priority",
		"status":     "r.status",
		"created_at": "r.created_at",
	}
	expr, ok := sortExpr[sortBy]
	if !ok {
		expr = "r.created_at"
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

	var list []RequestItem
	for rows.Next() {
		var item RequestItem
		if err := rows.Scan(
			&item.ID, &item.CompanyID, &item.CompanyName,
			&item.TypeID, &item.TypeName,
			&item.RequesterID, &item.RequesterName,
			&item.Code, &item.Title, &item.Description,
			&item.Priority, &item.Status, &item.CurrentStep,
			&item.CompletionNote, &item.CompletedBy, &item.CompletedAt,
			&item.CreatedAt,
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

func (r *repository) CreateRequest(companyID string, req RequestCreateRequest) error {
	query := `
		INSERT INTO "dat_request" (
			id, company_id, type_id, requester_id, code, title, description,
			priority, status, current_step, created_at, updated_at
		) VALUES (gen_random_uuid()::text,
			$1, $2, $3, $4, $5, $6, $7,
			'PENDING'::approval_info, 1, NOW(), NOW()
		)
	`
	_, err := r.db.ExecContext(context.Background(), query,
		companyID, req.TypeID, req.RequesterID, req.Code, req.Title,
		req.Description, req.Priority,
	)
	return err
}

func (r *repository) UpdateRequest(id string, req RequestUpdateRequest) error {
	query := `
		UPDATE "dat_request" SET
			type_id = $2, requester_id = $3, code = $4, title = $5,
			description = $6, priority = $7, updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.ExecContext(context.Background(), query,
		id, req.TypeID, req.RequesterID, req.Code, req.Title,
		req.Description, req.Priority,
	)
	return err
}

func (r *repository) DeleteRequest(id string) error {
	_, err := r.db.ExecContext(context.Background(),
		`DELETE FROM "dat_request" WHERE id = $1`, id)
	return err
}

func (r *repository) CountDocument(companyID, search string) (int, error) {
	query := `SELECT COUNT(*) FROM "dat_document" WHERE 1=1`
	args := []any{}
	argIdx := 1
	if companyID != "" {
		query += fmt.Sprintf(" AND company_id = $%d", argIdx)
		args = append(args, companyID)
		argIdx++
	}
	if search != "" {
		query += fmt.Sprintf(
			" AND (code ILIKE $%d OR title ILIKE $%d OR status ILIKE $%d)",
			argIdx, argIdx, argIdx)
		args = append(args, "%"+search+"%")
	}
	var total int
	err := r.db.QueryRowContext(context.Background(), query, args...).Scan(&total)
	return total, err
}

func (r *repository) ListDocument(companyID, search string, page, size int, sortBy, sortOrder string) (
	[]DocumentItem, error) {
	query := `
		SELECT d.id, d.company_id, c.name AS company_name,
		       COALESCE(d.category_id, ''),
		       COALESCE(dc.name, '') AS category_name,
		       d.code, d.title, COALESCE(d.description, ''),
		       d.content_type, COALESCE(d.content, ''),
		       COALESCE(d.file_name, ''),
		       COALESCE(d.file_size, 0),
		       d.version, d.status,
		       COALESCE(d.created_by, ''),
		       COALESCE(d.approved_by, ''),
		       COALESCE(TO_CHAR(d.approved_at, 'YYYY-MM-DD HH24:MI:SS'), ''),
		       d.is_active,
		       TO_CHAR(d.created_at, 'YYYY-MM-DD HH24:MI:SS')
		FROM   "dat_document" d
		JOIN   "dat_company" c ON c.id = d.company_id
		LEFT JOIN "dat_document_category" dc ON dc.id = d.category_id
		WHERE  1=1
	`
	args := []any{}
	argIdx := 1
	if companyID != "" {
		query += fmt.Sprintf(" AND d.company_id = $%d", argIdx)
		args = append(args, companyID)
		argIdx++
	}
	if search != "" {
		query += fmt.Sprintf(
			" AND (d.code ILIKE $%d OR d.title ILIKE $%d OR d.status ILIKE $%d)",
			argIdx, argIdx, argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}
	sortExpr := map[string]string{
		"code":          "d.code",
		"title":         "d.title",
		"status":        "d.status",
		"category_name": "dc.name",
		"created_at":    "d.created_at",
	}
	expr, ok := sortExpr[sortBy]
	if !ok {
		expr = "d.created_at"
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

	var list []DocumentItem
	for rows.Next() {
		var item DocumentItem
		if err := rows.Scan(
			&item.ID, &item.CompanyID, &item.CompanyName,
			&item.CategoryID, &item.CategoryName,
			&item.Code, &item.Title, &item.Description,
			&item.ContentType, &item.Content,
			&item.FileName, &item.FileSize,
			&item.Version, &item.Status,
			&item.CreatedBy, &item.ApprovedBy, &item.ApprovedAt,
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

func (r *repository) CreateDocument(companyID, createdBy string, req DocumentCreateRequest) error {
	var fileSize any
	if req.FileSize > 0 {
		fileSize = req.FileSize
	}
	query := `
		INSERT INTO "dat_document" (
			id, company_id, category_id, code, title, description,
			content_type, content, file_name, file_size,
			version, status, created_by, is_active, created_at, updated_at
		) VALUES (gen_random_uuid()::text,
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9,
			$10, $11, $12, TRUE, NOW(), NOW()
		)
	`
	_, err := r.db.ExecContext(context.Background(), query,
		companyID, mechanic.NullableString(req.CategoryID), req.Code, req.Title, req.Description,
		req.ContentType, mechanic.NullableString(req.Content), mechanic.NullableString(req.FileName), fileSize,
		req.Version, req.Status, createdBy,
	)
	return err
}

func (r *repository) UpdateDocument(id string, req DocumentUpdateRequest) error {
	var fileSize any
	if req.FileSize > 0 {
		fileSize = req.FileSize
	}
	query := `
		UPDATE "dat_document" SET
			category_id = $2, code = $3, title = $4, description = $5,
			content_type = $6, content = $7, file_name = $8, file_size = $9,
			version = $10, status = $11, updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.ExecContext(context.Background(), query,
		id, mechanic.NullableString(req.CategoryID), req.Code, req.Title, req.Description,
		req.ContentType, mechanic.NullableString(req.Content), mechanic.NullableString(req.FileName), fileSize,
		req.Version, req.Status,
	)
	return err
}

func (r *repository) DeleteDocument(id string) error {
	_, err := r.db.ExecContext(context.Background(),
		`DELETE FROM "dat_document" WHERE id = $1`, id)
	return err
}

func (r *repository) ListDocumentCategory(companyID string) ([]DocumentCategoryItem, error) {
	query := `SELECT id, name FROM "dat_document_category" WHERE is_active = TRUE`
	args := []any{}
	argIdx := 1
	if companyID != "" {
		query += fmt.Sprintf(" AND company_id = $%d", argIdx)
		args = append(args, companyID)
	}
	query += ` ORDER BY name`
	rows, err := r.db.QueryContext(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []DocumentCategoryItem
	for rows.Next() {
		var item DocumentCategoryItem
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) CountDocumentRevision(documentID, companyID, search string) (int, error) {
	query := `SELECT COUNT(*) FROM "dat_document_version" WHERE document_id = $1`
	args := []any{documentID}
	argIdx := 2
	if companyID != "" {
		query += fmt.Sprintf(" AND company_id = $%d", argIdx)
		args = append(args, companyID)
		argIdx++
	}
	if search != "" {
		query += fmt.Sprintf(
			" AND (version ILIKE $%d OR status ILIKE $%d OR COALESCE(note, '') ILIKE $%d)",
			argIdx, argIdx, argIdx)
		args = append(args, "%"+search+"%")
	}
	var total int
	err := r.db.QueryRowContext(context.Background(), query, args...).Scan(&total)
	return total, err
}

func (r *repository) ListDocumentRevision(documentID, companyID, search string, page, size int, sortBy, sortOrder string) (
	[]DocumentRevisionItem, error) {
	query := `
		SELECT v.id, v.document_id, v.version, COALESCE(v.content, ''),
		       COALESCE(v.file_path, ''), v.status, COALESCE(v.note, ''),
		       COALESCE(v.created_by, ''),
		       TO_CHAR(v.created_at, 'YYYY-MM-DD HH24:MI:SS')
		FROM   "dat_document_version" v
		WHERE  v.document_id = $1
	`
	args := []any{documentID}
	argIdx := 2
	if companyID != "" {
		query += fmt.Sprintf(" AND v.company_id = $%d", argIdx)
		args = append(args, companyID)
		argIdx++
	}
	if search != "" {
		query += fmt.Sprintf(
			" AND (v.version ILIKE $%d OR v.status ILIKE $%d OR COALESCE(v.note, '') ILIKE $%d)",
			argIdx, argIdx, argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}
	sortExpr := map[string]string{
		"version":    "v.version",
		"status":     "v.status",
		"created_at": "v.created_at",
	}
	expr, ok := sortExpr[sortBy]
	if !ok {
		expr = "v.created_at"
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

	var list []DocumentRevisionItem
	for rows.Next() {
		var item DocumentRevisionItem
		if err := rows.Scan(
			&item.ID, &item.DocumentID, &item.Version, &item.Content,
			&item.FilePath, &item.Status, &item.Note,
			&item.CreatedBy, &item.CreatedAt,
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

func (r *repository) CreateDocumentRevision(documentID, companyID, createdBy string, req DocumentRevisionCreateRequest) error {
	query := `
		INSERT INTO "dat_document_version" (
			id, document_id, company_id, version, content, file_path,
			status, note, created_by, created_at
		) VALUES (gen_random_uuid()::text,
			$1, $2, $3, $4, $5,
			$6, $7, $8, NOW()
		)
	`
	_, err := r.db.ExecContext(context.Background(), query,
		documentID, companyID, req.Version,
		mechanic.NullableString(req.Content), mechanic.NullableString(req.FilePath),
		req.Status, mechanic.NullableString(req.Note), createdBy,
	)
	return err
}

func (r *repository) UpdateDocumentRevision(id string, req DocumentRevisionUpdateRequest) error {
	query := `
		UPDATE "dat_document_version" SET
			version = $2, content = $3, file_path = $4,
			status = $5, note = $6
		WHERE id = $1
	`
	_, err := r.db.ExecContext(context.Background(), query,
		id, req.Version,
		mechanic.NullableString(req.Content), mechanic.NullableString(req.FilePath),
		req.Status, mechanic.NullableString(req.Note),
	)
	return err
}

func (r *repository) DeleteDocumentRevision(id string) error {
	_, err := r.db.ExecContext(context.Background(),
		`DELETE FROM "dat_document_version" WHERE id = $1`, id)
	return err
}

func (r *repository) CountDocumentEvidence(documentID, companyID, search string) (int, error) {
	query := `SELECT COUNT(*) FROM "dat_document_approval" WHERE document_id = $1`
	args := []any{documentID}
	argIdx := 2
	if companyID != "" {
		query += fmt.Sprintf(" AND company_id = $%d", argIdx)
		args = append(args, companyID)
		argIdx++
	}
	if search != "" {
		query += fmt.Sprintf(
			" AND (action ILIKE $%d OR COALESCE(note, '') ILIKE $%d)",
			argIdx, argIdx)
		args = append(args, "%"+search+"%")
	}
	var total int
	err := r.db.QueryRowContext(context.Background(), query, args...).Scan(&total)
	return total, err
}

func (r *repository) ListDocumentEvidence(documentID, companyID, search string, page, size int, sortBy, sortOrder string) (
	[]DocumentEvidenceItem, error) {
	query := `
		SELECT e.id, e.document_id, e.action, COALESCE(e.note, ''),
		       e.user_id, COALESCE(u.fullname, '') AS user_name,
		       TO_CHAR(e.created_at, 'YYYY-MM-DD HH24:MI:SS')
		FROM   "dat_document_approval" e
		JOIN   "dat_user" u ON u.id = e.user_id
		WHERE  e.document_id = $1
	`
	args := []any{documentID}
	argIdx := 2
	if companyID != "" {
		query += fmt.Sprintf(" AND e.company_id = $%d", argIdx)
		args = append(args, companyID)
		argIdx++
	}
	if search != "" {
		query += fmt.Sprintf(
			" AND (e.action ILIKE $%d OR COALESCE(e.note, '') ILIKE $%d OR u.fullname ILIKE $%d)",
			argIdx, argIdx, argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}
	sortExpr := map[string]string{
		"action":     "e.action",
		"user_name":  "u.fullname",
		"created_at": "e.created_at",
	}
	expr, ok := sortExpr[sortBy]
	if !ok {
		expr = "e.created_at"
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

	var list []DocumentEvidenceItem
	for rows.Next() {
		var item DocumentEvidenceItem
		if err := rows.Scan(
			&item.ID, &item.DocumentID, &item.Action, &item.Note,
			&item.UserID, &item.UserName, &item.CreatedAt,
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

func (r *repository) CreateDocumentEvidence(documentID, companyID, userID string, req DocumentEvidenceCreateRequest) error {
	query := `
		INSERT INTO "dat_document_approval" (
			id, document_id, company_id, action, note, user_id, created_at
		) VALUES (gen_random_uuid()::text,
			$1, $2, $3, $4, $5, NOW()
		)
	`
	_, err := r.db.ExecContext(context.Background(), query,
		documentID, companyID, req.Action, mechanic.NullableString(req.Note), userID,
	)
	return err
}

func (r *repository) UpdateDocumentEvidence(id string, req DocumentEvidenceUpdateRequest) error {
	query := `
		UPDATE "dat_document_approval" SET
			action = $2, note = $3
		WHERE id = $1
	`
	_, err := r.db.ExecContext(context.Background(), query,
		id, req.Action, mechanic.NullableString(req.Note),
	)
	return err
}

func (r *repository) DeleteDocumentEvidence(id string) error {
	_, err := r.db.ExecContext(context.Background(),
		`DELETE FROM "dat_document_approval" WHERE id = $1`, id)
	return err
}

func (r *repository) GetRequestMeta(requestID string) (string, string, int, error) {
	var typeID, status string
	var currentStep int
	query := `SELECT type_id, status::text, current_step FROM "dat_request" WHERE id = $1`
	err := r.db.QueryRowContext(context.Background(), query, requestID).Scan(&typeID, &status, &currentStep)
	return typeID, status, currentStep, err
}

func (r *repository) CountSignatureForms(requestID string) (int, error) {
	var total int
	err := r.db.QueryRowContext(context.Background(),
		`SELECT COUNT(*) FROM "dat_signature_form" WHERE request_id = $1`, requestID).Scan(&total)
	return total, err
}

func (r *repository) ListSignatureForms(requestID, companyID string) ([]SignatureFormItem, error) {
	query := `
		SELECT f.id, f.step, f.request_id, f.condition::text, f.status::text,
		       TO_CHAR(f.created_at, 'YYYY-MM-DD HH24:MI:SS')
		FROM   "dat_signature_form" f
		JOIN   "dat_request" r ON r.id = f.request_id
		WHERE  f.request_id = $1
	`
	args := []any{requestID}
	argIdx := 2
	if companyID != "" {
		query += fmt.Sprintf(" AND r.company_id = $%d", argIdx)
		args = append(args, companyID)
	}
	query += ` ORDER BY f.step`

	rows, err := r.db.QueryContext(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []SignatureFormItem
	for rows.Next() {
		var item SignatureFormItem
		if err := rows.Scan(
			&item.ID, &item.Step, &item.RequestID, &item.Condition, &item.Status,
			&item.CreatedAt,
		); err != nil {
			return nil, err
		}
		flags, err := r.ListSignatureFlags(item.ID)
		if err == nil {
			item.Flags = flags
		}
		list = append(list, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) ListSignatureFlags(formID string) ([]SignatureFlagItem, error) {
	query := `
		SELECT f.id, f.form_id, f.user_id, u.fullname AS user_name,
		       f.status::text, COALESCE(f.comment, ''),
		       TO_CHAR(f.created_at, 'YYYY-MM-DD HH24:MI:SS')
		FROM   "dat_signature_flag" f
		JOIN   "dat_user" u ON u.id = f.user_id
		WHERE  f.form_id = $1
		ORDER  BY f.created_at`
	rows, err := r.db.QueryContext(context.Background(), query, formID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []SignatureFlagItem
	for rows.Next() {
		var item SignatureFlagItem
		if err := rows.Scan(
			&item.ID, &item.FormID, &item.UserID, &item.UserName,
			&item.Status, &item.Comment, &item.CreatedAt,
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

func (r *repository) GetSignatureForm(formID string) (string, string, int, error) {
	var requestID, condition string
	var step int
	query := `SELECT request_id, condition::text, step FROM "dat_signature_form" WHERE id = $1`
	err := r.db.QueryRowContext(context.Background(), query, formID).Scan(&requestID, &condition, &step)
	return requestID, condition, step, err
}

func (r *repository) GetSignatureFlag(flagID string) (string, string, error) {
	var formID, status string
	query := `SELECT form_id, status::text FROM "dat_signature_flag" WHERE id = $1`
	err := r.db.QueryRowContext(context.Background(), query, flagID).Scan(&formID, &status)
	return formID, status, err
}

func (r *repository) InsertSignatureForm(requestID string, step int, condition string) (string, error) {
	var id string
	query := `
		INSERT INTO "dat_signature_form" (
			id, step, request_id, condition, status, created_at, updated_at
		) VALUES (gen_random_uuid()::text,
			$1, $2, $3::approval_flag, 'PENDING'::approval_info, NOW(), NOW()
		) RETURNING id`
	err := r.db.QueryRowContext(context.Background(), query, step, requestID, condition).Scan(&id)
	return id, err
}

func (r *repository) InsertSignatureFlag(formID, userID string) error {
	query := `
		INSERT INTO "dat_signature_flag" (
			id, form_id, user_id, status, created_at, updated_at
		) VALUES (gen_random_uuid()::text,
			$1, $2, 'PENDING'::approval_info, NOW(), NOW()
		) ON CONFLICT (form_id, user_id) DO NOTHING`
	_, err := r.db.ExecContext(context.Background(), query, formID, userID)
	return err
}

func (r *repository) UpdateSignatureFlag(flagID, status, comment string) error {
	query := `UPDATE "dat_signature_flag" SET status = $2::approval_info, comment = $3, updated_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(context.Background(), query, flagID, status, comment)
	return err
}

func (r *repository) UpdateSignatureFormStatus(formID, status string) error {
	query := `UPDATE "dat_signature_form" SET status = $2::approval_info, updated_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(context.Background(), query, formID, status)
	return err
}

func (r *repository) UpdateRequestStatus(id, status string) error {
	query := `UPDATE "dat_request" SET status = $2::approval_info, updated_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(context.Background(), query, id, status)
	return err
}

func (r *repository) UpdateRequestStep(id string, step int) error {
	query := `UPDATE "dat_request" SET current_step = $2, updated_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(context.Background(), query, id, step)
	return err
}

func (r *repository) UpdateRequestCompletion(id, status, completedBy, completionNote string) error {
	query := `
		UPDATE "dat_request" SET
			status = $2::approval_info,
			completed_by = $3,
			completion_note = $4,
			completed_at = NOW(),
			updated_at = NOW()
		WHERE id = $1`
	_, err := r.db.ExecContext(context.Background(), query,
		id, status, mechanic.NullableString(completedBy), mechanic.NullableString(completionNote))
	return err
}
