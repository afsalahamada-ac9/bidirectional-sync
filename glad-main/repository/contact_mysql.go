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

// ContactMySQL mysql repo
type ContactMySQL struct {
	db *sql.DB
}

// NewContactMySQL create new repository
func NewContactMySQL(db *sql.DB) *ContactMySQL {
	return &ContactMySQL{
		db: db,
	}
}

// Create a contact
func (r *ContactMySQL) Create(e *entity.Contact) error {
	stmt, err := r.db.Prepare(`
		INSERT INTO contact (id, tenant_id, account_id, name, handle, is_stale) 
			VALUES(?,?,?,?,?, ?)`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		e.ID,
		e.TenantID,
		e.AccountID,
		e.Name,
		e.Handle,
		e.IsStale,
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

// GetByHandle retrieves a contact for a given tenant and account using the handle
func (r *ContactMySQL) GetByHandle(tenantID, accountID entity.ID, handle string) (*entity.Contact, error) {
	stmt, err := r.db.Prepare(`
		SELECT id FROM contact
			WHERE tenant_id = ? AND account_id = ? AND handle = ?`)
	if err != nil {
		return nil, err
	}

	var contactID entity.ID
	err = stmt.QueryRow(tenantID, accountID, handle).Scan(&contactID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return getContact(tenantID, contactID, r.db)
}

// GetByID retrieves a contact for the given tenant using the contact identifier
func (r *ContactMySQL) GetByID(tenantID, contactID entity.ID) (*entity.Contact, error) {
	return getContact(tenantID, contactID, r.db)
}

func getContact(tenantID, contactID entity.ID, db *sql.DB) (*entity.Contact, error) {
	stmt, err := db.Prepare(`
		SELECT name, account_id, handle, is_stale FROM contact
			WHERE tenant_id = ? AND id = ?`)
	if err != nil {
		return nil, err
	}

	// Get the contact
	var c entity.Contact
	var name sql.NullString
	var handle sql.NullString

	err = stmt.QueryRow(tenantID, contactID).Scan(&name, &c.AccountID, &handle, &c.IsStale)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	c.TenantID = tenantID
	c.Handle = handle.String
	c.Name = name.String

	// Get the labels
	stmt, err = db.Prepare(`
		SELECT label_id FROM label_contact
			WHERE tenant_id = ? AND contact_id = ?`)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(tenantID, contactID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var labelID entity.ID
		if err = rows.Scan(&labelID); err != nil {
			return nil, err
		}
		c.Labels = append(c.Labels, labelID)
	}
	return &c, nil
}

// Update a contact
func (r *ContactMySQL) Update(c *entity.Contact) error {
	// Update the contact
	_, err := r.db.Exec(`
		UPDATE contact SET account_id = ?, name = ?, handle = ?, is_stale = ?
			WHERE tenant_id = ? AND id = ?`,
		c.AccountID, c.Name, c.Handle, c.IsStale, c.TenantID, c.ID)
	if err != nil {
		return err
	}

	// TODO-perf: We could create a label hash and update the database only if the hash has changed

	// Update the labels
	_, err = r.db.Exec("DELETE FROM label_contact WHERE tenant_id = ? AND contact_id = ?", c.TenantID, c.ID)
	if err != nil {
		return err
	}

	// TODO-perf: https://stackoverflow.com/questions/21108084/how-to-insert-multiple-data-at-once

	for _, labelID := range c.Labels {
		_, err := r.db.Exec(`INSERT INTO label_contact (label_id, contact_id, tenant_id)
			VALUES(?,?,?)`,
			labelID, c.ID, c.TenantID)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetMulti gets multiple contacts
func (r *ContactMySQL) GetMulti(tenantID entity.ID, page, page_size int) ([]*entity.Contact, error) {
	stmt, err := r.db.Prepare(`
		SELECT id, account_id, name, handle, is_stale FROM contact
			WHERE tenant_id = ?
			LIMIT ?,?`)
	if err != nil {
		return nil, err
	}
	var contacts []*entity.Contact
	rows, err := stmt.Query(tenantID, page*page_size, page_size)
	if err != nil {
		return nil, err
	}

	var name sql.NullString
	var handle sql.NullString
	for rows.Next() {
		var t entity.Contact
		err = rows.Scan(&t.ID, &t.AccountID, &name, &handle, &t.IsStale)
		if err != nil {
			return nil, err
		}
		t.TenantID = tenantID
		t.Name = name.String
		t.Handle = handle.String
		contacts = append(contacts, &t)
	}
	return contacts, nil
}

// SetStaleByAccount sets all the contacts for the given tenant and account as stale
func (r *ContactMySQL) SetStaleByAccount(tenantID, accountID entity.ID, value bool) error {
	_, err := r.db.Exec(`
		UPDATE contact SET is_stale = ?
			WHERE tenant_id = ? AND account_id = ?`,
		value, tenantID, accountID)
	if err != nil {
		return err
	}
	return nil
}

// DeleteStaleByAccount deletes all stale contacts for the given tenant and the account
func (r *ContactMySQL) DeleteStaleByAccount(tenantID, accountID entity.ID) error {
	res, err := r.db.Exec(`DELETE FROM contact
		WHERE tenant_id = ? AND account_id = ? AND is_stale = ?`,
		tenantID, accountID, true)
	if err != nil {
		return err
	}

	if cnt, _ := res.RowsAffected(); cnt == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// TODO: DeleteByID()
// Delete a contact for the given tenant and the account
func (r *ContactMySQL) DeleteByAccount(tenantID, accountID entity.ID) error {
	res, err := r.db.Exec(`DELETE FROM contact
		WHERE tenant_id = ? AND account_id = ?`,
		tenantID, accountID)
	if err != nil {
		return err
	}

	if cnt, _ := res.RowsAffected(); cnt == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Get total contacts
func (r *ContactMySQL) GetCount(tenantID entity.ID) (int, error) {
	stmt, err := r.db.Prepare(`SELECT count(*) FROM contact WHERE tenant_id = ?`)
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
