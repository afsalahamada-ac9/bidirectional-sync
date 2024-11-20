/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package presenter

import (
	"sudhagar/glad/entity"
)

// Tenant data
type Tenant struct {
	ID      entity.ID `json:"id"`
	Name    string    `json:"name"`
	Country string    `json:"country"`
	// Do not return password
	// AuthToken is returned at login
	AuthToken string `json:"token,omitempty"`
}
