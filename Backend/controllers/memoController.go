package controllers

import (
	"log"
	"project-gin/initializers"
	"project-gin/models"
	"time"

	"github.com/gin-gonic/gin"
)

type MemoRequest struct {
	Tanggal string `json:"tanggal"`
	NoMemo  string `json:"no_memo"`
	Perihal string `json:"perihal"`
	Pic     string `json:"pic"`
}

func MemoIndex(c *gin.Context) {

	var memos []models.Memo

	initializers.DB.Find(&memos)

	c.JSON(200, gin.H{
		"memos": memos,
	})

}

func MemoCreate(c *gin.Context) {

	var requestBody MemoRequest

	if err := c.BindJSON(&requestBody); err != nil {
		c.Status(400)
		c.Error(err) // log the error
		return
	}

	// Add some logging to see what's being received
	log.Println("Received request body:", requestBody)

	// Parse the date string
	tanggalString := requestBody.Tanggal
	tanggal, err := time.Parse("2006-01-02", tanggalString)
	if err != nil {
		c.Status(400)
		c.JSON(400, gin.H{"error": "Invalid date format: " + err.Error()})
		return
	}

	memo := models.Memo{
		Tanggal: tanggal,
		NoMemo:  requestBody.NoMemo,
		Perihal: requestBody.Perihal,
		Pic:     requestBody.Pic,
	}

	result := initializers.DB.Create(&memo)

	if result.Error != nil {
		c.Status(400)
		c.JSON(400, gin.H{"error": "Failed to create Memo: " + result.Error.Error()})
		return
	}

	c.JSON(201, gin.H{
		"memo": memo,
	})

}

func MemoShow(c *gin.Context) {

	id := c.Params.ByName("id")

	var memo models.Memo

	initializers.DB.First(&memo, id)

	c.JSON(200, gin.H{
		"memo": memo,
	})

}

func MemoUpdate(c *gin.Context) {

	var requestBody MemoRequest

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	tanggalString := requestBody.Tanggal
	tanggal, err := time.Parse("2006-01-02", tanggalString)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid date format: " + err.Error()})
		return
	}

	id := c.Params.ByName("id")

	var memo models.Memo

	if err := initializers.DB.First(&memo, id); err.Error != nil {
		c.JSON(404, gin.H{"error": "Memo not found"})
		return
	}

	if err := initializers.DB.Model(&memo).Updates(models.Memo{
		Tanggal: tanggal,
		NoMemo:  requestBody.NoMemo,
		Perihal: requestBody.Perihal,
		Pic:     requestBody.Pic,
	}). Error; err != nil {
		c.JSON(400, gin.H{"error": "Failed to update Memo: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"memo": memo,
	})

}

func MemoDelete(c *gin.Context) {

	id := c.Params.ByName("id")

	var memo models.Memo

	if err := initializers.DB.First(&memo, id); err.Error != nil {
		c.JSON(404, gin.H{"error": "Memo not found"})
		return
	}

	if err := initializers.DB.Delete(&memo).Error; err != nil {
		c.JSON(400, gin.H{"error": "Failed to delete Memo: " + err.Error()})
		return
	}

	c.Status(204)

}
