/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package course

import (
	"strings"

	"sudhagar/glad/entity"
)

// inmem in memory repo
type inmem struct {
	m map[entity.ID]*entity.Course
}

// newInmem create new repository
func newInmem() *inmem {
	var m = map[entity.ID]*entity.Course{}
	return &inmem{
		m: m,
	}
}

// Create a course
func (r *inmem) Create(e *entity.Course) (entity.ID, error) {
	r.m[e.ID] = e
	return e.ID, nil
}

// Get a course
func (r *inmem) Get(id entity.ID) (*entity.Course, error) {
	if r.m[id] == nil {
		return nil, entity.ErrNotFound
	}
	return r.m[id], nil
}

// Update a course
func (r *inmem) Update(e *entity.Course) error {
	_, err := r.Get(e.ID)
	if err != nil {
		return err
	}
	r.m[e.ID] = e
	return nil
}

// Search courses
func (r *inmem) Search(tenantID entity.ID,
	query string,
) ([]*entity.Course, error) {
	var d []*entity.Course
	for _, j := range r.m {
		if j.TenantID == tenantID &&
			strings.Contains(strings.ToLower(j.Name), query) {
			d = append(d, j)
		}
	}
	return d, nil
}

// List courses
func (r *inmem) List(tenantID entity.ID) ([]*entity.Course, error) {
	var d []*entity.Course
	for _, j := range r.m {
		if j.TenantID == tenantID {
			d = append(d, j)
		}
	}
	return d, nil
}

// Delete a course
func (r *inmem) Delete(id entity.ID) error {
	if r.m[id] == nil {
		return entity.ErrNotFound
	}
	r.m[id] = nil
	delete(r.m, id)
	return nil
}

// GetCount gets total courses for a given tenant
func (r *inmem) GetCount(tenantID entity.ID) (int, error) {
	count := 0
	for _, j := range r.m {
		if j.TenantID == tenantID {
			count++
		}
	}
	return count, nil
}
