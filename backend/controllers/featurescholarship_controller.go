package controllers

import (
	"backend/config"
	"backend/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

// getall
func GetAllFeaturescholarship(c *gin.Context) {
	var items []entity.Featurescholarship
	if err := config.DB.Preload("Scholarship").Preload("Typefeature").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// getby id
func GetFeaturescholarshipByID(c *gin.Context) {
	id := c.Param("id")
	var item entity.Featurescholarship
	if err := config.DB.Preload("Scholarship").Preload("Typefeature").First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// post
func CreateFeaturescholarship(c *gin.Context) {
	var item entity.Featurescholarship
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
func UpdateFeaturescholarship(c *gin.Context) {
	id := c.Param("id")
	var item entity.Featurescholarship
	if err := config.DB.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Featurescholarship not found"})
		return
	}

	var input entity.Featurescholarship
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Model(&item).Updates(input)
	c.JSON(http.StatusOK, item)
}

// delete
func DeleteFeaturescholarship(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&entity.Featurescholarship{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Featurescholarship deleted"})
}
