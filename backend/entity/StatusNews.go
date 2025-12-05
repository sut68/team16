package entity
import "gorm.io/gorm"

type StatusNews struct{
	gorm.Model

	StatusNews string `json:"status_news" gorm:"not null"`
}