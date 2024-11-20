/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package center

import (
	"strings"

	"sudhagar/glad/entity"
)

// inmem in memory repo
type inmem struct {
	m map[entity.ID]*entity.Center
}

// newInmem create new repository
func newInmem() *inmem {
	var m = map[entity.ID]*entity.Center{}
	return &inmem{
		m: m,
	}
}

// Create a center
func (r *inmem) Create(e *entity.Center) (entity.ID, error) {
	r.m[e.ID] = e
	return e.ID, nil
}

// Get a center
func (r *inmem) Get(id entity.ID) (*entity.Center, error) {
	if r.m[id] == nil {
		return nil, entity.ErrNotFound
	}
	return r.m[id], nil
}

// Update a center
func (r *inmem) Update(e *entity.Center) error {
	_, err := r.Get(e.ID)
	if err != nil {
		return err
	}
	r.m[e.ID] = e
	return nil
}

// Search centers
func (r *inmem) Search(tenantID entity.ID,
	query string,
) ([]*entity.Center, error) {
	var d []*entity.Center
	for _, j := range r.m {
		if j.TenantID == tenantID &&
			strings.Contains(strings.ToLower(j.Name), query) {
			d = append(d, j)
		}
	}
	return d, nil
}

// List centers
func (r *inmem) List(tenantID entity.ID) ([]*entity.Center, error) {
	var d []*entity.Center
	for _, j := range r.m {
		if j.TenantID == tenantID {
			d = append(d, j)
		}
	}
	return d, nil
}

// Delete a center
func (r *inmem) Delete(id entity.ID) error {
	if r.m[id] == nil {
		return entity.ErrNotFound
	}
	r.m[id] = nil
	delete(r.m, id)
	return nil
}

// GetCount gets total centers for a given tenant
func (r *inmem) GetCount(tenantID entity.ID) (int, error) {
	count := 0
	for _, j := range r.m {
		if j.TenantID == tenantID {
			count++
		}
	}
	return count, nil
}
