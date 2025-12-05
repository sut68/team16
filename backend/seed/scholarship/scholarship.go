package scholarship

import (
	"backend/entity"
	"gorm.io/gorm"
	"time"
)

func SeedScholarships(db *gorm.DB) error {
	if err := db.First(&entity.Scholarship{}).Error; err == gorm.ErrRecordNotFound {
		scholarships := []entity.Scholarship{
			{
				ScholarshipName:     "Beasiswa Anak Bangsa",
				Description:         "Beasiswa untuk anak bangsa yang berprestasi",
				OpenDate:            time.Now().String(),
				CloseDate:           time.Now().AddDate(0, 1, 0).String(),
				StatusscholarshipID: 1,
				TypescholarshipID:   1,
			},
		}
		if err := db.Create(&scholarships).Error; err != nil {
			return err
		}
	}
	return nil
}
