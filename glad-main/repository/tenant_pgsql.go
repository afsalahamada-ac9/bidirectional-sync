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

// TenantPGSQL mysql repo
type TenantPGSQL struct {
	db *sql.DB
}

// NewTenantPGSQL create new repository
func NewTenantPGSQL(db *sql.DB) *TenantPGSQL {
	return &TenantPGSQL{
		db: db,
	}
}

// Create a Tenant
func (r *TenantPGSQL) Create(e *entity.Tenant) (entity.ID, error) {

	stmt, err := r.db.Prepare(`
		INSERT INTO tenant (id, name, country, created_at) 
		VALUES($1, $2, $3, $4)`)
	if err != nil {
		return e.ID, err
	}
	_, err = stmt.Exec(
		e.ID,
		e.Name,
		e.Country,
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

// Get a Tenant
func (r *TenantPGSQL) Get(id entity.ID) (*entity.Tenant, error) {
	stmt, err := r.db.Prepare(`
		SELECT id, name, country, created_at FROM tenant WHERE id = $1;`)
	if err != nil {
		return nil, err
	}
	var t entity.Tenant
	var token sql.NullString

	err = stmt.QueryRow(id).Scan(&t.ID, &t.Name, &t.Country, &t.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	t.AuthToken = token.String
	return &t, nil
}

// Get a Tenant by username
func (r *TenantPGSQL) GetByName(name string) (*entity.Tenant, error) {
	stmt, err := r.db.Prepare(`
		SELECT id, name, country, created_at FROM tenant WHERE username = $1;`)
	if err != nil {
		return nil, err
	}
	var t entity.Tenant
	var token sql.NullString

	err = stmt.QueryRow(name).Scan(&t.ID, &t.Name, &t.Country, &t.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	t.AuthToken = token.String
	return &t, nil
}

// Update a Tenant
func (r *TenantPGSQL) Update(e *entity.Tenant) error {
	e.UpdatedAt = time.Now()
	_, err := r.db.Exec(`UPDATE tenant SET name = $1, country = $2, updated_at = $3 WHERE id = $4;`,
		e.Name, e.Country, e.UpdatedAt.Format("2006-01-02"), e.ID)
	if err != nil {
		return err
	}
	return nil
}

// List Tenants
func (r *TenantPGSQL) List() ([]*entity.Tenant, error) {
	stmt, err := r.db.Prepare(`SELECT id, name, country, created_at FROM tenant;`)
	if err != nil {
		return nil, err
	}
	var Tenants []*entity.Tenant
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var t entity.Tenant
		var country sql.NullString
		err = rows.Scan(&t.ID, &t.Name, &country, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		t.Country = country.String
		Tenants = append(Tenants, &t)
	}
	return Tenants, nil
}

// Delete a Tenant
func (r *TenantPGSQL) Delete(id entity.ID) error {
	res, err := r.db.Exec(`DELETE FROM tenant WHERE id = $1;`, id)
	if err != nil {
		return err
	}

	if cnt, _ := res.RowsAffected(); cnt == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Get total tenants
func (r *TenantPGSQL) GetCount() (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT count(*) FROM tenant;`).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
