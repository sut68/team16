package entity

import (
	"gorm.io/gorm"
)

type Assistance struct {
	gorm.Model
	Massage string `json:"massage"`
	
	ChatroomID uint `json:"chatroom_id"`
	Chatroom Chatroom `gorm:"foreignKey:ChatroomID" json:"chatroom"`
	
	SenderID uint `json:"sender_id"`
	Sender User	`gorm:"foreignKey:SenderID" json:"sender"`
	
}
