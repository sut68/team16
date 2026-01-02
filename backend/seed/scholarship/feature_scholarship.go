package scholarship

import (
	"backend/entity"

	"gorm.io/gorm"
)

func SeedFeatureScholarships(db *gorm.DB) error {
	features := []entity.Featurescholarship{
		{TypefeatureID: 1, Operator: ">=", Value: "3.5", ScholarshipID: 1},
		{TypefeatureID: 2, Operator: "<=", Value: "250000", ScholarshipID: 1},
		{TypefeatureID: 3, Operator: "<=", Value: "2", ScholarshipID: 1},
		{TypefeatureID: 1, Operator: ">=", Value: "3.0", ScholarshipID: 2},
		{TypefeatureID: 2, Operator: "<=", Value: "300000", ScholarshipID: 2},
		{TypefeatureID: 4, Operator: "<=", Value: "6", ScholarshipID: 2},
	}

	for _, fs := range features {
		var exist entity.Featurescholarship
		// เช็คว่า row นี้มีอยู่แล้วหรือยัง
		if err := db.Where("scholarship_id = ? AND typefeature_id = ? AND operator = ? AND value = ?",
			fs.ScholarshipID, fs.TypefeatureID, fs.Operator, fs.Value).First(&exist).Error; err == gorm.ErrRecordNotFound {

			// ใช้ FirstOrCreate จะปลอดภัยต่อ foreign key
			if err := db.FirstOrCreate(&fs, entity.Featurescholarship{
				ScholarshipID: fs.ScholarshipID,
				TypefeatureID: fs.TypefeatureID,
				Operator:      fs.Operator,
				Value:         fs.Value,
			}).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
