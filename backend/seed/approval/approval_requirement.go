package approval

import (
	"backend/entity"
	"fmt"

	"gorm.io/gorm"
)

func SeedApprovalRequirements(db *gorm.DB) error {
	// Check if data already exists
	var count int64
	db.Model(&entity.ApprovalRequirement{}).Count(&count)
	if count > 0 {
		return nil // Data already seeded
	}

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

		var requirementsToCreate []entity.ApprovalRequirement

		for _, reqName := range reqNames {
			// Find the requirement by name
			var requirement entity.Requirement
			if err := db.Where("name = ?", reqName).First(&requirement).Error; err != nil {
				return fmt.Errorf("requirement '%s' not found: %w", reqName, err)
			}

			// Prepare the association
			requirementsToCreate = append(requirementsToCreate, entity.ApprovalRequirement{
				ScholarshipID: scholarship.ID,
				RequirementID: requirement.ID,
			})
		}

		// Create the associations for the current scholarship
		if len(requirementsToCreate) > 0 {
			if err := db.Create(&requirementsToCreate).Error; err != nil {
				return fmt.Errorf("failed to create approval requirements for scholarship '%s': %w", scholarshipName, err)
			}
		}
	}

	return nil
}