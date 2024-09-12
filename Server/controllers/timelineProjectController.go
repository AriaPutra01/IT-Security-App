package controllers

import (
	"log"
	"net/http"
	"project-its/initializers"
	"project-its/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetEventsTimeline retrieves all timeline events
func GetEventsProject(c *gin.Context) {
	var events []models.TimelineProject
	if err := initializers.DB.Find(&events).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"events": events})
}

// CreateEventTimeline creates a new timeline event
func CreateEventProject(c *gin.Context) {
	var event models.TimelineProject
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parsing waktu untuk notifikasi
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error loading location"})
		return
	}

	// Ubah format parsing sesuai dengan format yang diterima
	startTime, err := time.ParseInLocation("2006-01-02 15:04:05", event.Start, loc) // Ubah format di sini
	if err != nil {
		log.Printf("Error parsing start time: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing start time"})
		return
	}

	// Panggil fungsi SetNotification
	SetNotification(event.Title, startTime, "TimelineProject")

	if err := initializers.DB.Create(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, event)
}

// DeleteEventTimeline deletes a timeline event by ID
func DeleteEventProject(c *gin.Context) {
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

	// Pastikan ID dikonversi ke tipe data yang sesuai
	if err := initializers.DB.Where("id = ?", uint(id)).Delete(&models.TimelineProject{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// Resources

// GetResources retrieves all resources
func GetResourcesProject(c *gin.Context) {
	var resources []models.ResourceProject
	if err := initializers.DB.Find(&resources).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"resources": resources})
}

// CreateResource creates a new resource
func CreateResourceProject(c *gin.Context) {
	var resource models.ResourceProject
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
func DeleteResourceProject(c *gin.Context) {
	idParam := c.Param("id")
	if idParam == "" || idParam == "undefined" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID harus disertakan dan valid"})
		return
	}

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	log.Printf("Attempting to delete ResourceProject with ID: %d", id)

	if err := initializers.DB.Where("id = ?", uint(id)).Delete(&models.ResourceProject{}).Error; err != nil {
		log.Printf("Error deleting ResourceProject with ID: %d, error: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Successfully deleted ResourceProject with ID: %d", id)
	c.Status(http.StatusNoContent)
}
