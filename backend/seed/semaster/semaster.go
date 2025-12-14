package semaster

import (
	"backend/entity"
	"gorm.io/gorm"
	"time"
)

func CreateSemasters(db *gorm.DB) error {
	semasters := []entity.Semaster{
		{
			AcademicYear: "2567",
			Term:         "2",
			Round:        "1",
			StartDate:    time.Now().Format(time.RFC3339),
			EndDate:      time.Now().AddDate(0, 4, 0).Format(time.RFC3339),
			IsActive:     true,
		},
	}

	for _, semaster := range semasters {
		if err := db.Create(&semaster).Error; err != nil {
			return err
		}
	}
	return nil
}
