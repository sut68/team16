package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null" json:"username"`
	Password string `gorm:"not null" json:"-"`

	RoleID *uint  `gorm:"index" json:"role_id,omitempty"`

	Role *Role `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"role,omitempty"`

	AdminProfiles []AdminProfile `gorm:"foreignKey:UserID" json:"admin_profile"`
	StudentProfiles []StudentProfile `gorm:"foreignKey:UserID" json:"student_profile"`
	Applications []Application `gorm:"foreignKey:UserID" json:"applications"`
}