package entity

import "gorm.io/gorm"

type ApprovalRequirement struct {
	gorm.Model

	ScholarshipID uint        `json:"scholarship_id" valid:"required~Scholarship ID is required"`
	Scholarship   Scholarship `gorm:"foreignKey:ScholarshipID" json:"scholarship"`

	RequirementID uint        `json:"requirement_id" valid:"required~Requirement ID is required"`
	Requirement   Requirement `gorm:"foreignKey:RequirementID" json:"requirement"`
}
