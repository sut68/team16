package test

import (
	"strings"
	"testing"

	"backend/entity"
	"backend/validators"

	. "github.com/onsi/gomega"
)

// unit test Sponsor Contact

// case ปกติ
func TestSponsorContactValidation_AllValid(t *testing.T) {
	g := NewWithT(t)

	contact := entity.SponsorContact{
		Name:  "Somchai Rattanakul",
		Email: "somchai@test.com",
		Phone: "0812345678",
	}

	err := validators.ValidateStruct(&contact)

	// ผ่าน
	g.Expect(err).To(BeNil())
}

// case Name ว่าง
func TestSponsorContactValidation_NameRequired(t *testing.T) {
	g := NewWithT(t)

	contact := entity.SponsorContact{
		Name:  "",
		Email: "somchai@test.com",
		Phone: "0812345678",
	}

	err := validators.ValidateStruct(&contact)

	// ไม่ผ่าน
	g.Expect(err).ToNot(BeNil())
}

// case Name ไม่ถึง 2
func TestSponsorContactValidation_NameTooShort(t *testing.T) {
	g := NewWithT(t)

	contact := entity.SponsorContact{
		Name:  "A",
		Email: "somchai@test.com",
		Phone: "0812345678",
	}

	err := validators.ValidateStruct(&contact)

	// ไม่ผ่าน
	g.Expect(err).ToNot(BeNil())
}

// case Email ผิด format
func TestSponsorContactValidation_InvalidEmail(t *testing.T) {
	g := NewWithT(t)

	contact := entity.SponsorContact{
		Name:  "Somchai",
		Email: "invalid-email",
		Phone: "0812345678",
	}

	err := validators.ValidateStruct(&contact)

	// ไม่ผ่าน
	g.Expect(err).ToNot(BeNil())
}

// case Phone มีตัวอักษร
func TestSponsorContactValidation_PhoneNotNumeric(t *testing.T) {
	g := NewWithT(t)

	contact := entity.SponsorContact{
		Name:  "Somchai",
		Email: "somchai@test.com",
		Phone: "08A123456",
	}

	err := validators.ValidateStruct(&contact)

	// ไม่ผ่าน
	g.Expect(err).ToNot(BeNil())
}

// case Phone ไม่ถึง 9
func TestSponsorContactValidation_PhoneTooShort(t *testing.T) {
	g := NewWithT(t)

	contact := entity.SponsorContact{
		Name:  "Somchai",
		Email: "somchai@test.com",
		Phone: "08123456",
	}

	err := validators.ValidateStruct(&contact)

	// ไม่ผ่าน
	g.Expect(err).ToNot(BeNil())
}

// case Phone เกิน 10
func TestSponsorContactValidation_PhoneTooLong(t *testing.T) {
	g := NewWithT(t)

	contact := entity.SponsorContact{
		Name:  "Somchai",
		Email: "somchai@test.com",
		Phone: "081234567890",
	}

	err := validators.ValidateStruct(&contact)

	// ไม่ผ่าน
	g.Expect(err).ToNot(BeNil())
}

// case Position เป็น nil หรือ ไม่ส่ง
func TestSponsorContactValidation_PositionOptional(t *testing.T) {
	g := NewWithT(t)

	contact := entity.SponsorContact{
		Name: "Somchai",
		Email: "somchai@test.com",
		Phone: "0812345678",
	}

	err := validators.ValidateStruct(&contact)

	// ผ่าน
	g.Expect(err).To(BeNil())
}

// case Position เกิน 50
func TestSponsorContactValidation_PositionTooLong(t *testing.T) {
	g := NewWithT(t)

	longPos := strings.Repeat("P", 51)

	contact := entity.SponsorContact{
		Name:     "Somchai",
		Email:    "somchai@test.com",
		Phone:    "0812345678",
		Position: &longPos,
	}

	err := validators.ValidateStruct(&contact)

	// ไม่ผ่าน
	g.Expect(err).ToNot(BeNil())
}