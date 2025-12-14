package test

import (
	"backend/entity"
	"backend/validators"
	"strings"
	"testing"

	. "github.com/onsi/gomega"
)

// unit test Sponsor

// case ปกติ
func TestSponsorValidation_AllValid(t *testing.T) {
	g := NewWithT(t)

	website := "https://www.pttplc.com"
	desc := "Scholarship sponsor"

	sponsor := entity.Sponsor {
		CompanyName: "PTT Public Company Limited",
		Website:     &website,
		Status:      "active",
		Description: &desc,
	}

	err := validators.ValidateStruct(&sponsor)

	// ผ่าน
	g.Expect(err).To(BeNil())
}

// case CompanyName ว่าง
func TestSponsorValidation_CompanyNameEmpty(t *testing.T) {
	g := NewWithT(t)

	sponsor := entity.Sponsor{
		CompanyName: "",
		Status:      "active",
	}

	err := validators.ValidateStruct(&sponsor)

	// ไม่ผ่าน
	g.Expect(err).ToNot(BeNil())
}

// case CompanyName สั้นเกิน
func TestSponsorValidation_CompanyNameTooShort(t *testing.T) {
	g := NewWithT(t)

	sponsor := entity.Sponsor{
		CompanyName: "A",
		Status:      "active",
	}

	err := validators.ValidateStruct(&sponsor)

	// ไม่ผ่าน
	g.Expect(err).ToNot(BeNil())
}

// case CompanyName ยาวเกิน 100 ตัว
func TestSponsorValidation_CompanyNameTooLong(t *testing.T) {
	g := NewWithT(t)

	longName := strings.Repeat("A", 101)

	sponsor := entity.Sponsor{
		CompanyName: longName,
		Status:      "active",
	}

	err := validators.ValidateStruct(&sponsor)

	// ไม่ผ่าน
	g.Expect(err).ToNot(BeNil())
}

// case Status inactive
func TestSponsorValidation_StatusInactiveValid(t *testing.T) {
	g := NewWithT(t)

	sponsor := entity.Sponsor{
		CompanyName: "Valid Company",
		Status:      "inactive",
	}

	err := validators.ValidateStruct(&sponsor)

	// ผ่าน
	g.Expect(err).To(BeNil())
}

// case Status ที่ไม่มี
func TestSponsorValidation_StatusInvalid(t *testing.T) {
	g := NewWithT(t)

	sponsor := entity.Sponsor{
		CompanyName: "Valid Company",
		Status:      "deleted",
	}

	err := validators.ValidateStruct(&sponsor)

	// ไม่ผ่าน
	g.Expect(err).ToNot(BeNil())
}

// case Website ไม่ใช่ URL
func TestSponsorValidation_InvalidWebsite(t *testing.T) {
	g := NewWithT(t)

	website := "not-a-url"

	sponsor := entity.Sponsor{
		CompanyName: "Valid Company",
		Status:      "active",
		Website:     &website,
	}

	err := validators.ValidateStruct(&sponsor)

	// ไม่ผ่าน
	g.Expect(err).ToNot(BeNil())
}

// case ไม่ส่ง Website
func TestSponsorValidation_WebsiteOptional(t *testing.T) {
	g := NewWithT(t)

	sponsor := entity.Sponsor{
		CompanyName: "Valid Company",
		Status:      "active",
	}

	err := validators.ValidateStruct(&sponsor)

	// ผ่าน
	g.Expect(err).To(BeNil())
}

// case Description ไม่เกิน 500
func TestSponsorValidation_DescriptionValid(t *testing.T) {
	g := NewWithT(t)

	desc := "OK"

	sponsor := entity.Sponsor{
		CompanyName: "Valid Company",
		Status:      "active",
		Description: &desc,
	}

	err := validators.ValidateStruct(&sponsor)

	// ผ่าน
	g.Expect(err).To(BeNil())
}

// case Description ยาวเกิน 500
func TestSponsorValidation_DescriptionTooLong(t *testing.T) {
	g := NewWithT(t)

	longDesc := strings.Repeat("D", 501)

	sponsor := entity.Sponsor{
		CompanyName: "Valid Company",
		Status:      "active",
		Description: &longDesc,
	}

	err := validators.ValidateStruct(&sponsor)

	// ไม่ผ่าน
	g.Expect(err).ToNot(BeNil())
}