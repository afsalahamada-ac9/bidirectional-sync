/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package labeler

import (
	"sudhagar/glad/entity"
)

// UseCase interface
type UseCase interface {
	SetLabel(tenantID entity.ID, contactID entity.ID, labelID entity.ID) error
	RemoveLabel(tenantID entity.ID, contactID entity.ID, labelID entity.ID) error
}
