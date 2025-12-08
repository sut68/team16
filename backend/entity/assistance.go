package entity

import (
	"gorm.io/gorm"
	"gorm.io/datatypes"
)

type Assistance struct {
	gorm.Model
	Massage datatypes.JSON `gorm:"type:jsonb"`
	
	AdminID uint `json:"admin_id"`
	AdminProfile AdminProfile `gorm:"foreignKey:AdminID"`
	

	StudentID uint `json:"student_id"`
	StudentProfile StudentProfile `gorm:"foreignKey:StudentID"`
	

	
}
