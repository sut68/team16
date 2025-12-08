
package approval

import (
	"backend/entity"
	"gorm.io/gorm"
)

func SeedApplicationDocuments(db *gorm.DB) error {
	if err := db.First(&entity.ApplicationDocument{}).Error; err == gorm.ErrRecordNotFound {
		// Find the application scholarship that should be in 'pending' state
		var appScholarship entity.ApplicationScholarship
		if err := db.Where("status = ?", "pending").First(&appScholarship).Error; err != nil {
			return err
		}

		// Find the student user to use for the 'UploadedBy' field
		var user entity.User
		if err := db.Where("username = ?", "user").First(&user).Error; err != nil {
			return err
		}

		documents := []entity.ApplicationDocument{
			{
				FileName:                 "transkrip_nilai_seed.pdf",
				FilePath:                 "test.pdf",
				FileType:                 "application/pdf",
				UploadedBy:               "user",
				ApplicationScholarshipID: appScholarship.ID,
			},
			{
				FileName:                 "surat_rekomendasi_seed.pdf",
				FilePath:                 "test.pdf",
				FileType:                 "application/pdf",
				UploadedBy:               "user",
				ApplicationScholarshipID: appScholarship.ID,
			},
		}

		if err := db.Create(&documents).Error; err != nil {
			return err
		}
	}
	return nil
}
