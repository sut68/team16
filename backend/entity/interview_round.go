package entity

import (
	"time"

	"gorm.io/gorm"
)

type InterviewRound struct {

	gorm.Model

	Name          string    `gorm:"unique;not null" json:"name" valid:"required~Name is required"`

	Description   string    `json:"description"`

	StartDateTime time.Time `json:"start_date_time"`

	EndDateTime   time.Time `json:"end_date_time"`

	SlotDuration  uint      `json:"slot_duration" valid:"required,range(1|60)~SlotDuration is required and must be between 1 and 60"`

	ScholarshipID uint      `json:"scholarship_id" valid:"required~ScholarshipID is required"`

	Scholarship   Scholarship `gorm:"foreignKey:ScholarshipID" json:"scholarship" valid:"-"`

	AdminProfileID uint      `json:"admin_profile_id" valid:"required~AdminProfileID is required"`

	AdminProfile  AdminProfile `gorm:"foreignKey:AdminProfileID" json:"admin_profile" valid:"-"`



	InterviewModeID *uint         `json:"interview_mode_id" valid:"required~InterviewModeID is required"`

	InterviewMode   *InterviewMode `gorm:"foreignKey:InterviewModeID" json:"interview_mode" valid:"-"`

	LocationID      *uint         `json:"location_id"`

	Location        *Location     `gorm:"foreignKey:LocationID" json:"location" valid:"-"`



	MeetingLink string `json:"meeting_link"`



	Slots                  []Slot                    `gorm:"foreignKey:InterviewRoundID" json:"slots" valid:"-"`

	InterviewRoundCriteria []InterviewRoundCriterion `gorm:"foreignKey:InterviewRoundID" json:"interview_round_criteria" valid:"-"`

	Evaluations            []Evaluation              `gorm:"foreignKey:InterviewRoundID" json:"evaluations" valid:"-"`

}
