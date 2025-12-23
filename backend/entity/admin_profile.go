package entity

import "gorm.io/gorm"

type AdminProfile struct {
	gorm.Model
	AdminFirstname string `gorm:"not null" json:"admin_firstname" valid:"required~Firstname is required"`
	AdminLastname  string `gorm:"not null" json:"admin_lastname" valid:"required~Lastname is required"`
	Position       string `gorm:"not null" json:"position" valid:"required~Position is required"`
	Email          string `gorm:"not null" json:"email" valid:"required~Email is required,email~Invalid email format"`
	Phone          string `gorm:"not null" json:"phone" valid:"required~Phone is required,numeric~Phone must handle numbers only,stringlength(10|10)~Phone must be 10 digits"`

	UserID uint `json:"user_id" valid:"required~User ID is required"`
	User   User `gorm:"foreignKey:UserID;references:ID" json:"user" valid:"-"`

	InterviewRounds []InterviewRound `gorm:"foreignKey:AdminProfileID" json:"interview_rounds" valid:"-"`
	ApprovalDecisions []ApprovalDecision `gorm:"foreignKey:AdminID" json:"approval_decisions" valid:"-"`
}