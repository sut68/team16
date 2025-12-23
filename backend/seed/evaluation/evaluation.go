package evaluation

import (
	"backend/entity"
	"fmt"
	"math/rand"

	"gorm.io/gorm"
)

func SeedEvaluations(db *gorm.DB) error {
	// ตรวจสอบว่ามี Evaluation อยู่แล้วหรือไม่
	var existingCount int64
	if err := db.Model(&entity.Evaluation{}).Count(&existingCount).Error; err != nil {
		return err
	}
	if existingCount > 0 {
		return nil // Already seeded
	}

	// ดึง ApplicationScholarship ที่มีสถานะ qualified หรือ pending
	var appScholarships []entity.ApplicationScholarship
	if err := db.Where("status IN ?", []string{"qualified", "pending"}).Find(&appScholarships).Error; err != nil {
		return fmt.Errorf("failed to fetch application scholarships: %v", err)
	}

	if len(appScholarships) == 0 {
		return nil // No applications to evaluate
	}

	// ดึง InterviewRound ทั้งหมด
	var interviewRounds []entity.InterviewRound
	if err := db.Find(&interviewRounds).Error; err != nil {
		return fmt.Errorf("failed to fetch interview rounds: %v", err)
	}

	if len(interviewRounds) == 0 {
		return nil // No interview rounds
	}

	// ดึง AdminProfile สำหรับเป็นผู้ประเมิน
	var admins []entity.AdminProfile
	if err := db.Find(&admins).Error; err != nil {
		return fmt.Errorf("failed to fetch admin profiles: %v", err)
	}

	if len(admins) == 0 {
		return nil // No admins
	}

	// สร้าง Evaluation สำหรับแต่ละ ApplicationScholarship
	statuses := []entity.EvaluationStatus{
		entity.EvaluationStatusPending,
		entity.EvaluationStatusInProgress,
		entity.EvaluationStatusCompleted,
		entity.EvaluationStatusApproved,
		entity.EvaluationStatusRejected,
	}

	remarks := []string{
		"ผู้สมัครมีคุณสมบัติเหมาะสม",
		"ต้องพิจารณาเพิ่มเติม",
		"ผลการสัมภาษณ์ดีมาก",
		"ผู้สมัครแสดงความมุ่งมั่นชัดเจน",
		"ข้อมูลครบถ้วน",
	}

	for i, appScholarship := range appScholarships {
		// เลือก InterviewRound (วนลูป)
		round := interviewRounds[i%len(interviewRounds)]
		// เลือก Admin (วนลูป)
		admin := admins[i%len(admins)]
		// เลือก Status (กระจายแบบสุ่ม)
		status := statuses[i%len(statuses)]

		// คำนวณคะแนนตาม status
		var totalScore float64
		switch status {
		case entity.EvaluationStatusCompleted, entity.EvaluationStatusApproved:
			totalScore = 70.0 + float64(rand.Intn(30)) // 70-99
		case entity.EvaluationStatusRejected:
			totalScore = 30.0 + float64(rand.Intn(30)) // 30-59
		case entity.EvaluationStatusInProgress:
			totalScore = float64(rand.Intn(50)) // 0-49 (ยังประเมินไม่เสร็จ)
		default:
			totalScore = 0 // pending
		}

		evaluation := entity.Evaluation{
			TotalScore:               totalScore,
			Status:                   status,
			Remark:                   remarks[i%len(remarks)],
			InterviewRoundID:         round.ID,
			ApplicationScholarshipID: appScholarship.ID,
			AdminID:                  admin.ID,
		}

		if err := db.Create(&evaluation).Error; err != nil {
			return fmt.Errorf("failed to seed evaluation: %v", err)
		}
	}

	return nil
}
