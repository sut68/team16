package entity

import "gorm.io/gorm"
type InterviewRound struct {
	gorm.Model
	Name string `gorm:"unique;not null" json:"name"`
	Description string `json:"description"`
	StartDateTime string `json:"start_date_time"`
	EndDateTime string `json:"end_date_time"`
	SlotDuration uint `json:"slot_duration"`
	ScholarshipID uint `json:"scholarship_id"`
	Scholarship Scholarship `gorm:"foreignKey:ScholarshipID"`
	AdminProfileID uint `json:"admin_profile_id"`
	AdminProfile AdminProfile `gorm:"foreignKey:AdminProfileID"`

	Slots []Slot `gorm:"foreignKey:InterviewRoundID" json:"slots"`
}