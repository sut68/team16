package scholarship

import (
	"backend/entity"
	"gorm.io/gorm"
	"time"
)

func SeedScholarships(db *gorm.DB) error {
	var scholarship entity.Scholarship
	// Check if "Beasiswa Anak Bangsa" scholarship already exists
	if err := db.Where("scholarship_name = ?", "Beasiswa Anak Bangsa").First(&scholarship).Error; err == gorm.ErrRecordNotFound {
		// If not found, create it
		newScholarship := entity.Scholarship{
			ScholarshipName:     "Beasiswa Anak Bangsa",
			Description:         "Beasiswa untuk anak bangsa yang berprestasi",
			OpenDate:            time.Now().String(),
			CloseDate:           time.Now().AddDate(0, 1, 0).String(),
			StatusscholarshipID: 1,
			TypescholarshipID:   1,
		}
		if err := db.Create(&newScholarship).Error; err != nil {
			return err
		}
	} else if err != nil {
		// Return any other error from the database query
		return err
	}
	return nil
}
