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

// AccountPGSQL mysql repo
type AccountPGSQL struct {
	db *sql.DB
}

// NewAccountPGSQL create new repository
func NewAccountPGSQL(db *sql.DB) *AccountPGSQL {
	return &AccountPGSQL{
		db: db,
	}
}

// Create a account
func (r *AccountPGSQL) Create(e *entity.Account) error {
	stmt, err := r.db.Prepare(`
		INSERT INTO account (id, ext_id, username, first_name, last_name, phone, email, type, created_at) 
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		e.ID,
		e.ExtID,
		e.Username,
		e.FirstName,
		e.LastName,
		e.Phone,
		e.Email,
		int(e.Type),
		time.Now().Format("2006-01-02"),
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

// TODO: GetByID()
// TODO: Need a filter by TenantID if/when TenantID support is added.
// Accounts are global in nature, but for storage purposes they will be assigned to some tenants.
// Get a account
func (r *AccountPGSQL) Get(username string) (*entity.Account, error) {
	stmt, err := r.db.Prepare(`
		SELECT id, ext_id, type, created_at FROM account WHERE username = $1;`)
	if err != nil {
		return nil, err
	}
	var t entity.Account
	err = stmt.QueryRow(username).Scan(&t.ID, &t.ExtID, &t.Type, &t.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	t.Username = username
	return &t, nil
}

// TODO: UpdateByID()
// Update a account
func (r *AccountPGSQL) Update(e *entity.Account) error {
	e.UpdatedAt = time.Now()
	_, err := r.db.Exec(`UPDATE account SET username = $1, type = $2, updated_at = $3 WHERE id = $4;`,
		e.Username, int(e.Type), e.UpdatedAt.Format("2006-01-02"), e.ID)
	if err != nil {
		return err
	}
	return nil
}

// List accounts
func (r *AccountPGSQL) List(tenantID entity.ID) ([]*entity.Account, error) {
	stmt, err := r.db.Prepare(`
		SELECT id, ext_id, username, type, created_at FROM account;`)
	if err != nil {
		return nil, err
	}
	var accounts []*entity.Account
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	var username sql.NullString
	for rows.Next() {
		var t entity.Account
		err = rows.Scan(&t.ID, &t.ExtID, &username, &t.Type, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		t.Username = username.String
		accounts = append(accounts, &t)
	}
	return accounts, nil
}

// TODO: DeleteByID()
// Delete a account
func (r *AccountPGSQL) Delete(username string) error {
	res, err := r.db.Exec(`DELETE FROM account WHERE username = $1;`, username)
	if err != nil {
		return err
	}

	if cnt, _ := res.RowsAffected(); cnt == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Get total accounts
func (r *AccountPGSQL) GetCount(tenantID entity.ID) (int, error) {
	stmt, err := r.db.Prepare(`SELECT count(*) FROM account;`)
	if err != nil {
		return 0, err
	}

	var count int
	err = stmt.QueryRow().Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
