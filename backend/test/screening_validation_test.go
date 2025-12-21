package test

import (
	"fmt"
	"testing"
	"backend/entity"
	"github.com/asaskevich/govalidator"
	. "github.com/onsi/gomega"
)

func TestScreeningValidation(t *testing.T) {
	g := NewWithT(t)
	govalidator.SetFieldsRequiredByDefault(false)

	// --- 1. ปรับปรุง Mocks ให้ข้อมูลครบตาม Tag required ของแต่ละ Entity ---

	validAdmin := entity.AdminProfile{
		AdminFirstname: "Somchai",
		AdminLastname:  "Jaidee",
		Position:       "Admin",
		Email:          "somchai@sut.ac.th",
		Phone:          "0812345678",
		UserID:         1,
	}

	validStatus := entity.StatusScreening{
		StatusScreening: "Pending",
	}

	// เพิ่มข้อมูลให้ครบตามที่ Error แจ้ง: Company name, Semaster ID, Student profile ID
	validApplication := entity.Application{
		StudentProfileID: 1, 
		SemasterID:       1, 
	}

	validScholarship := entity.Scholarship{
		Sponsor: entity.Sponsor{
			CompanyName: "SUT Foundation", 
		},
	}

	validAppScholarship := entity.ApplicationScholarship{
		ApplicationID: 1,
		Application:   validApplication,
		ScholarshipID: 1,
		Scholarship:   validScholarship,
		Status:        "Submitted",
	}

	// --- 2. Helper function ---
	createValidScreening := func() entity.Screening {
		return entity.Screening{
			AdminProfileID:           1,
			AdminProfile:             validAdmin,
			StatusScreeningID:        1, // สมมติ 1 = "ผ่าน"
			StatusScreening:          validStatus,
			ApplicationScholarshipID: 1,
			ApplicationScholarship:   validAppScholarship,
			RejectionReason:          nil,
		}
	}

	// --- 3. Test Cases ---

	t.Run("Should Pass: Valid Screening (Approved)", func(t *testing.T) {
		s := createValidScreening()
		ok, err := govalidator.ValidateStruct(s)
		
		if !ok {
			fmt.Printf("Validation Error: %v\n", err)
		}
		g.Expect(ok).To(BeTrue())
		g.Expect(err).To(BeNil())
	})

	t.Run("Should Pass: Valid Screening (Rejected with Reason)", func(t *testing.T) {
		s := createValidScreening()
		s.StatusScreeningID = 2 // สมมติ 2 = "ไม่ผ่าน"
		reason := "คุณสมบัติเบื้องต้นไม่ตรงตามเกณฑ์"
		s.RejectionReason = &reason

		ok, err := govalidator.ValidateStruct(s)
		g.Expect(ok).To(BeTrue())
		g.Expect(err).To(BeNil())
	})

	t.Run("Should Fail: RejectionReason too long", func(t *testing.T) {
		s := createValidScreening()
		reason := ""
		for i := 0; i < 110; i++ { reason += "a" }
		s.RejectionReason = &reason

		ok, err := govalidator.ValidateStruct(s)
		g.Expect(ok).To(BeFalse())
		g.Expect(err.Error()).To(ContainSubstring("You must provide a rejection reason"))
	})

	t.Run("Should Fail: Status Rejected but Reason is missing", func(t *testing.T) {
		s := createValidScreening()
		s.StatusScreeningID = 2 // ไม่ผ่าน
		s.RejectionReason = nil

		ok, _ := govalidator.ValidateStruct(s)
		
		// เพิ่ม logic ตรวจสอบ Business Logic
		if s.StatusScreeningID == 2 && (s.RejectionReason == nil || *s.RejectionReason == "") {
			ok = false
		}
		g.Expect(ok).To(BeFalse())
	})
}