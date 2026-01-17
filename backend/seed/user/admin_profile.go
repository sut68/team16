package user

import (
	"backend/entity"

	"gorm.io/gorm"
)

func SeedAdminProfiles(db *gorm.DB) error {
	var adminUser entity.User
	if err := db.Where("username = ?", "admin").First(&adminUser).Error; err != nil {
		return err
	}

	adminProfile := entity.AdminProfile{
		AdminFirstname: "Super",
		AdminLastname:  "Admin",
		Position:       "SuperAdmin",
		Email:          "admin@example.com",
		UserID:         adminUser.ID,
		Phone:          "0812345678",
	}

	// ใช้ UserID เป็นตัวเช็คว่ามีโปรไฟล์ของคนนี้หรือยัง
	if err := db.Where("user_id = ?", adminUser.ID).FirstOrCreate(&adminProfile).Error; err != nil {
		return err
	}
	return nil
}
