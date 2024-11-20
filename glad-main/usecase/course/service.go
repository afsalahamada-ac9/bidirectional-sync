/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package course

import (
	"strings"
	"time"

	"sudhagar/glad/entity"
)

// Service course usecase
type Service struct {
	repo Repository
}

// NewService create new service
func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

// CreateCourse creates a course
func (s *Service) CreateCourse(tenantID entity.ID,
	extID string,
	centerID entity.ID,
	name, notes, timezone string,
	location entity.CourseLocation,
	status entity.CourseStatus,
	ctype entity.CourseType,
	maxAttendees, numAttendees int32,
	isAutoApprove bool,
) (entity.ID, error) {
	c, err := entity.NewCourse(tenantID, extID, centerID,
		name, notes, timezone,
		location, status, ctype,
		maxAttendees, numAttendees, isAutoApprove)
	if err != nil {
		return c.ID, err
	}
	return s.repo.Create(c)
}

// GetCourse retrieves a course
func (s *Service) GetCourse(id entity.ID) (*entity.Course, error) {
	t, err := s.repo.Get(id)
	if t == nil {
		return nil, entity.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return t, nil
}

// SearchCourses search course
func (s *Service) SearchCourses(tenantID entity.ID,
	query string,
) ([]*entity.Course, error) {
	courses, err := s.repo.Search(tenantID, strings.ToLower(query))
	if err != nil {
		return nil, err
	}
	if len(courses) == 0 {
		return nil, entity.ErrNotFound
	}
	return courses, nil
}

// ListCourses list course
func (s *Service) ListCourses(tenantID entity.ID) ([]*entity.Course, error) {
	courses, err := s.repo.List(tenantID)
	if err != nil {
		return nil, err
	}
	if len(courses) == 0 {
		return nil, entity.ErrNotFound
	}
	return courses, nil
}

// DeleteCourse Delete a course
func (s *Service) DeleteCourse(id entity.ID) error {
	t, err := s.GetCourse(id)
	if t == nil {
		return entity.ErrNotFound
	}
	if err != nil {
		return err
	}

	return s.repo.Delete(id)
}

// UpdateCourse Update a course
func (s *Service) UpdateCourse(c *entity.Course) error {
	err := c.Validate()
	if err != nil {
		return err
	}
	c.UpdatedAt = time.Now()
	return s.repo.Update(c)
}

// GetCount gets total course count
func (s *Service) GetCount(tenantID entity.ID) int {
	count, err := s.repo.GetCount(tenantID)
	if err != nil {
		return 0
	}

	return count
}
