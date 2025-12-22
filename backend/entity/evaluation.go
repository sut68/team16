package entity

import "gorm.io/gorm"

// สร้าง Custom Type จำกัดค่าของ สถานะ
type EvaluationStatus string

const (
	EvaluationStatusPending    EvaluationStatus = "pending"     // รอประเมิน
	EvaluationStatusInProgress EvaluationStatus = "in_progress" // กำลังประเมิน
	EvaluationStatusCompleted  EvaluationStatus = "completed"   // ประเมินเสร็จ
	EvaluationStatusApproved   EvaluationStatus = "approved"    // อนุมัติ
	EvaluationStatusRejected   EvaluationStatus = "rejected"    // ไม่อนุมัติ
)

// การประเมิน
type Evaluation struct {
	gorm.Model
	TotalScore float64          `json:"total_score" valid:"optional,range(0|100)~TotalScore must be between 0 and 100"`
	Status     EvaluationStatus `gorm:"type:varchar(20);default:'pending'" json:"status" valid:"optional,in(pending|in_progress|completed|approved|rejected)~Invalid status"`
	Remark     string           `json:"remark" valid:"optional,stringlength(0|1000)~Remark too long"`

	// Foreign Keys
	InterviewRoundID uint           `json:"interview_round_id" valid:"required~Interview Round ID is required"`
	InterviewRound   InterviewRound `gorm:"foreignKey:InterviewRoundID" json:"interview_round" valid:"-"`

	ApplicationScholarshipID uint                   `json:"application_scholarship_id" valid:"required~Application Scholarship ID is required"`
	ApplicationScholarship   ApplicationScholarship `gorm:"foreignKey:ApplicationScholarshipID" json:"application_scholarship" valid:"-"`

	AdminID      uint         `json:"admin_id" valid:"required~Admin ID is required"`
	AdminProfile AdminProfile `gorm:"foreignKey:AdminID" json:"admin_profile" valid:"-"`

	// Relations
	EvaluationScores []EvaluationScore `gorm:"foreignKey:EvaluationID" json:"evaluation_scores" valid:"-"`
}
