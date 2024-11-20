//go:build prod

/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package config

const (
	DB_USER                = "glad_user"
	DB_PASSWORD            = "glad1234"
	DB_DATABASE            = "glad"
	DB_HOST                = "127.0.0.1"
	DB_PORT                = 5432 /* 3306 for MySQL */
	DB_SSLMODE             = "verify-full"
	API_PORT               = 8080
	PROMETHEUS_PUSHGATEWAY = "http://localhost:9091/"
)
