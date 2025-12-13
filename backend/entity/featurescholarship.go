package entity

import(
	"gorm.io/gorm"
)

type Featurescholarship struct{
	gorm.Model

	//แก้ให้ชื่อมันถูกต้อง - เติม e เข้ามา
	Featurescholarshipname string `json:"feature_scholarship_name"`

	Operator string `json:"operator"`
	Value string `json:"value"`

	ScholarshipID uint `json:"scholarship_id"`
	Scholarship Scholarship `gorm:"foreignKey:ScholarshipID"`

	TypefeatureID uint `json:"typefeature_id"`
	Typefeature Typefeature `gorm:"foreignKey:TypefeatureID"`
	
}