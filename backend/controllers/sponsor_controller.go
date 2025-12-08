package controllers

import (
	"backend/config"
	"backend/entity"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

// GET /sponsors
func GetSponsors(ctx *gin.Context) {
	var sponsors []entity.Sponsor

	// SELECT ALL
	if err := config.DB.Preload("Industry").Preload("Contacts").Find(&sponsors).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, sponsors)
}

// GET /sponsors/:id 
func GetSponsorsByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid sponsor id"})
		return
	}

	var sponsor entity.Sponsor

	// SELECT FIRST ROW
	if err := config.DB.Preload("Industry").Preload("Contacts").First(&sponsor, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, sponsor)
}

// POST /sponsors
func CreateSponsor(ctx *gin.Context) {
	var inputValues struct {
		CompanyName		string		`json:"company_name" binding:"required"`
		IndustryID		*uint     `json:"industry_id"`
		Website				*string		`json:"website"`
		Status				string		`json:"status"`
		Description		*string		`json:"description"`
		Contacts			[]entity.SponsorContact		`json:"contacts"`
	}

	if err := ctx.ShouldBindJSON(&inputValues); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sponsor := entity.Sponsor{
		CompanyName:		inputValues.CompanyName,
		IndustryID:			inputValues.IndustryID,
		Website:				inputValues.Website,
		Status:					inputValues.Status,
		Description:		inputValues.Description,
		Contacts:			inputValues.Contacts,
	}

	// INSERT
	if err := config.DB.Create(&sponsor).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, sponsor)
}

// PATCH /sponsors/:id
func UpdateSponsor(ctx *gin.Context) {
	idParam := ctx.Param("id")
	idUint, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	id := uint(idUint)

	var sponsor entity.Sponsor
	// หา sponsor ว่ามีมั้ย กันส่ง มาไม่มีอยู่
	if err := config.DB.Preload("Contacts").First(&sponsor, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "sponsor not found"})
		return
	}

	var inputValues struct {
		CompanyName		*string		`json:"company_name"`
		IndustryID    *string		`json:"industry_id"`
		Website				*string		`json:"website"`
		Status				*string		`json:"status"`
		Description		*string		`json:"description"`
	}

	if err := ctx.ShouldBindJSON(&inputValues); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// เริ่มต้นการดำเนินการฐานข้อมูล
	tx := config.DB.Begin()
	if tx.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": tx.Error.Error()})
		return
	}
	// defer เพื่อจัดการ Rollback หากเกิด Panic
	defer func ()  {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	updates := map[string]interface{}{}

	if inputValues.CompanyName != nil {
		updates["company_name"] = *inputValues.CompanyName
	}
	if inputValues.IndustryID != nil {
		updates["industry_id"] = *inputValues.IndustryID
	}
	if inputValues.Website != nil {
		updates["website"] = *inputValues.Website
	}
	if inputValues.Status != nil {
		updates["status"] = *inputValues.Status
	}
	if inputValues.Description != nil {
		updates["description"] = *inputValues.Description
	}

	if len(updates) > 0 {
		// UPDATE
		if err := config.DB.Model(&sponsor).Updates(updates).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if err := config.DB.Preload("Contacts").First(&sponsor, sponsor.ID).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "reload failed"})
		return
	}

	ctx.JSON(http.StatusOK, sponsor)
}

// PATCH /sponsors/:id/contacts
type BatchContactsPayload struct {
	Upsert			[]entity.SponsorContact	 `json:"upsert"`
	DeleteIDs		[]uint									 `json:"delete_ids"`
}
// จำกัดผู้ติดต่อ 10 คน / บริษัท
const MaxContactsBatch = 10

func UpdateSponsorContacts(ctx *gin.Context) {
	idParam := ctx.Param("id")
	idUint, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid sponsor id"})
		return
	}
	sponsorID := uint(idUint)

	var payload BatchContactsPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// เช็ค limit ผู้ติดต่อ
	if len(payload.Upsert) > MaxContactsBatch || len(payload.DeleteIDs) > MaxContactsBatch {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "too many items in batch"})
		return
	}

	// เริ่มต้นการดำเนินการฐานข้อมูล
	tx := config.DB.Begin()
	if tx.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": tx.Error.Error()})
		return
	}
	// defer เพื่อจัดการ Rollback หากเกิด Panic
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	var sponsor entity.Sponsor
	if err := tx.First(&sponsor, sponsorID).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusNotFound, gin.H{"error": "sponsor not found"})
		return
	}

	var existingCount int64
	if err := tx.Model(&entity.SponsorContact{}).
		Where("sponsor_id = ?", sponsorID).
		Count(&existingCount).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newCreates := 0
	for i := range payload.Upsert {
		// บังคับใช้ sponsorID จาก URL (override ค่า client)
		payload.Upsert[i].SponsorID = sponsorID

		// ถ้าเป็น create (ไม่มี ID) ให้บวก newCreates
		if payload.Upsert[i].ID == 0 {
			newCreates++
		}
	}

	// คำนวณผลลัพธ์หลัง apply: existing - deletes + newCreates
	deleteCount := len(payload.DeleteIDs)
	resultTotal := int(existingCount) - deleteCount + newCreates
	if resultTotal < 0 {
		resultTotal = 0
	}
	if resultTotal > MaxContactsBatch {
		tx.Rollback()
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "resulting contacts exceed limit"})
		return
	}

	if len(payload.Upsert) > 0 {
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"name", "email", "phone", "position", "sponsor_id", "updated_at"}),
		}).Create(&payload.Upsert).Error; err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if len(payload.DeleteIDs) > 0 {
		if err := tx.Where("id IN ? AND sponsor_id = ?", payload.DeleteIDs, sponsorID).
			Delete(&entity.SponsorContact{}).Error; err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var contacts []entity.SponsorContact
	if err := config.DB.Where("sponsor_id = ?", sponsorID).Find(&contacts).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"contacts": contacts})
}

// DELETE /sponsors/:id
func DeleteSponsor(ctx *gin.Context) {
	id := ctx.Param("id")

	// DELETE
	if err := config.DB.Delete(&entity.Sponsor{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Sponsor delete successfully"})
}