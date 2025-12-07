package scholarship

import (
	"backend/entity"
	"gorm.io/gorm"
	"time"
)

func SeedScholarships(db *gorm.DB) error {
	// Define scholarships to be seeded
	scholarships := []entity.Scholarship{
		{
			ScholarshipName:     "Beasiswa Anak Bangsa",
			Description:         "Beasiswa untuk anak bangsa yang berprestasi di bidang akademik.",
			OpenDate:            time.Now().Format("2006-01-02"),
			CloseDate:           time.Now().AddDate(0, 1, 0).Format("2006-01-02"),
			StatusscholarshipID: 1, // Assumes "Open" is ID 1
			TypescholarshipID:   1, // Assumes "Full" is ID 1
		},
		{
			ScholarshipName:     "Beasiswa Olahraga Nasional",
			Description:         "Dukungan untuk atlet muda berprestasi tingkat nasional.",
			OpenDate:            time.Now().Format("2006-01-02"),
			CloseDate:           time.Now().AddDate(0, 2, 0).Format("2006-01-02"),
			StatusscholarshipID: 1, // Assumes "Open" is ID 1
			TypescholarshipID:   2, // Assumes "Partial" is ID 2
		},
	}

	for _, s := range scholarships {
		// Use FirstOrCreate to prevent duplicates on subsequent runs
		if err := db.Where("scholarship_name = ?", s.ScholarshipName).FirstOrCreate(&s).Error; err != nil {
			return err
		}
	}

	return nil
}
