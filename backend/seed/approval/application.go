package approval
import (
	"backend/entity"
	"gorm.io/gorm"
)

func SeedApplications(db *gorm.DB) error {
	applications := []entity.Application{
		{
			StudentProfileID: 1,
			SemasterID: 1,
			ApplicationScholarships: []entity.ApplicationScholarship{},
		},
	}
	for _, application := range applications {
		if err := db.Create(&application).Error; err != nil {
			return err
		}
	}
	return nil
}