
package approval

import (
	"backend/entity"
	"gorm.io/gorm"
)

func SeedApplications(db *gorm.DB) error {
	if err := db.First(&entity.Application{}).Error; err == gorm.ErrRecordNotFound {
		var studentProfile entity.StudentProfile
		if err := db.Where("email = ?", "test.student@example.com").First(&studentProfile).Error; err != nil {
			return err
		}

		applications := []entity.Application{
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
