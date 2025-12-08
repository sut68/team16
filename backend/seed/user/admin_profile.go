package user

import (
	"backend/entity"
	"gorm.io/gorm"
)

func SeedAdminProfiles(db *gorm.DB) error {
	if err := db.First(&entity.AdminProfile{}).Error; err == gorm.ErrRecordNotFound {
		var adminUser entity.User
		if err := db.Where("username = ?", "admin").First(&adminUser).Error; err != nil {
			return err
		}

		adminProfile := entity.AdminProfile{
			AdminFirstname: "Admin",
			AdminLastname:  "User",
			Position:       1,
			Department:     1,
			Email:          "admin@example.com",
			UserID:         adminUser.ID,
		}

		if err := db.Create(&adminProfile).Error; err != nil {
			return err
		}
	}
	return nil
}
