package approval

import (
	"backend/entity"
	"fmt"
	"gorm.io/gorm"
)

func SeedApplications(db *gorm.DB) error {
	if err := db.First(&entity.Application{}).Error; err == gorm.ErrRecordNotFound {
		var studentProfile entity.StudentProfile
		if err := db.Where("email = ?", "test.student@example.com").First(&studentProfile).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("in SeedApplications: prerequisite student profile with email 'test.student@example.com' not found")
			}
			return err
		}

		applications := []entity.Application{
			{
				StudentProfileID: studentProfile.ID,
			},
			{
				StudentProfileID: studentProfile.ID,
			},
		}

		if err := db.Create(&applications).Error; err != nil {
			return err
		}
	}
	return nil
}
