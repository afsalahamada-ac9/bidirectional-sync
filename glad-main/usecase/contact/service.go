/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package contact

import (
	"sudhagar/glad/entity"
)

// Service contact usecase
type Service struct {
	repo Repository
}

// NewService create new service
func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

// CreateOrUpdate create a contact, if it doesn't exist or update the contact
func (s *Service) CreateOrUpdate(
	tenantID, accountID entity.ID,
	handle string,
	name string,
) error {
	contact, err := s.repo.GetByHandle(tenantID, accountID, handle)
	switch err {
	case entity.ErrNotFound:
		// create contact
		contact, err = entity.NewContact(tenantID, accountID, handle, name)
		if err == nil {
			err = s.repo.Create(contact)
		}
	case nil:
		// update
		contact.Name = name
		err = s.repo.Update(contact)
	default:
		return err
	}

	return err
}

// SetStaleByAccount sets all the contacts belonging to the handle as stale
func (s *Service) SetStaleByAccount(tenantID, accountID entity.ID) error {
	return s.repo.SetStaleByAccount(tenantID, accountID, true)
}

// ResetStaleByAccount sets all the contacts belonging to the handle as not stale
func (s *Service) ResetStaleByAccount(tenantID, accountID entity.ID) error {
	return s.repo.SetStaleByAccount(tenantID, accountID, false)
}

// DeleteStaleByAccount deletes all the contacts belonging to the handle as stale
func (s *Service) DeleteStaleByAccount(tenantID, accountID entity.ID) error {
	return s.repo.DeleteStaleByAccount(tenantID, accountID)
}

// Get gets a contact using the given identifier
func (s *Service) Get(tenantID, contactID entity.ID) (*entity.Contact, error) {
	return s.repo.GetByID(tenantID, contactID)
}

// Update updates the given contact
func (s *Service) Update(c *entity.Contact) error {
	return s.repo.Update(c)
}

// GetMulti get contacts in batches or pages
func (s *Service) GetMulti(tenantID entity.ID, page, page_size int) ([]*entity.Contact, error) {
	contacts, err := s.repo.GetMulti(tenantID, page, page_size)
	if err != nil {
		return nil, err
	}
	if len(contacts) == 0 {
		return nil, entity.ErrNotFound
	}
	return contacts, nil
}

// DeleteByAccount deletes all contact(s) for the given handle
func (s *Service) DeleteByAccount(tenantID, accountID entity.ID) error {
	return s.repo.DeleteByAccount(tenantID, accountID)
}

// GetCount gets total contact count
func (s *Service) GetCount(tenantID entity.ID) (int, error) {
	count, err := s.repo.GetCount(tenantID)
	if err != nil {
		return 0, err
	}

	return count, err
}
