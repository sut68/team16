package entity

import "gorm.io/gorm"

type ApplicationScholarship struct {
	gorm.Model
	Status string `json:"status"`

	ApplicationID uint        `json:"application_id"`
	Application   Application `gorm:"foreignKey:ApplicationID"`

	ScholarshipID uint        `json:"scholarship_id"`
	Scholarship   Scholarship `gorm:"foreignKey:ScholarshipID"`

	ApplicationDocuments []ApplicationDocument `gorm:"foreignKey:ApplicationScholarshipID"`
}
