package controllers
import (
	"net/http"
	"github.com/gin-gonic/gin"
	"backend/config"
	"backend/entity"
)

// GET /student_favs/my_favs/:id
// ฟังก์ชันนี้จะดึง "รายการข่าวทั้งหมด" ที่นักเรียนคนนี้กดถูกใจไว้
func GetStudentFavsByStudentID(c *gin.Context) {
    id := c.Param("id") // รับ id ของนักเรียน (StudentProfileID)
    var favorites []entity.StudentFavNews

    // ค้นหาในตาราง fav ว่ามี record ไหนบ้างที่เป็นของ student_id นี้
    if err := config.DB.Preload("NewsPost").Where("student_profile_id = ?", id).Find(&favorites).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Favorites not found"})
        return
    }

    c.JSON(http.StatusOK, favorites)
}

// POST /student_favs/toggle
func ToggleStudentFav(c *gin.Context) {
    var input struct {
        StudentProfileID uint `json:"student_profile_id"`
        NewsPostID       uint `json:"news_post_id"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var existingFav entity.StudentFavNews
    // เช็คว่ามีอยู่แล้วไหม?
    result := config.DB.Where("student_profile_id = ? AND news_post_id = ?", input.StudentProfileID, input.NewsPostID).First(&existingFav)

    if result.RowsAffected > 0 {
        // ถ้ามี -> ลบออก (Un-fav)
        config.DB.Delete(&existingFav) // หรือ Unscoped().Delete() ถ้าอยากลบถาวร
        c.JSON(http.StatusOK, gin.H{"message": "Unfavorited", "status": false})
    } else {
        // ถ้าไม่มี -> สร้างใหม่ (Fav)
        newFav := entity.StudentFavNews{
            StudentProfileID: input.StudentProfileID,
            NewsPostID:       input.NewsPostID,
        }
        config.DB.Create(&newFav)
        c.JSON(http.StatusCreated, gin.H{"message": "Favorited", "status": true})
    }
}