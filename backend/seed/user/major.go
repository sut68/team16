package user

import (
	"backend/entity"

	"gorm.io/gorm"
)

func SeedMajors(db *gorm.DB) error {
	majors := []entity.Major{
		{MajorName: "Computer Engineering"},
		{MajorName: "Electrical Engineering"},
		{MajorName: "Mechanical Engineering"},
		{MajorName: "Civil Engineering"},
		{MajorName: "Chemical Engineering"},
		{MajorName: "Ceramic Engineering"},
		{MajorName: "Polimer Engineering"},
		{MajorName: "Logistics Engineering"},
		{MajorName: "Environmental Engineering"},
		{MajorName: "Agricultural Engineering"},
		{MajorName: "Industrial Engineering"},
		{MajorName: "Metallurgical Engineering"},
		{MajorName: "Telecommunication Engineering"},
		{MajorName: "Transportation Engineering"},
		{MajorName: "Geological Engineering"},
		{MajorName: "Electronics Engineering"},
		{MajorName: "Design Technology"},
	}

	for _, m := range majors {
		if err := db.Where("major_name = ?", m.MajorName).FirstOrCreate(&m).Error; err != nil {
			return err
		}
	}
	return nil
}
