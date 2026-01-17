package user

import (
	"backend/entity"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func SeedStudentProfiles(db *gorm.DB) error {
	// Find prerequisite user 'user'
	var user entity.User
	if err := db.Where("username = ?", "user").First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("in SeedStudentProfiles: prerequisite user 'user' not found")
		}
		return err
	}

	// Find Computer Engineering major
	var major entity.Major
	if err := db.Where("major_name = ?", "Computer Engineering").First(&major).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("in SeedStudentProfiles: prerequisite major 'Computer Engineering' not found")
		}
		return err
	}

	students := []entity.StudentProfile{
		{
			StudentID:        "65010001",
			FirstNameTH:      "นักศึกษา",
			LastNameTH:       "ทดสอบ",
			FirstNameEN:      "Test",
			LastNameEN:       "Student",
			NationalID:       "1234567890123",
			BirthDate:        time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
			CurrentYear:      3,
			GPAX:             3.50,
			AdvisorName:      "อ.ทดสอบ",
			Phone:            "0812345678",
			Email:            "test.student@example.com",
			PermanentAddress: "123/456 ถ. ทดสอบ ต. ทดสอบ อ. ทดสอบ จ. ทดสอบ 12345",
			CurrentAddress:   "123/456 ถ. ทดสอบ ต. ทดสอบ อ. ทดสอบ จ. ทดสอบ 12345",
			Province:         "บุรีรัมย์",
			SiblingsCount:    1,
			UserID:           user.ID,
			MajorID:          major.ID,
		},
	}

	for _, s := range students {
		var existing entity.StudentProfile
		err := db.Where("student_id = ?", s.StudentID).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(&s).Error; err != nil {
				return fmt.Errorf("failed to seed student profile '%s': %v", s.StudentID, err)
			}
		} else if err != nil {
			return err
		}
	}

	return nil
}
