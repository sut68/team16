package config

import (
	"os"
	"time"
)

var (
	// Cookie Config
	CookieDomain   = os.Getenv("COOKIE_DOMAIN")           // Set in production via env or build args
	AccessTTL      = 24 * 60 * 60 * time.Second           // 24 hours
	CookieSecure   = os.Getenv("COOKIE_SECURE") == "true" // Set to true in production with HTTPS
	CookieSameSite = "lax"
)
