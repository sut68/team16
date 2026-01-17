package semaster

import (
	"backend/entity"
	"time"

	"gorm.io/gorm"
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
		var existing entity.Semaster
		err := db.Where("academic_year = ? AND term = ? AND round = ?", semaster.AcademicYear, semaster.Term, semaster.Round).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(&semaster).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}
	return nil
}
