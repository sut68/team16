package test

import (
	"testing"
	"backend/entity"
	"github.com/asaskevich/govalidator"
	. "github.com/onsi/gomega"
)

func TestNewsPostValidation(t *testing.T) {
	g := NewWithT(t)
	govalidator.SetFieldsRequiredByDefault(false)

	// --- Base Mock ---
	validAdmin := entity.AdminProfile{
		AdminFirstname: "Somchai",
		AdminLastname:  "Jaidee",
		Position:       "Admin",
		Email:          "somchai@sut.ac.th",
		Phone:          "0812345678",
		UserID:         1,
	}

	validScholarship := entity.Scholarship{
		Sponsor: entity.Sponsor{
			CompanyName: "SUT Foundation",
		},
	}

	createValidNewsPost := func() entity.NewsPost {
		return entity.NewsPost{
			Title:         "ทุนการศึกษา 2568",
			FilePath:      "uploads/news/file.pdf",
			PostDetail:    "รายละเอียดทุนการศึกษาที่มีความยาวมากกว่า 10 ตัวอักษร",
			AdminID:       1,
			Admin:         validAdmin,
			ScholarshipID: 1,
			Scholarship:   validScholarship,
			StatusNewsID:  1,
		}
	}

	// --- Test Case 1: Valid Data ---
	t.Run("Valid Data", func(t *testing.T) {
		news := createValidNewsPost()
		ok, err := govalidator.ValidateStruct(news)
		if ok {
			t.Logf("PASS: Valid Data")
		} else {
			t.Logf("FAIL: Valid Data - %v", err)
		}
		g.Expect(ok).To(BeTrue())
		g.Expect(err).To(BeNil())
	})

	// --- Test Case 2: Title too short ---
	t.Run("Title too short", func(t *testing.T) {
		news := createValidNewsPost()
		news.Title = "Abcd"
		ok, err := govalidator.ValidateStruct(news)
		if !ok {
			t.Logf("PASS (expected fail): Title too short - %v", err)
		} else {
			t.Logf("FAIL: Title too short")
		}
		g.Expect(ok).To(BeFalse())
		g.Expect(err).ToNot(BeNil())
	})

	// --- Test Case 3: Title empty ---
	t.Run("Title empty", func(t *testing.T) {
		news := createValidNewsPost()
		news.Title = ""
		ok, err := govalidator.ValidateStruct(news)
		if !ok {
			t.Logf("PASS (expected fail): Title empty - %v", err)
		} else {
			t.Logf("FAIL: Title empty")
		}
		g.Expect(ok).To(BeFalse())
		g.Expect(err).ToNot(BeNil())
	})

	// --- Test Case 4: FilePath empty ---
	t.Run("FilePath empty", func(t *testing.T) {
		news := createValidNewsPost()
		news.FilePath = ""
		ok, err := govalidator.ValidateStruct(news)
		if !ok {
			t.Logf("PASS (expected fail): FilePath empty - %v", err)
		} else {
			t.Logf("FAIL: FilePath empty")
		}
		g.Expect(ok).To(BeFalse())
		g.Expect(err).ToNot(BeNil())
	})

	// --- Test Case 5: PostDetail too short ---
	t.Run("PostDetail too short", func(t *testing.T) {
		news := createValidNewsPost()
		news.PostDetail = "Short"
		ok, err := govalidator.ValidateStruct(news)
		if !ok {
			t.Logf("PASS (expected fail): PostDetail too short - %v", err)
		} else {
			t.Logf("FAIL: PostDetail too short")
		}
		g.Expect(ok).To(BeFalse())
		g.Expect(err).ToNot(BeNil())
	})

	// --- Test Case 6: StatusNewsID out of range ---
	t.Run("StatusNewsID out of range", func(t *testing.T) {
		news := createValidNewsPost()
		news.StatusNewsID = 9
		ok, err := govalidator.ValidateStruct(news)
		if !ok {
			t.Logf("PASS (expected fail): StatusNewsID out of range - %v", err)
		} else {
			t.Logf("FAIL: StatusNewsID out of range")
		}
		g.Expect(ok).To(BeFalse())
		g.Expect(err).ToNot(BeNil())
	})

	// --- Test Case 7: AdminID missing ---
	t.Run("AdminID missing", func(t *testing.T) {
		news := createValidNewsPost()
		news.AdminID = 0
		ok, err := govalidator.ValidateStruct(news)
		if !ok {
			t.Logf("PASS (expected fail): AdminID missing - %v", err)
		} else {
			t.Logf("FAIL: AdminID missing")
		}
		g.Expect(ok).To(BeFalse())
		g.Expect(err).ToNot(BeNil())
	})
}
