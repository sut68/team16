package evaluation

import (
	"backend/entity"
	"fmt"

	"gorm.io/gorm"
)

func SeedEvaluationCriteria(db *gorm.DB) error {
	criteria := []entity.EvaluationCriterion{
		{
			Name:        "ความเหมาะสมของผู้สมัคร",
			Description: "ประเมินความเหมาะสมโดยรวมของผู้สมัครกับทุนการศึกษา",
			ScoreType:   entity.ScoreTypeNumeric,
			MaxScore:    100,
			Weight:      1.5,
			IsActive:    true,
		},
		{
			Name:        "ผลการเรียน",
			Description: "ประเมินจากเกรดเฉลี่ยสะสมและผลการเรียนในรายวิชาที่เกี่ยวข้อง",
			ScoreType:   entity.ScoreTypeNumeric,
			MaxScore:    100,
			Weight:      2.0,
			IsActive:    true,
		},
		{
			Name:        "ฐานะทางการเงิน",
			Description: "ประเมินความจำเป็นทางการเงินของผู้สมัคร",
			ScoreType:   entity.ScoreTypeNumeric,
			MaxScore:    100,
			Weight:      1.5,
			IsActive:    true,
		},
		{
			Name:        "ทักษะการสื่อสาร",
			Description: "ประเมินความสามารถในการสื่อสาร การนำเสนอ และการตอบคำถาม",
			ScoreType:   entity.ScoreTypeNumeric,
			MaxScore:    100,
			Weight:      1.0,
			IsActive:    true,
		},
		{
			Name:        "ความมุ่งมั่นและเป้าหมาย",
			Description: "ประเมินความตั้งใจ เป้าหมายในอนาคต และแรงจูงใจในการขอรับทุน",
			ScoreType:   entity.ScoreTypeNumeric,
			MaxScore:    100,
			Weight:      1.0,
			IsActive:    true,
		},
		{
			Name:        "กิจกรรมและจิตอาสา",
			Description: "ประเมินการมีส่วนร่วมในกิจกรรมนอกหลักสูตรและงานจิตอาสา",
			ScoreType:   entity.ScoreTypeNumeric,
			MaxScore:    100,
			Weight:      0.5,
			IsActive:    true,
		},
		{
			Name:        "ความประพฤติ",
			Description: "ประเมินความประพฤติและจริยธรรมของผู้สมัคร",
			ScoreType:   entity.ScoreTypePassFail,
			MaxScore:    1,
			Weight:      1.0,
			IsActive:    true,
		},
		{
			Name:        "เอกสารครบถ้วน",
			Description: "ตรวจสอบความครบถ้วนของเอกสารประกอบการสมัคร",
			ScoreType:   entity.ScoreTypePassFail,
			MaxScore:    1,
			Weight:      1.0,
			IsActive:    true,
		},
	}

	// INSERT
	for _, criterion := range criteria {
		var existing entity.EvaluationCriterion
		err := db.Where("name = ?", criterion.Name).First(&existing).Error

		if err == nil {
			continue
		}

		if err != gorm.ErrRecordNotFound {
			return err
		}

		if err := db.Create(&criterion).Error; err != nil {
			return fmt.Errorf("failed to seed evaluation criterion %s: %v", criterion.Name, err)
		}
	}

	return nil
}
