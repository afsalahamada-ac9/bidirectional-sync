/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package presenter

import (
	"sudhagar/glad/entity"
)

// Account data - TenantID is returned in the HTTP header (may be not, as account is global?)
// X-GLAD-TenantID
type Account struct {
	ID entity.ID `json:"id"`
	// TenantID entity.ID            `json:"tenantId"`
	ExtID     string             `json:"extId"`
	Username  string             `json:"username"`
	FirstName string             `json:"firstName"`
	LastName  string             `json:"lastName"`
	Phone     string             `json:"phone"`
	Email     string             `json:"email"`
	Type      entity.AccountType `json:"type"`
}
