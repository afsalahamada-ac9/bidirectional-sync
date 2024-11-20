/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package repository

import (
	"database/sql"
	"time"

	"sudhagar/glad/entity"
)

// CenterPGSQL mysql repo
type CenterPGSQL struct {
	db *sql.DB
}

// NewCenterPGSQL create new repository
func NewCenterPGSQL(db *sql.DB) *CenterPGSQL {
	return &CenterPGSQL{
		db: db,
	}
}

// Create creates a center
func (r *CenterPGSQL) Create(e *entity.Center) (entity.ID, error) {
	stmt, err := r.db.Prepare(`
		INSERT INTO center (id, tenant_id, ext_id, name, location, geo_location, capacity, mode, webpage, is_national_center, created_at) 
		VALUES( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`)
	if err != nil {
		return e.ID, err
	}
	_, err = stmt.Exec(
		e.ID,
		e.TenantID,
		e.ExtID,
		e.Name,
		e.Location,    // TODO: to be converted into json
		e.GeoLocation, // TODO: to be converted into json
		e.Capacity,
		int(e.Mode),
		e.WebPage,
		e.IsNationalCenter,
		time.Now().Format("2006-01-02"),
	)
	if err != nil {
		return e.ID, err
	}
	err = stmt.Close()
	if err != nil {
		return e.ID, err
	}
	return e.ID, nil
}

// Get retrieves a center
// Not all fields are required for v1
func (r *CenterPGSQL) Get(id entity.ID) (*entity.Center, error) {
	stmt, err := r.db.Prepare(`
		SELECT id, tenant_id, ext_id, name, mode, created_at FROM center WHERE id = $1;`)
	if err != nil {
		return nil, err
	}
	var c entity.Center
	var ext_id sql.NullString
	var name sql.NullString
	err = stmt.QueryRow(id).Scan(&c.ID, &c.TenantID, &ext_id, &name, &c.Mode, &c.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	c.ExtID = ext_id.String
	c.Name = name.String

	return &c, nil
}

// Update updates a center
func (r *CenterPGSQL) Update(e *entity.Center) error {
	e.UpdatedAt = time.Now()
	_, err := r.db.Exec(`
		UPDATE center SET name = $1, mode = $2, updated_at = $3 WHERE id = $4;`,
		e.Name, int(e.Mode), e.UpdatedAt.Format("2006-01-02"), e.ID)
	if err != nil {
		return err
	}
	return nil
}

// Search searches centers
func (r *CenterPGSQL) Search(tenantID entity.ID,
	query string,
) ([]*entity.Center, error) {
	stmt, err := r.db.Prepare(`
		SELECT id, tenant_id, ext_id, name, mode, created_at FROM center
		WHERE tenant_id = $1 AND name LIKE $2;`)
	if err != nil {
		return nil, err
	}
	var centers []*entity.Center
	rows, err := stmt.Query(tenantID, "%"+query+"%")
	if err != nil {
		return nil, err
	}

	var ext_id sql.NullString
	var name sql.NullString
	for rows.Next() {
		var c entity.Center
		err = rows.Scan(&c.ID, &c.TenantID, &ext_id, &name, &c.Mode, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		c.ExtID = ext_id.String
		c.Name = name.String
		centers = append(centers, &c)
	}

	return centers, nil
}

// List lists centers
func (r *CenterPGSQL) List(tenantID entity.ID) ([]*entity.Center, error) {
	stmt, err := r.db.Prepare(`
		SELECT id, tenant_id, ext_id, name, mode, created_at FROM center WHERE tenant_id = $1;`)
	if err != nil {
		return nil, err
	}
	var centers []*entity.Center
	rows, err := stmt.Query(tenantID)
	if err != nil {
		return nil, err
	}

	var ext_id sql.NullString
	var name sql.NullString
	for rows.Next() {
		var c entity.Center
		err = rows.Scan(&c.ID, &c.TenantID, &ext_id, &name, &c.Mode, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		c.ExtID = ext_id.String
		c.Name = name.String
		centers = append(centers, &c)
	}
	return centers, nil
}

// Delete deletes a center
func (r *CenterPGSQL) Delete(id entity.ID) error {
	res, err := r.db.Exec(`DELETE FROM center WHERE id = $1;`, id)
	if err != nil {
		return err
	}

	if cnt, _ := res.RowsAffected(); cnt == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Get total centers
func (r *CenterPGSQL) GetCount(tenantID entity.ID) (int, error) {
	stmt, err := r.db.Prepare(`SELECT count(*) FROM center WHERE tenant_id = $1;`)
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
