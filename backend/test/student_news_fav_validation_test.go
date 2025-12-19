package test

import (
	"backend/entity"
	"github.com/asaskevich/govalidator"
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

func TestStudentFavNewsValidation(t *testing.T) {
	g := NewWithT(t)
	govalidator.SetFieldsRequiredByDefault(false)
	// --- Base Mocks ---

	validNewsPost := entity.NewsPost{
		Title:         "ทุนการศึกษา 2568",
		FilePath:      "uploads/news/file.pdf",
		PostDetail:    "รายละเอียดทุนการศึกษาที่มีความยาวมากกว่า 10 ตัวอักษร",
		AdminID:       1,
		ScholarshipID: 1,
		StatusNewsID:  1,
	}
	validStudent := entity.StudentProfile{
		StudentID: "B6630652", // รหัสถูกต้อง
		FirstNameTH: "สมชาย", LastNameTH: "เรียนดี",
		FirstNameEN: "Somchai", LastNameEN: "Reandee",
		NationalID: "1100012345678",
		BirthDate:  time.Now(), CurrentYear: 1, GPAX: 3.50,
		AdvisorName: "Dr.Smith", Phone: "0812345678", Email: "b6630652@g.sut.ac.th",
		PermanentAddress: "Korat", CurrentAddress: "SUT Dorm", Province: "Nakhon Ratchasima",
		UserID: 1, MajorID: 1,
	}

	// --- Helper function ---
	createValidStudentFavNews := func() entity.StudentFavNews {
		return entity.StudentFavNews{
			NewsPostID:       1,
			NewsPost:         validNewsPost,
			StudentProfileID: 1,
			Student:          validStudent,
		}
	}
	// --- Test Case 1: Valid Data ---
	t.Run("Valid Data", func(t *testing.T) {
		studentFavNews := createValidStudentFavNews()
		ok, err := govalidator.ValidateStruct(studentFavNews)
		if ok {
			t.Logf("PASS: Valid Data")
		} else {
			t.Logf("FAIL: Valid Data - %v", err)
		}
		//g.Expect(ok).To(BeTrue())
		g.Expect(err).To(BeNil())
	})
	// --- Test Case 2: Missing NewsPostID ---
	t.Run("Missing NewsPostID", func(t *testing.T) {
		studentFavNews := createValidStudentFavNews()
		studentFavNews.NewsPostID = 0
		ok, err := govalidator.ValidateStruct(studentFavNews)
		if ok {
			t.Logf("FAIL: Missing NewsPostID - should be invalid")
		} else {
			t.Logf("PASS: Missing NewsPostID - %v", err)
		}
		g.Expect(ok).To(BeFalse())
		g.Expect(err).ToNot(BeNil())
	})
	// --- Test Case 3: Missing StudentProfileID ---
	t.Run("Missing StudentProfileID", func(t *testing.T) {
		studentFavNews := createValidStudentFavNews()
		studentFavNews.StudentProfileID = 0
		ok, err := govalidator.ValidateStruct(studentFavNews)
		if ok {
			t.Logf("FAIL: Missing StudentProfileID - should be invalid")
		} else {
			t.Logf("PASS: Missing StudentProfileID - %v", err)
		}
		g.Expect(ok).To(BeFalse())
		g.Expect(err).ToNot(BeNil())
	})

	// --- Test Case 4: Invalid StudentID format ---
    t.Run("Invalid StudentID format", func(t *testing.T) {
        studentFavNews := createValidStudentFavNews()
        studentFavNews.Student.StudentID = "A123" // ผิดกฎ (ต้อง B/C/M/D และ 8 หลัก)
        ok, err := govalidator.ValidateStruct(studentFavNews)
        
        g.Expect(ok).To(BeFalse())
        g.Expect(err.Error()).To(ContainSubstring("Invalid Student ID format"))
    })
}
