package entity

import "gorm.io/gorm"

type Location struct {
	gorm.Model
	Name string `gorm:"unique;not null" json:"name"`
	Building string `gorm:"not null" json:"building"`
	Room string `gorm:"not null" json:"room"`
	Floor uint `gorm:"not null" json:"floor"`
	Description string `json:"description"`

	Slots []Slot `gorm:"foreignKey:LocationID" json:"slots"`
}