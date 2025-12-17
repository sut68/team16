package test

import (
	"testing"
	"backend/entity"
	"backend/validators"
	. "github.com/onsi/gomega"
)

func TestFamilyInfoValidation(t *testing.T) {
	g := NewWithT(t)

	// Mock Data
	validFamily := entity.FamilyInfo{
		FatherName: "Father", FatherOccupation: "Worker", FatherIncome: 20000,
		MotherName: "Mother", MotherOccupation: "Housewife", MotherIncome: 0,
		GuardianName: "Uncle", GuardianOccupation: "Business", GuardianIncome: 50000, GuardianRelation: "Uncle",
		ProfileID: 1,
	}

	// -----------------------------------------
	// กลุ่ม Test: รายได้ (Income)
	// -----------------------------------------
	t.Run("Check Income Logic", func(t *testing.T) {
		// Case: รายได้พ่อติดลบ
		f1 := validFamily
		f1.FatherIncome = -100
		err := validators.ValidateStruct(&f1)
		g.Expect(err).ToNot(BeNil())
		g.Expect(err.Error()).To(ContainSubstring("Father Income cannot be negative"))

		// Case: รายได้แม่ติดลบ
		f2 := validFamily
		f2.MotherIncome = -1
		g.Expect(validators.ValidateStruct(&f2)).ToNot(BeNil())

		// Case: รายได้เป็น 0 (ต้องผ่าน)
		f3 := validFamily
		f3.FatherIncome = 0
		g.Expect(validators.ValidateStruct(&f3)).To(BeNil())
	})

	// -----------------------------------------
	// กลุ่ม Test: ข้อมูลที่จำเป็น (Required)
	// -----------------------------------------
	t.Run("Check Required Fields", func(t *testing.T) {
		// Case: ชื่อพ่อว่าง
		f := validFamily
		f.FatherName = ""
		err := validators.ValidateStruct(&f)
		g.Expect(err).ToNot(BeNil())
		g.Expect(err.Error()).To(ContainSubstring("Father Name is required"))
	})
}