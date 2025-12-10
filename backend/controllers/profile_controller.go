package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"backend/config"
	"backend/entity"
)

// GetMyStudentProfile - GET /api/profile/student/me
func GetMyStudentProfile(c *gin.Context) {
	uid, ok := getUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var sp entity.StudentProfile
	if err := config.DB.Where("user_id = ?", uid).Preload("Major").First(&sp).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "student profile not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query profile"})
		return
	}
	c.JSON(http.StatusOK, sp)
}

// UpdateMyStudentProfile - PUT /api/profile/student/me
// Allowed: personal fields only (not StudentID, GPAX, MajorID unless you want)
func UpdateMyStudentProfile(c *gin.Context) {
	uid, ok := getUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var sp entity.StudentProfile
	if err := config.DB.Where("user_id = ?", uid).First(&sp).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "student profile not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query profile"})
		return
	}

	var payload struct {
		FirstNameTH      *string `json:"first_name_th"`
		LastNameTH       *string `json:"last_name_th"`
		FirstNameEN      *string `json:"first_name_en"`
		LastNameEN       *string `json:"last_name_en"`
		AdvisorName      *string `json:"advisor_name"`
		Phone            *string `json:"phone"`
		Email            *string `json:"email"`
		PermanentAddress *string `json:"permanent_address"`
		CurrentAddress   *string `json:"current_address"`
		Province         *string `json:"province"`
		SiblingsCount    *int    `json:"siblings_count"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if payload.FirstNameTH != nil {
		sp.FirstNameTH = *payload.FirstNameTH
	}
	if payload.LastNameTH != nil {
		sp.LastNameTH = *payload.LastNameTH
	}
	if payload.FirstNameEN != nil {
		sp.FirstNameEN = *payload.FirstNameEN
	}
	if payload.LastNameEN != nil {
		sp.LastNameEN = *payload.LastNameEN
	}
	if payload.AdvisorName != nil {
		sp.AdvisorName = *payload.AdvisorName
	}
	if payload.Phone != nil {
		sp.Phone = *payload.Phone
	}
	if payload.Email != nil {
		sp.Email = *payload.Email
	}
	if payload.PermanentAddress != nil {
		sp.PermanentAddress = *payload.PermanentAddress
	}
	if payload.CurrentAddress != nil {
		sp.CurrentAddress = *payload.CurrentAddress
	}
	if payload.Province != nil {
		sp.Province = *payload.Province
	}
	if payload.SiblingsCount != nil {
		sp.SiblingsCount = *payload.SiblingsCount
	}

	if err := config.DB.Save(&sp).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save profile"})
		return
	}
	c.JSON(http.StatusOK, sp)
}

// Admin profile handlers
// GetMyAdminProfile - GET /api/profile/admin/me
func GetMyAdminProfile(c *gin.Context) {
	uid, ok := getUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var ap entity.AdminProfile
	if err := config.DB.Where("user_id = ?", uid).First(&ap).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "admin profile not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query admin profile"})
		return
	}
	c.JSON(http.StatusOK, ap)
}

// UpdateMyAdminProfile - PUT /api/profile/admin/me
// Position kept as string (per your request)
func UpdateMyAdminProfile(c *gin.Context) {
	uid, ok := getUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var ap entity.AdminProfile
	if err := config.DB.Where("user_id = ?", uid).First(&ap).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "admin profile not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query admin profile"})
		return
	}

	var payload struct {
		AdminFirstName *string `json:"admin_firstname"`
		AdminLastName  *string `json:"admin_lastname"`
		Position       *string `json:"position"` // <-- string
		Department     *string `json:"department"`
		Email          *string `json:"email"`
		Phone          *string `json:"phone"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if payload.AdminFirstName != nil {
		ap.AdminFirstname = *payload.AdminFirstName
	}
	if payload.AdminLastName != nil {
		ap.AdminLastname = *payload.AdminLastName
	}
	if payload.Position != nil {
		ap.Position = *payload.Position
	}
	if payload.Email != nil {
		ap.Email = *payload.Email
	}
	if payload.Phone != nil {
		ap.Phone = *payload.Phone
	}

	if err := config.DB.Save(&ap).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save admin profile"})
		return
	}
	c.JSON(http.StatusOK, ap)
}
