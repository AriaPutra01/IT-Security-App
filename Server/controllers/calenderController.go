package controllers

import (
	"net/http"
	"project-gin/initializers"
	"project-gin/models"
	"log"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RuangRapat struct {
	ID     string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Title  string
	Start  string
	End    string
	AllDay bool
}

func generateUUID() string {
	return uuid.New().String()
}

// Create a new event
func GetEvents(c *gin.Context) {
	var events []models.RuangRapat
	if err := initializers.DB.Find(&events).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, events)
}

// Example of using generated UUID
func CreateEvent(c *gin.Context) {
	var event RuangRapat
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

func DeleteEvent(c *gin.Context) {
	id := c.Query("id")
	if err := initializers.DB.Delete(&RuangRapat{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}