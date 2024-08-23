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
	tanggal, err := time.Parse("2006-01-02", requestBody.Tanggal)
    if err != nil {
        log.Printf("Error parsing date: %v", err)
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
		"BeritaAcara": berita_acara,
	})

}

func BeritaAcaraIndex(c *gin.Context) {

	// Get models from DB
	var berita_acara []models.BeritaAcara
	initializers.DB.Find(&berita_acara)

	//Respond with them
	c.JSON(200, gin.H{
		"BeritaAcara": berita_acara,
	})
}

func BeritaAcaraShow(c *gin.Context) {

	id := c.Params.ByName("id")
	// Get models from DB
	var berita_acara models.BeritaAcara

	initializers.DB.First(&berita_acara, id)

	//Respond with them
	c.JSON(200, gin.H{
		"BeritaAcara": berita_acara,
	})
}

func BeritaAcaraUpdate(c *gin.Context) {

	var requestBody beritaRequest

	if err := c.BindJSON(&requestBody); err != nil {
		c.Status(400)
		c.Error(err) // log the error
		return
	}

    id := c.Params.ByName("id")

    var berita_acara models.BeritaAcara
    initializers.DB.First(&berita_acara, id)

    if err := initializers.DB.First(&berita_acara, id).Error; err != nil {
        c.JSON(404, gin.H{"error": "Berita Acara tidak ditemukan"})
        return
    }

    if requestBody.Tanggal != "" {
        tanggal, err := time.Parse("2006-01-02", requestBody.Tanggal)
        if err != nil {
            c.JSON(400, gin.H{"error": "Format tanggal tidak valid: " + err.Error()})
            return
        }
		berita_acara.Tanggal = tanggal
    }

    if requestBody.NoSurat != "" {
		berita_acara.NoSurat = requestBody.NoSurat
    } else {
		berita_acara.NoSurat = berita_acara.NoSurat // gunakan nilai yang ada dari database
    }

	if requestBody.Perihal != "" {
		berita_acara.Perihal = requestBody.Perihal
	} else {
		berita_acara.Perihal = berita_acara.Perihal // gunakan nilai yang ada dari database
	}
	
	if requestBody.Pic != "" {
		berita_acara.Pic = requestBody.Pic
	} else {
		berita_acara.Pic = berita_acara.Pic // gunakan nilai yang ada dari database
	}

    initializers.DB.Model(&berita_acara).Updates(berita_acara)

    c.JSON(200, gin.H{
        "BeritaAcara": berita_acara,
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
