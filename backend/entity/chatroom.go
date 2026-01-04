package entity

import (
	"gorm.io/gorm"
)

type Chatroom struct {
	gorm.Model
	
	UserID uint `json:"user_id"`
	User User `gorm:"foreignKey:UserID" json:"user"`

	Statuschatroom string `json:"status_chatroom"`  
}
