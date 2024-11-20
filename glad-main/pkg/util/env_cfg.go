/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package util

import (
	"os"
	"strconv"
)

func GetStrEnvOrConfig(key string, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func GetIntEnvOrConfig(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if port, err := strconv.Atoi(value); err == nil {
			return port
		}
	}
	return fallback
}
