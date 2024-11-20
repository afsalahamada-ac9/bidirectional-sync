/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package label

import (
	"sudhagar/glad/entity"
)

// Reader interface
type Reader interface {
	Get(tenantID, id entity.ID) (*entity.Label, error)
	GetMulti(tenantID entity.ID, page, page_size int) ([]*entity.Label, error)
	GetCount(tenantId entity.ID) (int, error)
	// TODO: Maybe implement Search by name for the given tenant
}

// Writer interface
type Writer interface {
	Create(e *entity.Label) error
	Update(e *entity.Label) error
	Delete(tenantID, id entity.ID) error
}

// Repository interface
type Repository interface {
	Reader
	Writer
}

// UseCase interface
type UseCase interface {
	Create(
		tenantID entity.ID,
		name string,
		color uint32) (entity.ID, error)

	Update(e *entity.Label) error
	Delete(tenantID, id entity.ID) error

	Get(tenantID, id entity.ID) (*entity.Label, error)
	GetMulti(tenantID entity.ID, page, page_size int) ([]*entity.Label, error)
	GetCount(tenantId entity.ID) (int, error)
}
