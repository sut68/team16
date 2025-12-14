package approval

import (
	"backend/entity"
	"fmt"
	"gorm.io/gorm"
)

func SeedApplicationScholarships(db *gorm.DB) error {
	if err := db.First(&entity.ApplicationScholarship{}).Error; err == gorm.ErrRecordNotFound {
		// Find the first application (which belongs to the test student)
		var application entity.Application
		if err := db.First(&application, 1).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("in SeedApplicationScholarships: prerequisite Application with ID=1 not found")
			}
			return err
		}

		// Find the first two scholarships
		var scholarship1, scholarship2 entity.Scholarship
		if err := db.First(&scholarship1, 1).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("in SeedApplicationScholarships: prerequisite Scholarship with ID=1 not found")
			}
			return err
		}
		if err := db.First(&scholarship2, 2).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("in SeedApplicationScholarships: prerequisite Scholarship with ID=2 not found")
			}
			return err
		}

		appScholarships := []entity.ApplicationScholarship{
			{
				// This one is in 'qualified' state, waiting for documents
				Status:        "qualified",
				ApplicationID: application.ID,
				ScholarshipID: scholarship1.ID,
			},
			{
				// This one will have documents and be in 'pending' state
				Status:        "pending",
				ApplicationID: application.ID,
				ScholarshipID: scholarship2.ID,
			},
		}

		if err := db.Create(&appScholarships).Error; err != nil {
			return err
		}
	}
	return nil
}
