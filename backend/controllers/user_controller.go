package controllers

import (
	"backend/config"
	"backend/entity"
	"backend/services"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Helper: Get User ID from Token
func getUserID(c *gin.Context) (uint, error) {
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

// ---------------------------------------------------
//  DROPDOWN DATA
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
//  USER MANAGEMENT (ADMIN ONLY)
// ---------------------------------------------------

func ListUsers(c *gin.Context) {
	var users []entity.User
	// Preload ข้อมูลที่จำเป็นสำหรับการแสดงผลในตาราง
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
	// Struct รับค่าจาก Frontend (Admin Form)
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
		RoleID   uint   `json:"role_id"`

		// --- Student Official Data ---
		StudentID   string  `json:"student_id"`
		NationalID  string  `json:"national_id"`
		FirstNameTH string  `json:"first_name_th"`
		LastNameTH  string  `json:"last_name_th"`
		FirstNameEN string  `json:"first_name_en"` // เพิ่ม
		LastNameEN  string  `json:"last_name_en"`  // เพิ่ม
		CurrentYear int     `json:"current_year"`  // เพิ่ม
		MajorID     uint    `json:"major_id"`
		GPAX        float64 `json:"gpax"`
		AdvisorName string  `json:"advisor_name"`

		// --- Admin Data ---
		AdminFirstname string `json:"admin_firstname"`
		AdminLastname  string `json:"admin_lastname"`
		Position       string `json:"position"`
		Email          string `json:"email"` // Admin ยังใช้อีเมลในการติดต่อได้
		Phone          string `json:"phone"` // Admin ยังใช้เบอร์โทรได้
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check Permission
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

	// 1. Create User Account
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

	// 2. Create Profile based on Role
	if targetRole.Name == "student" {
		student := entity.StudentProfile{
			UserID:      user.ID,
			StudentID:   input.StudentID,
			NationalID:  input.NationalID,
			FirstNameTH: input.FirstNameTH,
			LastNameTH:  input.LastNameTH,
			FirstNameEN: input.FirstNameEN, // บันทึกชื่ออังกฤษ
			LastNameEN:  input.LastNameEN,  // บันทึกนามสกุลอังกฤษ
			CurrentYear: input.CurrentYear, // บันทึกชั้นปี
			MajorID:     input.MajorID,
			GPAX:        input.GPAX,
			AdvisorName: input.AdvisorName,
			// Contact Info ปล่อยว่างหรือใส่ Default ให้นิสิตไปกรอกเอง
			Email:            "-",
			Phone:            "-",
			PermanentAddress: "-",
			CurrentAddress:   "-",
			Province:         "-",
		}
		if err := tx.Create(&student).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Create Student Profile Failed: " + err.Error()})
			return
		}

		// สร้าง FamilyInfo เปล่าๆ ไว้รอเลยก็ได้ เพื่อกัน Error เวลา Query แบบ Preload
		family := entity.FamilyInfo{ProfileID: student.ID}
		tx.Create(&family)

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
			c.JSON(http.StatusBadRequest, gin.H{"error": "Create Admin Profile Failed"})
			return
		}
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")

	// Input Structure
	var input struct {
		Username string `json:"username"` // ✅ เพิ่มตรงนี้: รับค่า Username มาแก้ไข
		Password string `json:"password"`
		
		// Student Official Data Only
		StudentID 	string 	`json:"student_id"`
		NationalID 	string 	`json:"national_id"`
		FirstNameTH string 	`json:"first_name_th"`
		LastNameTH 	string 	`json:"last_name_th"`
		FirstNameEN string 	`json:"first_name_en"`
		LastNameEN 	string 	`json:"last_name_en"`
		CurrentYear int 	`json:"current_year"`
		MajorID 	uint 	`json:"major_id"`
		GPAX 		float64 `json:"gpax"`
		AdvisorName string 	`json:"advisor_name"`
		// Admin Data
		AdminFirstname string `json:"admin_firstname"`
		AdminLastname 	string `json:"admin_lastname"`
		Position 	 	string `json:"position"`
		Email 		 	string `json:"email"`
		Phone 		 	string `json:"phone"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Permission Check...
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

	// -------------------------------------------------------------
	// 1. Update User Account (Username & Password)
	// -------------------------------------------------------------
	
	// Logic แก้ไข Username
	if input.Username != "" && input.Username != user.Username {
		// เช็คว่า Username ใหม่ซ้ำกับคนอื่นไหม
		var checkUser entity.User
		if err := tx.Where("username = ?", input.Username).First(&checkUser).Error; err == nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username นี้มีผู้ใช้งานแล้ว กรุณาใช้ชื่ออื่น"})
			return
		}
		user.Username = input.Username
	}

	// Logic แก้ไข Password
	if input.Password != "" {
		hashed, _ := services.HashPassword(input.Password)
		user.Password = hashed
	}

	// บันทึกการเปลี่ยนแปลง User (ทั้ง Username และ Password)
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Update User failed: " + err.Error()})
		return
	}

	// -------------------------------------------------------------
	// 2. Update Profile
	// -------------------------------------------------------------
	if user.Role.Name == "student" {
		var student entity.StudentProfile
		if err := tx.Where("user_id = ?", user.ID).First(&student).Error; err == nil {
			// Update Only Official Fields
			student.StudentID = input.StudentID
			student.NationalID = input.NationalID
			student.FirstNameTH = input.FirstNameTH
			student.LastNameTH = input.LastNameTH
			student.FirstNameEN = input.FirstNameEN
			student.LastNameEN = input.LastNameEN
			student.CurrentYear = input.CurrentYear
			student.MajorID = input.MajorID
			student.GPAX = input.GPAX
			student.AdvisorName = input.AdvisorName
			// หมายเหตุ: ไม่ Update Email/Phone ที่นี่ เพื่อไม่ให้ทับข้อมูลที่นิสิตกรอกเอง
			tx.Save(&student)
		}
	} else if user.Role.Name == "admin" {
		var admin entity.AdminProfile
		if err := tx.Where("user_id = ?", user.ID).First(&admin).Error; err == nil {
			admin.AdminFirstname = input.AdminFirstname
			admin.AdminLastname = input.AdminLastname
			admin.Position = input.Position
			admin.Email = input.Email
			admin.Phone = input.Phone
			tx.Save(&admin)
		}
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	// 1. ตรวจสอบสิทธิ์ (Code เดิม)
	currentUserID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// เริ่ม Transaction
	tx := config.DB.Begin()

	// 2. ค้นหา User (ใช้ Unscoped เพื่อหาคนที่ถูก Soft Delete ไปแล้วด้วย)
	var user entity.User
	if err := tx.Unscoped().Preload("Role").First(&user, id).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 3. กฎห้ามลบ Super Admin / Admin ข้ามสิทธิ์
	if user.ID == 1 {
		tx.Rollback()
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot delete Super Admin"})
		return
	}
	if user.Role.Name == "admin" && currentUserID != 1 {
		tx.Rollback()
		c.JSON(http.StatusForbidden, gin.H{"error": "Only Super Admin can delete other Admins"})
		return
	}

	// =================================================================
	// 🧹 Manual Cleanup: ล้างบาง (เหลน -> หลาน -> ลูก -> พ่อ)
	// =================================================================

	// A. กรณีเป็น STUDENT
	if user.Role.Name == "student" {
		var studentProfiles []entity.StudentProfile
		// หา Profile ทั้งหมดของ User นี้
		if err := tx.Unscoped().Where("user_id = ?", user.ID).Find(&studentProfiles).Error; err == nil {
			
			for _, sp := range studentProfiles {
				
				// -------------------------------------------------------------
				// STEP 1: ลบ "เหลน" (Screenings) ** จุดที่เพิ่มใหม่ **
				// -------------------------------------------------------------
				// ต้องหา ID ของ Application ที่ผูกกับ Profile นี้ก่อน
				var appIDs []uint
				tx.Model(&entity.Application{}).Unscoped().Where("student_profile_id = ?", sp.ID).Pluck("id", &appIDs)
				
				if len(appIDs) > 0 {
					// สั่งลบ screenings ที่ application_id อยู่ในลิสต์ที่เราหามา
					if err := tx.Exec("DELETE FROM screenings WHERE application_id IN ?", appIDs).Error; err != nil {
						tx.Rollback()
						c.JSON(http.StatusBadRequest, gin.H{"error": "ลบ Screenings ไม่ผ่าน: " + err.Error()})
						return
					}
				}

				// -------------------------------------------------------------
				// STEP 2: ลบ "หลาน" (Applications & Family)
				// -------------------------------------------------------------
				
				// 2.1 ลบ Applications (ใบสมัคร)
				if err := tx.Unscoped().Where("student_profile_id = ?", sp.ID).Delete(&entity.Application{}).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusBadRequest, gin.H{"error": "ลบ Applications ไม่ผ่าน: " + err.Error()})
					return
				}

				// 2.2 ลบ FamilyInfo (ครอบครัว)
				if err := tx.Unscoped().Where("profile_id = ?", sp.ID).Delete(&entity.FamilyInfo{}).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusBadRequest, gin.H{"error": "ลบ FamilyInfo ไม่ผ่าน: " + err.Error()})
					return
				}

				// -------------------------------------------------------------
				// STEP 3: ลบ "ลูก" (StudentProfile)
				// -------------------------------------------------------------
				if err := tx.Unscoped().Delete(&sp).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusBadRequest, gin.H{"error": "ลบ StudentProfile ไม่ผ่าน: " + err.Error()})
					return
				}
			}
		}

	// B. กรณีเป็น ADMIN
	} else if user.Role.Name == "admin" {
		if err := tx.Unscoped().Where("user_id = ?", user.ID).Delete(&entity.AdminProfile{}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "ลบ AdminProfile ไม่ผ่าน"})
			return
		}
	}

	// -------------------------------------------------------------
	// STEP 4: ลบ "พ่อ" (User)
	// -------------------------------------------------------------
	if err := tx.Unscoped().Delete(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "ลบ User ไม่ผ่าน: " + err.Error()})
		return
	}

	// ยืนยันการลบ
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "ลบข้อมูลผู้ใช้งานและประวัติทั้งหมด (Screenings, Applications, Family) เรียบร้อยแล้ว"})
}