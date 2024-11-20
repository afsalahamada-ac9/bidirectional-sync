/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package entity

// Label data
type Label struct {
	ID       ID
	TenantID ID
	Name     string
	Color    uint32 // store the actual color itself (RGB)

	RefCount uint32 // Number of contacts using this label
	// no meta data such as create/update time are required
}

// NewLabel create a new Label
func NewLabel(tenantID ID, name string, color uint32) (*Label, error) {
	t := &Label{
		ID:       NewID(),
		TenantID: tenantID,
		Name:     name,
		Color:    color,
	}
	err := t.Validate()
	if err != nil {
		return nil, ErrInvalidEntity
	}
	return t, nil
}

// Validate validate Label
func (t *Label) Validate() error {
	if t.Name == "" {
		return ErrInvalidEntity
	}
	return nil
}
