package scholarship

import (
	"backend/entity"
	"gorm.io/gorm"
)

func SeedStatusScholarships(db *gorm.DB) error {
	if err := db.First(&entity.Statusscholarship{}).Error; err == gorm.ErrRecordNotFound {
		statuses := []entity.Statusscholarship{
			{
				Statusname: "Open",
			},
			{
				Statusname: "Closed",
			},
		}
		if err := db.Create(&statuses).Error; err != nil {
			return err
		}
	}
	return nil
}
