package entity

import (
	"time"
	"gorm.io/gorm"
)
type InterviewRound struct {
	gorm.Model
	Name string `gorm:"unique;not null" json:"name"`
	Description string `json:"description"`
	StartDateTime time.Time `json:"start_date_time"`
	EndDateTime time.Time `json:"end_date_time"`
	SlotDuration uint `json:"slot_duration"`
	ScholarshipID uint `json:"scholarship_id"`
	Scholarship Scholarship `gorm:"foreignKey:ScholarshipID" json:"scholarship"`
	AdminProfileID uint `json:"admin_profile_id"`
	AdminProfile AdminProfile `gorm:"foreignKey:AdminProfileID" json:"admin_profile"`

	InterviewModeID *uint `json:"interview_mode_id"`
	InterviewMode *InterviewMode `gorm:"foreignKey:InterviewModeID" json:"interview_mode"`
	LocationID *uint `json:"location_id"`
	Location   *Location `gorm:"foreignKey:LocationID" json:"location"`

	MeetingLink string `json:"meeting_link"`

	Slots []Slot `gorm:"foreignKey:InterviewRoundID" json:"slots"`
}