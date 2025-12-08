package user

import (
	"backend/entity"
	"gorm.io/gorm"
)

func SeedMajors(db *gorm.DB) error {
	if err := db.First(&entity.Major{}).Error; err == gorm.ErrRecordNotFound {
		majors := []entity.Major{
			{
				MajorName: "Computer Engineering",
			},
			{
				MajorName: "Electrical Engineering",
			},
		}
		if err := db.Create(&majors).Error; err != nil {
			return err
		}
	}
	return nil
}
