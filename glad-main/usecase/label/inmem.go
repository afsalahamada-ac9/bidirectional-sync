/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package label

import (
	"sudhagar/glad/entity"
)

// inmem in memory repo
type inmem struct {
	m map[entity.ID]*entity.Label
}

// newInmem create new repository
func newInmem() *inmem {
	var m = map[entity.ID]*entity.Label{}
	return &inmem{
		m: m,
	}
}

// Create a label
func (r *inmem) Create(e *entity.Label) error {
	r.m[e.ID] = e
	return nil
}

// Update a label
func (r *inmem) Update(e *entity.Label) error {
	label := r.m[e.ID]
	if label == nil {
		return entity.ErrNotFound
	}

	label.TenantID = e.TenantID
	label.Name = e.Name
	label.Color = e.Color

	r.m[e.ID] = label
	return nil
}

// Delete deletes label for the given label id.
func (r *inmem) Delete(tenantID, id entity.ID) error {
	if label, ok := r.m[id]; ok {
		if label.TenantID == tenantID {
			delete(r.m, id)
			return nil
		}
	}

	return entity.ErrNotFound
}

// Get a label
func (r *inmem) Get(tenantID, id entity.ID) (*entity.Label, error) {
	if label, ok := r.m[id]; ok {
		if label.TenantID == tenantID {
			return label, nil
		}
	}

	return nil, entity.ErrNotFound
}

// Get labels in batches or pages
func (r *inmem) GetMulti(tenantID entity.ID, page, page_size int) ([]*entity.Label, error) {
	var d []*entity.Label

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

// GetCount gets total labels for a given tenant
func (r *inmem) GetCount(tenantID entity.ID) (int, error) {
	count := 0
	for _, j := range r.m {
		if j.TenantID == tenantID {
			count++
		}
	}
	return count, nil
}
