package approval

import (
	"time"

	"gorm.io/gorm"

	"backend/entity"
)

func SeedApprovalDecisions(db *gorm.DB) error {
	if err := db.First(&entity.ApprovalDecision{}).Error; err == gorm.ErrRecordNotFound {
		var doc1, doc2 entity.ApplicationDocument
		if err := db.Where("file_name = ?", "transkrip_nilai_seed.pdf").First(&doc1).Error; err != nil {
			return err
		}
		if err := db.Where("file_name = ?", "surat_rekomendasi_seed.pdf").First(&doc2).Error; err != nil {
			return err
		}

		var task1, task2 entity.ApprovalTask
		if err := db.Where("document_id = ?", doc1.ID).First(&task1).Error; err != nil {
			return err
		}
		if err := db.Where("document_id = ?", doc2.ID).First(&task2).Error; err != nil {
			return err
		}

		decisions := []entity.ApprovalDecision{
			{
				DecisionAt: time.Now(),
				Decision:   "Approved",
				Comment:    "Dokumen valid.",
				TaskID:     task1.ID,
			},
			{
				DecisionAt: time.Now(),
				Decision:   "Rejected",
				Comment:    "Tanda tangan tidak sesuai.",
				TaskID:     task2.ID,
			},
		}

		if err := db.Create(&decisions).Error; err != nil {
			return err
		}
	}
	return nil
}
