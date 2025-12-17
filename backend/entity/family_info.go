package entity

import "gorm.io/gorm"

type FamilyInfo struct {
	gorm.Model
	FatherName       string  `json:"father_name" valid:"required~Father Name is required"`
	FatherOccupation string  `json:"father_occupation" valid:"required~Father Occupation is required"`
	// รายได้พ่อ ห้ามติดลบ (สมมติ max 10 ล้านต่อเดือน)
	FatherIncome     float64 `gorm:"type:decimal(10,2)" json:"father_income" valid:"range(0|10000000)~Father Income cannot be negative"`

	MotherName       string  `json:"mother_name" valid:"required~Mother Name is required"`
	MotherOccupation string  `json:"mother_occupation" valid:"required~Mother Occupation is required"`
	MotherIncome     float64 `gorm:"type:decimal(10,2)" json:"mother_income" valid:"range(0|10000000)~Mother Income cannot be negative"`

	GuardianName       string  `json:"guardian_name" valid:"required~Guardian Name is required"`
	GuardianOccupation string  `json:"guardian_occupation" valid:"required~Guardian Occupation is required"`
	GuardianIncome     float64 `gorm:"type:decimal(10,2)" json:"guardian_income" valid:"range(0|10000000)~Guardian Income cannot be negative"`
	GuardianRelation   string  `json:"guardian_relation" valid:"required~Guardian Relation is required"`

	ProfileID uint           `json:"profile_id" valid:"required~Profile ID is required"`
	Profile   StudentProfile `gorm:"foreignKey:ProfileID;references:ID" json:"student_profile" valid:"-"`
}