package controllers

import (
	"log"
	"net/http"
	"project-its/initializers"
	"project-its/models"
	"time"

	"github.com/gin-gonic/gin"
)

func SetNotification(title string, startTime time.Time, category string) {
	// Set lokasi ke WIB
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Printf("Error loading location: %v", err)
		return
	}

	// Parse waktu mulai event ke WIB
	startTime, err = time.ParseInLocation(time.RFC3339, startTime.Format(time.RFC3339), loc) // Ubah ini
	if err != nil {
		log.Printf("Error parsing start time: %v", err)
		return
	}
	log.Printf("Parsed start time in WIB: %v", startTime)

	// Tentukan waktu notifikasi 24 jam sebelum event
	notificationTime24 := startTime.Add(-24 * time.Hour)
	log.Printf("24-hour notification scheduled for %s", notificationTime24)

	// Tentukan waktu notifikasi 1 jam sebelum event
	notificationTime0 := startTime.Add(time.Hour)
	log.Printf("1-hour notification  for %s", notificationTime0)

	// Simulasi pengiriman notifikasi 24 jam sebelum event
	go func() {
		time.Sleep(time.Until(notificationTime24))
		log.Printf("24-hour notification sent for event %s at %s", title, notificationTime24) // Ubah ini
	}()

	// Simulasi pengiriman notifikasi 1 jam sebelum event
	go func() {
		time.Sleep(time.Until(notificationTime0))
		log.Printf("1-hour notification sent for event %s at %s", title, notificationTime0) // Ubah ini
	}()

	notification := models.Notification{
		Title:    title,
		Start:    startTime,
		Category: category,
	}
	if err := initializers.DB.Create(&notification).Error; err != nil {
		log.Printf("Error creating notification: %v", err)
	}
	log.Printf("Notification created with category: %s, title: %s", category, title) // Tambahkan log ini

}

func GetNotifications(c *gin.Context) {
	var notifications []models.Notification
	if err := initializers.DB.Find(&notifications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, notifications)
}

func DeleteNotification(c *gin.Context) {
	id := c.Param("id")
	log.Printf("ID yang diterima untuk dihapus: %s", id) // Tambahkan log ini
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID harus disertakan"})
		return
	}

	// Menghapus notifikasi berdasarkan ID
	if err := initializers.DB.Where("id = ?", id).Delete(&models.Notification{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent) // Mengembalikan status 204 No Content
}
