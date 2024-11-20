/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package presenter

import (
	"sudhagar/glad/entity"
)

// Course data - TenantID is returned in the HTTP header
// X-GLAD-TenantID
type Course struct {
	ID           entity.ID             `json:"id"`
	TenantID     entity.ID             `json:"tenantId"`
	ExtID        string                `json:"extId"`
	CenterID     entity.ID             `json:"centerId"`
	Name         string                `json:"name"`
	Notes        string                `json:"notes"`
	Timezone     string                `json:"timezone"`
	Location     entity.CourseLocation `json:"location"` // TODO: To be defined in here in presenter
	Status       entity.CourseStatus   `json:"status"`
	CType        entity.CourseType     `json:"ctype"`
	MaxAttendees int32                 `json:"maxAttendees"`
	NumAttendees int32                 `json:"numAttendees"`
}
