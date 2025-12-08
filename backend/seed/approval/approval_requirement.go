package approval

import (
	"backend/entity"
	"errors"

	"gorm.io/gorm"
)

func SeedApprovalRequirements(db *gorm.DB) error {
	if err := db.First(&entity.ApprovalRequirement{}).Error; err == gorm.ErrRecordNotFound {
		var scholarship entity.Scholarship
		err := db.Where("scholarship_name = ?", "Beasiswa Anak Bangsa").First(&scholarship).Error
		if err == gorm.ErrRecordNotFound {
			return errors.New("Scholarship 'Beasiswa Anak Bangsa' not found for seeding ApprovalRequirements")
		} else if err != nil {
			return err
		}

		requirements := []entity.ApprovalRequirement{
			{
				Description: "Transkrip nilai semester terakhir",
				ScholarshipID: scholarship.ID,
			},
			{
				Description: "Surat rekomendasi dari dekan fakultas",
				ScholarshipID: scholarship.ID,
			},
		}

		if err := db.Create(&requirements).Error; err != nil {
			return err
		}
	}
	return nil
}
