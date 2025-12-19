package test

import (
	"testing"

	"backend/entity"
	"github.com/asaskevich/govalidator"
	. "github.com/onsi/gomega"
)
// 1. กรณี Positive (ข้อมูลถูกต้องทั้งหมด)
func TestFamilyInfoCorrect(t *testing.T) {
	g := NewGomegaWithT(t)

	family := entity.FamilyInfo{
		FatherName:         "Mr. Somchai",
		FatherOccupation:   "Engineer",
		FatherIncome:       50000.00,      // ถูกต้อง: > 0
		MotherName:         "Mrs. Somying",
		MotherOccupation:   "Doctor",
		MotherIncome:       60000.00,      // ถูกต้อง: > 0
		GuardianName:       "Mr. Somchai",
		GuardianOccupation: "Engineer",
		GuardianIncome:     50000.00,      // ถูกต้อง: > 0
		GuardianRelation:   "Father",
		ProfileID:          1,
	}

	// ตรวจสอบ
	ok, err := govalidator.ValidateStruct(family)

	// คาดหวัง: ผ่าน (true) และไม่มี error
	g.Expect(ok).To(BeTrue())
	g.Expect(err).To(BeNil())
}

// 2. กรณี Negative 1: รายได้บิดาติดลบ
func TestFatherIncomeNegative(t *testing.T) {
	g := NewGomegaWithT(t)

	family := entity.FamilyInfo{
		FatherName:       "Mr. Somchai",
		FatherOccupation: "Engineer",
		FatherIncome:     -5000.00, // ผิด: ติดลบไม่ได้ (valid:"range(0|10000000)")
		MotherName:       "Mrs. Somying",
		MotherOccupation: "Doctor",
		MotherIncome:     60000.00,
		GuardianName:     "Mr. Somchai",
		GuardianOccupation: "Engineer",
		GuardianIncome:   50000.00,
		GuardianRelation: "Father",
		ProfileID:        1,
	}

	ok, err := govalidator.ValidateStruct(family)

	// คาดหวัง: ไม่ผ่าน (false) และเจอ Error message ที่กำหนด
	g.Expect(ok).To(BeFalse())
	g.Expect(err).ToNot(BeNil())
	g.Expect(err.Error()).To(ContainSubstring("Father Income cannot be negative"))
}

// 3. กรณี Negative 2: รายได้มารดาติดลบ
func TestMotherIncomeNegative(t *testing.T) {
	g := NewGomegaWithT(t)

	family := entity.FamilyInfo{
		FatherName:       "Mr. Somchai",
		FatherOccupation: "Engineer",
		FatherIncome:     50000.00,
		MotherName:       "Mrs. Somying",
		MotherOccupation: "Doctor",
		MotherIncome:     -100.00, // ผิด: ติดลบไม่ได้
		GuardianName:     "Mr. Somchai",
		GuardianOccupation: "Engineer",
		GuardianIncome:   50000.00,
		GuardianRelation: "Father",
		ProfileID:        1,
	}

	ok, err := govalidator.ValidateStruct(family)

	// คาดหวัง: ไม่ผ่าน (false)
	g.Expect(ok).To(BeFalse())
	g.Expect(err).ToNot(BeNil())
	g.Expect(err.Error()).To(ContainSubstring("Mother Income cannot be negative"))
}

// 4. กรณี Negative 3: ข้อมูลจำเป็นเป็นค่าว่าง (Required)
func TestGuardianNameRequired(t *testing.T) {
	g := NewGomegaWithT(t)

	family := entity.FamilyInfo{
		FatherName:       "Mr. Somchai",
		FatherOccupation: "Engineer",
		FatherIncome:     50000.00,
		MotherName:       "Mrs. Somying",
		MotherOccupation: "Doctor",
		MotherIncome:     60000.00,
		GuardianName:     "", // ผิด: ห้ามเป็นค่าว่าง (valid:"required")
		GuardianOccupation: "Engineer",
		GuardianIncome:   50000.00,
		GuardianRelation: "Father",
		ProfileID:        1,
	}

	ok, err := govalidator.ValidateStruct(family)

	// คาดหวัง: ไม่ผ่าน (false)
	g.Expect(ok).To(BeFalse())
	g.Expect(err).ToNot(BeNil())
	g.Expect(err.Error()).To(ContainSubstring("Guardian Name is required"))
}