package news

import (
    "backend/entity"
    "gorm.io/gorm"
    "os"
    "path/filepath"
)

// ฟังก์ชันสร้างไฟล์หลอก
func createDummyFile(filePath string) error {
    dir := filepath.Dir(filePath)
    if _, err := os.Stat(dir); os.IsNotExist(err) {
        if err := os.MkdirAll(dir, 0755); err != nil {
            return err
        }
    }

    file, err := os.Create(filePath)
    if err != nil {
        return err
    }
    defer file.Close()
    return nil
}

func SeedNewsPosts(db *gorm.DB) error {

    // ===== 1) CHECK ว่ามี NEWS อยู่หรือยัง =====
    if err := db.First(&entity.NewsPost{}).Error; err != gorm.ErrRecordNotFound {
        return nil // มีแล้ว ไม่ seed ซ้ำ
    }

    // ===== 2) ดึง admin =====
    var admin entity.AdminProfile
    if err := db.First(&admin).Error; err != nil {
        return err
    }

    // ===== 3) ดึง StatusNews =====
    var statusPub entity.StatusNews
    if err := db.First(&statusPub).Error; err != nil {
        return err
    }

    // ===== 4) ดึง Scholarship =====
    var scholarship entity.Scholarship
    if err := db.Preload("Featurescholarships").First(&scholarship).Error; err != nil {
        return err
    }

    // ===== 5) เตรียม file path =====
    file1 := "uploads/news/scholarship_announcement_2024.pdf"
    file2 := "uploads/news/new_scholarship_opportunities.pdf"

    _ = createDummyFile(file1)
    _ = createDummyFile(file2)

    // ===== 6) Seed =====
    posts := []entity.NewsPost{
        {
            Title:         "Scholarship Announcement 2024",
            PostDetail:    "We are excited to announce...",
            FilePath:      file1,
            AdminID:       admin.ID,
            ScholarshipID: scholarship.ID,
            StatusNewsID:  statusPub.ID,
        },
        {
            Title:         "New Scholarship Opportunities",
            PostDetail:    "Discover new scholarship...",
            FilePath:      file2,
            AdminID:       admin.ID,
            ScholarshipID: scholarship.ID,
            StatusNewsID:  statusPub.ID,
        },
    }

    if err := db.Create(&posts).Error; err != nil {
        return err
    }

    return nil
}
