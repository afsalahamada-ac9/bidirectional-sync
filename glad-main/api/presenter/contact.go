/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package presenter

import (
	"sudhagar/glad/entity"
)

// Contact data - TenantID is returned in the HTTP header
// X-GLAD-TenantID
type Contact struct {
	ID        entity.ID `json:"id"`
	TenantID  entity.ID `json:"tenantId"`
	AccountID entity.ID `json:"accountId"`
	Handle    string    `json:"handle"`
	Name      string    `json:"name"`
}
