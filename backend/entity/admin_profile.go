package entity

import "gorm.io/gorm"

type AdminProfile struct {
	gorm.Model
	AdminFirstname string `gorm:"not null" json:"admin_firstname"`
	AdminLastname  string `gorm:"not null" json:"admin_lastname"`
	Position       string `gorm:"not null" json:"position"`
	Email          string `gorm:"not null" json:"email"`
	Phone          string `gorm:"not null" json:"phone"`

	//เพิ่ม refferences:ID เข้ามา
	UserID uint `json:"user_id"`
	User   User `gorm:"foreignKey:UserID;references:ID" json:"user"`

	ApprovalTasks []ApprovalTask `gorm:"foreignKey:AdminID" json:"approval_tasks"`
	InterviewRounds []InterviewRound `gorm:"foreignKey:AdminProfileID" json:"interview_rounds"`
}
