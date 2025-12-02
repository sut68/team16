package entity

import (
	"gorm.io/gorm"
)

type Sponsor struct {
	gorm.Model
	CompanyName		string						`json:"company_name" gorm:"not null" valid:"required~Company is required"`

	IndustryID		*uint							`json:"industry_id"`
	Industry			*SponsorIndustry	`json:"industry"`

	Website				*string						`json:"website"`
	Status				string						`json:"status" gorm:"default:'active'"`
	Description		*string						`json:"description"`

	Contacts			[]SponsorContact	`json:"contacts" gorm:"foreignKey:SponsorID"`
}