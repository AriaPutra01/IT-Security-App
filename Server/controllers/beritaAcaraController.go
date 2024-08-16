package controllers

import (
	"project-gin/initializers"
	"project-gin/models"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type beritaRequest struct {
	Tanggal string `json:"tanggal"`
	NoSurat string `json:"no_surat"`
	Perihal string `json:"perihal"`
	Pic     string `json:"pic"`
}

func BeritaAcaraCreate(c *gin.Context) {

	// Get data off req body
	var requestBody beritaRequest

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

	berita_acara := models.BeritaAcara{
		Tanggal: tanggal,
		NoSurat: requestBody.NoSurat,
		Perihal: requestBody.Perihal,
		Pic:     requestBody.Pic,
	}

	result := initializers.DB.Create(&berita_acara)

	if result.Error != nil {
		c.Status(400)
		return
	}

	// Return it
	c.JSON(200, gin.H{
		"Berita Acara": berita_acara,
	})

}

func BeritaAcaraIndex(c *gin.Context) {

	// Get models from DB
	var berita_acara []models.BeritaAcara
	initializers.DB.Find(&berita_acara)

	//Respond with them
	c.JSON(200, gin.H{
		"Berita Acara": berita_acara,
	})
}

func BeritaAcaraShow(c *gin.Context) {

	id := c.Params.ByName("id")
	// Get models from DB
	var berita_acara models.BeritaAcara

	initializers.DB.First(&berita_acara, id)

	//Respond with them
	c.JSON(200, gin.H{
		"Berita Acara": berita_acara,
	})
}

func BeritaAcaraUpdate(c *gin.Context) {

	var requestBody beritaRequest

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

	// find the berita_acara
	var berita_acara models.BeritaAcara
	initializers.DB.First(&berita_acara, id)

	if err := initializers.DB.First(&berita_acara, id); err.Error != nil {
		c.JSON(404, gin.H{"error": "berita_acara not found"})
		return
	}

	// update it
	initializers.DB.Model(&berita_acara).Updates(models.BeritaAcara{
		Tanggal: tanggal,
		NoSurat: requestBody.NoSurat,
		Perihal: requestBody.Perihal,
		Pic:     requestBody.Pic,
	})

	//Respond with them
	c.JSON(200, gin.H{
		"Berita Acara": berita_acara,
	})

}

func BeritaAcaraDelete(c *gin.Context) {

	//get id
	id := c.Params.ByName("id")

	// find the berita_acara
	var berita_acara models.BeritaAcara

	if err := initializers.DB.First(&berita_acara, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "berita_acara not found"})
		return
	}

	/// delete it
	if err := initializers.DB.Delete(&berita_acara).Error; err != nil {
		c.JSON(404, gin.H{"error": "berita_acara Failed to Delete"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Deleted",
	})
}
