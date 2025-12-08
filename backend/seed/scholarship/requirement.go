package scholarship

import (
	"backend/entity"

	"gorm.io/gorm"
)

func SeedRequirements(db *gorm.DB) error {
	requirements := []entity.Requirement{
		{Name: "Transkrip nilai semester terakhir"},
		{Name: "Surat rekomendasi dari dekan fakultas"},
		{Name: "Esai motivasi"},
		{Name: "Sertifikat kegiatan ekstrakurikuler"},
	}

	for _, req := range requirements {
		// Use FirstOrCreate to avoid duplicate entries on subsequent runs
		if err := db.FirstOrCreate(&req, entity.Requirement{Name: req.Name}).Error; err != nil {
			return err
		}
	}

	return nil
}
