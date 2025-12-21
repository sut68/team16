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
			{
				MajorName: "Mechanical Engineering",
			},
			{
				MajorName: "Civil Engineering",
			},
			{
				MajorName: "Chemical Engineering",
			},
			{
				MajorName: "Ceramic Engineering",
			},
			{
				MajorName: "Polimer Engineering",
			},
			{
				MajorName: "Logistics Engineering",
			},
			{
				MajorName: "Environmental Engineering",
			},
			{
				MajorName: "Agricultural Engineering",
			},
			{
				MajorName: "Industrial Engineering",
			},
			{
				MajorName: "Metallurgical Engineering",
			},
			{
				MajorName: "Telecommunication Engineering",
			},
			{
				MajorName: "Transportation Engineering",
			},
			{
				MajorName: "Geological Engineering",
			},
		}
		if err := db.Create(&majors).Error; err != nil {
			return err
		}
	}
	return nil
}
