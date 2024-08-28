package controllers

import (
	"log"
	"net/http"
	"project-its/initializers"
	"project-its/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RuangRapat struct {
	ID     string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()";json:"id"`
	Title  string `json:"title"`
	Start  string `json:"start"`
	End    string `json:"end"`
	AllDay bool   `json:"allDay"`
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

	// Parse Start time
	startTime, err := time.Parse(time.RFC3339, event.Start)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start time format"})
		return
	}

	// Generate a new UUID if not provided
	if event.ID == "" {
		event.ID = generateUUID()
	}

	// Parse UUID
	eventUUID, err := uuid.Parse(event.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	// Log the generated UUID for debugging
	log.Printf("Generated UUID: %s", event.ID)

	if err := initializers.DB.Create(&event).Error; err != nil {
		log.Printf("Error creating event: %v", err) // Add this line
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Add notifications
	userID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000") // Ganti dengan UUID yang valid
	notifyAtOneDayBefore := startTime.Add(-24 * time.Hour)
	addNotification(userID, eventUUID, notifyAtOneDayBefore, "Notifikasi sehari sebelum event")

	notifyAtOneHourBefore := startTime.Add(-1 * time.Hour)
	addNotification(userID, eventUUID, notifyAtOneHourBefore, "Notifikasi satu jam sebelum event")

	c.Status(http.StatusNoContent)
}

func DeleteEvent(c *gin.Context) {
	id := c.Param("id") // Menggunakan c.Param jika UUID dikirim sebagai bagian dari URL
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID harus disertakan"})
		return
	}
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}
	if err := initializers.DB.Where("id = ?", id).Delete(&RuangRapat{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// Mendapatkan semua notifikasi
func GetAllNotifications(c *gin.Context) {
	var notifications []models.Notification

	if err := initializers.DB.Find(&notifications).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
			return
	}

	c.JSON(http.StatusOK, gin.H{"notifications": notifications})
}

// Mendapatkan semua notifikasi untuk pengguna tertentu
func GetNotifications(c *gin.Context) {
	userID := c.Param("user_id") // asumsikan user_id dikirim sebagai parameter URL
	var notifications []models.Notification

	if err := initializers.DB.Where("user_id = ?", userID).Find(&notifications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"notifications": notifications})
}

func addNotification(userID uuid.UUID, eventID uuid.UUID, notifyAt time.Time, message string) {
	notification := models.Notification{
		UserID:   userID,
		EventID:  eventID, // Pastikan ini disertakan
		Message:  message,
		NotifyAt: notifyAt,
		Sent:     false,
	}
	if err := initializers.DB.Create(&notification).Error; err != nil {
		log.Printf("Failed to schedule notification: %v", err)
		return
	}
	log.Printf("Notification scheduled for event %s at %s", eventID, notifyAt.Format("2006-01-02 15:04:05"))
}
