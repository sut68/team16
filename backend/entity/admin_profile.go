package entity

import "gorm.io/gorm"

type AdminProfile struct {
	gorm.Model
	AdminFirstname string `gorm:"not null" json:"admin_first_name"`
	AdminLastname  string `gorm:"not null" json:"admin_last_name"`
	Position       string `gorm:"not null" json:"position"`
	Email          string `gorm:"not null" json:"email"`
	Phone          string `gorm:"not null" json:"phone"`

	UserID uint `json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"user"`

	ApprovalTasks []ApprovalTask `gorm:"foreignKey:AdminID" json:"approval_tasks"`
}
