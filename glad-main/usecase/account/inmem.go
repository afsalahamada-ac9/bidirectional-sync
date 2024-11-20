/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package account

import (
	"sudhagar/glad/entity"
)

// inmem in memory repo
type inmem struct {
	m map[entity.ID]*entity.Account
}

// newInmem create new repository
func newInmem() *inmem {
	var m = map[entity.ID]*entity.Account{}
	return &inmem{
		m: m,
	}
}

// Create an account
func (r *inmem) Create(e *entity.Account) error {
	r.m[e.ID] = e
	return nil
}

// Get an account
func (r *inmem) Get(username string) (*entity.Account, error) {
	for _, j := range r.m {
		if j.Username == username {
			return r.m[j.ID], nil
		}
	}

	return nil, entity.ErrNotFound
}

// Update an account
func (r *inmem) Update(e *entity.Account) error {
	account := r.m[e.ID]
	if account == nil {
		return entity.ErrNotFound
	}

	account.ExtID = e.ExtID
	account.Username = e.Username
	account.Type = e.Type
	account.CreatedAt = e.CreatedAt
	account.UpdatedAt = e.UpdatedAt

	r.m[e.ID] = account
	return nil
}

// List accounts
func (r *inmem) List(tenantID entity.ID) ([]*entity.Account, error) {
	var d []*entity.Account
	for _, j := range r.m {
		// TenantID check removed
		d = append(d, j)
	}
	return d, nil
}

// Delete an account
func (r *inmem) Delete(username string) error {
	account, err := r.Get(username)
	if err != nil {
		return err
	}

	r.m[account.ID] = nil
	delete(r.m, account.ID)
	return nil
}

// GetCount gets total accounts for a given tenant
func (r *inmem) GetCount(tenantID entity.ID) (int, error) {
	count := 0
	for _, _ = range r.m {
		// TenantID check removed
		count++
	}
	return count, nil
}
