package controllers

import (
	"log"
	"net/http"
	"project-its/initializers"
	"project-its/models"
	"time"

	"github.com/gin-gonic/gin"
)

// Create a new event
func GetEventsBookingRapat(c *gin.Context) {
	var events []models.BookingRapat
	if err := initializers.DB.Find(&events).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"booking": events})
}

// Example of using generated UUID
func CreateEventBookingRapat(c *gin.Context) {
	var event models.BookingRapat
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set notification menggunakan fungsi dari notificationController
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Printf("Error loading location: %v", err)
		return
	}

	var startTime time.Time
	if event.AllDay {
		// Jika AllDay = true, set waktu ke tengah malam
		startTime, err = time.ParseInLocation("2006-01-02T15:04:05", event.Start+"T00:00:00", loc)
	} else {
		// Jika tidak, parse dengan format RFC3339
		startTime, err = time.ParseInLocation(time.RFC3339, event.Start, loc)
	}

	if err != nil {
		log.Printf("Error parsing start time: %v", err)
		return
	}

	SetNotification(event.Title, startTime, "BookingRapat") // Panggil fungsi SetNotification

	if err := initializers.DB.Create(&event).Error; err != nil {
		log.Printf("Error creating event: %v", err) // Add this line
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, event)
}

func DeleteEventBookingRapat(c *gin.Context) {
	id := c.Param("id") // Menggunakan c.Param jika UUID dikirim sebagai bagian dari URL
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID harus disertakan"})
		return
	}
	if err := initializers.DB.Where("id = ?", id).Delete(&models.BookingRapat{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
