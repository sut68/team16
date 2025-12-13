package scholarship
import (
	"backend/entity"
	"gorm.io/gorm"
)

func SeedFeatureScholarships(db *gorm.DB) error {
	featureScholarships := []entity.Featurescholarship{
		{
			TypefeatureID: 1,
			Operator: ">=",
			Value:    "3.5",
			ScholarshipID: 1,
		},
		{
			TypefeatureID: 2,
			Operator: "<=",
			Value:    "300000",
			ScholarshipID: 1,
		},
		{
			TypefeatureID: 3,
			Operator: ">=",
			Value:    "2",
			ScholarshipID: 1,
		},
		{
			TypefeatureID: 1,
			Operator: ">=",
			Value:    "2.0",
			ScholarshipID: 2,
		},
		{
			TypefeatureID: 4,
			Operator: "<=",
			Value:    "6",
			ScholarshipID: 2,
		},
		{
			TypefeatureID: 5,
			Operator: "=",
			Value:    "20",
			ScholarshipID: 2,
		},
	}
	for _, fs := range featureScholarships {
		var exist entity.Featurescholarship
		if err := db.Where("typefeature_id = ? AND scholarship_id = ?", fs.TypefeatureID, fs.ScholarshipID).First(&exist).Error; err == gorm.ErrRecordNotFound {
			if err := db.Create(&fs).Error; err != nil {
				return err
			}
		}
	}
	return nil
}