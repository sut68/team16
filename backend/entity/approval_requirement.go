package entity

import "gorm.io/gorm"

type ApprovalRequirement struct {
	gorm.Model
	Description string `gorm:"not null" json:"description" valid:"required~Description is required"`

	ScholarshipID uint        `json:"scholarship_id" valid:"required~Scholarship ID is required"`
	Scholarship   Scholarship `gorm:"foreignKey:ScholarshipID" json:"scholarship"`

	ApplicationDocuments []ApplicationDocument `gorm:"foreignKey:RequirementID" json:"application_documents"`
}
