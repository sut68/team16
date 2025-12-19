package news
import (
    "backend/entity"
    "gorm.io/gorm"
)

func SeedStudentFavs(db *gorm.DB) error{
	studentFavs := []entity.StudentFavNews{
		{StudentProfileID: 1, NewsPostID: 1},
		{StudentProfileID: 1, NewsPostID: 2},
	}
	for _, sf := range studentFavs {
		db.FirstOrCreate(&sf, entity.StudentFavNews{StudentProfileID: sf.StudentProfileID, NewsPostID: sf.NewsPostID})
	}
	return nil
}
