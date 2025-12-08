
package approval

import (
	"backend/entity"
	"gorm.io/gorm"
)

func SeedApprovalTasks(db *gorm.DB) error {
	if err := db.First(&entity.ApprovalTask{}).Error; err == gorm.ErrRecordNotFound {
		var adminProfile entity.AdminProfile
		if err := db.Where("email = ?", "admin@example.com").First(&adminProfile).Error; err != nil {
			return err
		}

		var doc1, doc2 entity.ApplicationDocument
		if err := db.Where("file_name = ?", "transkrip_semester_5.pdf").First(&doc1).Error; err != nil {
			return err
		}
		if err := db.Where("file_name = ?", "surat_rekomendasi_dekan.pdf").First(&doc2).Error; err != nil {
			return err
		}

		tasks := []entity.ApprovalTask{
			{
				Status:        "Pending",
				AdminID:       adminProfile.ID,
				DocumentID:    doc1.ID,
				ApplicationID: doc1.ApplicationID,
				RequirementID: doc1.RequirementID,
			},
			{
				Status:        "Pending",
				AdminID:       adminProfile.ID,
				DocumentID:    doc2.ID,
				ApplicationID: doc2.ApplicationID,
				RequirementID: doc2.RequirementID,
			},
		}

		if err := db.Create(&tasks).Error; err != nil {
			return err
		}
	}
	return nil
}
