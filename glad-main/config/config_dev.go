//go:build dev

/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package config

const (
	// Database
	DB_USER     = "glad_user"
	DB_PASSWORD = "glad1234"
	DB_DATABASE = "glad"
	DB_HOST     = "127.0.0.1"
	DB_PORT     = 5432 /* 3306 for MySQL */
	DB_SSLMODE  = "disable"

	// Defaults
	DEFAULT_TENANT = "5306526529902621696"

	// API port
	API_PORT = 8080

	// Metrics
	PROMETHEUS_PUSHGATEWAY = "http://localhost:9091/"
)
