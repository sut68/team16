package validators

import "github.com/asaskevich/govalidator"

// ใช้ตรวจสอบ struct ตาม tag valid
func ValidateStruct(v interface{}) error {
	_, err := govalidator.ValidateStruct(v)
	return err
}