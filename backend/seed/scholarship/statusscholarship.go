package scholarship

import (
	"backend/entity"

	"gorm.io/gorm"
)

func SeedStatusScholarships(db *gorm.DB) error {
	statuses := []entity.Statusscholarship{
		{Statusname: "Open"},
		{Statusname: "Closed"},
	}

	for _, s := range statuses {
		if err := db.Where("statusname = ?", s.Statusname).FirstOrCreate(&s).Error; err != nil {
			return err
		}
	}
	return nil
}
