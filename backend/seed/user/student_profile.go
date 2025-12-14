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
// package user

// import (
// 	"backend/entity"
// 	"fmt"
// 	"gorm.io/gorm"
// 	"time"
// )

// func SeedStudentProfiles(db *gorm.DB) error {
// 	// ตรวจสอบว่ามีนักเรียนแล้วหรือยัง
// 	var count int64
// 	db.Model(&entity.StudentProfile{}).Count(&count)
// 	if count > 0 {
// 		return nil
// 	}

// 	// หาผู้ใช้ตัวอย่าง
// 	var user entity.User
// 	if err := db.Where("username = ?", "user").First(&user).Error; err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			return fmt.Errorf("in SeedStudentProfiles: prerequisite user 'user' not found")
// 		}
// 		return err
// 	}

// 	// หาวิชาเอก / สร้างถ้าไม่เจอ
// 	major := entity.Major{MajorName: "Computer Engineering"}
// 	if err := db.Where(entity.Major{MajorName: "Computer Engineering"}).FirstOrCreate(&major).Error; err != nil {
// 		return fmt.Errorf("in SeedStudentProfiles: failed to find or create major: %w", err)
// 	}

// 	// สร้างข้อมูลนักเรียนหลายคน
// 	students := []entity.StudentProfile{
// 		{
// 			StudentID:        "65010001",
// 			FirstNameTH:      "นักศึกษา",
// 			LastNameTH:       "ทดสอบ",
// 			FirstNameEN:      "Test",
// 			LastNameEN:       "Student",
// 			NationalID:       "1234567890123",
// 			BirthDate:        time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
// 			CurrentYear:      3,
// 			GPAX:             3.50,
// 			AdvisorName:      "อ. ทดสอบ",
// 			Phone:            "0812345678",
// 			Email:            "test.student@example.com",
// 			PermanentAddress: "123/456 ถ. ทดสอบ ต. ทดสอบ อ. ทดสอบ จ. ทดสอบ 12345",
// 			CurrentAddress:   "123/456 ถ. ทดสอบ ต. ทดสอบ อ. ทดสอบ จ. ทดสอบ 12345",
// 			Province:         "ทดสอบ",
// 			SiblingsCount:    5,
// 			UserID:           user.ID,
// 			MajorID:          major.ID,
// 		},
// 		{
// 			StudentID:        "65010002",
// 			FirstNameTH:      "สมชาย",
// 			LastNameTH:       "ใจดี",
// 			FirstNameEN:      "Somchai",
// 			LastNameEN:       "Jaidee",
// 			NationalID:       "9876543210987",
// 			BirthDate:        time.Date(2001, 5, 15, 0, 0, 0, 0, time.UTC),
// 			CurrentYear:      4,
// 			GPAX:             3.80,
// 			AdvisorName:      "อ. สมชาย",
// 			Phone:            "0898765432",
// 			Email:            "somchai.jaidee@example.com",
// 			PermanentAddress: "789/012 ถ. ตัวอย่าง ต. ตัวอย่าง อ. ตัวอย่าง จ. ตัวอย่าง 54321",
// 			CurrentAddress:   "789/012 ถ. ตัวอย่าง ต. ตัวอย่าง อ. ตัวอย่าง จ. ตัวอย่าง 54321",
// 			Province:         "ตัวอย่าง",
// 			SiblingsCount:    5,
// 			UserID:           user.ID,
// 			MajorID:          major.ID,
// 		},
// 	}

// 	// Insert ทีละคน ถ้าไม่ซ้ำ student_id
// 	for _, s := range students {
// 		var exist entity.StudentProfile
// 		if err := db.Where("student_id = ?", s.StudentID).First(&exist).Error; err == gorm.ErrRecordNotFound {
// 			if err := db.Create(&s).Error; err != nil {
// 				return err
// 			}
// 		}
// 	}

// 	return nil
// }
