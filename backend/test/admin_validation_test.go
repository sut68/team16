package test

import (
	"testing"

	"backend/entity"
	"github.com/asaskevich/govalidator"
	. "github.com/onsi/gomega"
)
// 1. กรณี Positive (ข้อมูลถูกต้องทั้งหมด)
func TestAdminProfileCorrect(t *testing.T) {
	g := NewGomegaWithT(t)

	admin := entity.AdminProfile{
		AdminFirstname: "Admin",
		AdminLastname:  "System",
		Position:       "HR Manager",
		Email:          "admin@sut.ac.th", // ถูกต้อง: format email
		Phone:          "0891234567",      // ถูกต้อง: ตัวเลข 10 หลัก
		UserID:         1,
	}

	ok, err := govalidator.ValidateStruct(admin)
	g.Expect(ok).To(BeTrue())
	g.Expect(err).To(BeNil())
}

// 2. กรณี Negative 1: Email ผิดรูปแบบ
func TestAdminEmailInvalid(t *testing.T) {
	g := NewGomegaWithT(t)

	admin := entity.AdminProfile{
		AdminFirstname: "Admin",
		AdminLastname:  "System",
		Position:       "HR",
		Phone:          "0891234567",
		Email:          "admin-no-at-sign.com", // ผิด: ไม่มี @
	}

	ok, err := govalidator.ValidateStruct(admin)

	g.Expect(ok).To(BeFalse())
	g.Expect(err).ToNot(BeNil())
	g.Expect(err.Error()).To(ContainSubstring("Invalid email format"))
}

// 3. กรณี Negative 2: เบอร์โทรศัพท์ไม่ใช่ตัวเลข
func TestAdminPhoneNotNumeric(t *testing.T) {
	g := NewGomegaWithT(t)

	admin := entity.AdminProfile{
		AdminFirstname: "Admin",
		Phone:          "089-123-45", // ผิด: มีขีด (-) หรือตัวหนังสือ
		Email:          "test@email.com",
	}

	ok, err := govalidator.ValidateStruct(admin)

	g.Expect(ok).To(BeFalse())
	g.Expect(err).ToNot(BeNil())
	g.Expect(err.Error()).To(ContainSubstring("Phone must handle numbers only"))
}

// 4. กรณี Negative 3: ข้อมูลจำเป็นเป็นค่าว่าง (Required)
func TestAdminFirstnameRequired(t *testing.T) {
	g := NewGomegaWithT(t)

	admin := entity.AdminProfile{
		AdminFirstname: "", // ผิด: เป็นค่าว่างไม่ได้ (valid:"required")
		AdminLastname:  "System",
		Position:       "HR",
		Email:          "test@email.com",
		Phone:          "0891234567",
	}

	ok, err := govalidator.ValidateStruct(admin)

	g.Expect(ok).To(BeFalse())
	g.Expect(err).ToNot(BeNil())
	g.Expect(err.Error()).To(ContainSubstring("Firstname is required"))
}