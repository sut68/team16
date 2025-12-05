
package approval

import (
	"backend/entity"
	"gorm.io/gorm"
)

func SeedApplicationDocuments(db *gorm.DB) error {
	if err := db.First(&entity.ApplicationDocument{}).Error; err == gorm.ErrRecordNotFound {
		var user entity.User
		if err := db.Where("username = ?", "user").First(&user).Error; err != nil {
			return err
		}

		var studentProfile entity.StudentProfile
		if err := db.Where("user_id = ?", user.ID).First(&studentProfile).Error; err != nil {
			return err
		}

		var application entity.Application
		if err := db.Where("student_profile_id = ?", studentProfile.ID).First(&application).Error; err != nil {
			return err
		}

		var req1, req2 entity.ApprovalRequirement
		if err := db.Where("name = ?", "Transkrip Nilai").First(&req1).Error; err != nil {
			return err
		}
		if err := db.Where("name = ?", "Surat Rekomendasi").First(&req2).Error; err != nil {
			return err
		}

		documents := []entity.ApplicationDocument{
			{
				FileName:      "transkrip_semester_5.pdf",
				UploadedBy:    "student1",
				ApplicationID: application.ID,
				RequirementID: req1.ID,
			},
			{
				FileName:      "surat_rekomendasi_dekan.pdf",
				UploadedBy:    "student1",
				ApplicationID: application.ID,
				RequirementID: req2.ID,
			},
		}

		if err := db.Create(&documents).Error; err != nil {
			return err
		}
	}
	return nil
}
