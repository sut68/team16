package evaluation

import (
	"backend/entity"
	"fmt"
	"math/rand"

	"gorm.io/gorm"
)

func SeedEvaluationScores(db *gorm.DB) error {
	// ตรวจสอบว่ามี EvaluationScore อยู่แล้วหรือไม่
	var existingCount int64
	if err := db.Model(&entity.EvaluationScore{}).Count(&existingCount).Error; err != nil {
		return err
	}
	if existingCount > 0 {
		return nil // Already seeded
	}

	// ดึง Evaluation ที่มีสถานะ in_progress, completed, approved, rejected (ไม่รวม pending)
	var evaluations []entity.Evaluation
	if err := db.Where("status != ?", entity.EvaluationStatusPending).
		Preload("InterviewRound").
		Find(&evaluations).Error; err != nil {
		return fmt.Errorf("failed to fetch evaluations: %v", err)
	}

	if len(evaluations) == 0 {
		return nil // No evaluations to score
	}

	// ดึง AdminProfile สำหรับเป็นผู้ให้คะแนน
	var admins []entity.AdminProfile
	if err := db.Find(&admins).Error; err != nil {
		return fmt.Errorf("failed to fetch admin profiles: %v", err)
	}

	if len(admins) == 0 {
		return nil // No admins
	}

	comments := []string{
		"ดีมาก",
		"พอใช้",
		"ต้องปรับปรุง",
		"โดดเด่น",
		"เป็นไปตามเกณฑ์",
		"ยอดเยี่ยม",
		"",
	}

	for _, evaluation := range evaluations {
		// ดึงเกณฑ์ที่ผูกกับ InterviewRound นี้
		var roundCriteria []entity.InterviewRoundCriterion
		if err := db.Where("interview_round_id = ? AND is_enabled = ?", evaluation.InterviewRoundID, true).
			Preload("EvaluationCriterion").
			Find(&roundCriteria).Error; err != nil {
			return fmt.Errorf("failed to fetch round criteria: %v", err)
		}

		// สร้างคะแนนสำหรับแต่ละเกณฑ์
		for _, rc := range roundCriteria {
			criterion := rc.EvaluationCriterion

			// คำนวณคะแนนตามประเภท
			var scoreValue float64
			switch criterion.ScoreType {
			case entity.ScoreTypeNumeric:
				// คะแนน 50-100% ของ MaxScore
				minScore := criterion.MaxScore * 0.5
				scoreRange := criterion.MaxScore - minScore
				scoreValue = minScore + (scoreRange * float64(rand.Intn(100)) / 100)
			case entity.ScoreTypePassFail:
				// 80% โอกาสผ่าน
				if rand.Intn(100) < 80 {
					scoreValue = 1
				} else {
					scoreValue = 0
				}
			case entity.ScoreTypeGrade:
				// คะแนน 60-100
				scoreValue = 60 + float64(rand.Intn(40))
			default:
				scoreValue = criterion.MaxScore * 0.7
			}

			// เลือก Admin (วนลูป)
			admin := admins[rand.Intn(len(admins))]

			evalScore := entity.EvaluationScore{
				ScoreValue:            scoreValue,
				Comment:               comments[rand.Intn(len(comments))],
				EvaluationID:          evaluation.ID,
				EvaluationCriterionID: criterion.ID,
				ScoringAdminID:        admin.ID,
			}

			if err := db.Create(&evalScore).Error; err != nil {
				return fmt.Errorf("failed to seed evaluation score: %v", err)
			}
		}
	}

	return nil
}
