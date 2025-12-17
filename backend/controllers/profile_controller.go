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

// Helper (Copy มาหรือ Import ก็ได้)
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
		// Preload FamilyInfo ด้วย เพื่อส่งไปให้ครบ
		if err := config.DB.Preload("Major").Preload("FamilyInfo").Where("user_id = ?", userID).First(&student).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
			return
		}
		// ส่งข้อมูลกลับแบบ Flatten เล็กน้อยเพื่อให้ Frontend ใช้ง่าย
		// หรือส่ง struct student ตรงๆ ก็ได้ เพราะเรา Preload FamilyInfo เข้าไปใน student แล้ว
		c.JSON(http.StatusOK, gin.H{
			"role": "student", 
			"data": student, 
			"family": student.FamilyInfo, // ส่งแยกด้วยเผื่อ Frontend เรียกใช้ง่ายๆ
		})

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

	// 1. STUDENT UPDATE
	if user.Role.Name == "student" {
		var student entity.StudentProfile
		if err := config.DB.Where("user_id = ?", userID).First(&student).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
			return
		}

		// Input struct รับเฉพาะข้อมูลส่วนตัวที่อนุญาตให้แก้
		var input struct {
			FirstNameEN      string            `json:"first_name_en"` // อนุญาตให้แก้ชื่ออังกฤษได้
			LastNameEN       string            `json:"last_name_en"`
			BirthDate        time.Time         `json:"birth_date"`
			Phone            string            `json:"phone"`
			Email            string            `json:"email"`
			PermanentAddress string            `json:"permanent_address"`
			CurrentAddress   string            `json:"current_address"`
			Province         string            `json:"province"`
			SiblingsCount    int               `json:"siblings_count"`
			FamilyInfo       entity.FamilyInfo `json:"family_info"` // รับ Object Family ทั้งก้อน
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		tx := config.DB.Begin()

		// Update Student Personal Data
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

		// Update/Create Family Info
		var family entity.FamilyInfo
		err := tx.Where("profile_id = ?", student.ID).First(&family).Error
		
		// Map Data from Input
		family.FatherName = input.FamilyInfo.FatherName
		family.FatherOccupation = input.FamilyInfo.FatherOccupation
		family.FatherIncome = input.FamilyInfo.FatherIncome
		family.MotherName = input.FamilyInfo.MotherName
		family.MotherOccupation = input.FamilyInfo.MotherOccupation
		family.MotherIncome = input.FamilyInfo.MotherIncome
		family.GuardianName = input.FamilyInfo.GuardianName
		family.GuardianRelation = input.FamilyInfo.GuardianRelation
		family.GuardianOccupation = input.FamilyInfo.GuardianOccupation
		family.GuardianIncome = input.FamilyInfo.GuardianIncome

		if err == gorm.ErrRecordNotFound {
			// Create New
			family.ProfileID = student.ID
			if err := tx.Create(&family).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Create family info failed"})
				return
			}
		} else {
			// Update Existing
			if err := tx.Save(&family).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Update family info failed"})
				return
			}
		}

		tx.Commit()
		c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})

	// 2. ADMIN UPDATE (Personal)
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