package controllers

import (
	"net/http"
	"project-its/initializers"
	"project-its/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetEventsTimeline retrieves all timeline events
func GetEventsTimeline(c *gin.Context) {
	var events []models.Timeline
	if err := initializers.DB.Find(&events).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"events": events})
}

// CreateEventTimeline creates a new timeline event
func CreateEventTimeline(c *gin.Context) {
	var event models.Timeline
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := initializers.DB.Create(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, event)
}

// DeleteEventTimeline deletes a timeline event by ID
func DeleteEventTimeline(c *gin.Context) {
	idParam := c.Param("id")
	if idParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID harus disertakan"})
		return
	}

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	if err := initializers.DB.Where("id = ?", uint(id)).Delete(&models.Timeline{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
