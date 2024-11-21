/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package middleware

import (
	"net/http"
	"sudhagar/glad/config"
	"sudhagar/glad/pkg/common"
	"sudhagar/glad/pkg/util"
)

// AddDefaultTenant Adds default tenant identifier
func AddDefaultTenant(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if r.Header.Get(common.HttpHeaderTenantID) == "" {
		r.Header.Set(common.HttpHeaderTenantID, util.GetStrEnvOrConfig("DEFAULT_TENANT", config.DEFAULT_TENANT))
	}
	next(w, r)
}
