package controllers

import (
	"log"
	"net/http"
	"project-its/initializers"
	"project-its/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Timeline struct {
	ID     string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()";json:"id"`
	Title  string `json:"title"`
	Start  string `json:"start"`
	End    string `json:"end"`
	AllDay bool   `json:"allDay"`
	Color  string `json:"color"`
}

// Create a new event
func GetEventsTimeline(c *gin.Context) {
	var events []models.Timeline
	if err := initializers.DB.Find(&events).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, events)
}

// Example of using generated UUID
func CreateEventTimeline(c *gin.Context) {
	var event Timeline
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate a new UUID if not provided
	if event.ID == "" {
		event.ID = generateUUID()
	}

	// Log the generated UUID for debugging
	log.Printf("Generated UUID: %s", event.ID)

	if err := initializers.DB.Create(&event).Error; err != nil {
		log.Printf("Error creating event: %v", err) // Add this line
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func DeleteEventTimeline(c *gin.Context) {
	id := c.Param("id") // Menggunakan c.Param jika UUID dikirim sebagai bagian dari URL
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID harus disertakan"})
		return
	}
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}
	if err := initializers.DB.Where("id = ?", id).Delete(&Timeline{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}