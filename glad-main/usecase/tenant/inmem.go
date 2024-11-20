/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package tenant

import (
	"sudhagar/glad/entity"
)

// inmem in memory repo
type inmem struct {
	m map[entity.ID]*entity.Tenant
}

// newInmem create new repository
func newInmem() *inmem {
	var m = map[entity.ID]*entity.Tenant{}
	return &inmem{
		m: m,
	}
}

// Create a tenant
func (r *inmem) Create(e *entity.Tenant) (entity.ID, error) {
	r.m[e.ID] = e
	return e.ID, nil
}

// Get a tenant
func (r *inmem) Get(id entity.ID) (*entity.Tenant, error) {
	if r.m[id] == nil {
		return nil, entity.ErrNotFound
	}
	return r.m[id], nil
}

// Get a tenant by name
func (r *inmem) GetByName(name string) (*entity.Tenant, error) {
	for _, j := range r.m {
		if j.Name == name {
			return j, nil
		}
	}
	return nil, entity.ErrNotFound
}

// Update a tenant
func (r *inmem) Update(e *entity.Tenant) error {
	_, err := r.Get(e.ID)
	if err != nil {
		return err
	}
	r.m[e.ID] = e
	return nil
}

// List tenants
func (r *inmem) List() ([]*entity.Tenant, error) {
	var d []*entity.Tenant
	for _, j := range r.m {
		d = append(d, j)
	}
	return d, nil
}

// Delete a tenant
func (r *inmem) Delete(id entity.ID) error {
	if r.m[id] == nil {
		return entity.ErrNotFound
	}
	r.m[id] = nil
	delete(r.m, id)
	return nil
}

// GetCount gets total tenants
func (r *inmem) GetCount() (int, error) {
	return len(r.m), nil
}
