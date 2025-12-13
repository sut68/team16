package controllers

import (
	"backend/config"
	"backend/entity"
	"backend/services"
	"fmt"
	"net/http"
	"strings"
	// "time"

	"github.com/gin-gonic/gin"
)

// Helper
func getUserID(c *gin.Context) (uint, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return 0, fmt.Errorf("no authorization header")
	}
	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
	claims, err := services.ValidateJWT(tokenString) // ใช้ Service
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}

// ---------------------------------------------------
//  DROPDOWN DATA (เพิ่มกลับมาให้แล้วครับ)
// ---------------------------------------------------

func ListRoles(c *gin.Context) {
	var roles []entity.Role
	config.DB.Find(&roles)
	c.JSON(http.StatusOK, gin.H{"data": roles})
}

func ListMajors(c *gin.Context) {
	var majors []entity.Major
	config.DB.Find(&majors)
	c.JSON(http.StatusOK, gin.H{"data": majors})
}

// ---------------------------------------------------
//  USER MANAGEMENT
// ---------------------------------------------------

func ListUsers(c *gin.Context) {
	var users []entity.User
	if err := config.DB.Preload("Role").
		Preload("StudentProfiles.Major").
		Preload("AdminProfiles").
		Find(&users).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}

func CreateUser(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
		RoleID   uint   `json:"role_id"`
		// Student
		StudentID   string    `json:"student_id"`
		FirstNameTH string    `json:"first_name_th"`
		LastNameTH  string    `json:"last_name_th"`
		NationalID  string    `json:"national_id"`
		MajorID     uint      `json:"major_id"`
		GPAX        float64   `json:"gpax"`
		AdvisorName string    `json:"advisor_name"`
		// Admin
		AdminFirstname string `json:"admin_first_name"`
		AdminLastname  string `json:"admin_last_name"`
		Position       string `json:"position"`
		// Common
		Email string `json:"email"`
		Phone string `json:"phone"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentUserID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var targetRole entity.Role
	config.DB.First(&targetRole, input.RoleID)

	if targetRole.Name == "admin" && currentUserID != 1 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only Super Admin can create Admins"})
		return
	}

	tx := config.DB.Begin()

	hashedPassword, _ := services.HashPassword(input.Password)
	user := entity.User{
		Username: input.Username,
		Password: hashedPassword,
		RoleID:   &input.RoleID,
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		if strings.Contains(err.Error(), "Duplicate") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	if targetRole.Name == "student" {
		student := entity.StudentProfile{
			UserID:      user.ID,
			StudentID:   input.StudentID,
			FirstNameTH: input.FirstNameTH,
			LastNameTH:  input.LastNameTH,
			NationalID:  input.NationalID,
			MajorID:     input.MajorID,
			GPAX:        input.GPAX,
			AdvisorName: input.AdvisorName,
			Email:       input.Email,
			Phone:       input.Phone,
			// Default values to prevent errors
			FirstNameEN: "-", LastNameEN: "-", PermanentAddress: "-", CurrentAddress: "-", Province: "-",
		}
		if err := tx.Create(&student).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Create Student Failed"})
			return
		}
	} else if targetRole.Name == "admin" {
		admin := entity.AdminProfile{
			UserID:         user.ID,
			AdminFirstname: input.AdminFirstname,
			AdminLastname:  input.AdminLastname,
			Position:       input.Position,
			Email:          input.Email,
			Phone:          input.Phone,
		}
		if err := tx.Create(&admin).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Create Admin Failed"})
			return
		}
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	// รับเฉพาะข้อมูลที่ Admin แก้ไขให้คนอื่นได้
	var input struct {
		// Student
		StudentID   string    `json:"student_id"`
		FirstNameTH string    `json:"first_name_th"`
		LastNameTH  string    `json:"last_name_th"`
		NationalID  string    `json:"national_id"`
		MajorID     uint      `json:"major_id"`
		GPAX        float64   `json:"gpax"`
		AdvisorName string    `json:"advisor_name"`
		// Admin
		AdminFirstname string `json:"admin_first_name"`
		AdminLastname  string `json:"admin_last_name"`
		Position       string `json:"position"`
		// Common
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentUserID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user entity.User
	if err := config.DB.Preload("Role").First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.Role.Name == "admin" && currentUserID != 1 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only Super Admin can update Admins"})
		return
	}

	tx := config.DB.Begin()

	if input.Password != "" {
		hashed, _ := services.HashPassword(input.Password)
		user.Password = hashed
		tx.Save(&user)
	}

	if user.Role.Name == "student" {
		var student entity.StudentProfile
		if err := tx.Where("user_id = ?", user.ID).First(&student).Error; err == nil {
			student.StudentID = input.StudentID
			student.FirstNameTH = input.FirstNameTH
			student.LastNameTH = input.LastNameTH
			student.NationalID = input.NationalID
			student.MajorID = input.MajorID
			student.GPAX = input.GPAX
			student.AdvisorName = input.AdvisorName
			student.Email = input.Email
			student.Phone = input.Phone
			tx.Save(&student)
		}
	} else if user.Role.Name == "admin" {
		var admin entity.AdminProfile
		// ใช้ Where user_id เพื่อความชัวร์
		if err := tx.Where("user_id = ?", user.ID).First(&admin).Error; err == nil {
			admin.AdminFirstname = input.AdminFirstname
			admin.AdminLastname = input.AdminLastname
			admin.Position = input.Position
			admin.Email = input.Email
			admin.Phone = input.Phone
			tx.Save(&admin)
		} else {
			// สร้างใหม่ถ้า Admin ยังไม่มี Profile
			newAdmin := entity.AdminProfile{
				UserID: user.ID,
				AdminFirstname: input.AdminFirstname,
				AdminLastname: input.AdminLastname,
				Position: input.Position,
				Email: input.Email,
				Phone: input.Phone,
			}
			tx.Create(&newAdmin)
		}
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	currentUserID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user entity.User
	if err := config.DB.Preload("Role").First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.ID == 1 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot delete Super Admin"})
		return
	}
	if user.Role.Name == "admin" && currentUserID != 1 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only Super Admin can delete other Admins"})
		return
	}

	if err := config.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Delete failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}