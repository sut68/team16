package entity

import "gorm.io/gorm"

type ApprovalTask struct {
	gorm.Model
	Status string `gorm:"not null" json:"status" valid:"required~Status is required"`

	AdminID uint         `json:"admin_id" valid:"required~Admin ID is required"`
	Admin   AdminProfile `gorm:"foreignKey:AdminID" json:"admin_profile"`

	DocumentID          uint                `json:"document_id" valid:"required~Document ID is required"`
	ApplicationDocument ApplicationDocument `gorm:"foreignKey:DocumentID" json:"application_document"`

	ApprovalDecisions []ApprovalDecision `gorm:"foreignKey:TaskID" json:"approval_decisions"`
}
