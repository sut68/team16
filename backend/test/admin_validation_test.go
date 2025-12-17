package test

import (
	"testing"
	"backend/entity"
	"backend/validators"
	. "github.com/onsi/gomega"
)

func TestAdminProfileValidation(t *testing.T) {
	g := NewWithT(t)

	// Mock Data
	validAdmin := entity.AdminProfile{
		AdminFirstname: "Admin", AdminLastname: "System",
		Position: "Registrar",
		Email: "admin@sut.ac.th", Phone: "0442233445",
		UserID: 99,
	}
	t.Run("Check Email Format", func(t *testing.T) {
		// Case: อีเมลผิดรูปแบบ (ไม่มี @)
		a := validAdmin
		a.Email = "adminsut.ac.th"
		err := validators.ValidateStruct(&a)
		g.Expect(err).ToNot(BeNil())
		g.Expect(err.Error()).To(ContainSubstring("Invalid email format"))
	})

	t.Run("Check Phone Format", func(t *testing.T) {
		// Case: เบอร์โทรมีตัวหนังสือ
		a := validAdmin
		a.Phone = "081-CallMe"
		err := validators.ValidateStruct(&a)
		g.Expect(err).ToNot(BeNil())
		
		// Case: เบอร์สั้นไป
		a2 := validAdmin
		a2.Phone = "02"
		g.Expect(validators.ValidateStruct(&a2)).ToNot(BeNil())
	})

	// -----------------------------------------
	// กลุ่ม Test: ข้อมูลบังคับ (Required)
	// -----------------------------------------
	t.Run("Check Position Required", func(t *testing.T) {
		// Case: ตำแหน่งว่าง
		a := validAdmin
		a.Position = ""
		err := validators.ValidateStruct(&a)
		g.Expect(err).ToNot(BeNil())
		g.Expect(err.Error()).To(ContainSubstring("Position is required"))
	})
}