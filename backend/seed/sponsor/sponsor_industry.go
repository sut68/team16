package sponsor

import (
	"backend/entity"
	"fmt"

	"gorm.io/gorm"
)

func SeedIndustries(db *gorm.DB) error {
	industries := []entity.SponsorIndustry{
    {Name: "Technology"},
    {Name: "Manufacturing"},
    {Name: "Energy"},
    {Name: "Construction"},
    {Name: "Automotive"},
    {Name: "Aerospace"},
    {Name: "Electronics"},
    {Name: "Telecommunications"},
    {Name: "Robotics"},
    {Name: "Biomedical Engineering"},
    {Name: "Civil Engineering"},
    {Name: "Mechanical Engineering"},
    {Name: "Electrical Engineering"},
    {Name: "Chemical Engineering"},
    {Name: "Industrial Engineering"},
    {Name: "Engineering Consulting"},
    {Name: "R&D / Innovation"},
    {Name: "Education"},
    {Name: "Non-Profit"},
	}

	for _, index := range industries {
		var existing entity.SponsorIndustry
		if err := db.Where("name = ?", index.Name).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&index).Error; err != nil {
					return fmt.Errorf("failed seed industry %s: %v", index.Name, err)
				}
			} else {
				return err
			}
		}
	}

	return nil
}