package controllers

import (
	"backend/config"
	"backend/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

// getall
func GetAllAssistance(c *gin.Context) {
	var items []entity.Assistance
	if err := config.DB.Preload("Chatroom").Preload("Sender").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// getby id
func GetAssistanceByID(c *gin.Context) {
	id := c.Param("id")
	var item entity.Assistance
	if err := config.DB.Preload("Chatroom").Preload("Sender").First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// post
func CreateAssistance(c *gin.Context) {
	var item entity.Assistance
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// put
func UpdateAssistance(c *gin.Context) {
	id := c.Param("id")
	var item entity.Assistance
	if err := config.DB.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Assistance not found"})
		return
	}

	var input entity.Assistance
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Model(&item).Updates(input)
	c.JSON(http.StatusOK, item)
}

// delete
func DeleteAssistance(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&entity.Assistance{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Assistance deleted"})
}
