package entity

import "gorm.io/gorm"
type Semaster struct {
	gorm.Model
	AcademicYear string `gorm:"not null" json:"academic_year"`
	Term		 string `gorm:"not null" json:"term"`
	Round		 string `gorm:"not null" json:"round"`
	StartDate	 string `gorm:"not null" json:"start_date"`
	EndDate	 string `gorm:"not null" json:"end_date"`
	IsActive	 bool   `gorm:"not null" json:"is_active"`

	Applications []Application `gorm:"foreignKey:SemasterID" json:"applications"`
	Scholarships  []Scholarship `gorm:"foreignKey:SemasterID" json:"scholarships"`
}