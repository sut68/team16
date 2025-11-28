package entity

import "gorm.io/gorm"

type ApprovalDecision struct {
	gorm.Model
	DecisionAt string `gorm:"not null" json:"decision_at"`
	Decision	 string `gorm:"not null" json:"decision"`
	Comment	 string `gorm:"not null" json:"comment"`

	TaskID uint `json:"task_id"`
	ApprovalTask ApprovalTask `gorm:"foreignKey:TaskID" json:"approval_task"`
}
