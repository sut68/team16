package seed

import (
	"backend/seed/scholarship"
	"gorm.io/gorm"
	"backend/seed/approval"
	"backend/seed/semaster"
	"backend/seed/sponsor"
	"backend/seed/user"
	"backend/seed/screening"
	"backend/seed/news"
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
		return err
	}

	if err := user.SeedMajors(db); err != nil {
		return err
	}

	if err := user.SeedStudentProfiles(db); err != nil {
		return err
	}

	if err := user.SeedAdminProfiles(db); err != nil {
		return err
	}

	 if err := user.SeedFamilyInfos(db); err != nil {
		 return err
	 }

	if err := semaster.CreateSemasters(db); err != nil {
		return err
	}

	if err := scholarship.SeedStatusScholarships(db); err != nil {
		return err
	}

	if err := scholarship.SeedTypeScholarships(db); err != nil {
		return err
	}

	if err := scholarship.SeedTypeFeatures(db); err != nil {
		return err
	}

	if err := scholarship.SeedScholarships(db); err != nil {
		return err
	}

	if err := scholarship.SeedFeatureScholarships(db); err != nil {
		return err
	}

	if err := scholarship.SeedRequirements(db); err != nil {
		return err
	}

	if err := approval.SeedApplications(db); err != nil {
		return err
	}
	// if err := approval.SeedApplicationScholarships(db); err != nil {
	// 	return err
	// }

	// if err := approval.SeedApplicationDocuments(db); err != nil {
	// 	return err
	// }

	if err := approval.SeedApprovalRequirements(db); err != nil {
		return err
	}


	if err := screening.SeedStatusScreenings(db); err != nil {
		return err
	}

	if err := screening.SeedScreenings(db); err != nil {
		return err
	}

	if err := news.SeedStatusNews(db); err != nil {
		return err
	}

	if err := news.SeedNewsPosts(db); err != nil {
		return err
	}

	return nil
}