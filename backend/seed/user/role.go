package user

import (
	"backend/entity"
	"log"

	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) error {
	roles := []entity.Role{
		{Name: "student"},
		{Name: "admin"},
	}

	for _, r := range roles {
		var existing entity.Role
		if err := db.Where("role_name = ?", r.Name).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&r).Error; err != nil {
					return err
				}
				log.Printf("Seeded role '%s'", r.Name)
			} else {
				return err
			}
		} else {
			log.Printf("Role '%s' already exists, skipping...", r.Name)
		}
	}

	return nil
}