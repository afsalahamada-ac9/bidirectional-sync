/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package repository

import (
	"database/sql"

	"sudhagar/glad/entity"
)

// LabelMySQL mysql repo
type LabelMySQL struct {
	db *sql.DB
}

// NewLabelMySQL create new repository
func NewLabelMySQL(db *sql.DB) *LabelMySQL {
	return &LabelMySQL{
		db: db,
	}
}

// Create a label
func (r *LabelMySQL) Create(e *entity.Label) error {
	stmt, err := r.db.Prepare(`
		INSERT INTO label (id, tenant_id, name, color, ref_count)
		VALUES(?,?,?,?,?)`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		e.ID,
		e.TenantID,
		e.Name,
		e.Color,
		e.RefCount,
	)
	if err != nil {
		return err
	}
	err = stmt.Close()
	if err != nil {
		return err
	}
	return nil
}

// Get a label
func (r *LabelMySQL) Get(tenantID, id entity.ID) (*entity.Label, error) {
	stmt, err := r.db.Prepare(`
		SELECT name, color, ref_count FROM label WHERE tenant_id = ? AND id = ?`)
	if err != nil {
		return nil, err
	}
	var l entity.Label
	var name sql.NullString
	err = stmt.QueryRow(tenantID, id).Scan(&name, &l.Color, &l.RefCount)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	l.ID = id
	l.TenantID = tenantID
	l.Name = name.String
	return &l, nil
}

// Update a label
func (r *LabelMySQL) Update(e *entity.Label) error {
	_, err := r.db.Exec(`UPDATE label SET name = ?, color = ?, ref_count = ?
			WHERE tenant_id = ? AND id = ?`,
		e.Name, e.Color, e.RefCount, e.TenantID, e.ID)
	if err != nil {
		return err
	}
	return nil
}

// GetMulti retrieves multiple labels from the database
func (r *LabelMySQL) GetMulti(tenantID entity.ID, page, page_size int) ([]*entity.Label, error) {
	stmt, err := r.db.Prepare(`
		SELECT id, tenant_id, name, color, ref_count FROM label
			WHERE tenant_id = ?
			LIMIT ?,?`)
	if err != nil {
		return nil, err
	}
	var labels []*entity.Label
	rows, err := stmt.Query(tenantID, page*page_size, page_size)
	if err != nil {
		return nil, err
	}

	var name sql.NullString
	for rows.Next() {
		var l entity.Label
		err = rows.Scan(&l.ID, &l.TenantID, &name, &l.Color, &l.RefCount)
		if err != nil {
			return nil, err
		}
		l.TenantID = tenantID
		l.Name = name.String
		labels = append(labels, &l)
	}
	return labels, nil
}

// Delete a label
func (r *LabelMySQL) Delete(tenantID, id entity.ID) error {
	res, err := r.db.Exec(`DELETE FROM label WHERE tenant_id = ? AND id = ?`, tenantID, id)
	if err != nil {
		return err
	}

	if cnt, _ := res.RowsAffected(); cnt == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Get total labels
func (r *LabelMySQL) GetCount(tenantID entity.ID) (int, error) {
	stmt, err := r.db.Prepare(`SELECT count(*) FROM label WHERE tenant_id = ?`)
	if err != nil {
		return 0, err
	}

	var count int
	err = stmt.QueryRow(tenantID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
