
package approval

import (
	"backend/entity"
	"gorm.io/gorm"
)

func SeedApprovalRequirements(db *gorm.DB) error {
	if err := db.First(&entity.ApprovalRequirement{}).Error; err == gorm.ErrRecordNotFound {
		requirements := []entity.ApprovalRequirement{
			{
				Name:        "Transkrip Nilai",
				Description: "Transkrip nilai semester terakhir",
				ScholarshipID: 1,
			},
			{
				Name:        "Surat Rekomendasi",
				Description: "Surat rekomendasi dari dekan fakultas",
				ScholarshipID: 1,
			},
		}

		if err := db.Create(&requirements).Error; err != nil {
			return err
		}
	}
	return nil
}
