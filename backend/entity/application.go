package entity

import "gorm.io/gorm"

type Application struct {
	gorm.Model

	StudentProfileID uint           `json:"student_profile_id" valid:"required~Student profile ID is required"`
	StudentProfile   *StudentProfile `gorm:"foreignKey:StudentProfileID" json:"student_profile"`

	SemasterID uint     `json:"semaster_id" valid:"required~Semaster ID is required"`
	Semaster   Semaster `gorm:"foreignKey:SemasterID" json:"semaster"`
	
	ApplicationScholarships []ApplicationScholarship `gorm:"foreignKey:ApplicationID" json:"application_scholarships"`
}
