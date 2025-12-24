package test

import (
	"backend/entity"
	"backend/validators"
	"strings"
	"testing"

	. "github.com/onsi/gomega"
)

// unit test Evaluation

// case ปกติ - ข้อมูลครบถ้วน
func TestEvaluationValidation_AllValid(t *testing.T) {
	g := NewWithT(t)

	evaluation := entity.Evaluation{
		TotalScore:               85.5,
		Status:                   entity.EvaluationStatusCompleted,
		Remark:                   "ผู้สมัครมีคุณสมบัติเหมาะสม",
		InterviewRoundID:         1,
		ApplicationScholarshipID: 1,
		AdminID:                  1,
	}

	err := validators.ValidateStruct(&evaluation)

	// ผ่าน
	g.Expect(err).To(BeNil())
}

// case TotalScore เท่ากับ 0 (ผ่าน)
func TestEvaluationValidation_TotalScoreZero(t *testing.T) {
	g := NewWithT(t)

	evaluation := entity.Evaluation{
		TotalScore:               0,
		Status:                   entity.EvaluationStatusPending,
		InterviewRoundID:         1,
		ApplicationScholarshipID: 1,
		AdminID:                  1,
	}

	err := validators.ValidateStruct(&evaluation)

	// ผ่าน
	g.Expect(err).To(BeNil())
}

// case TotalScore เท่ากับ 100 (ผ่าน)
func TestEvaluationValidation_TotalScoreMax(t *testing.T) {
	g := NewWithT(t)

	evaluation := entity.Evaluation{
		TotalScore:               100,
		Status:                   entity.EvaluationStatusApproved,
		InterviewRoundID:         1,
		ApplicationScholarshipID: 1,
		AdminID:                  1,
	}

	err := validators.ValidateStruct(&evaluation)

	// ผ่าน
	g.Expect(err).To(BeNil())
}

// case TotalScore มากกว่า 100 (ไม่ผ่าน)
func TestEvaluationValidation_TotalScoreExceeds100(t *testing.T) {
	g := NewWithT(t)

	evaluation := entity.Evaluation{
		TotalScore:               101,
		Status:                   entity.EvaluationStatusCompleted,
		InterviewRoundID:         1,
		ApplicationScholarshipID: 1,
		AdminID:                  1,
	}

	err := validators.ValidateStruct(&evaluation)

	// ไม่ผ่าน
	g.Expect(err).ToNot(BeNil())
}

// case TotalScore ติดลบ (ไม่ผ่าน)
func TestEvaluationValidation_TotalScoreNegative(t *testing.T) {
	g := NewWithT(t)

	evaluation := entity.Evaluation{
		TotalScore:               -1,
		Status:                   entity.EvaluationStatusPending,
		InterviewRoundID:         1,
		ApplicationScholarshipID: 1,
		AdminID:                  1,
	}

	err := validators.ValidateStruct(&evaluation)

	// ไม่ผ่าน
	g.Expect(err).ToNot(BeNil())
}

// case Status pending (ผ่าน)
func TestEvaluationValidation_StatusPending(t *testing.T) {
	g := NewWithT(t)

	evaluation := entity.Evaluation{
		TotalScore:               0,
		Status:                   entity.EvaluationStatusPending,
		InterviewRoundID:         1,
		ApplicationScholarshipID: 1,
		AdminID:                  1,
	}

	err := validators.ValidateStruct(&evaluation)

	// ผ่าน
	g.Expect(err).To(BeNil())
}

// case Status in_progress (ผ่าน)
func TestEvaluationValidation_StatusInProgress(t *testing.T) {
	g := NewWithT(t)

	evaluation := entity.Evaluation{
		TotalScore:               50,
		Status:                   entity.EvaluationStatusInProgress,
		InterviewRoundID:         1,
		ApplicationScholarshipID: 1,
		AdminID:                  1,
	}

	err := validators.ValidateStruct(&evaluation)

	// ผ่าน
	g.Expect(err).To(BeNil())
}

// case Status completed (ผ่าน)
func TestEvaluationValidation_StatusCompleted(t *testing.T) {
	g := NewWithT(t)

	evaluation := entity.Evaluation{
		TotalScore:               80,
		Status:                   entity.EvaluationStatusCompleted,
		InterviewRoundID:         1,
		ApplicationScholarshipID: 1,
		AdminID:                  1,
	}

	err := validators.ValidateStruct(&evaluation)

	// ผ่าน
	g.Expect(err).To(BeNil())
}

// case Status approved (ผ่าน)
func TestEvaluationValidation_StatusApproved(t *testing.T) {
	g := NewWithT(t)

	evaluation := entity.Evaluation{
		TotalScore:               90,
		Status:                   entity.EvaluationStatusApproved,
		InterviewRoundID:         1,
		ApplicationScholarshipID: 1,
		AdminID:                  1,
	}

	err := validators.ValidateStruct(&evaluation)

	// ผ่าน
	g.Expect(err).To(BeNil())
}

// case Status rejected (ผ่าน)
func TestEvaluationValidation_StatusRejected(t *testing.T) {
	g := NewWithT(t)

	evaluation := entity.Evaluation{
		TotalScore:               40,
		Status:                   entity.EvaluationStatusRejected,
		InterviewRoundID:         1,
		ApplicationScholarshipID: 1,
		AdminID:                  1,
	}

	err := validators.ValidateStruct(&evaluation)

	// ผ่าน
	g.Expect(err).To(BeNil())
}

// case Status ไม่ถูกต้อง (ไม่ผ่าน)
func TestEvaluationValidation_StatusInvalid(t *testing.T) {
	g := NewWithT(t)

	evaluation := entity.Evaluation{
		TotalScore:               50,
		Status:                   "invalid_status",
		InterviewRoundID:         1,
		ApplicationScholarshipID: 1,
		AdminID:                  1,
	}

	err := validators.ValidateStruct(&evaluation)

	// ไม่ผ่าน
	g.Expect(err).ToNot(BeNil())
}

// case Remark ไม่เกิน 1000 ตัว (ผ่าน)
func TestEvaluationValidation_RemarkValid(t *testing.T) {
	g := NewWithT(t)

	evaluation := entity.Evaluation{
		TotalScore:               75,
		Status:                   entity.EvaluationStatusCompleted,
		Remark:                   "ผู้สมัครมีคุณสมบัติดี มีความมุ่งมั่นสูง",
		InterviewRoundID:         1,
		ApplicationScholarshipID: 1,
		AdminID:                  1,
	}

	err := validators.ValidateStruct(&evaluation)

	// ผ่าน
	g.Expect(err).To(BeNil())
}

// case Remark ยาวเกิน 1000 ตัว (ไม่ผ่าน)
func TestEvaluationValidation_RemarkTooLong(t *testing.T) {
	g := NewWithT(t)

	longRemark := strings.Repeat("R", 1001)

	evaluation := entity.Evaluation{
		TotalScore:               75,
		Status:                   entity.EvaluationStatusCompleted,
		Remark:                   longRemark,
		InterviewRoundID:         1,
		ApplicationScholarshipID: 1,
		AdminID:                  1,
	}

	err := validators.ValidateStruct(&evaluation)

	// ไม่ผ่าน
	g.Expect(err).ToNot(BeNil())
}

// case ไม่ส่ง InterviewRoundID (ไม่ผ่าน)
func TestEvaluationValidation_InterviewRoundIDRequired(t *testing.T) {
	g := NewWithT(t)

	evaluation := entity.Evaluation{
		TotalScore:               75,
		Status:                   entity.EvaluationStatusCompleted,
		ApplicationScholarshipID: 1,
		AdminID:                  1,
	}

	err := validators.ValidateStruct(&evaluation)

	// ไม่ผ่าน
	g.Expect(err).ToNot(BeNil())
}

// case ไม่ส่ง ApplicationScholarshipID (ไม่ผ่าน)
func TestEvaluationValidation_ApplicationScholarshipIDRequired(t *testing.T) {
	g := NewWithT(t)

	evaluation := entity.Evaluation{
		TotalScore:       75,
		Status:           entity.EvaluationStatusCompleted,
		InterviewRoundID: 1,
		AdminID:          1,
	}

	err := validators.ValidateStruct(&evaluation)

	// ไม่ผ่าน
	g.Expect(err).ToNot(BeNil())
}

// case ไม่ส่ง AdminID (ไม่ผ่าน)
func TestEvaluationValidation_AdminIDRequired(t *testing.T) {
	g := NewWithT(t)

	evaluation := entity.Evaluation{
		TotalScore:               75,
		Status:                   entity.EvaluationStatusCompleted,
		InterviewRoundID:         1,
		ApplicationScholarshipID: 1,
	}

	err := validators.ValidateStruct(&evaluation)

	// ไม่ผ่าน
	g.Expect(err).ToNot(BeNil())
}

// case ไม่ส่ง Remark (optional - ผ่าน)
func TestEvaluationValidation_RemarkOptional(t *testing.T) {
	g := NewWithT(t)

	evaluation := entity.Evaluation{
		TotalScore:               75,
		Status:                   entity.EvaluationStatusCompleted,
		InterviewRoundID:         1,
		ApplicationScholarshipID: 1,
		AdminID:                  1,
	}

	err := validators.ValidateStruct(&evaluation)

	// ผ่าน
	g.Expect(err).To(BeNil())
}
