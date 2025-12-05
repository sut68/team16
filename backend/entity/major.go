package entity

import "gorm.io/gorm"

type Major struct {
	gorm.Model
	MajorName       string           `gorm:"not null" json:"major_name"`
	StudentProfiles []StudentProfile `gorm:"foreignKey:MajorID" json:"student_profiles"`
}