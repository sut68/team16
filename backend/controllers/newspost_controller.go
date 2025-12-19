package controllers

import (
	"backend/config"
	"backend/entity"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm" // นำเข้า gorm เพื่อใช้ DB.Where และจัดการกับ error
)

// GET /newsposts
func GetAllNewsPosts(c *gin.Context) {
	fmt.Println("\n[GET] GetAllNewsPosts Called...")
	var items []entity.NewsPost
	if err := config.DB.
		Preload("Admin").
		Preload("Admin.User").
		Preload("Scholarship").
		Preload("StatusNews").
		Find(&items).Error; err != nil {
		fmt.Printf("[GET] Error fetching news: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("[GET] Found %d items\n", len(items))
	c.JSON(http.StatusOK, items)
}

// =================================================================
// 1. แก้ไข: GetNewsPostByID (ดึง Features และส่งโครงสร้าง JSON ใหม่)
// =================================================================
// GET /newsposts/:id
func GetNewsPostByID(c *gin.Context) {
	id := c.Param("id")
	fmt.Printf("\n[GET] GetNewsPostByID Called with ID: %s\n", id)

	var item entity.NewsPost
	// ดึง NewsPost พร้อม Preload ข้อมูลพื้นฐาน
	if err := config.DB.
		Preload("Admin").
		Preload("Admin.User").
		Preload("Scholarship").
		Preload("Scholarship.Typescholarship"). // เพิ่ม Typescholarship สำหรับ Frontend Preview
		Preload("StatusNews").
		First(&item, id).Error; err != nil {
		fmt.Printf("[GET] Error finding news ID %s: %v\n", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "News post not found"})
		return
	}

	// 2. ดึง Featurescholarships แยกต่างหาก (สำหรับ Frontend Preview)
	var features []entity.Featurescholarship
	if item.ScholarshipID != 0 {
		if err := config.DB.
			Where("scholarship_id = ?", item.ScholarshipID).
			Preload("Typefeature"). // Preload Typefeature ตามที่ Frontend ต้องการ
			Find(&features).Error; err != nil {
			fmt.Printf("[GET] Warning: Error finding features for Scholarship ID %d: %v\n", item.ScholarshipID, err)
			// ดำเนินการต่อได้แม้จะมี Error
		}
	}

	// 3. ส่ง JSON ในโครงสร้างที่ Frontend คาดหวัง
	response := gin.H{
		"news_post": item,
		"features":  features,
	}

	fmt.Printf("[GET] Found item: %s, Features Count: %d\n", item.Title, len(features))
	c.JSON(http.StatusOK, response)
}


//POST /newsposts
func CreateNewsPost(c *gin.Context) {
	fmt.Println("\n[POST] CreateNewsPost Called...")

	title := c.PostForm("title")
	postDetail := c.PostForm("post_detail")
	
	// ใช้ ParseUint อย่างระมัดระวัง (AdminID ต้องมี, ScholarshipID/StatusNewsID อาจไม่มี)
	adminID, _ := strconv.ParseUint(c.PostForm("admin_id"), 10, 64)
	scholarshipID, _ := strconv.ParseUint(c.PostForm("scholarship_id"), 10, 64)
	statusNewsID, _ := strconv.ParseUint(c.PostForm("status_news_id"), 10, 64)

	var filePath string = ""

	// 2. พยายามดึงไฟล์
	file, err := c.FormFile("file_path")

	if err == nil {
		fmt.Printf("[POST] File received: %s (Size: %d bytes)\n", file.Filename, file.Size)

		extension := filepath.Ext(file.Filename)
		fileName := fmt.Sprintf("news-%d%s", time.Now().UnixNano(), extension)
		filePath = "uploads/news/" + fileName

		if _, err := os.Stat("uploads/news"); os.IsNotExist(err) {
			os.MkdirAll("uploads/news", 0755)
		}

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			fmt.Printf("[POST] Error saving file to disk: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image file"})
			return
		}
		fmt.Printf("[POST] File saved at: %s\n", filePath)

	} else {
		fmt.Println("[POST] No file uploaded (or error getting file). Proceeding without image.")
	}

	// 4. บันทึก DB
	item := entity.NewsPost{
		Title:           title,
		PostDetail:      postDetail,
		FilePath:        filePath,
		AdminID:         uint(adminID),
		ScholarshipID:   uint(scholarshipID),
		StatusNewsID:    uint(statusNewsID),
	}
    
    // ตรวจสอบว่ามี AdminID หรือไม่ (ถ้าเป็น 0 อาจไม่ผ่าน Validation/FK)
	if item.AdminID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "AdminID is required"})
		return
	}

	if err := config.DB.Create(&item).Error; err != nil {
		fmt.Printf("[POST] Database Create Error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("[POST] Success! Created News ID: %d\n", item.ID)
	// ส่ง item กลับไปให้ Frontend เพื่อให้มี ID และ CreatedAt
	c.JSON(http.StatusCreated, item)
}

// =================================================================
// 2. แก้ไข: UpdateNewsPost (จัดการกับไฟล์ใหม่/เก่า และ ID ที่เป็น 0)
// =================================================================
// PUT /newsposts/:id
func UpdateNewsPost(c *gin.Context) {
	fmt.Println("\n[PUT] UpdateNewsPost Called...")
	
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// 2. Parse Form Data
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form data"})
		return
	}
    
    // 3. ดึงข้อมูลเก่าเพื่อใช้ FilePath เดิม
    var oldData entity.NewsPost
    if err := config.DB.First(&oldData, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "News post not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        }
        return
    }
    
    // 4. เตรียมข้อมูลสำหรับ Update (ใช้ Map เพื่อระบุ field ชัดเจน)
	updateData := map[string]interface{}{}
    
    // ดึงค่า String (ถ้ามี)
    if title := c.PostForm("title"); title != "" {
        updateData["title"] = title
    }
    if postDetail := c.PostForm("post_detail"); postDetail != "" {
        updateData["post_detail"] = postDetail
    }
    
    // ดึงค่า ID และจัดการ 0
    if adminIDStr := c.PostForm("admin_id"); adminIDStr != "" {
        if adminID, err := strconv.ParseUint(adminIDStr, 10, 64); err == nil {
            updateData["admin_id"] = uint(adminID)
        }
    }
    
    // 🔥 แก้ไข: หาก scholarship_id เป็น 0 หรือไม่ได้ส่งมา (ว่าง) ไม่ต้องอัปเดต
    if scholarshipIDStr := c.PostForm("scholarship_id"); scholarshipIDStr != "" {
        if scholarshipID, err := strconv.ParseUint(scholarshipIDStr, 10, 64); err == nil && scholarshipID != 0 {
            updateData["scholarship_id"] = uint(scholarshipID)
        }
    } else {
        // หาก Frontend ต้องการอนุญาตให้ตั้งค่าเป็น NULL ได้ ต้องใช้ gorm.DeletedAt/NULLABLE
        // แต่ถ้าไม่อยากให้มัน override ค่าเดิมถ้าไม่ได้ส่งมา ให้ข้าม field นี้ไปเลย
    }
    
    if statusNewsIDStr := c.PostForm("status_news_id"); statusNewsIDStr != "" {
        if statusNewsID, err := strconv.ParseUint(statusNewsIDStr, 10, 64); err == nil {
            updateData["status_news_id"] = uint(statusNewsID)
        }
    }
    
    // 5. จัดการไฟล์รูปภาพ
	file, err := c.FormFile("file_path") 

	if err == nil {
		// ✅ กรณีมีไฟล์ใหม่: บันทึกไฟล์ และอัปเดต FilePath
		extension := filepath.Ext(file.Filename)
		fileName := fmt.Sprintf("news-%d%s", time.Now().UnixNano(), extension)
		filePath := "uploads/news/" + fileName

		if _, err := os.Stat("uploads/news"); os.IsNotExist(err) {
			os.MkdirAll("uploads/news", 0755)
		}

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}
        updateData["file_path"] = filePath // อัปเดต Path ใหม่
        
        // (Optional) ลบไฟล์เก่าทิ้ง หากไม่ต้องการเก็บไฟล์ที่ถูกแทนที่
        if oldData.FilePath != "" {
            // os.Remove(oldData.FilePath)
        }

	} else {
		// ⚠️ กรณีไม่มีไฟล์ใหม่: ใช้ path เดิมจาก oldData
        // (ไม่ต้องทำอะไร เพราะ updateData จะไม่มี field file_path หากไม่มีการส่งไฟล์ใหม่)
        // หรือถ้า Frontend ส่ง "clear" มาเพื่อลบรูป ให้ใส่ updateData["file_path"] = ""
	}
    
    // 6. บันทึกลงฐานข้อมูล
	result := config.DB.Model(&entity.NewsPost{}).Where("id = ?", uint(id)).Updates(updateData)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
    
	fmt.Printf("[PUT] Success! Updated News ID: %d, Rows Affected: %d\n", id, result.RowsAffected)
	c.JSON(http.StatusOK, gin.H{
		"message": "Update successful",
		"data":    updateData,
	})
}


// DELETE /newsposts/:id

func DeleteNewsPost(c *gin.Context) {
	id := c.Param("id")
	fmt.Printf("\n[DELETE] DeleteNewsPost Called for ID: %s\n", id)

	if err := config.DB.Delete(&entity.NewsPost{}, id).Error; err != nil {
		fmt.Printf("[DELETE] Error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("[DELETE] Success!")
	c.JSON(http.StatusOK, gin.H{"message": "NewsPost deleted successfully"})
}