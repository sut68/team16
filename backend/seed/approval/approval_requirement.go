package approval

import (
	"backend/entity"
	"errors"

	"gorm.io/gorm"
)

func SeedApprovalRequirements(db *gorm.DB) error {
	if err := db.First(&entity.ApprovalRequirement{}).Error; err == gorm.ErrRecordNotFound {
		// Find scholarship
		var scholarship entity.Scholarship
		if err := db.Where("scholarship_name = ?", "Beasiswa Anak Bangsa").First(&scholarship).Error; err != nil {
			return errors.New("Scholarship 'Beasiswa Anak Bangsa' not found for seeding ApprovalRequirements")
		}

		// Find master requirements
		var req1, req2 entity.Requirement
		if err := db.Where("name = ?", "Transkrip nilai semester terakhir").First(&req1).Error; err != nil {
			return errors.New("Requirement 'Transkrip nilai semester terakhir' not found")
		}
		if err := db.Where("name = ?", "Surat rekomendasi dari dekan fakultas").First(&req2).Error; err != nil {
			return errors.New("Requirement 'Surat rekomendasi dari dekan fakultas' not found")
		}

		requirements := []entity.ApprovalRequirement{
			{
				ScholarshipID: scholarship.ID,
				RequirementID: req1.ID,
			},
			{
				ScholarshipID: scholarship.ID,
				RequirementID: req2.ID,
			},
		}

		if err := db.Create(&requirements).Error; err != nil {
			return err
		}
	}
	return nil
}
