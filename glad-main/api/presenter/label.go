/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package presenter

import (
	"sudhagar/glad/entity"
)

// Label data - TenantID is returned in the HTTP header
// X-GLAD-TenantID
type Label struct {
	ID       entity.ID `json:"id"`
	TenantID entity.ID `json:"tenantId"`
	Name     string    `json:"name"`
	Color    uint32    `json:"color"`
}
