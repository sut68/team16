package entity

import (
	"gorm.io/gorm"
)

type SponsorContact struct {
	gorm.Model
	Name			string			`json:"name"`
	Email			string			`json:"email"`
	Phone			string			`json:"phone"`
	Position	*string			`json:"position"`
	SponsorID	uint				`json:"sponsor_id"`
}