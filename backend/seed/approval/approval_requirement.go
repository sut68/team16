package approval

import (
	"backend/entity"
	"fmt"

	"gorm.io/gorm"
)

func SeedApprovalRequirements(db *gorm.DB) error {
	/*
		var count int64
		db.Model(&entity.ApprovalRequirement{}).Count(&count)
		if count > 0 {
			return nil // Data already seeded
		}
	*/

	// Map scholarship names to the requirements they need
	scholarshipRequirements := map[string][]string{
		"ทุนเรียนดี ประจำปีการศึกษา 2568": {
			"Transkrip nilai semester terakhir",
			"Surat rekomendasi dari dekan fakultas",
		},
		"ทุนช่วยเหลือนักศึกษาขาดแคลนทุนทรัพย์": {
			"Transkrip nilai semester terakhir",
			"Esai motivasi",
		},
		"ทุนนักกีฬาดีเด่น": {
			"Transkrip nilai semester terakhir",
			"Sertifikat kegiatan ekstrakurikuler",
		},
		"ทุนพัฒนาบุคลากรด้านเทคโนโลยี": {
			"Transkrip nilai semester terakhir",
			"Esai motivasi",
			"Surat rekomendasi dari dekan fakultas",
		},
	}

	for scholarshipName, reqNames := range scholarshipRequirements {
		// Find the scholarship by name
		var scholarship entity.Scholarship
		if err := db.Where("scholarship_name = ?", scholarshipName).First(&scholarship).Error; err != nil {
			return fmt.Errorf("scholarship '%s' not found for seeding ApprovalRequirements: %w", scholarshipName, err)
		}

		for _, reqName := range reqNames {
			// Find the requirement by name
			var requirement entity.Requirement
			if err := db.Where("name = ?", reqName).First(&requirement).Error; err != nil {
				return fmt.Errorf("requirement '%s' not found: %w", reqName, err)
			}

			// Check if this association already exists
			var existing entity.ApprovalRequirement
			err := db.Where("scholarship_id = ? AND requirement_id = ?", scholarship.ID, requirement.ID).First(&existing).Error
			if err == gorm.ErrRecordNotFound {
				// Create the association if it doesn't exist
				newReq := entity.ApprovalRequirement{
					ScholarshipID: scholarship.ID,
					RequirementID: requirement.ID,
				}
				if err := db.Create(&newReq).Error; err != nil {
					return fmt.Errorf("failed to create approval requirement for scholarship '%s' and requirement '%s': %w", scholarshipName, reqName, err)
				}
			} else if err != nil {
				return err
			}
		}
	}

	return nil
}
