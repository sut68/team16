package entity

import "gorm.io/gorm"

type ApplicationDocument struct {
	gorm.Model
	FileName   string `gorm:"not null" json:"file_name" valid:"required~File name is required"`
	FilePath   string `json:"file_path"` // Path to the file on the server relative to the uploads directory
	FileType   string `json:"file_type"` // MIME type of the file, e.g., "application/pdf"
	UploadedBy string `gorm:"not null" json:"uploaded_by" valid:"required~Uploaded by is required"`

	ApplicationID uint        `json:"application_id" valid:"required~Application ID is required"`
	Application   Application `gorm:"foreignKey:ApplicationID" json:"application"`

	RequirementID       uint                `json:"requirement_id" valid:"required~Requirement ID is required"`
	ApprovalRequirement ApprovalRequirement `gorm:"foreignKey:RequirementID" json:"approval_requirement"`

	ApprovalTasks []ApprovalTask `gorm:"foreignKey:DocumentID" json:"approval_tasks"`
}
