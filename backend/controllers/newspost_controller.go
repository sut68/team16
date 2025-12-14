package controllers

import (
	"backend/config"
	"backend/entity"
	"fmt" // ใช้สำหรับ print log
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GET /newsposts
func GetAllNewsPosts(c *gin.Context) {
	fmt.Println("\n[GET] GetAllNewsPosts Called...") // Log เริ่มฟังก์ชัน

	var items []entity.NewsPost
	if err := config.DB.
		Preload("Admin").
		Preload("Admin.User").
		Preload("Scholarship").
		Preload("StatusNews").
		Find(&items).Error; err != nil {
		fmt.Printf("[GET] Error fetching news: %v\n", err) // Log Error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("[GET] Found %d items\n", len(items)) // Log จำนวนข้อมูลที่เจอ
	c.JSON(http.StatusOK, items)
}

// GET /newsposts/:id
func GetNewsPostByID(c *gin.Context) {
	id := c.Param("id")
	fmt.Printf("\n[GET] GetNewsPostByID Called with ID: %s\n", id)

	var item entity.NewsPost
	if err := config.DB.
		Preload("Admin").
		Preload("Admin.User").
		Preload("Scholarship").
		Preload("Scholarship.Featurescholarships").
		Preload("Scholarship.Featurescholarships.Typefeature").
		Preload("Scholarship.Typescholarship").
		Preload("StatusNews").
		First(&item, id).Error; err != nil {
		fmt.Printf("[GET] Error finding news ID %s: %v\n", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("[GET] Found item: %s\n", item.Title)
	c.JSON(http.StatusOK, item)
}

// POST /newsposts
func CreateNewsPost(c *gin.Context) {
	fmt.Println("\n[POST] CreateNewsPost Called...")

	// 1. Log ค่า Text ที่รับมา (ส่วนนี้เหมือนเดิม)
	title := c.PostForm("title")
	postDetail := c.PostForm("post_detail")
	fmt.Printf("[POST] Received Title: %s\n", title)
	fmt.Printf("[POST] Received Detail (Length): %d chars\n", len(postDetail))

	adminID, _ := strconv.ParseUint(c.PostForm("admin_id"), 10, 64)
	scholarshipID, _ := strconv.ParseUint(c.PostForm("scholarship_id"), 10, 64)
	statusNewsID, _ := strconv.ParseUint(c.PostForm("status_news_id"), 10, 64)
	visibility := c.PostForm("visibility")
	fmt.Printf("[POST] IDs -> Admin: %d, Scholar: %d, Status: %d, Visibility: %s\n", adminID, scholarshipID, statusNewsID, visibility)

	var filePath string = "" // สร้างตัวแปรมารอรับ Path (ค่าเริ่มต้นเป็นว่าง)

	// 2. พยายามดึงไฟล์
	file, err := c.FormFile("file_path")

	if err == nil {
		// ✅ กรณีมีไฟล์ส่งมา: ทำการบันทึกไฟล์
		fmt.Printf("[POST] File received: %s (Size: %d bytes)\n", file.Filename, file.Size)

		// 3. เริ่มกระบวนการ Save
		extension := filepath.Ext(file.Filename)
		fileName := fmt.Sprintf("news-%d%s", time.Now().UnixNano(), extension)
		filePath = "uploads/news/" + fileName // อัปเดตตัวแปร filePath

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
		// ⚠️ กรณีไม่มีไฟล์ส่งมา หรือหาไม่เจอ: ให้ข้ามไป ไม่ต้อง return error
		fmt.Println("[POST] No file uploaded (or error getting file). Proceeding without image.")
		// filePath จะยังคงเป็น "" เหมือนเดิม
	}

	// ---------------------------------------------------------

	// 4. บันทึก DB (เหมือนเดิม)
	item := entity.NewsPost{
		Title:         title,
		PostDetail:    postDetail,
		FilePath:      filePath, // ค่านี้จะเป็น Path รูป หรือ "" ก็ได้
		AdminID:       uint(adminID),
		ScholarshipID: uint(scholarshipID),
		StatusNewsID:  uint(statusNewsID),
	}

	if err := config.DB.Create(&item).Error; err != nil {
		fmt.Printf("[POST] Database Create Error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("[POST] Success! Created News ID: %d\n", item.ID)
	c.JSON(http.StatusCreated, item)
}

// PUT /newsposts/:id
func UpdateNewsPost(c *gin.Context) {
	// 1. รับ ID และแปลงเป็นตัวเลข (แก้ปัญหา rows:0)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// 2. สั่ง Parse Form Data แบบระบุเจาะจง (แก้ปัญหาค่าว่าง)
	// 32MB คือ max memory
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form data"})
		return
	}

	// 3. ดึงค่าจาก Form (ตอนนี้จะได้ค่าชัวร์ๆ แล้ว)
	title := c.PostForm("title")
	postDetail := c.PostForm("post_detail")
	
	// แปลง String เป็นตัวเลขสำหรับ ID ต่างๆ
	adminID, _ := strconv.ParseUint(c.PostForm("admin_id"), 10, 64)
	scholarshipID, _ := strconv.ParseUint(c.PostForm("scholarship_id"), 10, 64)
	statusNewsID, _ := strconv.ParseUint(c.PostForm("status_news_id"), 10, 64)

	// 4. จัดการไฟล์รูปภาพ
	var filePath string
	file, err := c.FormFile("file_path") // ลองดึงไฟล์

	if err == nil {
		// ✅ กรณีมีไฟล์ใหม่: บันทึกไฟล์
		extension := filepath.Ext(file.Filename)
		// ตั้งชื่อไฟล์ใหม่ด้วย Timestamp เพื่อกันชื่อซ้ำ
		fileName := fmt.Sprintf("news-%d%s", time.Now().UnixNano(), extension)
		filePath = "uploads/news/" + fileName

		// สร้าง Folder ถ้ายังไม่มี
		if _, err := os.Stat("uploads/news"); os.IsNotExist(err) {
			os.MkdirAll("uploads/news", 0755)
		}

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}
	} else {
		// ⚠️ กรณีไม่มีไฟล์ใหม่: ต้องใช้ path เดิมจาก DB
		var oldData entity.NewsPost
		if err := config.DB.First(&oldData, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "News post not found"})
			return
		}
		filePath = oldData.FilePath
	}

	// 5. เตรียมข้อมูลสำหรับ Update (ใช้ Map เพื่อระบุ field ชัดเจน)
	updateData := map[string]interface{}{
		"title":          title,
		"post_detail":    postDetail,
		"admin_id":       uint(adminID),       // แปลงเป็น uint ให้ตรง model
		"scholarship_id": uint(scholarshipID), // แปลงเป็น uint ให้ตรง model
		"status_news_id": uint(statusNewsID),  // แปลงเป็น uint ให้ตรง model
		"file_path":      filePath,
	}

	// 6. บันทึกลงฐานข้อมูล
	// ใช้ id ที่เป็น uint ใน Where clause เพื่อความชัวร์
	result := config.DB.Model(&entity.NewsPost{}).Where("id = ?", uint(id)).Updates(updateData)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// เช็คว่ามีการแก้ไขจริงไหม
	if result.RowsAffected == 0 {
		// อาจจะเกิดขึ้นถ้าข้อมูลใหม่เหมือนข้อมูลเดิมเป๊ะๆ หรือหา ID ไม่เจอ (แต่เราเช็ค First ไปแล้ว)
		fmt.Println("[WARN] Update executed but 0 rows affected")
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Update successful",
		"data": updateData, // ส่งข้อมูลกลับไปดูว่า Backend เห็นเป็นอะไร
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
