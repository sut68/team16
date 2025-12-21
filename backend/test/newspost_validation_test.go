package test

import (
	"testing"

	"backend/entity"
	"github.com/asaskevich/govalidator"
	. "github.com/onsi/gomega"
)

//go test -v file_name.go -run TestFunctionName

func TestNewsPostValidation(t *testing.T) {
	govalidator.SetFieldsRequiredByDefault(false)

	createValidNewsPost := func() entity.NewsPost {
		return entity.NewsPost{
			Title:         "Test Title",
			FilePath:      "uploads/file.pdf",
			PostDetail:    "Some post detail",
			AdminID:       1,
			ScholarshipID: 1,
			StatusNewsID:  1,
		}
	}

	// ---------- Valid ----------
	t.Run("Valid data", func(t *testing.T) {
		g := NewWithT(t)

		news := createValidNewsPost()
		ok, err := govalidator.ValidateStruct(news)

		t.Log("PASS: Valid data")

		g.Expect(ok).To(BeTrue())
		g.Expect(err).To(BeNil())
	})

	// ---------- FilePath ----------
	t.Run("FilePath empty", func(t *testing.T) {
		g := NewWithT(t)

		news := createValidNewsPost()
		news.FilePath = ""

		ok, err := govalidator.ValidateStruct(news)

		t.Log("PASS: FilePath is required")

		g.Expect(ok).To(BeFalse())
		g.Expect(err).ToNot(BeNil())
	})

	// ---------- Title ----------
	t.Run("Title empty", func(t *testing.T) {
		g := NewWithT(t)

		news := createValidNewsPost()
		news.Title = ""

		t.Log("PASS: Title is required")

		ok, err := govalidator.ValidateStruct(news)
		g.Expect(ok).To(BeFalse())
		g.Expect(err).ToNot(BeNil())
	})

	// ---------- AdminID ----------
	t.Run("AdminID missing", func(t *testing.T) {
		g := NewWithT(t)

		news := createValidNewsPost()
		news.AdminID = 0

		t.Log("PASS: AdminID is required")

		ok, err := govalidator.ValidateStruct(news)
		g.Expect(ok).To(BeFalse())
		g.Expect(err).ToNot(BeNil())
	})

	// ---------- Status ----------
	t.Run("StatusNewsID out of range", func(t *testing.T) {
		g := NewWithT(t)

		news := createValidNewsPost()
		news.StatusNewsID = 9

		t.Log("PASS: StatusNewsID invalid")

		ok, err := govalidator.ValidateStruct(news)
		g.Expect(ok).To(BeFalse())
		g.Expect(err).ToNot(BeNil())
	})
}
