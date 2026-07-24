package SP01

import (
	"context"
	"database/sql"
)

type repository struct {
	db *sql.DB
}

func NRepo(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) ListUserCompany(userID string) ([]UserCompanyItem, error) {
	query := `
		SELECT c.id, c.name
		FROM "dat_company" c
		WHERE c.is_active = true AND (
			EXISTS (
				SELECT 1 FROM "dat_user_company" uc
				WHERE uc.user_id = $1 AND uc.company_id = c.id AND uc.is_active = true
			)
			OR c.id = (
				SELECT u.company_id FROM "dat_user" u
				WHERE u.id = $1 AND u.company_id != '')
		)
		ORDER BY c.name ASC
	`
	rows, err := r.db.QueryContext(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []UserCompanyItem
	for rows.Next() {
		var item UserCompanyItem
		if err := rows.Scan(&item.CompanyID, &item.Name); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) ListAllModuleTree() ([]ModuleTreeNode, error) {
	query := `
		SELECT id, COALESCE(parent_id,''), code, name, path, is_page
		FROM "dat_module"
		WHERE is_active = true
		ORDER BY code ASC
	`
	rows, err := r.db.QueryContext(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type flatModule struct {
		ID       string
		ParentID string
		Code     string
		Name     string
		Path     string
		IsPage   bool
	}
	var flat []flatModule
	for rows.Next() {
		var item flatModule
		if err := rows.Scan(
			&item.ID, &item.ParentID, &item.Code,
			&item.Name, &item.Path, &item.IsPage); err != nil {
			return nil, err
		}
		flat = append(flat, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	nodeMap := make(map[string]*ModuleTreeNode)
	for _, f := range flat {
		nodeMap[f.ID] = &ModuleTreeNode{
			ID:     f.ID,
			Code:   f.Code,
			Name:   f.Name,
			Path:   f.Path,
			IsPage: f.IsPage,
		}
	}

	for _, f := range flat {
		if f.ParentID != "" {
			if parent, ok := nodeMap[f.ParentID]; ok {
				parent.Children = append(parent.Children, *nodeMap[f.ID])
			}
		}
	}

	var roots []ModuleTreeNode
	for _, f := range flat {
		if f.ParentID == "" {
			roots = append(roots, *nodeMap[f.ID])
		}
	}
	return roots, nil
}

func (r *repository) ListUserModule(userID, companyID string) ([]ModuleTreeNode, error) {
	query := `
		WITH RECURSIVE "allowed" AS (
			SELECT m.id, m.parent_id, m.code, m.name, m.path, m.is_page, m.is_active, up.level::text AS level
			FROM "dat_user_company" uc
			JOIN "dat_company_module" cm ON cm.company_id = uc.company_id
			JOIN "dat_user_privilege" up ON up.user_company_id = uc.id AND up.module_id = cm.module_id
			JOIN "dat_module" m ON m.id = cm.module_id
			WHERE uc.user_id = $1
				AND uc.company_id = $2
				AND uc.is_active = true
				AND cm.is_active = true
				AND m.is_active = true
				AND up.level::text <> 'hide'
			UNION
			SELECT p.id, p.parent_id, p.code, p.name, p.path, p.is_page, p.is_active, NULL::text AS level
			FROM "dat_module" p
			JOIN "allowed" a ON a.parent_id = p.id
			WHERE p.is_active = true
		)
		SELECT
			id, COALESCE(parent_id,'') as parent_id, code, name, path, is_page, is_active,
			COALESCE(
				CASE MAX(CASE level WHEN 'post' THEN 2 WHEN 'view' THEN 1 ELSE 0 END)
					WHEN 2 THEN 'post'
					WHEN 1 THEN 'view'
					ELSE 'hide'
				END,
				'view'
			) AS level
		FROM "allowed"
		GROUP BY id, parent_id, code, name, path, is_page, is_active
		ORDER BY name
	`
	rows, err := r.db.QueryContext(context.Background(), query, userID, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type flatModule struct {
		ID       string
		ParentID string
		Code     string
		Name     string
		Path     string
		IsPage   bool
		IsActive bool
		Level    string
	}
	var flat []flatModule
	for rows.Next() {
		var item flatModule
		if err := rows.Scan(
			&item.ID, &item.ParentID, &item.Code, &item.Name,
			&item.Path, &item.IsPage, &item.IsActive, &item.Level); err != nil {
			return nil, err
		}
		flat = append(flat, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	nodeMap := make(map[string]*ModuleTreeNode)
	for _, f := range flat {
		nodeMap[f.ID] = &ModuleTreeNode{
			ID:     f.ID,
			Code:   f.Code,
			Name:   f.Name,
			Path:   f.Path,
			IsPage: f.IsPage,
		}
	}
	for _, f := range flat {
		if f.ParentID != "" {
			if parent, ok := nodeMap[f.ParentID]; ok {
				parent.Children = append(parent.Children, *nodeMap[f.ID])
			}
		}
	}
	var roots []ModuleTreeNode
	for _, f := range flat {
		if f.ParentID == "" {
			roots = append(roots, *nodeMap[f.ID])
		}
	}
	return roots, nil
}
