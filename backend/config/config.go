package config

import (
	"os"
	"strconv"
	"time"
)

var (
	JWTSecret		string
	AccessTTL		time.Duration
	CookieDomain	string
	CookieSecure	bool
	CookieSameSite	string
)

func init() {
	JWTSecret = os.Getenv("JWT_SECRET")
	if JWTSecret == "" {
		JWTSecret = "dev_secret"
	}

	at := 300
	if v := os.Getenv("ACCESS_TOKEN_TTL_MIN"); v != "" {
		if i, err := strconv.Atoi(v); err == nil { at = i }
	}
	AccessTTL = time.Duration(at) * time.Minute

	CookieDomain = os.Getenv("COOKIE_DOMAIN")
  CookieSecure = (os.Getenv("COOKIE_SECURE") == "true")
  CookieSameSite = os.Getenv("COOKIE_SAMESITE")
}