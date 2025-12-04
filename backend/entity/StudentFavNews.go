package entity
import "gorm.io/gorm"

type StudentFavNews struct{
	gorm.Model

	NewsPostID uint `json:"news_post_id" gorm:"not null"`
	NewsPost   NewsPost `json:"news_post" gorm:"foreignKey:NewsPostID"`

	StudentID uint `json:"student_id" gorm:"not null"`
	Student   StudentProfile `json:"student" gorm:"foreignKey:StudentID"`
	
}