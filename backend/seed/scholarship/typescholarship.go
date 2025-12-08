package scholarship

import (
	"backend/entity"
	"gorm.io/gorm"
)

func SeedTypeScholarships(db *gorm.DB) error {
	if err := db.First(&entity.Typescholarship{}).Error; err == gorm.ErrRecordNotFound {
		types := []entity.Typescholarship{
			{
				Typename: "Full",
			},
			{
				Typename: "Partial",
			},
		}
		if err := db.Create(&types).Error; err != nil {
			return err
		}
	}
	return nil
}
