package entity

import "gorm.io/gorm"

type Application struct {
	gorm.Model

	StudentProfileID uint           `json:"student_profile_id" valid:"required~Student profile is required"`
	StudentProfile   *StudentProfile `gorm:"foreignKey:StudentProfileID" json:"student_profile"`

	ApplicationDocuments []ApplicationDocument `gorm:"foreignKey:ApplicationID" json:"application_documents"`
}
