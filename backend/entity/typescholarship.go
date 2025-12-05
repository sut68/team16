package entity

import(
	"gorm.io/gorm"
	
)

type Typescholarship struct{
	gorm.Model
	
	Typename string `json:"type_name"`

	
}