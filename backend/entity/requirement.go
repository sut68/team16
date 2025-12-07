package entity

import "gorm.io/gorm"

// Requirement stores the master list of all possible requirement types.
type Requirement struct {
	gorm.Model
	Name string `gorm:"unique;not null" json:"name"`
	ApprovalRequirements []ApprovalRequirement `gorm:"foreignKey:RequirementID" json:"approval_requirements"`
}
