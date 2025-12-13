package entity

import "gorm.io/gorm"

type FamilyInfo struct {
	gorm.Model
	FatherName         string  `json:"father_name"`
	FatherOccupation   string  `json:"father_occupation"`
	FatherIncome       float64 `gorm:"type:decimal(10,2)" json:"father_income"`
	
	MotherName         string  `json:"mother_name"`
	MotherIncome       float64 `gorm:"type:decimal(10,2)" json:"mother_income"`
	MotherOccupation   string  `json:"mother_occupation"`
	
	GuardianName       string  `json:"guardian_name"`
	GuardianOccupation string  `json:"guardian_occupation"`
	GuardianIncome     float64 `gorm:"type:decimal(10,2)" json:"guardian_income"`
	GuardianRelation   string  `json:"guardian_relation"`

	// Foreign Key - เพิ่ม references:ID เข้ามา
	ProfileID uint           `json:"profile_id"`
	Profile   StudentProfile `gorm:"foreignKey:ProfileID;references:ID" json:"student_profile"`
}