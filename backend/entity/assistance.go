package entity

import (
	"gorm.io/gorm"
)

type Assistance struct {
	gorm.Model
	Massage string `json:"massage" valid:"required~massage is required"`
	
	ChatroomID uint `json:"chatroom_id" valid:"required~chatroom is required"`
	Chatroom Chatroom `gorm:"foreignKey:ChatroomID" json:"chatroom"`
	
	SenderID uint `json:"sender_id" valid:"required~sender is required"`
	Sender User	`gorm:"foreignKey:SenderID" json:"sender"`
	
}
