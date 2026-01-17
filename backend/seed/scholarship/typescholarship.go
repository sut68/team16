package scholarship

import (
	"backend/entity"

	"gorm.io/gorm"
)

func SeedTypeScholarships(db *gorm.DB) error {
	types := []entity.Typescholarship{
		{Typename: "Full"},
		{Typename: "Partial"},
	}

	for _, t := range types {
		if err := db.Where("typename = ?", t.Typename).FirstOrCreate(&t).Error; err != nil {
			return err
		}
	}
	return nil
}
