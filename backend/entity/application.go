package entity

import "gorm.io/gorm"

type Application struct {
	gorm.Model

	StudentProfileID uint           `json:"student_profile_id" valid:"required~Student profile ID is required"`
	StudentProfile   *StudentProfile `gorm:"foreignKey:StudentProfileID" json:"student_profile"`

	ApplicationScholarships []ApplicationScholarship `gorm:"foreignKey:ApplicationID" json:"application_scholarships"`
}
