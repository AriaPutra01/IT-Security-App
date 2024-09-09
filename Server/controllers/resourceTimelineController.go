package controllers

import (
	"net/http"
	"project-its/initializers"
	"project-its/models"

	"github.com/gin-gonic/gin"
)

// GetResources retrieves all resources
func GetResourcesTimeline(c *gin.Context) {
	var resources []models.ResourceTimeline
	if err := initializers.DB.Find(&resources).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"resources": resources})
}

// CreateResource creates a new resource
func CreateResourceTimeline(c *gin.Context) {
	var resource models.ResourceTimeline
	if err := c.ShouldBindJSON(&resource); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := initializers.DB.Create(&resource).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resource)
}

// DeleteResource deletes a resource by ID
func DeleteResourceTimeline(c *gin.Context) {
	id := c.Param("id")
	if err := initializers.DB.Where("id = ?", id).Delete(&models.ResourceTimeline{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
