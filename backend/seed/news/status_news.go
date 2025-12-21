package news

import (
    "backend/entity"
    "gorm.io/gorm"
)

func SeedStatusNews(db *gorm.DB) error {
    statusNewsList := []entity.StatusNews{
        {StatusNews: "Published"},
        {StatusNews: "Draft"},
        {StatusNews: "Archived"},
		{StatusNews: "Members Only"},
        {StatusNews: "Deleted"},
        
    }

    for _, status := range statusNewsList {
        var sn entity.StatusNews
        // ถ้า id ไม่มี จะสร้างใหม่
        if err := db.FirstOrCreate(&sn, entity.StatusNews{StatusNews: status.StatusNews}).Error; err != nil {
            return err
        }
        // อัปเดตชื่อให้ตรง
        if err := db.Model(&entity.StatusNews{}).Where("id = ?", sn.ID).Update("status_news", status.StatusNews).Error; err != nil {
            return err
        }
    }

    return nil
}
