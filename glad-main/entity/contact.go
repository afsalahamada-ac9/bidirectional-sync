/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package entity

// Contact data
type Contact struct {
	ID        ID
	TenantID  ID
	AccountID ID
	Handle    string
	Name      string

	// no meta data such as create/update time are required
	IsStale bool

	Labels []ID
}

const (
	maxLabelsPerContact int = 16
)

// NewContact create a new Contact
func NewContact(tenantID ID, accountID ID, handle, name string) (*Contact, error) {
	c := &Contact{
		ID:        NewID(),
		TenantID:  tenantID,
		AccountID: accountID,
		Handle:    handle,
		Name:      name,
		IsStale:   false,
	}
	err := c.Validate()
	if err != nil {
		return nil, ErrInvalidEntity
	}
	return c, nil
}

// Validate validate Contact
func (c *Contact) Validate() error {
	if c.Name == "" || c.Handle == "" {
		return ErrInvalidEntity
	}

	if len(c.Labels) > maxLabelsPerContact {
		return ErrTooManyLabels
	}

	return nil
}

// AddLabel adds label to the contact
func (c *Contact) AddLabel(id ID) error {
	if err := c.GetLabel(id); err == nil {
		return ErrLabelAlreadySet
	}

	if len(c.Labels)+1 > maxLabelsPerContact {
		return ErrTooManyLabels
	}

	c.Labels = append(c.Labels, id)
	return nil
}

// RemoveLabel removes label from the contact
func (c *Contact) RemoveLabel(id ID) error {
	for i, v := range c.Labels {
		if v == id {
			c.Labels = append(c.Labels[:i], c.Labels[i+1:]...)
			return nil
		}
	}

	return ErrNotFound
}

// GetLabel gets the label from the contact
func (c *Contact) GetLabel(id ID) error {
	for _, v := range c.Labels {
		if v == id {
			return nil
		}
	}

	return ErrNotFound
}
