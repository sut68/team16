package entity

import "gorm.io/gorm"

type ApprovalRequirement struct {
	gorm.Model
	Name        string `gorm:"not null" json:"status"`
	Description string `gorm:"not null" json:"description"`

	ScholarshipID uint        `json:"scholarship_id"`
	Scholarship   Scholarship `gorm:"foreignKey:ScholarshipID" json:"scholarship"`

	ApplicationDocuments []ApplicationDocument `gorm:"foreignKey:RequirementID" json:"application_document"`
}
