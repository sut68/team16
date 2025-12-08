package entity

import(
	"gorm.io/gorm"
)

type Featurescholarship struct{
	gorm.Model
	
	Featurscholarshipename string `json:"feature_scholarship_name"`
	Operator string `json:"operator"`
	Value string `json:"value"`

	ScholarshipID uint `json:"scholarship_id"`
	Scholarship Scholarship `gorm:"foreignKey:ScholarshipID"`

	TypefeatureID uint `json:"typefeature_id"`
	Typefeature Typefeature `gorm:"foreignKey:TypefeatureID"`


}