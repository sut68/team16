package user

import (
	"backend/entity"
	"fmt"
	"gorm.io/gorm"
	"time"
)

func SeedStudentProfiles(db *gorm.DB) error {
	if err := db.First(&entity.StudentProfile{}).Error; err == gorm.ErrRecordNotFound {
		var user entity.User
		if err := db.Where("username = ?", "user").First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("in SeedStudentProfiles: prerequisite user 'user' not found")
			}
			return err
		}

		// Use FirstOrCreate to make this seeder more robust.
		// It will find the major if it exists, or create it if it doesn't.
		major := entity.Major{MajorName: "Computer Engineering"}
		if err := db.Where(entity.Major{MajorName: "Computer Engineering"}).FirstOrCreate(&major).Error; err != nil {
			return fmt.Errorf("in SeedStudentProfiles: failed to find or create major: %w", err)
		}

		studentProfile := entity.StudentProfile{
			StudentID:        "65010001",
			FirstNameTH:      "นักศึกษา",
			LastNameTH:       "ทดสอบ",
			FirstNameEN:      "Test",
			LastNameEN:       "Student",
			NationalID:       "1234567890123",
			BirthDate:        time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
			CurrentYear:      3,
			GPAX:             3.50,
			AdvisorName:      "อ. ทดสอบ",
			Phone:            "0812345678",
			Email:            "test.student@example.com",
			PermanentAddress: "123/456 ถ. ทดสอบ ต. ทดสอบ อ. ทดสอบ จ. ทดสอบ 12345",
			CurrentAddress:   "123/456 ถ. ทดสอบ ต. ทดสอบ อ. ทดสอบ จ. ทดสอบ 12345",
			Province:         "ทดสอบ",
			SiblingsCount:    1,
			UserID:           user.ID,
			MajorID:          major.ID,
		}
		if err := db.Create(&studentProfile).Error; err != nil {
			return err
		}
	}
	return nil
}
