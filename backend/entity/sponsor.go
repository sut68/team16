package entity

import (
	"gorm.io/gorm"
)

type Sponsor struct {
	gorm.Model
	CompanyName		string	`json:"company_name" gorm:"not null" valid:"required~Company name is required, stringlength(2|100)~Company name must be 2-100 characters"`

	IndustryID		*uint							`json:"industry_id" valid:"optional"`
	Industry			*SponsorIndustry	`json:"industry"`

	Website				*string						`json:"website" valid:"optional,url~Website must be a valid URL"`
	Status				string						`json:"status" gorm:"default:'active'" valid:"in(active|inactive)~Invalid status"`
	Description		*string						`json:"description" valid:"optional,stringlength(1|500)~Description too long"`

	Contacts			[]SponsorContact	`json:"contacts" gorm:"foreignKey:SponsorID" valid:"optional"`
}