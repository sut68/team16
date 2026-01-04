package entity

import "gorm.io/gorm"

// สร้าง Custom Type จำกัดค่าของ ประเภท
type ScoreType string

const (
	ScoreTypeNumeric  ScoreType = "numeric"   // คะแนนตัวเลข
	ScoreTypeGrade    ScoreType = "grade"     // เกรด A-F
	ScoreTypePassFail ScoreType = "pass_fail" // ผ่าน/ไม่ผ่าน
)

// เกณฑ์การประเมิน
type EvaluationCriterion struct {
	gorm.Model
	Name        string    `gorm:"not null" json:"name" valid:"required~Name is required,stringlength(2|100)~Name must be 2-100 characters"`
	Description string    `json:"description" valid:"optional,stringlength(0|500)~Description too long"`
	ScoreType   ScoreType `gorm:"type:varchar(20);default:'numeric'" json:"score_type" valid:"optional,in(numeric|grade|pass_fail)~Invalid score type"`
	MaxScore    float64   `gorm:"default:100" json:"max_score" valid:"optional,range(0|1000)~MaxScore must be between 0 and 1000"`
	Weight      float64   `gorm:"default:1.0" json:"weight" valid:"optional,range(0|10)~Weight must be between 0 and 10"`
	IsActive    bool      `gorm:"default:true" json:"is_active" valid:"optional"`

	// Relations
	InterviewRoundCriteria []InterviewRoundCriterion `gorm:"foreignKey:EvaluationCriterionID" json:"interview_round_criteria" valid:"-"`
	EvaluationScores       []EvaluationScore         `gorm:"foreignKey:EvaluationCriterionID" json:"evaluation_scores" valid:"-"`
}
