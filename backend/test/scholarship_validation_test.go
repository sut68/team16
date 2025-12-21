package test

import (
	"backend/entity"
	"backend/validators"
	"testing"

	. "github.com/onsi/gomega"
)

func TestScholarshipFullValidation(t *testing.T) {
	g := NewGomegaWithT(t)

	scholarship := entity.Scholarship{
		ScholarshipName: "ทุนเรียนดี",
		Description:     "สำหรับนักศึกษาที่มีผลการเรียนดี",
		OpenDate:        "2025-01-01",
		CloseDate:       "2025-02-01",

		StatusscholarshipID: 1,
		TypescholarshipID:   1,
		SemasterID:          1,
		SponsorID:           1,
	}

	err := validators.ValidateStruct(&scholarship)

	g.Expect(err).To(BeNil())
}

func TestScholarshipNameValidation(t *testing.T) {
	g := NewGomegaWithT(t)

	t.Run("Scholarship name is required", func(t *testing.T) {
		s := entity.Scholarship{
			ScholarshipName: "",	//ผิด
			Description:     "supermegathon",
			OpenDate:        "2025-01-01",
			CloseDate:       "2025-02-01",
			StatusscholarshipID: 1,
			TypescholarshipID:   1,
			SemasterID:          1,
			SponsorID:           1,
		}

		err := validators.ValidateStruct(s)

		g.Expect(err).NotTo(BeNil())
		g.Expect(err.Error()).To(Equal("Scholarship name is required"))
	})
}

func TestDescriptionValidation(t *testing.T) {
	g := NewGomegaWithT(t)

	t.Run("Description is required", func(t *testing.T) {
		s := entity.Scholarship{
			ScholarshipName: "ทุน",
			Description:     "",	//ผิด
			OpenDate:        "2025-01-01",
			CloseDate:       "2025-02-01",
			StatusscholarshipID: 1,
			TypescholarshipID:   1,
			SemasterID:          1,
			SponsorID:           1,
		}

		err := validators.ValidateStruct(s)

		g.Expect(err).NotTo(BeNil())
		g.Expect(err.Error()).To(Equal("Description is required"))
	})
}

func TestDateValidation(t *testing.T) {
	g := NewGomegaWithT(t)

	t.Run("Open date is required", func(t *testing.T) {
		s := entity.Scholarship{
			ScholarshipName: "ทุน",
			Description:     "supermegathon",
			OpenDate:        "",	//ผิด
			CloseDate:       "2025-02-01",
			StatusscholarshipID: 1,
			TypescholarshipID:   1,
			SemasterID:          1,
			SponsorID:           1,
		}

		err := validators.ValidateStruct(s)

		g.Expect(err).NotTo(BeNil())
		g.Expect(err.Error()).To(Equal("Open date is required"))
	})
}

func TestScholarshipForeignKeyValidation(t *testing.T) {
	g := NewGomegaWithT(t)

	t.Run("Invalid semaster", func(t *testing.T) {
		s := entity.Scholarship{
			ScholarshipName: "ทุน",
			Description:     "supermegathon",
			OpenDate:        "2025-01-01",
			CloseDate:       "2025-02-01",

			StatusscholarshipID: 1,
			TypescholarshipID:   1,
			SemasterID:          0, //ผิด
			SponsorID:           1,
		}

		err := validators.ValidateStruct(s)

		g.Expect(err).ToNot(BeNil())
	})
}

