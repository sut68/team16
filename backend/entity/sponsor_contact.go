package entity

import (
	"gorm.io/gorm"
)

type SponsorContact struct {
	gorm.Model
	Name      string  `json:"name" valid:"required~Contact name is required,stringlength(2|100)~Contact name must be 2-100 characters"`
	Email     string  `json:"email" valid:"required~Email is required,email~Invalid email format"`
	Phone     string  `json:"phone" valid:"required~Phone is required,numeric~Phone must contain only numbers,stringlength(9|10)~Phone number must be 9-10 digits"`
	Position  *string `json:"position" valid:"optional,stringlength(2|50)~Position too long"`
	SponsorID uint    `json:"sponsor_id"`
}
