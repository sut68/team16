package scholarship
import (
	"backend/entity"
	"gorm.io/gorm"
)

func SeedTypeFeatures(db *gorm.DB) error {
	typeFeatures := []entity.Typefeature{
		{Typefeaturename: "เกรดเฉลี่ย"},
		{Typefeaturename: "รายได้ในครอบครัว"},
		{Typefeaturename: "จำนวนพี่น้องที่กำลังศึกษา"},
		{Typefeaturename: "ระยะเวลาการศึกษา"},
		{Typefeaturename: "กิจกรรมพิเศษ"},
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
