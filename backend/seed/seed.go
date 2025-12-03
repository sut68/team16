package seed

import (
	"gorm.io/gorm"

	"backend/seed/sponsor"
)

func SeedAll(db *gorm.DB) error {

	if err := sponsor.SeedIndustries(db); err != nil {
		return err
	}

	if err := sponsor.SeedSponsors(db); err != nil {
		return err
	}

	return nil
}