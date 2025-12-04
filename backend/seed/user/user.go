package user

import (
	"fmt"
	"log"

	"backend/entity"
	"backend/services"
	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) error {
	users := []struct {
		Username string
		Password string
		RoleName string
	}{
		{Username: "admin", Password: "admin123", RoleName: "admin"},
		{Username: "user", Password: "user123", RoleName: "student"},
	}

	for _, u := range users {
		// หา role id จาก role_name
		var role entity.Role
		if err := db.Where("role_name = ?", u.RoleName).First(&role).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("role '%s' not found, please seed roles first", u.RoleName)
			}
			return fmt.Errorf("failed query role '%s': %v", u.RoleName, err)
		}

		// เช็คว่ามี user นี้แล้วหรือยัง
		var existing entity.User
		err := db.Where("username = ?", u.Username).First(&existing).Error
		if err == nil {
			log.Printf("User '%s' already exists, skipping...", u.Username)
			continue
		}
		if err != gorm.ErrRecordNotFound {
			return fmt.Errorf("unable to query user '%s': %v", u.Username, err)
		}

		// Hash password
		hashedPassword, err := services.HashPassword(u.Password)
		if err != nil {
			return fmt.Errorf("failed to hash password for user '%s': %v", u.Username, err)
		}

		// Create new user
		newUser := entity.User{
			Username: u.Username,
			Password: hashedPassword,
			RoleID:   role.ID,
		}

		if err := db.Create(&newUser).Error; err != nil {
			return fmt.Errorf("failed to seed user '%s': %v", u.Username, err)
		}
		log.Printf("Seeded user: %s (role=%s)", u.Username, u.RoleName)
	}

	return nil
}
