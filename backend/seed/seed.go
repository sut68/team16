package seed

import (
	"gorm.io/gorm"

	"backend/seed/sponsor"
	"backend/seed/user"
)

func SeedAll(db *gorm.DB) error {

	if err := sponsor.SeedIndustries(db); err != nil {
		return err
	}

	if err := sponsor.SeedSponsors(db); err != nil {
		return err
	}

	if err := user.SeedRoles(db); err != nil {
		return err
	}

	if err := user.SeedUsers(db); err != nil {
		
	}

	return nil
}