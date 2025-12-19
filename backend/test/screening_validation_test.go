package test

import (
	"testing"
	"backend/entity"
	"github.com/asaskevich/govalidator"
	. "github.com/onsi/gomega"
)

func TestScreeningValidation(t *testing.T) {
	g := NewWithT(t)
	govalidator.SetFieldsRequiredByDefault(false)

	// --- Base Mocks ---

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

	validScholarship := entity.Scholarship{
		Sponsor: entity.Sponsor{
			CompanyName: "SUT Foundation",
		},
	}

	validApplication := entity.Application{
		StudentProfileID: 1,
		SemasterID:       1,
	}

	validAppScholarship := entity.ApplicationScholarship{
		Status:        "Submitted",
		ApplicationID: 1,
		Application:   validApplication,
		ScholarshipID: 1,
		Scholarship:   validScholarship,
	}

	// --- Helper function ---
	createValidScreening := func() entity.Screening {
		return entity.Screening{
			AdminProfileID:           1,
			AdminProfile:             validAdmin,
			StatusScreeningID:        1,
			StatusScreening:          validStatus,
			ApplicationScholarshipID: 1,
			ApplicationScholarship:   validAppScholarship,
			RejectionReason:          nil,
		}
	}

	// --- Test Cases ---

	// 1. Valid Screening
	t.Run("Valid Screening", func(t *testing.T) {
		s := createValidScreening()
		ok, err := govalidator.ValidateStruct(s)
		if ok {
			t.Logf("PASS: Valid Screening")
		} else {
			t.Logf("FAIL: Valid Screening - %v", err)
		}
		g.Expect(ok).To(BeTrue())
		g.Expect(err).To(BeNil())
	})

	// 2. Missing AdminProfileID
	t.Run("Missing AdminProfileID", func(t *testing.T) {
		s := createValidScreening()
		s.AdminProfileID = 0
		ok, err := govalidator.ValidateStruct(s)
		if !ok {
			t.Logf("PASS (expected fail): Missing AdminProfileID - %v", err)
		} else {
			t.Logf("FAIL: Missing AdminProfileID")
		}
		g.Expect(ok).To(BeFalse())
	})

	// 3. Missing StatusScreeningID
	t.Run("Missing StatusScreeningID", func(t *testing.T) {
		s := createValidScreening()
		s.StatusScreeningID = 0
		ok, err := govalidator.ValidateStruct(s)
		if !ok {
			t.Logf("PASS (expected fail): Missing StatusScreeningID - %v", err)
		} else {
			t.Logf("FAIL: Missing StatusScreeningID")
		}
		g.Expect(ok).To(BeFalse())
	})

	// 4. Missing ApplicationScholarshipID
	t.Run("Missing ApplicationScholarshipID", func(t *testing.T) {
		s := createValidScreening()
		s.ApplicationScholarshipID = 0
		ok, err := govalidator.ValidateStruct(s)
		if !ok {
			t.Logf("PASS (expected fail): Missing ApplicationScholarshipID - %v", err)
		} else {
			t.Logf("FAIL: Missing ApplicationScholarshipID")
		}
		g.Expect(ok).To(BeFalse())
	})

	// 5. RejectionReason too long
	t.Run("RejectionReason too long", func(t *testing.T) {
		reason := ""
		for i := 0; i < 300; i++ {
			reason += "a"
		}
		s := createValidScreening()
		s.RejectionReason = &reason
		ok, err := govalidator.ValidateStruct(s)
		if !ok {
			t.Logf("PASS (expected fail): RejectionReason too long - %v", err)
		} else {
			t.Logf("FAIL: RejectionReason too long")
		}
		g.Expect(ok).To(BeFalse())
	})
}
