package entity

import "gorm.io/gorm"

type ApprovalRequirement struct {
	gorm.Model
	Name        string `gorm:"not null" json:"name" valid:"required~Name is required"`
	Description string `gorm:"not null" json:"description" valid:"required~Description is required"`

	ScholarshipID uint        `json:"scholarship_id" valid:"required~Scholarship ID is required"`
	Scholarship   Scholarship `gorm:"foreignKey:ScholarshipID" json:"scholarship"`

	ApplicationDocuments []ApplicationDocument `gorm:"foreignKey:RequirementID" json:"application_documents"`
}
