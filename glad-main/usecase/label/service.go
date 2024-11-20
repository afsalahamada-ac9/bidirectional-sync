/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package label

import (
	"sudhagar/glad/entity"
)

// Service label usecase
type Service struct {
	repo Repository
}

// NewService create new service
func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

// Create create a label
func (s *Service) Create(
	tenantID entity.ID,
	name string,
	color uint32,
) (entity.ID, error) {
	// create label
	label, err := entity.NewLabel(tenantID, name, color)
	if err == nil {
		err = s.repo.Create(label)
	}

	return label.ID, err
}

// Update updates a label
func (s *Service) Update(label *entity.Label) error {
	return s.repo.Update(label)
}

// TODO: DeleteByTenant deletes all label(s) for the given tenant
// Delete deletes the label
func (s *Service) Delete(tenantID, id entity.ID) error {
	// TODO: Check reference count and do not allow to delete the label
	// if reference count is not zero (OR, may be remove all the references in label_contact)
	// On a second thought, latter might be better with a prompt to user at the UI when
	// they attempt to delete the label with non-zero reference count
	return s.repo.Delete(tenantID, id)
}

// Get retrieves the label from the storage
func (s *Service) Get(tenantID, id entity.ID) (*entity.Label, error) {
	return s.repo.Get(tenantID, id)
}

// GetMulti get labels in batches or pages
func (s *Service) GetMulti(tenantID entity.ID, page, page_size int) ([]*entity.Label, error) {
	labels, err := s.repo.GetMulti(tenantID, page, page_size)
	if err != nil {
		return nil, err
	}
	if len(labels) == 0 {
		return nil, entity.ErrNotFound
	}
	return labels, nil
}

// GetCount gets total label count
func (s *Service) GetCount(tenantID entity.ID) (int, error) {
	count, err := s.repo.GetCount(tenantID)
	if err != nil {
		return 0, err
	}

	return count, err
}
