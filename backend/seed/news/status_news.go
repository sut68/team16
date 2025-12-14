package news
import (
	"gorm.io/gorm"
	"backend/entity"
)

func SeedStatusNews(db *gorm.DB) error {
	if err := db.First(&entity.StatusNews{}).Error; err == gorm.ErrRecordNotFound {
		statusNews := []entity.StatusNews{
			{StatusNews: "Published"},
			{StatusNews: "Draft"},
			{StatusNews: "Archived"},
			{StatusNews: "Deleted"},
			{StatusNews: "Members Only"},
		}
		if err := db.Create(&statusNews).Error; err != nil {
			return err
		}
	}
	return nil
}