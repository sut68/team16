package entity

import(
	"gorm.io/gorm"
)

type Typefeature struct{
	gorm.Model
	
	Typefeaturename string `json:"type_feature_name"`

	
}