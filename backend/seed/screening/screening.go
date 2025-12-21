package screening
import (
	"backend/entity"
	"gorm.io/gorm"

)

func SeedScreenings(db *gorm.DB) error {
	screenings := []entity.Screening{
		{
			AdminProfileID:   1,
			StatusScreeningID: 1,
			ApplicationScholarshipID:    1,
			RejectionReason:  nil,
		},

	}
	for _, screening := range screenings {
		if err := db.Create(&screening).Error; err != nil {
			return err
		}
	}
	return nil
}