package entity

import "gorm.io/gorm"	

type InterviewerSlot struct {
	gorm.Model
	
	InterviewerID    uint        `json:"interviewer_id"`
	Interviewer      Interviewer `gorm:"foreignKey:InterviewerID" json:"interviewer"`
	SlotID           uint        `json:"slot_id"`
	Slot             Slot        `gorm:"foreignKey:SlotID" json:"slot"`
}
	