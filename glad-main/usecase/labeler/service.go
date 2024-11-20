/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package labeler

import (
	"sudhagar/glad/entity"

	"sudhagar/glad/usecase/contact"
	"sudhagar/glad/usecase/label"
)

// Service labeler usecase
type Service struct {
	contactService contact.UseCase
	labelService   label.UseCase
}

// NewService create new service
func NewService(c contact.UseCase, l label.UseCase) *Service {
	return &Service{
		contactService: c,
		labelService:   l,
	}
}

// SetLabel sets the label for the given contact
func (s *Service) SetLabel(
	tenantID, contactID, labelID entity.ID,
) error {

	// Get contact for the given tenant. If not found, return error (contact) not found
	c, err := s.contactService.Get(tenantID, contactID)
	if err != nil {
		return err
	}

	// Repeat for Label
	l, err := s.labelService.Get(tenantID, labelID)
	if err != nil {
		return err
	}

	// Update Contact with the Label ID
	// Add label does the validation for adding too many labels as well
	if err = c.AddLabel(labelID); err != nil {
		return err
	}

	// validate entity before persisting to the repository (database)
	if err = c.Validate(); err != nil {
		return err
	}

	// Increment label reference (ref) count
	l.RefCount += 1

	// persistent contacts and labels
	// TODO: Make this a transaction, so if label service update fails, then we can roll back
	if err = s.contactService.Update(c); err != nil {
		return err
	}

	if err = s.labelService.Update(l); err != nil {
		return err
	}

	return nil
}

// RemoveLabel removes the label for the given contact
func (s *Service) RemoveLabel(
	tenantID, contactID, labelID entity.ID,
) error {

	// Get contact for the given tenant. If not found, return error (contact) not found
	c, err := s.contactService.Get(tenantID, contactID)
	if err != nil {
		return err
	}

	// Repeat for Label
	l, err := s.labelService.Get(tenantID, labelID)
	if err != nil {
		return err
	}

	// Remove label for the Contact
	// Since we are removing the label, no need to check whether too many labels are applied.
	if err = c.RemoveLabel(labelID); err != nil {
		return err
	}

	// validate entity before persisting to the repository (database)
	if err = c.Validate(); err != nil {
		return err
	}

	// Decrement label reference (ref) count
	if l.RefCount > 0 {
		l.RefCount -= 1
	}

	// persistent contacts and labels
	// TODO: Make this a transaction, so if label service update fails, then we can roll back
	if err = s.contactService.Update(c); err != nil {
		return err
	}

	if err = s.labelService.Update(l); err != nil {
		return err
	}

	return nil
}
