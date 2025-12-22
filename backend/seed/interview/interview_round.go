package interview

import (
	"backend/entity"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func SeedInterviewRounds(db *gorm.DB) error {
	// ดึง Scholarship ที่มีอยู่
	var scholarship entity.Scholarship
	if err := db.First(&scholarship).Error; err != nil {
		return fmt.Errorf("no scholarship found for interview round: %v", err)
	}

	// ดึง Admin Profile ที่มีอยู่
	var adminProfile entity.AdminProfile
	if err := db.First(&adminProfile).Error; err != nil {
		return fmt.Errorf("no admin profile found for interview round: %v", err)
	}

	// ดึง Location ที่มีอยู่
	var location entity.Location
	if err := db.First(&location).Error; err != nil {
		return fmt.Errorf("no location found for interview round: %v", err)
	}

	// ดึง InterviewMode ที่มีอยู่
	var interviewMode entity.InterviewMode
	if err := db.First(&interviewMode).Error; err != nil {
		return fmt.Errorf("no interview mode found for interview round: %v", err)
	}

	rounds := []entity.InterviewRound{
		{
			Name:            "รอบสัมภาษณ์ครั้งที่ 1",
			Description:     "รอบสัมภาษณ์สำหรับทุนการศึกษาประจำปี 2567",
			StartDateTime:   time.Now().AddDate(0, 0, 7),  // 7 วันจากวันนี้
			EndDateTime:     time.Now().AddDate(0, 0, 14), // 14 วันจากวันนี้
			SlotDuration:    30,                           // 30 นาที
			ScholarshipID:   scholarship.ID,
			AdminProfileID:  adminProfile.ID,
			LocationID:      &location.ID,
			InterviewModeID: &interviewMode.ID,
			MeetingLink:     "https://meet.google.com/abc-defg-hij",
		},
		{
			Name:            "รอบสัมภาษณ์ครั้งที่ 2",
			Description:     "รอบสัมภาษณ์รอบสองสำหรับผู้ผ่านการคัดเลือก",
			StartDateTime:   time.Now().AddDate(0, 0, 21), // 21 วันจากวันนี้
			EndDateTime:     time.Now().AddDate(0, 0, 28), // 28 วันจากวันนี้
			SlotDuration:    45,                           // 45 นาที
			ScholarshipID:   scholarship.ID,
			AdminProfileID:  adminProfile.ID,
			LocationID:      &location.ID,
			InterviewModeID: &interviewMode.ID,
			MeetingLink:     "https://meet.google.com/xyz-uvwx-rst",
		},
	}

	for _, round := range rounds {
		var existing entity.InterviewRound
		err := db.Where("name = ?", round.Name).First(&existing).Error

		if err == nil {
			continue // already exists
		}

		if err != gorm.ErrRecordNotFound {
			return err
		}

		if err := db.Create(&round).Error; err != nil {
			return fmt.Errorf("failed to seed interview round %s: %v", round.Name, err)
		}
	}

	return nil
}
