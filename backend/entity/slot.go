package entity
import "gorm.io/gorm"

type Slot struct {
	gorm.Model
	Date 	string `gorm:"not null" json:"date"`
	StartTime string `gorm:"not null" json:"start_time"`
	EndTime   string `gorm:"not null" json:"end_time"`
	Mode 	string `gorm:"not null" json:"mode"`
	Capacity uint   `gorm:"not null" json:"capacity"`
	BookCount uint   `gorm:"not null" json:"book_count"`
	Status 	string `gorm:"not null" json:"status"`
	MeetingLink string `json:"meeting_link"`

	InterviewRoundID uint `json:"interview_round_id"`
	InterviewRound   InterviewRound `gorm:"foreignKey:InterviewRoundID" json:"interview_round"`

	LocationID uint `json:"location_id"`
	Location   Location `gorm:"foreignKey:LocationID" json:"location"`

	InterviewerSlots []InterviewerSlot `gorm:"foreignKey:SlotID" json:"interviewer_slots"`
}