package entity

import "gorm.io/gorm"

type ApplicationDocument struct {
	gorm.Model
	FileName   string `gorm:"not null" json:"file_name" valid:"required~File name is required"`
	FilePath   string `json:"file_path"`
	FileType   string `json:"file_type"`
	UploadedBy string `gorm:"not null" json:"uploaded_by" valid:"required~Uploaded by is required"`

	ApplicationScholarshipID uint                    `json:"application_scholarship_id" valid:"required~ApplicationScholarship ID is required"`
	ApplicationScholarship   *ApplicationScholarship `gorm:"foreignKey:ApplicationScholarshipID" json:"application_scholarship"`

	ApprovalTasks []ApprovalTask `gorm:"foreignKey:DocumentID" json:"approval_tasks"`
}
