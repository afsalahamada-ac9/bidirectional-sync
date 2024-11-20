/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package tenant

import (
	"time"

	"sudhagar/glad/entity"
)

// Service tenant usecase
type Service struct {
	repo Repository
}

// NewService create new service
func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

// CreateTenant create a tenant
func (s *Service) CreateTenant(name, country string) (entity.ID, error) {
	// TODO: Check whether tenant already exists with same name

	t, err := entity.NewTenant(name, country)
	if err != nil {
		return t.ID, err
	}
	return s.repo.Create(t)
}

// GetTenant get a tenant
func (s *Service) GetTenant(id entity.ID) (*entity.Tenant, error) {
	t, err := s.repo.Get(id)
	if t == nil {
		return nil, entity.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return t, nil
}

// GetTenant get a tenant by name
func (s *Service) GetTenantByName(name string) (*entity.Tenant, error) {
	t, err := s.repo.GetByName(name)
	if t == nil {
		return nil, entity.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return t, nil
}

// ListTenants list tenant
func (s *Service) ListTenants() ([]*entity.Tenant, error) {
	tenants, err := s.repo.List()
	if err != nil {
		return nil, err
	}
	if len(tenants) == 0 {
		return nil, entity.ErrNotFound
	}
	return tenants, nil
}

// DeleteTenant Delete a tenant
func (s *Service) DeleteTenant(id entity.ID) error {
	t, err := s.GetTenant(id)
	if t == nil {
		return entity.ErrNotFound
	}
	if err != nil {
		return err
	}

	return s.repo.Delete(id)
}

// UpdateTenant Update a tenant
func (s *Service) UpdateTenant(t *entity.Tenant) error {
	// retrieve and fill in empty values for mandatory fields
	if t.Country == "" {
		current, err := s.GetTenant(t.ID)
		if err != nil {
			return err
		}

		if t.Country == "" {
			t.Country = current.Country
		}

		if t.AuthToken == "" {
			t.AuthToken = current.AuthToken
		}
	}

	err := t.Validate()
	if err != nil {
		return err
	}
	t.UpdatedAt = time.Now()
	return s.repo.Update(t)
}

// UNUSED: Login Validates credentials, generates token and update the DB
func (s *Service) Login(name, password string) (*entity.Tenant, error) {
	// Get tenant by name
	t, err := s.GetTenantByName(name)
	if err != nil {
		return nil, err
	}

	// // Validate credentials
	// if t.ValidatePassword(password) != nil {
	// 	return nil, entity.ErrAuthFailure
	// }

	// // Generate token
	// if t.GenToken() != nil {
	// 	return nil, entity.ErrCreateToken
	// }

	// Update tenant: store token to database
	if err = s.UpdateTenant(t); err != nil {
		return nil, err
	}

	return t, nil
}

// GetCount gets total tenant count
func (s *Service) GetCount() int {
	count, err := s.repo.GetCount()
	if err != nil {
		return 0
	}

	return count
}
