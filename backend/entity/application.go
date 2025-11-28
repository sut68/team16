package entity

import "gorm.io/gorm"

type Application struct {
	gorm.Model
	
	UserID uint `json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"user"`
}