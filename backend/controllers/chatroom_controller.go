package controllers

import (
	"backend/config"
	"backend/entity"
	"backend/ws"
	"net/http"

	"github.com/gin-gonic/gin"
)

// getall
func GetAllChatroom(c *gin.Context) {
	var items []entity.Chatroom
	if err := config.DB.Preload("User").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// getby id
func GetChatroomByID(c *gin.Context) {
	id := c.Param("id")
	var item entity.Chatroom
	if err := config.DB.Preload("User").First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// post
func CreateChatroom(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	// 🔒 check open room ก่อน
	var existing entity.Chatroom
	if err := config.DB.
		Where("user_id = ? AND statuschatroom = ?", userID, "open").
		First(&existing).Error; err == nil {
		// มีอยู่แล้ว → ส่งกลับไปเลย
		c.JSON(http.StatusOK, existing)
		return
	}

	room := entity.Chatroom{
		UserID:         userID,
		Statuschatroom: "open",
	}

	config.DB.Create(&room)

	// Broadcast to admins (via ApprovalHub reuse)
	if ws.ApprovalHubInstance != nil {
		ws.ApprovalHubInstance.BroadcastUpdate("chatroom_updated", room)
	}

	c.JSON(http.StatusCreated, room)
}

// put
func UpdateChatroom(c *gin.Context) {
	id := c.Param("id")
	var item entity.Chatroom
	if err := config.DB.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chatroom not found"})
		return
	}

	var input entity.Chatroom
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Model(&item).Updates(input)
	c.JSON(http.StatusOK, item)
}

// delete
func DeleteChatroom(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&entity.Chatroom{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Chatroom deleted"})
}

func GetMyOpenChatroom(c *gin.Context) {
	userIDAny, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDAny.(uint)

	var room entity.Chatroom
	err := config.DB.
		Where("user_id = ? AND statuschatroom = ?", userID, "open").
		First(&room).Error

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": room})
}

func GetAllOpenChatrooms(c *gin.Context) {
	var rooms []entity.Chatroom
	if err := config.DB.Preload("User.StudentProfiles").Where("statuschatroom = ?", "open").Find(&rooms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rooms})
}
