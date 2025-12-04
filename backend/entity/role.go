package entity

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	RoleName string `gorm:"unique;not null" json:"role_name"`

	// 1 Role มี User ได้หลายคน (One-to-Many)
	Users []User `gorm:"foreignKey:RoleID" json:"users"` 
}
