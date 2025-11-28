package models

import "gorm.io/gorm"

type StudentProfile struct {
	gorm.Model
	StudentFirstname string `gorm:"not null" json:"student_firstname"`
	StudentLastname string `gorm:"not null" json:"student_lastname"`
	Grade uint	`gorm:"not null" json:"grade"`
	Class uint	`gorm:"not null" json:"class"`
	Email string `gorm:"not null" json:"email"`

	UserID uint `json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"user"`
}