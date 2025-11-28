package entity

import(
	"gorm.io/gorm"
	
)

type Statusscholarship struct{
	gorm.Model
	
	Statusname string `json:"status_name"`

	
}