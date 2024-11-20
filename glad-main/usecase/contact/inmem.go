/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package contact

import (
	"sudhagar/glad/entity"
)

// inmem in memory repo
type inmem struct {
	m map[entity.ID]*entity.Contact
}

// newInmem create new repository
func newInmem() *inmem {
	var m = map[entity.ID]*entity.Contact{}
	return &inmem{
		m: m,
	}
}

// Create a contact
func (r *inmem) Create(e *entity.Contact) error {
	r.m[e.ID] = e
	return nil
}

// GetByID get a contact by ID
func (r *inmem) GetByID(tenantID, contactID entity.ID) (*entity.Contact, error) {
	for _, j := range r.m {
		if j.TenantID == tenantID && j.ID == contactID {
			return r.m[j.ID], nil
		}
	}

	return nil, entity.ErrNotFound
}

// GetByHandle get a contact by handle
func (r *inmem) GetByHandle(tenantID, accountID entity.ID, handle string) (*entity.Contact, error) {
	for _, j := range r.m {
		if j.TenantID == tenantID && j.AccountID == accountID && j.Handle == handle {
			return r.m[j.ID], nil
		}
	}

	return nil, entity.ErrNotFound
}

// Update a contact
func (r *inmem) Update(e *entity.Contact) error {
	contact := r.m[e.ID]
	if contact == nil {
		return entity.ErrNotFound
	}

	contact.TenantID = e.TenantID
	contact.AccountID = e.AccountID
	contact.Handle = e.Handle
	contact.Name = e.Name

	r.m[e.ID] = contact
	return nil
}

// Set all contacts for the given tenant and the account id to the given stale value
func (r *inmem) SetStaleByAccount(tenantID, accountID entity.ID, stale bool) error {
	for _, j := range r.m {
		if j.TenantID == tenantID && j.AccountID == accountID {
			j.IsStale = stale
		}
	}

	return nil
}

// Delete all stale contacts for the given tenant and the account id
func (r *inmem) DeleteStaleByAccount(tenantID, accountID entity.ID) error {
	var staleIDs []entity.ID
	for _, j := range r.m {
		if j.TenantID == tenantID && j.AccountID == accountID && j.IsStale {
			staleIDs = append(staleIDs, j.ID)
		}
	}

	for _, id := range staleIDs {
		delete(r.m, id)
	}

	return nil
}

// Get contacts in batches or pages
func (r *inmem) GetMulti(tenantID entity.ID, page, page_size int) ([]*entity.Contact, error) {
	var d []*entity.Contact

	count := 0
	for _, j := range r.m {
		if j.TenantID == tenantID {
			if count < page*page_size || count >= (page+1)*page_size {
				break
			}

			d = append(d, j)
			count++
		}
	}
	return d, nil
}

// Deletes all contact(s) for the given account id
func (r *inmem) DeleteByAccount(tenantID, accountID entity.ID) error {
	var IDs []entity.ID
	for _, j := range r.m {
		if j.TenantID == tenantID && j.AccountID == accountID {
			IDs = append(IDs, j.ID)
		}
	}

	for _, id := range IDs {
		delete(r.m, id)
	}

	return nil
}

// GetCount gets total contacts for a given tenant
func (r *inmem) GetCount(tenantID entity.ID) (int, error) {
	count := 0
	for _, j := range r.m {
		if j.TenantID == tenantID {
			count++
		}
	}
	return count, nil
}
