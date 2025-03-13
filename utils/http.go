package utils

import (
	"os"
	"strings"
)

// GetTrustedProxies returns list of trusted proxies.
func GetTrustedProxies() []string {
	return split("TRUSTED_PROXIES")
}

// GetAllowedOrigins returns list of allowed origins.
func GetAllowedOrigins() []string {
	return split("ALLOWED_ORIGINS")
}

func split(val string) []string {
	l := os.Getenv(val)
	if l == "" {
		return []string{}
	}
	return strings.Split(l, ",")
}
