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

func GetNotifications(c *gin.Context) {
	var notifications []models.Notification
	if err := initializers.DB.Find(&notifications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"notifications": notifications})
}

func DeleteNotification(c *gin.Context) {
	id := c.Param("id")
	log.Printf("Received ID: %s", id)

	notificationID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	if err := initializers.DB.Where("id = ?", notificationID).Delete(&models.Notification{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func setNotification(event interface{}) {
	var startTime time.Time
	var title string

	switch e := event.(type) {
	case *models.JadwalRapat:
		loc, err := time.LoadLocation("Asia/Jakarta")
		if err != nil {
			log.Printf("Error loading location: %v", err)
			return
		}
		startTime, err = time.ParseInLocation(time.RFC3339, e.Start, loc)
		if err != nil {
			log.Printf("Error parsing start time: %v", err)
			return
		}
		title = e.Title
	case *models.JadwalCuti:
		loc, err := time.LoadLocation("Asia/Jakarta")
		if err != nil {
			log.Printf("Error loading location: %v", err)
			return
		}
		startTime, err = time.ParseInLocation(time.RFC3339, e.Start, loc)
		if err != nil {
			log.Printf("Error parsing start time: %v", err)
			return
		}
		title = e.Title
	case *models.BookingRapat:
		loc, err := time.LoadLocation("Asia/Jakarta")
		if err != nil {
			log.Printf("Error loading location: %v", err)
			return
		}
		startTime, err = time.ParseInLocation(time.RFC3339, e.Start, loc)
		if err != nil {
			log.Printf("Error parsing start time: %v", err)
			return
		}
		title = e.Title
	case *models.TimelineProject:
		loc, err := time.LoadLocation("Asia/Jakarta")
		if err != nil {
			log.Printf("Error loading location: %v", err)
			return
		}
		startTime, err = time.ParseInLocation(time.RFC3339, e.Start, loc)
		if err != nil {
			log.Printf("Error parsing start time: %v", err)
			return
		}
		title = e.Title
	case *models.TimelineDesktop:
		loc, err := time.LoadLocation("Asia/Jakarta")
		if err != nil {
			log.Printf("Error loading location: %v", err)
			return
		}
		startTime, err = time.ParseInLocation(time.RFC3339, e.Start, loc)
		if err != nil {
			log.Printf("Error parsing start time: %v", err)
			return
		}
		title = e.Title
	default:
		log.Printf("Unknown event type")
		return
	}

	notificationTime24 := startTime.Add(-24 * time.Hour)
	notificationTime1 := startTime.Add(-1 * time.Hour)

	go func() {
		time.Sleep(time.Until(notificationTime24))
		log.Printf("24-hour notification sent for event %s at %s", title, notificationTime24)
	}()

	go func() {
		time.Sleep(time.Until(notificationTime1))
		log.Printf("1-hour notification sent for event %s at %s", title, notificationTime1)
	}()

	notification := models.Notification{
		Title: title,
		Start: startTime.Add(-1 * time.Hour),
	}
	if err := initializers.DB.Create(&notification).Error; err != nil {
		log.Printf("Error creating notification: %v", err)
	}
}
