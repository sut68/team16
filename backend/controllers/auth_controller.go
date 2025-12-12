package controllers

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"backend/config"
	"backend/entity"
	"backend/services"
)

// AuthInput reused
type AuthInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var input AuthInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username := strings.TrimSpace(input.Username)
	password := strings.TrimSpace(input.Password)

	if username == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username and password are required"})
		return
	}

	var role entity.Role
	if err := config.DB.FirstOrCreate(&role, entity.Role{Name: "user"}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to ensure default role 'user'"})
		return
	}

	hashedPassword, err := services.HashPassword(password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	user := entity.User{
		Username: username,
		Password: hashedPassword,
		RoleID:   &role.ID, 
	}

	if err := config.DB.Create(&user).Error; err != nil {
		errMsg := strings.ToLower(err.Error())
		if strings.Contains(errMsg, "duplicate") ||
			strings.Contains(errMsg, "unique") ||
			strings.Contains(errMsg, "already exists") ||
			strings.Contains(errMsg, "violates unique constraint") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username already exists"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "registration successful",
		"id":       user.ID,
		"username": user.Username,
		"role_id":  user.RoleID,
	})
}

func makeCSRFToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func Login(c *gin.Context) {
	var input AuthInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username := strings.TrimSpace(input.Username)
	password := strings.TrimSpace(input.Password)

	if username == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username and password are required"})
		return
	}

	var user entity.User
	if err := config.DB.Preload("Role").Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	if !services.CheckPasswordHash(password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	var roleName string
	if user.Role != nil {
		roleName = user.Role.Name
	} else {
		roleName = "user"
	}

	// Generate JWT
	token, err := services.GenerateJWT(user.ID, roleName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	accessCookie := &http.Cookie{
		Name:			"access_token",
		Value:    token,
		Path:     "/",
		Domain:   config.CookieDomain,
		MaxAge:   int(config.AccessTTL.Seconds()),
		HttpOnly: true,
		Secure:   config.CookieSecure,	
	}
	switch strings.ToLower(config.CookieSameSite) {
	case "strict":
		accessCookie.SameSite = http.SameSiteStrictMode
	case "none":
		accessCookie.SameSite = http.SameSiteNoneMode
	default:
		accessCookie.SameSite = http.SameSiteLaxMode
	}
	http.SetCookie(c.Writer, accessCookie)

	csrfToken, err := makeCSRFToken(32)
	if err != nil {
		csrfToken = ""
	}
	csrfCookie := &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Path:     "/",
		Domain:   config.CookieDomain,
		MaxAge:   int(config.AccessTTL.Seconds()),
		HttpOnly: false,
		Secure:   config.CookieSecure,
	}
	csrfCookie.SameSite = accessCookie.SameSite
	http.SetCookie(c.Writer, csrfCookie)

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"user": gin.H{
			"id": user.ID,
			"username": user.Username,
			"role_id": user.RoleID,
			"role": roleName,
		},
	})
}

func Logout(c *gin.Context) {
	clearAccess := &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		Domain:   config.CookieDomain,
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   config.CookieSecure,
	}
	switch strings.ToLower(config.CookieSameSite) {
	case "strict":
		clearAccess.SameSite = http.SameSiteStrictMode
	case "none":
		clearAccess.SameSite = http.SameSiteNoneMode
	default:
		clearAccess.SameSite = http.SameSiteLaxMode
	}
	http.SetCookie(c.Writer, clearAccess)

	clearCSRF := &http.Cookie{
		Name:     "csrf_token",
		Value:    "",
		Path:     "/",
		Domain:   config.CookieDomain,
		MaxAge:   -1,
		HttpOnly: false,
		Secure:   config.CookieSecure,
	}
	clearCSRF.SameSite = clearAccess.SameSite
	http.SetCookie(c.Writer, clearCSRF)

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}