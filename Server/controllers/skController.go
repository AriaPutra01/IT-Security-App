package controllers

import (
	"project-gin/initializers"
	"project-gin/models"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type skRequest struct {
	Tanggal string `json:"tanggal"`
	NoSurat string `json:"no_surat"`
	Perihal string `json:"perihal"`
	Pic     string `json:"pic"`
}

func SkCreate(c *gin.Context) {

	// Get data off req body
	var requestBody skRequest

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

	sk := models.Sk{
		Tanggal: tanggal,
		NoSurat: requestBody.NoSurat,
		Perihal: requestBody.Perihal,
		Pic:     requestBody.Pic,
	}

	result := initializers.DB.Create(&sk)

	if result.Error != nil {
		c.Status(400)
		return
	}

	// Return it
	c.JSON(200, gin.H{
		"SK": sk ,
	})

}

func SkIndex(c *gin.Context) {

	// Get models from DB
	var sk []models.Sk
	initializers.DB.Find(&sk)

	//Respond with them
	c.JSON(200, gin.H{
		"SK": sk,
	})
}

func SkShow(c *gin.Context) {

	id := c.Params.ByName("id")
	// Get models from DB
	var sk models.Sk

	initializers.DB.First(&sk, id)

	//Respond with them
	c.JSON(200, gin.H{
		"SK": sk,
	})
}

func SkUpdate(c *gin.Context) {

	var requestBody skRequest

	if err := c.BindJSON(&requestBody); err != nil {
		c.Status(400)
		c.Error(err) // log the error
		return
	}

	tanggalString := requestBody.Tanggal
	tanggal, err := time.Parse("2006-01-02", tanggalString)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid date format: " + err.Error()})
		return
	}

	//get id
	id := c.Params.ByName("id")

	// find the sk
	var sk models.Sk
	initializers.DB.First(&sk, id)

	if err := initializers.DB.First(&sk, id); err.Error != nil {
		c.JSON(404, gin.H{"error": "sk not found"})
		return
	}

	// update it
	initializers.DB.Model(&sk).Updates(models.Sk{
		Tanggal: tanggal,
		NoSurat: requestBody.NoSurat,
		Perihal: requestBody.Perihal,
		Pic:     requestBody.Pic,
	})

	//Respond with them
	c.JSON(200, gin.H{
		"SK": sk,
	})

}

func SkDelete(c *gin.Context) {

	//get id
	id := c.Params.ByName("id")

	// find the sk
	var sk models.Sk

	if err := initializers.DB.First(&sk, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "sk not found"})
		return
	}

	/// delete it
	if err := initializers.DB.Delete(&sk).Error; err != nil {
		c.JSON(404, gin.H{"error": "sk Failed to Delete"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Deleted",
	})
}


