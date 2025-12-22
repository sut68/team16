package entity

import "gorm.io/gorm"

// คะแนนรายข้อ
type EvaluationScore struct {
	gorm.Model
	ScoreValue float64 `json:"score_value" valid:"required~Score value is required"`
	Comment    string  `json:"comment" valid:"optional,stringlength(0|500)~Comment too long"`

	// Foreign Keys
	EvaluationID uint       `json:"evaluation_id" valid:"required~Evaluation ID is required"`
	Evaluation   Evaluation `gorm:"foreignKey:EvaluationID" json:"evaluation" valid:"-"`

	EvaluationCriterionID uint                `json:"evaluation_criterion_id" valid:"required~Evaluation Criterion ID is required"`
	EvaluationCriterion   EvaluationCriterion `gorm:"foreignKey:EvaluationCriterionID" json:"evaluation_criterion" valid:"-"`

	ScoringAdminID uint         `json:"scoring_admin_id" valid:"required~Scoring Admin ID is required"`
	ScoringAdmin   AdminProfile `gorm:"foreignKey:ScoringAdminID" json:"scoring_admin" valid:"-"`
}
