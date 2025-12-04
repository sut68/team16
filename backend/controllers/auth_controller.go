package controllers

import (
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
	if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	if !services.CheckPasswordHash(password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	token, err := services.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":      token,
		"token_type": "Bearer",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role_id":  user.RoleID,
		},
	})
}