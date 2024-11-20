/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package center

import (
	"sudhagar/glad/entity"
)

// Reader interface
type Reader interface {
	Get(id entity.ID) (*entity.Center, error)
	Search(tenantID entity.ID, query string) ([]*entity.Center, error)
	List(tenantID entity.ID) ([]*entity.Center, error)
	GetCount(id entity.ID) (int, error)
}

// Writer center writer
type Writer interface {
	Create(e *entity.Center) (entity.ID, error)
	Update(e *entity.Center) error
	Delete(id entity.ID) error
}

// Repository interface
type Repository interface {
	Reader
	Writer
}

// UseCase interface
type UseCase interface {
	GetCenter(id entity.ID) (*entity.Center, error)
	SearchCenters(tenantID entity.ID, query string) ([]*entity.Center, error)
	ListCenters(tenantID entity.ID) ([]*entity.Center, error)
	CreateCenter(tenantID entity.ID, extID, name string, mode entity.CenterMode) (entity.ID, error)
	UpdateCenter(e *entity.Center) error
	DeleteCenter(id entity.ID) error
	GetCount(id entity.ID) int
}
