package entity

import (
	"time"

	"gorm.io/gorm"
)

type ApprovalDecision struct {
	gorm.Model
	DecisionAt time.Time `gorm:"not null" json:"decision_at" valid:"required~Decision date is required"`
	Decision     string `gorm:"not null" json:"decision" valid:"required~Decision is required"`
	Comment     string `gorm:"not null" json:"comment"`

	TaskID       uint         `json:"task_id" valid:"required~Task ID is required"`
	ApprovalTask ApprovalTask `gorm:"foreignKey:TaskID" json:"approval_task" valid:"-"`

	AdminID uint         `json:"admin_id" valid:"required~Admin ID is required"`
	Admin   AdminProfile `gorm:"foreignKey:AdminID" json:"admin_profile" valid:"-"`
}
