package evaluation

import (
	"backend/entity"
	"fmt"

	"gorm.io/gorm"
)

func SeedInterviewRoundCriteria(db *gorm.DB) error {
	// ดึง InterviewRound ทั้งหมด
	var rounds []entity.InterviewRound
	if err := db.Find(&rounds).Error; err != nil {
		return fmt.Errorf("failed to fetch interview rounds: %v", err)
	}

	if len(rounds) == 0 {
		return nil // No interview rounds to assign criteria
	}

	// ดึง EvaluationCriterion ทั้งหมด
	var criteria []entity.EvaluationCriterion
	if err := db.Where("is_active = ?", true).Find(&criteria).Error; err != nil {
		return fmt.Errorf("failed to fetch evaluation criteria: %v", err)
	}

	if len(criteria) == 0 {
		return nil // No criteria to assign
	}

	// สำหรับแต่ละ InterviewRound ให้เพิ่มเกณฑ์ทั้งหมด
	for _, round := range rounds {
		for _, criterion := range criteria {
			roundCriterion := entity.InterviewRoundCriterion{
				InterviewRoundID:      round.ID,
				EvaluationCriterionID: criterion.ID,
				Weight:                criterion.Weight, // ใช้ weight จากเกณฑ์หลัก
				IsEnabled:             true,
			}

			// ตรวจสอบว่ามีอยู่แล้วหรือไม่
			var existing entity.InterviewRoundCriterion
			err := db.Where("interview_round_id = ? AND evaluation_criterion_id = ?",
				round.ID, criterion.ID).First(&existing).Error

			if err == nil {
				continue // already exists
			}

			if err != gorm.ErrRecordNotFound {
				return err
			}

			if err := db.Create(&roundCriterion).Error; err != nil {
				return fmt.Errorf("failed to seed interview round criterion: %v", err)
			}
		}
	}

	return nil
}
