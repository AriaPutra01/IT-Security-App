package controllers

import (
	"fmt"
	"log"
	"project-its/initializers"
	"project-its/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/robfig/cron"
)

// Struktur data untuk request notifikasi
type NotificationRequest struct {
	UserID   uuid.UUID `json:"user_id"` // Ubah tipe data dari int menjadi uint
	Message  string    `json:"message"`
	NotifyAt time.Time `json:"notify_at"`
}

// Menambahkan notifikasi
func AddNotification(c *gin.Context) {
	var data NotificationRequest
	if err := c.BindJSON(&data); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	log.Println("Attempting to create notification for user:", data.UserID)
	notification := models.Notification{
		UserID:   data.UserID,
		Message:  data.Message,
		NotifyAt: data.NotifyAt,
	}

	// Format NotifyAt untuk output
	formattedNotifyAt := notification.NotifyAt.Format("2006-01-02 jam 15.04")
	log.Println("Notification will be sent at:", formattedNotifyAt)

	if err := initializers.DB.Create(&notification).Error; err != nil {
		log.Printf("Failed to create notification: %s", err)
		c.JSON(500, gin.H{"error": "Could not create notification"})
		return
	}

	c.JSON(200, gin.H{"message": "Notification scheduled", "notify_at": formattedNotifyAt})
}

// Scheduler untuk mengirim notifikasi
func SendNotifications() {
	var notifications []models.Notification
	now := time.Now()

	initializers.DB.Where("notify_at <= ? AND sent = false", now).Find(&notifications)

	if len(notifications) == 0 {
		log.Println("No notifications to send at", now)
	} else {
		log.Printf("Found %d notifications to send", len(notifications))
	}

	for _, notification := range notifications {
		fmt.Println("Sending notification to user:", notification.UserID, "Message:", notification.Message)
		initializers.DB.Model(&notification).Update("sent", true)
	}
}

// Inisialisasi scheduler
func InitNotificationScheduler() {
	c := cron.New()
	c.AddFunc("@every 1m", SendNotifications)
	c.Start()
}
