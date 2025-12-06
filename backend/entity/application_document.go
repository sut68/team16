package entity

import "gorm.io/gorm"

type ApplicationDocument struct {
	gorm.Model
	FileName   string `gorm:"not null" json:"file_name" valid:"required~File name is required"`
	UploadedBy string `gorm:"not null" json:"uploaded_by" valid:"required~Uploaded by is required"`

	ApplicationID uint        `json:"application_id"`
	Application   Application `gorm:"foreignKey:ApplicationID" json:"application"`

	RequirementID       uint                `json:"requirement_id"`
	ApprovalRequirement ApprovalRequirement `gorm:"foreignKey:RequirementID" json:"approval_requirement"`

	ApprovalTasks []ApprovalTask `gorm:"foreignKey:DocumentID" json:"approval_tasks"`
}
