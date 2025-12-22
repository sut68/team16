package entity

import "gorm.io/gorm"

// เกณฑ์ประจำรอบสัมภาษณ์
type InterviewRoundCriterion struct {
	gorm.Model
	Weight    float64 `gorm:"default:1.0" json:"weight" valid:"optional,range(0|10)~Weight must be between 0 and 10"`
	IsEnabled bool    `gorm:"default:true" json:"is_enabled" valid:"optional"`

	// Foreign Keys
	InterviewRoundID uint           `json:"interview_round_id" valid:"required~Interview Round ID is required"`
	InterviewRound   InterviewRound `gorm:"foreignKey:InterviewRoundID" json:"interview_round" valid:"-"`

	EvaluationCriterionID uint                `json:"evaluation_criterion_id" valid:"required~Evaluation Criterion ID is required"`
	EvaluationCriterion   EvaluationCriterion `gorm:"foreignKey:EvaluationCriterionID" json:"evaluation_criterion" valid:"-"`
}
