package entity

import "gorm.io/gorm"

type ApprovalTask struct {
	gorm.Model
	Status string `gorm:"not null" json:"status"`

	AdminID uint         `json:"admin_id"`
	Admin   AdminProfile `gorm:"foreignKey:AdminID" json:"admin_profile"`

	DocumentID          uint                `json:"document_id"`
	ApplicationDocument ApplicationDocument `gorm:"foreignKey:DocumentID" json:"application_document"`

	ApplicationID uint        `json:"application_id"`
	Application   Application `gorm:"foreignKey:ApplicationID" json:"application"`

	RequirementID       uint                `json:"requirement_id"`
	ApprovalRequirement ApprovalRequirement `gorm:"foreignKey:RequirementID" json:"approval_requirement"`
}
