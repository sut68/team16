package test

import (
	"testing"
	"time"

	"backend/entity"
	"backend/validators"
	. "github.com/onsi/gomega"
)

func TestStudentProfileValidation(t *testing.T) {
	g := NewWithT(t)

	// Mock ข้อมูลที่ "ถูกต้อง" ไว้ก่อน (ใช้เป็นฐาน)
	validStudent := entity.StudentProfile{
		StudentID: "B6630409", // รหัสถูกต้อง
		FirstNameTH: "สมชาย", LastNameTH: "เรียนดี",
		FirstNameEN: "Somchai", LastNameEN: "Reandee",
		NationalID: "1100012345678",
		BirthDate:  time.Now(), CurrentYear: 1, GPAX: 3.50,
		AdvisorName: "Dr.Smith", Phone: "0812345678", Email: "b6630409@g.sut.ac.th",
		PermanentAddress: "Korat", CurrentAddress: "SUT Dorm", Province: "Nakhon Ratchasima",
		UserID: 1, MajorID: 1,
	}

	// -----------------------------------------
	// กลุ่ม Test: รหัสนักศึกษา (Student ID)
	// -----------------------------------------
	t.Run("Check Student ID Format", func(t *testing.T) {
		// Case 1: รหัสถูกต้อง (B - ป.ตรี)
		s1 := validStudent
		s1.StudentID = "B6600001"
		g.Expect(validators.ValidateStruct(&s1)).To(BeNil())

		// Case 2: รหัสถูกต้อง (M - ป.โท)
		s2 := validStudent
		s2.StudentID = "M6500001"
		g.Expect(validators.ValidateStruct(&s2)).To(BeNil())

		// Case 3: ผิด - ขึ้นต้นด้วยตัวอื่น (A)
		s3 := validStudent
		s3.StudentID = "A6600001"
		err := validators.ValidateStruct(&s3)
		g.Expect(err).ToNot(BeNil())
		g.Expect(err.Error()).To(ContainSubstring("Invalid Student ID format"))

		// Case 4: ผิด - ความยาวไม่ครบ
		s4 := validStudent
		s4.StudentID = "B66"
		g.Expect(validators.ValidateStruct(&s4)).ToNot(BeNil())
	})

	// -----------------------------------------
	// กลุ่ม Test: เกรดเฉลี่ย (GPAX)
	// -----------------------------------------
	t.Run("Check GPAX Range", func(t *testing.T) {
		// Case: GPAX เกิน 4.00
		s := validStudent
		s.GPAX = 4.01
		err := validators.ValidateStruct(&s)
		g.Expect(err).ToNot(BeNil())
		g.Expect(err.Error()).To(ContainSubstring("GPAX must be between 0.00 and 4.00"))

		// Case: GPAX ติดลบ
		s2 := validStudent
		s2.GPAX = -0.50
		g.Expect(validators.ValidateStruct(&s2)).ToNot(BeNil())
	})

	// -----------------------------------------
	// กลุ่ม Test: ชื่อภาษาอังกฤษ (Name EN)
	// -----------------------------------------
	t.Run("Check Name EN Characters", func(t *testing.T) {
		// Case: ชื่อมีภาษาไทยปน
		s := validStudent
		s.FirstNameEN = "Somchaiสมชาย"
		err := validators.ValidateStruct(&s)
		g.Expect(err).ToNot(BeNil())
		g.Expect(err.Error()).To(ContainSubstring("First Name EN must be English letters"))

		// Case: ชื่อมีตัวเลข
		s2 := validStudent
		s2.LastNameEN = "Reandee123"
		g.Expect(validators.ValidateStruct(&s2)).ToNot(BeNil())
	})

	// -----------------------------------------
	// กลุ่ม Test: บัตรประชาชน (National ID)
	// -----------------------------------------
	t.Run("Check National ID", func(t *testing.T) {
		// Case: ไม่ครบ 13 หลัก
		s := validStudent
		s.NationalID = "123456789"
		err := validators.ValidateStruct(&s)
		g.Expect(err).ToNot(BeNil())

		// Case: มีตัวอักษรปน
		s2 := validStudent
		s2.NationalID = "110001234567A"
		g.Expect(validators.ValidateStruct(&s2)).ToNot(BeNil())
	})
}