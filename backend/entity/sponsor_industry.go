package entity

import (
	"gorm.io/gorm"
)

type SponsorIndustry struct {
	gorm.Model
	Name		string		`json:"name" valid:"required~Industry is required"`
}