package screening
import (
	"backend/entity"
	"gorm.io/gorm"
)

func SeedStatusScreenings(db *gorm.DB) error {
	// สร้างข้อมูล StatusScreening เบื้องต้น
	statusScreenings := []entity.StatusScreening{
		{Model: gorm.Model{ID: 1}, StatusScreening: "รอตรวจสอบ"},   // ID 1
		{Model: gorm.Model{ID: 2}, StatusScreening: "ผ่านการคัดกรอง"},   // ID 1
		{Model: gorm.Model{ID: 3}, StatusScreening: "ไม่ผ่านการคัดกรอง"}, // ID 2
	}
	// วนลูปสร้างข้อมูล
	for _, status := range statusScreenings {
		var exist entity.StatusScreening
		if err := db.Where("status_screening = ?", status.StatusScreening).First(&exist).Error; err == gorm.ErrRecordNotFound {
			if err := db.Create(&status).Error; err != nil {
				return err
			}
		}
	}
	return nil
}