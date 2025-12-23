
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
		if err := db.Where("file_name = ?", "transkrip_nilai_seed.pdf").First(&doc1).Error; err != nil {
			return err
		}
		if err := db.Where("file_name = ?", "surat_rekomendasi_seed.pdf").First(&doc2).Error; err != nil {
			return err
		}

		tasks := []entity.ApprovalTask{
			{
				Status:     "Pending",
				DocumentID: doc1.ID,
			},
			{
				Status:     "Pending",
				DocumentID: doc2.ID,
			},
		}

		if err := db.Create(&tasks).Error; err != nil {
			return err
		}
	}
	return nil
}
