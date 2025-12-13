package controllers

import (
	"backend/config"
	"backend/entity"
	"backend/services"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func getUserIDFromToken(c *gin.Context) (uint, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return 0, fmt.Errorf("no authorization header")
	}
	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
	claims, err := services.ValidateJWT(tokenString)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}

// GET /profile/me
func GetMyProfile(c *gin.Context) {
	userID, err := getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user entity.User
	if err := config.DB.Preload("Role").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.Role.Name == "student" {
		var student entity.StudentProfile
		if err := config.DB.Preload("Major").Where("user_id = ?", userID).First(&student).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
			return
		}
		var family entity.FamilyInfo
		config.DB.Where("profile_id = ?", student.ID).First(&family)
		c.JSON(http.StatusOK, gin.H{"role": "student", "data": student, "family": family})

	} else if user.Role.Name == "admin" {
		var admin entity.AdminProfile
		if err := config.DB.Where("user_id = ?", userID).First(&admin).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"role": "admin", "data": admin})
	}
}

// PUT /profile/me
func UpdateMyProfile(c *gin.Context) {
	userID, err := getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user entity.User
	config.DB.Preload("Role").First(&user, userID)

	// 1. STUDENT
	if user.Role.Name == "student" {
		var student entity.StudentProfile
		if err := config.DB.Where("user_id = ?", userID).First(&student).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
			return
		}

		var input struct {
			// Fields ที่ Student แก้ได้
			FirstNameEN      string            `json:"first_name_en"`
			LastNameEN       string            `json:"last_name_en"`
			BirthDate        time.Time         `json:"birth_date"`
			Phone            string            `json:"phone"`
			Email            string            `json:"email"`
			PermanentAddress string            `json:"permanent_address"`
			CurrentAddress   string            `json:"current_address"`
			Province         string            `json:"province"`
			SiblingsCount    int               `json:"siblings_count"`
			Family           entity.FamilyInfo `json:"family_info"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		tx := config.DB.Begin()

		// **สำคัญ: Map ค่าให้ครบ**
		student.FirstNameEN = input.FirstNameEN
		student.LastNameEN = input.LastNameEN
		student.BirthDate = input.BirthDate
		student.Phone = input.Phone
		student.Email = input.Email
		student.PermanentAddress = input.PermanentAddress
		student.CurrentAddress = input.CurrentAddress
		student.Province = input.Province
		student.SiblingsCount = input.SiblingsCount

		if err := tx.Save(&student).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Update profile failed"})
			return
		}

		// Family Info
		var family entity.FamilyInfo
		err := tx.Where("profile_id = ?", student.ID).First(&family).Error
		if err == gorm.ErrRecordNotFound {
			input.Family.ProfileID = student.ID
			tx.Create(&input.Family)
		} else {
			// Update Existing
			family.FatherName = input.Family.FatherName
			family.FatherOccupation = input.Family.FatherOccupation
			family.FatherIncome = input.Family.FatherIncome
			family.MotherName = input.Family.MotherName
			family.MotherOccupation = input.Family.MotherOccupation
			family.MotherIncome = input.Family.MotherIncome
			family.GuardianName = input.Family.GuardianName
			family.GuardianOccupation = input.Family.GuardianOccupation
			family.GuardianIncome = input.Family.GuardianIncome
			family.GuardianRelation = input.Family.GuardianRelation
			tx.Save(&family)
		}

		tx.Commit()
		c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})

	// 2. ADMIN
	} else if user.Role.Name == "admin" {
		var admin entity.AdminProfile
		if err := config.DB.Where("user_id = ?", userID).First(&admin).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
			return
		}

		var input struct {
			AdminFirstname string `json:"admin_firstname"`
			AdminLastname  string `json:"admin_lastname"`
			Phone          string `json:"phone"`
			Email          string `json:"email"`
			Position       string `json:"position"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		admin.AdminFirstname = input.AdminFirstname
		admin.AdminLastname = input.AdminLastname
		admin.Phone = input.Phone
		admin.Email = input.Email
		admin.Position = input.Position

		config.DB.Save(&admin)
		c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
	}
}