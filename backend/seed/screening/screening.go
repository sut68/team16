package screening
import (
	"backend/entity"
	"gorm.io/gorm"

)

func SeedScreenings(db *gorm.DB) error {
	screenings := []entity.Screening{
		{
			AdminProfileID:   1,
			//StudentProfileID: 2,
			ApplicationID:    1,
			StatusScreeningID: 1,
			ScholarshipID:    1,
			RejectionReason:  nil,
		},
		{
			AdminProfileID:   1,
			ApplicationID: 2,
			StatusScreeningID: 1,
			ScholarshipID:    2,
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