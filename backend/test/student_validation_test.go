package test

import (
	"testing"
	"time"

	"backend/entity"
	"github.com/asaskevich/govalidator"
	. "github.com/onsi/gomega"
)
// 1. กรณี Positive (ข้อมูลถูกต้องทั้งหมด)
func TestStudentProfileCorrect(t *testing.T) {
	g := NewGomegaWithT(t)

	student := entity.StudentProfile{
		StudentID:        "B6630409",      // ถูกต้อง: ขึ้นต้น B ตามด้วยเลข 7 ตัว
		FirstNameTH:      "สมชาย",
		LastNameTH:       "ใจดี",
		FirstNameEN:      "Somchai",       // ถูกต้อง: ภาษาอังกฤษล้วน
		LastNameEN:       "Jaidee",
		NationalID:       "1234567890123", // ถูกต้อง: 13 หลัก
		BirthDate:        time.Now(),
		CurrentYear:      2,
		GPAX:             3.50,            // ถูกต้อง: 0.00 - 4.00
		AdvisorName:      "Dr.Smith",
		Phone:            "0812345678",
		Email:            "test@student.com",
		PermanentAddress: "Address 1",
		CurrentAddress:   "Address 2",
		Province:         "Nakhon Ratchasima",
		SiblingsCount:    1,
		UserID:           1,
		MajorID:          1,
	}

	// ตรวจสอบ
	ok, err := govalidator.ValidateStruct(student)
	
	// คาดหวัง: ผ่าน (true) และไม่มี error
	g.Expect(ok).To(BeTrue())
	g.Expect(err).To(BeNil())
}

// 2. กรณี Negative 1: Format รหัสนักศึกษาผิด
func TestStudentIDInvalidFormat(t *testing.T) {
	g := NewGomegaWithT(t)

	student := entity.StudentProfile{
		StudentID:   "K6630409", // ผิด: ขึ้นต้นด้วย K ไม่ได้ (ต้อง B,C,M,D)
		FirstNameTH: "สมชาย",
		FirstNameEN: "Somchai",
		LastNameEN:  "Jaidee",
		NationalID:  "1234567890123",
		Phone:       "0812345678",
		Email:       "test@student.com",
		GPAX:        3.50,
	}

	ok, err := govalidator.ValidateStruct(student)

	// คาดหวัง: ไม่ผ่าน (false) และเจอ Error message ที่กำหนด
	g.Expect(ok).To(BeFalse())
	g.Expect(err).ToNot(BeNil())
	g.Expect(err.Error()).To(ContainSubstring("Invalid Student ID format"))
}

// 3. กรณี Negative 2: GPAX เกินช่วงที่กำหนด
func TestGPAXOutOfRange(t *testing.T) {
	g := NewGomegaWithT(t)

	student := entity.StudentProfile{
		StudentID: "B6630409",
		GPAX:      4.50, // ผิด: เกรดเกิน 4.00
	}

	ok, err := govalidator.ValidateStruct(student)

	// คาดหวัง: ไม่ผ่าน (false)
	g.Expect(ok).To(BeFalse())
	g.Expect(err).ToNot(BeNil())
	g.Expect(err.Error()).To(ContainSubstring("GPAX must be between 0.00 and 4.00"))
}

// 4. กรณี Negative 3: ชื่อภาษาอังกฤษมีตัวเลขปน
func TestFirstNameENMustBeAlpha(t *testing.T) {
	g := NewGomegaWithT(t)

	student := entity.StudentProfile{
		StudentID:   "B6630409",
		FirstNameEN: "Somchai123", // ผิด: มีตัวเลข (ต้องเป็นตัวอักษรเท่านั้น)
	}

	ok, err := govalidator.ValidateStruct(student)

	// คาดหวัง: ไม่ผ่าน (false)
	g.Expect(ok).To(BeFalse())
	g.Expect(err).ToNot(BeNil())
	g.Expect(err.Error()).To(ContainSubstring("First Name EN must be English letters"))
}