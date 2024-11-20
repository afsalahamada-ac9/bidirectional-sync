/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package tenant

import (
	"sudhagar/glad/entity"
)

// Reader interface
type Reader interface {
	Get(id entity.ID) (*entity.Tenant, error)
	GetByName(username string) (*entity.Tenant, error)
	List() ([]*entity.Tenant, error)
	GetCount() (int, error)
}

// Writer tenant writer
type Writer interface {
	Create(e *entity.Tenant) (entity.ID, error)
	Update(e *entity.Tenant) error
	Delete(id entity.ID) error
}

// Repository interface
type Repository interface {
	Reader
	Writer
}

// UseCase interface
type UseCase interface {
	GetTenant(id entity.ID) (*entity.Tenant, error)
	ListTenants() ([]*entity.Tenant, error)
	CreateTenant(username, country string) (entity.ID, error)
	UpdateTenant(e *entity.Tenant) error
	DeleteTenant(id entity.ID) error
	Login(username, password string) (*entity.Tenant, error)
	GetCount() int
	// Thoughts: Need to validate token; use tenant id and token to validate
}
