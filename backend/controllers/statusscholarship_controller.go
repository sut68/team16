package controllers

import (
	"backend/config"
	"backend/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

// getall
func GetAllStatusscholarship(c *gin.Context) {
	var items []entity.Statusscholarship
	if err := config.DB.Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// getby id
func GetStatusscholarshipByID(c *gin.Context) {
	id := c.Param("id")
	var item entity.Statusscholarship
	if err := config.DB.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// post
func CreateStatusscholarship(c *gin.Context) {
	var item entity.Statusscholarship
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
func UpdateStatusscholarship(c *gin.Context) {
	id := c.Param("id")
	var item entity.Statusscholarship
	if err := config.DB.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Statusscholarship not found"})
		return
	}

	var input entity.Statusscholarship
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Model(&item).Updates(input)
	c.JSON(http.StatusOK, item)
}

// delete
func DeleteStatusscholarship(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&entity.Statusscholarship{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Statusscholarship deleted"})
}