/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package contact

import (
	"sudhagar/glad/entity"
)

// Reader interface
type Reader interface {
	GetByHandle(tenantID entity.ID, accountID entity.ID, handle string) (*entity.Contact, error)
	GetMulti(tenantID entity.ID, page, page_size int) ([]*entity.Contact, error)
	GetCount(tenantId entity.ID) (int, error)

	GetByID(tenantID, contactID entity.ID) (*entity.Contact, error)
}

// Writer interface
type Writer interface {
	Create(e *entity.Contact) error
	Update(e *entity.Contact) error
	SetStaleByAccount(tenantID entity.ID, accountID entity.ID, value bool) error
	DeleteStaleByAccount(tenantID entity.ID, accountID entity.ID) error
	DeleteByAccount(tenantID entity.ID, accountID entity.ID) error
}

// Repository interface
type Repository interface {
	Reader
	Writer
}

// UseCase interface
type UseCase interface {
	// TODO: Wondering whether ID for contact must be returned here
	// or Contact entity itself.
	// TODO: Need to change this interface to Create the contact
	// and use Update to update the contact
	CreateOrUpdate(
		tenantID entity.ID,
		accountID entity.ID,
		handle string,
		name string) error

	SetStaleByAccount(tenantID entity.ID, accountID entity.ID) error
	ResetStaleByAccount(tenantID entity.ID, accountID entity.ID) error
	DeleteStaleByAccount(tenantID entity.ID, accountID entity.ID) error

	Get(tenantID, contactID entity.ID) (*entity.Contact, error)
	Update(c *entity.Contact) error

	GetMulti(tenantID entity.ID, page, page_size int) ([]*entity.Contact, error)
	DeleteByAccount(tenantID entity.ID, accountID entity.ID) error
	GetCount(tenantId entity.ID) (int, error)
}
