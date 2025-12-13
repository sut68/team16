package entity

import "gorm.io/gorm"

type AdminProfile struct {
	gorm.Model
	AdminFirstname string `gorm:"not null" json:"admin_firstname"`
	AdminLastname  string `gorm:"not null" json:"admin_lastname"`
	Position       uint   `gorm:"not null" json:"position"`
	Department     uint   `gorm:"not null" json:"department"`
	Email          string `gorm:"not null" json:"email"`

	UserID uint `json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"user"`

	ApprovalTasks []ApprovalTask `gorm:"foreignKey:AdminID" json:"approval_tasks"`
	InterviewRounds []InterviewRound `gorm:"foreignKey:AdminProfileID" json:"interview_rounds"`
}
