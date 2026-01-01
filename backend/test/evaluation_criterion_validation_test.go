package test

import (
	"backend/entity"
	"backend/validators"
	"strings"
	"testing"

	. "github.com/onsi/gomega"
)

// unit test EvaluationCriterion

// case ปกติ - ข้อมูลครบถ้วน
func TestEvaluationCriterionValidation_AllValid(t *testing.T) {
	g := NewWithT(t)

	criterion := entity.EvaluationCriterion{
		Name:        "ความเหมาะสมของผู้สมัคร",
		Description: "ประเมินความเหมาะสมโดยรวมของผู้สมัครกับทุนการศึกษา",
		ScoreType:   entity.ScoreTypeNumeric,
		MaxScore:    100,
		Weight:      1.5,
		IsActive:    true,
	}

	err := validators.ValidateStruct(&criterion)

	// ผ่าน
	g.Expect(err).To(BeNil())
}

// case Name ว่าง (ไม่ผ่าน)
func TestEvaluationCriterionValidation_NameEmpty(t *testing.T) {
	g := NewWithT(t)

	criterion := entity.EvaluationCriterion{
		Name:      "",
		ScoreType: entity.ScoreTypeNumeric,
		MaxScore:  100,
		Weight:    1.0,
	}

	err := validators.ValidateStruct(&criterion)

	// ไม่ผ่าน
	g.Expect(err).ToNot(BeNil())
}

// case Name สั้นเกิน (1 ตัว) (ไม่ผ่าน)
func TestEvaluationCriterionValidation_NameTooShort(t *testing.T) {
	g := NewWithT(t)

	criterion := entity.EvaluationCriterion{
		Name:      "A",
		ScoreType: entity.ScoreTypeNumeric,
		MaxScore:  100,
		Weight:    1.0,
	}

	err := validators.ValidateStruct(&criterion)

	// ไม่ผ่าน
	g.Expect(err).ToNot(BeNil())
}

// case Name ยาวเกิน 100 ตัว (ไม่ผ่าน)
func TestEvaluationCriterionValidation_NameTooLong(t *testing.T) {
	g := NewWithT(t)

	longName := strings.Repeat("A", 101)

	criterion := entity.EvaluationCriterion{
		Name:      longName,
		ScoreType: entity.ScoreTypeNumeric,
		MaxScore:  100,
		Weight:    1.0,
	}

	err := validators.ValidateStruct(&criterion)

	// ไม่ผ่าน
	g.Expect(err).ToNot(BeNil())
}

// case Name พอดี 2 ตัว (ผ่าน)
func TestEvaluationCriterionValidation_NameMinLength(t *testing.T) {
	g := NewWithT(t)

	criterion := entity.EvaluationCriterion{
		Name:      "AB",
		ScoreType: entity.ScoreTypeNumeric,
		MaxScore:  100,
		Weight:    1.0,
	}

	err := validators.ValidateStruct(&criterion)

	// ผ่าน
	g.Expect(err).To(BeNil())
}

// case Name พอดี 100 ตัว (ผ่าน)
func TestEvaluationCriterionValidation_NameMaxLength(t *testing.T) {
	g := NewWithT(t)

	maxName := strings.Repeat("A", 100)

	criterion := entity.EvaluationCriterion{
		Name:      maxName,
		ScoreType: entity.ScoreTypeNumeric,
		MaxScore:  100,
		Weight:    1.0,
	}

	err := validators.ValidateStruct(&criterion)

	// ผ่าน
	g.Expect(err).To(BeNil())
}

// case Description ไม่เกิน 500 ตัว (ผ่าน)
func TestEvaluationCriterionValidation_DescriptionValid(t *testing.T) {
	g := NewWithT(t)

	criterion := entity.EvaluationCriterion{
		Name:        "ผลการเรียน",
		Description: "ประเมินจากเกรดเฉลี่ยสะสมและผลการเรียนในรายวิชาที่เกี่ยวข้อง",
		ScoreType:   entity.ScoreTypeNumeric,
		MaxScore:    100,
		Weight:      2.0,
	}

	err := validators.ValidateStruct(&criterion)

	// ผ่าน
	g.Expect(err).To(BeNil())
}

// case Description ยาวเกิน 500 ตัว (ไม่ผ่าน)
func TestEvaluationCriterionValidation_DescriptionTooLong(t *testing.T) {
	g := NewWithT(t)

	longDesc := strings.Repeat("D", 501)

	criterion := entity.EvaluationCriterion{
		Name:        "ผลการเรียน",
		Description: longDesc,
		ScoreType:   entity.ScoreTypeNumeric,
		MaxScore:    100,
		Weight:      2.0,
	}

	err := validators.ValidateStruct(&criterion)

	// ไม่ผ่าน
	g.Expect(err).ToNot(BeNil())
}

// case ScoreType numeric (ผ่าน)
func TestEvaluationCriterionValidation_ScoreTypeNumeric(t *testing.T) {
	g := NewWithT(t)

	criterion := entity.EvaluationCriterion{
		Name:      "คะแนนทดสอบ",
		ScoreType: entity.ScoreTypeNumeric,
		MaxScore:  100,
		Weight:    1.0,
	}

	err := validators.ValidateStruct(&criterion)

	// ผ่าน
	g.Expect(err).To(BeNil())
}

// case ScoreType grade (ผ่าน)
func TestEvaluationCriterionValidation_ScoreTypeGrade(t *testing.T) {
	g := NewWithT(t)

	criterion := entity.EvaluationCriterion{
		Name:      "เกรดเฉลี่ย",
		ScoreType: entity.ScoreTypeGrade,
		MaxScore:  4,
		Weight:    1.0,
	}

	err := validators.ValidateStruct(&criterion)

	// ผ่าน
	g.Expect(err).To(BeNil())
}

// case ScoreType pass_fail (ผ่าน)
func TestEvaluationCriterionValidation_ScoreTypePassFail(t *testing.T) {
	g := NewWithT(t)

	criterion := entity.EvaluationCriterion{
		Name:      "ความประพฤติ",
		ScoreType: entity.ScoreTypePassFail,
		MaxScore:  1,
		Weight:    1.0,
	}

	err := validators.ValidateStruct(&criterion)

	// ผ่าน
	g.Expect(err).To(BeNil())
}

// case ScoreType ไม่ถูกต้อง (ไม่ผ่าน)
func TestEvaluationCriterionValidation_ScoreTypeInvalid(t *testing.T) {
	g := NewWithT(t)

	criterion := entity.EvaluationCriterion{
		Name:      "คะแนน",
		ScoreType: "invalid_type",
		MaxScore:  100,
		Weight:    1.0,
	}

	err := validators.ValidateStruct(&criterion)

	// ไม่ผ่าน
	g.Expect(err).ToNot(BeNil())
}

// case MaxScore เท่ากับ 0 (ผ่าน)
func TestEvaluationCriterionValidation_MaxScoreZero(t *testing.T) {
	g := NewWithT(t)

	criterion := entity.EvaluationCriterion{
		Name:      "ข้อสังเกต",
		ScoreType: entity.ScoreTypeNumeric,
		MaxScore:  0,
		Weight:    1.0,
	}

	err := validators.ValidateStruct(&criterion)

	// ผ่าน
	g.Expect(err).To(BeNil())
}

// case MaxScore เท่ากับ 1000 (ผ่าน)
func TestEvaluationCriterionValidation_MaxScoreMax(t *testing.T) {
	g := NewWithT(t)

	criterion := entity.EvaluationCriterion{
		Name:      "คะแนนรวม",
		ScoreType: entity.ScoreTypeNumeric,
		MaxScore:  1000,
		Weight:    1.0,
	}

	err := validators.ValidateStruct(&criterion)

	// ผ่าน
	g.Expect(err).To(BeNil())
}

// case MaxScore มากกว่า 1000 (ไม่ผ่าน)
func TestEvaluationCriterionValidation_MaxScoreExceeds(t *testing.T) {
	g := NewWithT(t)

	criterion := entity.EvaluationCriterion{
		Name:      "คะแนนรวม",
		ScoreType: entity.ScoreTypeNumeric,
		MaxScore:  1001,
		Weight:    1.0,
	}

	err := validators.ValidateStruct(&criterion)

	// ไม่ผ่าน
	g.Expect(err).ToNot(BeNil())
}

// case MaxScore ติดลบ (ไม่ผ่าน)
func TestEvaluationCriterionValidation_MaxScoreNegative(t *testing.T) {
	g := NewWithT(t)

	criterion := entity.EvaluationCriterion{
		Name:      "คะแนนรวม",
		ScoreType: entity.ScoreTypeNumeric,
		MaxScore:  -1,
		Weight:    1.0,
	}

	err := validators.ValidateStruct(&criterion)

	// ไม่ผ่าน
	g.Expect(err).ToNot(BeNil())
}

// case Weight เท่ากับ 0 (ผ่าน)
func TestEvaluationCriterionValidation_WeightZero(t *testing.T) {
	g := NewWithT(t)

	criterion := entity.EvaluationCriterion{
		Name:      "ข้อมูลเสริม",
		ScoreType: entity.ScoreTypeNumeric,
		MaxScore:  100,
		Weight:    0,
	}

	err := validators.ValidateStruct(&criterion)

	// ผ่าน
	g.Expect(err).To(BeNil())
}

// case Weight เท่ากับ 10 (ผ่าน)
func TestEvaluationCriterionValidation_WeightMax(t *testing.T) {
	g := NewWithT(t)

	criterion := entity.EvaluationCriterion{
		Name:      "เกณฑ์สำคัญมาก",
		ScoreType: entity.ScoreTypeNumeric,
		MaxScore:  100,
		Weight:    10,
	}

	err := validators.ValidateStruct(&criterion)

	// ผ่าน
	g.Expect(err).To(BeNil())
}

// case Weight มากกว่า 10 (ไม่ผ่าน)
func TestEvaluationCriterionValidation_WeightExceeds(t *testing.T) {
	g := NewWithT(t)

	criterion := entity.EvaluationCriterion{
		Name:      "เกณฑ์สำคัญมาก",
		ScoreType: entity.ScoreTypeNumeric,
		MaxScore:  100,
		Weight:    11,
	}

	err := validators.ValidateStruct(&criterion)

	// ไม่ผ่าน
	g.Expect(err).ToNot(BeNil())
}

// case Weight ติดลบ (ไม่ผ่าน)
func TestEvaluationCriterionValidation_WeightNegative(t *testing.T) {
	g := NewWithT(t)

	criterion := entity.EvaluationCriterion{
		Name:      "เกณฑ์",
		ScoreType: entity.ScoreTypeNumeric,
		MaxScore:  100,
		Weight:    -1,
	}

	err := validators.ValidateStruct(&criterion)

	// ไม่ผ่าน
	g.Expect(err).ToNot(BeNil())
}

// case Description optional (ผ่าน)
func TestEvaluationCriterionValidation_DescriptionOptional(t *testing.T) {
	g := NewWithT(t)

	criterion := entity.EvaluationCriterion{
		Name:      "เกณฑ์ทั่วไป",
		ScoreType: entity.ScoreTypeNumeric,
		MaxScore:  100,
		Weight:    1.0,
	}

	err := validators.ValidateStruct(&criterion)

	// ผ่าน
	g.Expect(err).To(BeNil())
}
