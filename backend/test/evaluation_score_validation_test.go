package test

import (
	"backend/entity"
	"backend/validators"
	"strings"
	"testing"

	. "github.com/onsi/gomega"
)

// unit test EvaluationScore

// case ปกติ - ข้อมูลครบถ้วน
func TestEvaluationScoreValidation_AllValid(t *testing.T) {
	g := NewWithT(t)

	score := entity.EvaluationScore{
		ScoreValue:            85.5,
		Comment:               "ผู้สมัครทำได้ดี",
		EvaluationID:          1,
		EvaluationCriterionID: 1,
		ScoringAdminID:        1,
	}

	err := validators.ValidateStruct(&score)

	// ผ่าน
	g.Expect(err).To(BeNil())
}

// case ScoreValue เท่ากับ 0 (ไม่ผ่าน - เป็น required และ zero value ถือว่าว่างเปล่า)
func TestEvaluationScoreValidation_ScoreValueZero(t *testing.T) {
	g := NewWithT(t)

	score := entity.EvaluationScore{
		ScoreValue:            0,
		EvaluationID:          1,
		EvaluationCriterionID: 1,
		ScoringAdminID:        1,
	}

	err := validators.ValidateStruct(&score)

	// ไม่ผ่าน (required field ไม่รับ zero value)
	g.Expect(err).ToNot(BeNil())
}

// case ScoreValue ค่าบวก (ผ่าน)
func TestEvaluationScoreValidation_ScoreValuePositive(t *testing.T) {
	g := NewWithT(t)

	score := entity.EvaluationScore{
		ScoreValue:            100,
		EvaluationID:          1,
		EvaluationCriterionID: 1,
		ScoringAdminID:        1,
	}

	err := validators.ValidateStruct(&score)

	// ผ่าน
	g.Expect(err).To(BeNil())
}

// case ScoreValue ค่าทศนิยม (ผ่าน)
func TestEvaluationScoreValidation_ScoreValueDecimal(t *testing.T) {
	g := NewWithT(t)

	score := entity.EvaluationScore{
		ScoreValue:            87.5,
		EvaluationID:          1,
		EvaluationCriterionID: 1,
		ScoringAdminID:        1,
	}

	err := validators.ValidateStruct(&score)

	// ผ่าน
	g.Expect(err).To(BeNil())
}

// case Comment ไม่เกิน 500 ตัว (ผ่าน)
func TestEvaluationScoreValidation_CommentValid(t *testing.T) {
	g := NewWithT(t)

	score := entity.EvaluationScore{
		ScoreValue:            80,
		Comment:               "ผู้สมัครมีทักษะการนำเสนอที่ดี สามารถตอบคำถามได้ชัดเจน",
		EvaluationID:          1,
		EvaluationCriterionID: 1,
		ScoringAdminID:        1,
	}

	err := validators.ValidateStruct(&score)

	// ผ่าน
	g.Expect(err).To(BeNil())
}

// case Comment ยาวเกิน 500 ตัว (ไม่ผ่าน)
func TestEvaluationScoreValidation_CommentTooLong(t *testing.T) {
	g := NewWithT(t)

	longComment := strings.Repeat("C", 501)

	score := entity.EvaluationScore{
		ScoreValue:            80,
		Comment:               longComment,
		EvaluationID:          1,
		EvaluationCriterionID: 1,
		ScoringAdminID:        1,
	}

	err := validators.ValidateStruct(&score)

	// ไม่ผ่าน
	g.Expect(err).ToNot(BeNil())
}

// case Comment พอดี 500 ตัว (ผ่าน)
func TestEvaluationScoreValidation_CommentMaxLength(t *testing.T) {
	g := NewWithT(t)

	maxComment := strings.Repeat("C", 500)

	score := entity.EvaluationScore{
		ScoreValue:            80,
		Comment:               maxComment,
		EvaluationID:          1,
		EvaluationCriterionID: 1,
		ScoringAdminID:        1,
	}

	err := validators.ValidateStruct(&score)

	// ผ่าน
	g.Expect(err).To(BeNil())
}

// case Comment optional (ผ่าน)
func TestEvaluationScoreValidation_CommentOptional(t *testing.T) {
	g := NewWithT(t)

	score := entity.EvaluationScore{
		ScoreValue:            75,
		EvaluationID:          1,
		EvaluationCriterionID: 1,
		ScoringAdminID:        1,
	}

	err := validators.ValidateStruct(&score)

	// ผ่าน
	g.Expect(err).To(BeNil())
}

// case ไม่ส่ง EvaluationID (ไม่ผ่าน)
func TestEvaluationScoreValidation_EvaluationIDRequired(t *testing.T) {
	g := NewWithT(t)

	score := entity.EvaluationScore{
		ScoreValue:            80,
		EvaluationCriterionID: 1,
		ScoringAdminID:        1,
	}

	err := validators.ValidateStruct(&score)

	// ไม่ผ่าน
	g.Expect(err).ToNot(BeNil())
}

// case ไม่ส่ง EvaluationCriterionID (ไม่ผ่าน)
func TestEvaluationScoreValidation_EvaluationCriterionIDRequired(t *testing.T) {
	g := NewWithT(t)

	score := entity.EvaluationScore{
		ScoreValue:     80,
		EvaluationID:   1,
		ScoringAdminID: 1,
	}

	err := validators.ValidateStruct(&score)

	// ไม่ผ่าน
	g.Expect(err).ToNot(BeNil())
}

// case ไม่ส่ง ScoringAdminID (ไม่ผ่าน)
func TestEvaluationScoreValidation_ScoringAdminIDRequired(t *testing.T) {
	g := NewWithT(t)

	score := entity.EvaluationScore{
		ScoreValue:            80,
		EvaluationID:          1,
		EvaluationCriterionID: 1,
	}

	err := validators.ValidateStruct(&score)

	// ไม่ผ่าน
	g.Expect(err).ToNot(BeNil())
}

// case ไม่ส่ง required fields ทั้งหมด (ไม่ผ่าน)
func TestEvaluationScoreValidation_AllRequiredFieldsMissing(t *testing.T) {
	g := NewWithT(t)

	score := entity.EvaluationScore{
		Comment: "ไม่มีข้อมูลสำคัญ",
	}

	err := validators.ValidateStruct(&score)

	// ไม่ผ่าน
	g.Expect(err).ToNot(BeNil())
}

// case ส่ง foreign keys ทั้งหมดถูกต้อง (ผ่าน)
func TestEvaluationScoreValidation_AllForeignKeysValid(t *testing.T) {
	g := NewWithT(t)

	score := entity.EvaluationScore{
		ScoreValue:            95,
		Comment:               "ยอดเยี่ยม",
		EvaluationID:          10,
		EvaluationCriterionID: 5,
		ScoringAdminID:        3,
	}

	err := validators.ValidateStruct(&score)

	// ผ่าน
	g.Expect(err).To(BeNil())
}
