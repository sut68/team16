package entity

import "gorm.io/gorm"

type Scholarship struct {
	gorm.Model
	ScholarshipName string `gorm:"not null" json:"scholarship_name"`
	Description string `gorm:"not null" json:"description"`
	OpenDate string `gorm:"not null" json:"open_date"`
	CloseDate string `gorm:"not null" json:"close_date"`
	
	ApprovalRequirements []ApprovalRequirement `gorm:"foreignKey:ScholarshipID" json:"approval_requirements"`
}