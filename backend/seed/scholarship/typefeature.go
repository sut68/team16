package scholarship
import (
	"backend/entity"
	"gorm.io/gorm"
)

func SeedTypeFeatures(db *gorm.DB) error {
	typeFeatures := []entity.Typefeature{
		{Model: gorm.Model{ID: 1}, Typefeaturename: "เกรดเฉลี่ย"},
		{Model: gorm.Model{ID: 2}, Typefeaturename: "รายได้ในครอบครัว"},
		{Model: gorm.Model{ID: 3}, Typefeaturename: "จำนวนพี่น้องที่กำลังศึกษา"},
		{Model: gorm.Model{ID: 4}, Typefeaturename: "ระยะเวลาการศึกษา"},
		{Model: gorm.Model{ID: 5}, Typefeaturename: "กิจกรรมพิเศษ"},
	}
	for _, tf := range typeFeatures {
		var exist entity.Typefeature
		if err := db.Where("typefeaturename = ?", tf.Typefeaturename).First(&exist).Error; err == gorm.ErrRecordNotFound {
			if err := db.Create(&tf).Error; err != nil {
				return err
			}
		}
	}
	return nil
}