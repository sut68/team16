package entity

import (
	"gorm.io/gorm"
)

type SponsorContact struct {
	gorm.Model
	Name      string  `json:"name"`
	Email     string  `json:"email" valid:"required~Email is required, email~Invalid email format"`
	Phone     string  `json:"phone" valid:"required~Phone is required,  numeric~Phone must contain only numbers, length(9|10)~Phone number length must be 9-10 digits"`
	Position  *string `json:"position"`
	SponsorID uint    `json:"sponsor_id"`
}
