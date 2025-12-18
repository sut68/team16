package entity

import "gorm.io/gorm"

type InterviewMode struct {
	gorm.Model
	Name string `gorm:"unique;not null" json:"name"`

	InterviewRounds []InterviewRound `gorm:"foreignKey:InterviewModeID" json:"interview_rounds"`
}
