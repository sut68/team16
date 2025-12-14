package user
import (
	"backend/entity"
	"gorm.io/gorm"
	"fmt"
)

func SeedFamilyInfos(db *gorm.DB) error {
	if err := db.First(&entity.FamilyInfo{}).Error; err == gorm.ErrRecordNotFound {
		var studentProfile entity.StudentProfile
		if err := db.Where("student_id = ?", "65010001").First(&studentProfile).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("in SeedFamilyInfos: prerequisite student profile '65010001' not found")
			}
			return err
		}
		familyInfo := entity.FamilyInfo{
			FatherName:         "สมชาย ทดสอบ",
			FatherOccupation:   "วิศวกร",
			FatherIncome:       123456.0,
			
			MotherName:         "สมหญิง ทดสอบ",
			MotherOccupation:   "ครู",
			MotherIncome:       78901.0,

			GuardianName:       "สมปอง ทดสอบ",
			GuardianOccupation: "พ่อค้า",
			GuardianIncome:     45678.0,
			GuardianRelation:   "ลุง",
			
			ProfileID:          studentProfile.ID,
		}
		if err := db.Create(&familyInfo).Error; err != nil {
			return err
		}
	}
	return nil
}
