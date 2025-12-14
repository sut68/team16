package entity

import "gorm.io/gorm"

type Interviewer struct {
	gorm.Model
	InterviewerFirstname string `gorm:"not null" json:"interviewer_firstname"`
	InterviewerLastname  string `gorm:"not null" json:"interviewer_lastname"`
	Email                string `gorm:"not null" json:"email"`

	InterviewerSlots []InterviewerSlot `gorm:"foreignKey:InterviewerID" json:"interviewer_slots"`
}