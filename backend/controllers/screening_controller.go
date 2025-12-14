package controllers

import (
	"backend/config"
	"backend/entity"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ---------------------- CREATE SCREENING ----------------------
func CreateScreening(c *gin.Context) {
	var screening entity.Screening
	if err := c.ShouldBindJSON(&screening); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Create(&screening).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": screening})
}

// ---------------------- GET ALL ----------------------
func GetAllScreenings(c *gin.Context) {
    var screenings []entity.Screening

    tx := config.DB.
        Preload("StatusScreening").
        Preload("Scholarship").
		Preload("Application.Semaster").
		Preload("Scholarship.Semaster").
        Preload("Application.StudentProfile")

    if err := tx.Find(&screenings).Error; err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, gin.H{"data": screenings})
}


// ---------------------- GET BY ID ----------------------
func GetScreeningByID(c *gin.Context) {
    var screening entity.Screening
    id := c.Param("id")

    tx := config.DB.
        Preload("AdminProfile").
        Preload("StatusScreening").
        Preload("Scholarship").
        Preload("Scholarship.Featurescholarships").
        Preload("Scholarship.Featurescholarships.Typefeature"). // ดึงชื่อเกณฑ์ (GPA, Income, etc.)
        Preload("Application").                               // (Optional) Preload แม่ข่ายให้ชัดเจน
        Preload("Application.StudentProfile").                // ดึง GPAX, Siblings
        Preload("Application.StudentProfile.FamilyInfo")      // ดึงรายได้พ่อแม่

    // เปลี่ยนจาก err != nil เฉยๆ เป็นการเช็ค RecordNotFound เพื่อส่ง 404
    if err := tx.First(&screening, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Data not found"}) // ใช้ 404 ถ้าหาไม่เจอ
        return
    }

    c.JSON(http.StatusOK, gin.H{"data": screening})
}


// ---------------------- UPDATE STATUS ----------------------
func UpdateScreeningStatus(c *gin.Context) {
	var screening entity.Screening
	id := c.Param("id")

	if err := config.DB.First(&screening, id).Error; err != nil {
		c.JSON(400, gin.H{"error": "Screening not found"})
		return
	}

	var input struct {
		StatusScreeningID uint    `json:"status_screening_id"`
		RejectionReason   *string `json:"rejection_reason"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	screening.StatusScreeningID = input.StatusScreeningID
	screening.RejectionReason = input.RejectionReason

	if err := config.DB.Model(&screening).Updates(map[string]interface{}{
		"status_screening_id": input.StatusScreeningID,
		"rejection_reason":    input.RejectionReason,
	}).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": screening})
}

// ---------------------- DELETE ----------------------
func DeleteScreening(c *gin.Context) {
	var screening entity.Screening
	id := c.Param("id")

	if err := config.DB.First(&screening, id).Error; err != nil {
		c.JSON(400, gin.H{"error": "Screening not found"})
		return
	}

	if err := config.DB.Delete(&screening).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": true})
}
