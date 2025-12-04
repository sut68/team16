package entity

import "gorm.io/gorm"

type StatusScreening struct {
	gorm.Model

	StatusScreening string `json:"status_screening" gorm:"not null"`
}