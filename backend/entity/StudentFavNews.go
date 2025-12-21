package entity
import "gorm.io/gorm"

type StudentFavNews struct {
    gorm.Model
    
    NewsPostID uint `json:"news_post_id" gorm:"not null;uniqueIndex:idx_student_news" valid:"required~NewsPostID is required"`
    NewsPost   NewsPost `json:"news_post" gorm:"foreignKey:NewsPostID"`

    StudentProfileID uint `json:"student_profile_id" gorm:"not null;uniqueIndex:idx_student_news" valid:"required~StudentProfileID is required"`
    Student          StudentProfile `json:"student" gorm:"foreignKey:StudentProfileID"`
}

//เพิ่ม Unique Index เพื่อให้แต่ละนักเรียนสามารถเพิ่มข่าวสารโปรดได้เพียงครั้งเดียว