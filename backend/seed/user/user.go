package user

import (
	"fmt"
	"log"
	"os"

	"backend/entity"
	"backend/services"

	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) error {
	// ดึงค่าจาก Environment Variables ถ้าไม่มีให้ใช้ค่า Default
	adminUsername := os.Getenv("SEED_ADMIN_USERNAME")
	if adminUsername == "" {
		adminUsername = "admin"
	}
	adminPassword := os.Getenv("SEED_ADMIN_PASSWORD")
	if adminPassword == "" {
		adminPassword = "admin123"
	}

	userUsername := os.Getenv("SEED_USER_USERNAME")
	if userUsername == "" {
		userUsername = "user"
	}
	userPassword := os.Getenv("SEED_USER_PASSWORD")
	if userPassword == "" {
		userPassword = "user123"
	}

	users := []struct {
		Username string
		Password string
		RoleName string
	}{
		{Username: adminUsername, Password: adminPassword, RoleName: "admin"},
		{Username: userUsername, Password: userPassword, RoleName: "student"},
	}

	for _, u := range users {
		var role entity.Role
		if err := db.Where("name = ?", u.RoleName).First(&role).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("role '%s' not found, please seed roles first", u.RoleName)
			}
			return fmt.Errorf("failed query role '%s': %v", u.RoleName, err)
		}

		// Hash password
		hashedPassword, err := services.HashPassword(u.Password)
		if err != nil {
			return fmt.Errorf("failed to hash password for user '%s': %v", u.Username, err)
		}

		// เช็คว่ามี user นี้แล้วหรือยัง
		var existing entity.User
		err = db.Where("username = ?", u.Username).First(&existing).Error
		if err == nil {
			// ถ้ามีแล้ว ให้ Update รหัสผ่านและ Role (เพื่อให้ค่าจาก .env ล่าสุดถูกนำไปใช้)
			log.Printf("User '%s' already exists, updating password...", u.Username)
			if err := db.Model(&existing).Updates(entity.User{
				Password: hashedPassword,
				RoleID:   &role.ID,
			}).Error; err != nil {
				return fmt.Errorf("failed to update user '%s': %v", u.Username, err)
			}
			continue
		}

		if err != gorm.ErrRecordNotFound {
			return fmt.Errorf("unable to query user '%s': %v", u.Username, err)
		}

		// Create new user
		newUser := entity.User{
			Username: u.Username,
			Password: hashedPassword,
			RoleID:   &role.ID,
		}

		if err := db.Create(&newUser).Error; err != nil {
			return fmt.Errorf("failed to seed user '%s': %v", u.Username, err)
		}
		log.Printf("Seeded user: %s (role=%s)", u.Username, u.RoleName)
	}

	return nil
}
