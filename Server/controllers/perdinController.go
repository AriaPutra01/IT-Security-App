package controllers

import (
	"project-gin/initializers"
	"project-gin/models"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type perdinRequest struct {
	NoPerdin  string `json:"no_perdin"`
	Tanggal string `json:"tanggal"`
	Hotel string `json:"hotel"`
	Transport     string `json:"transport"`
}

func PerdinCreate(c *gin.Context) {

	// Get data off req body
	var requestBody perdinRequest

	if err := c.BindJSON(&requestBody); err != nil {
		c.Status(400)
		c.Error(err) // log the error
		return
	}

	// Add some logging to see what's being received
	log.Println("Received request body:", requestBody)


	tanggal, err := time.Parse("2006-01-02", requestBody.Tanggal)
    if err != nil {
        log.Printf("Error parsing date: %v", err)
        c.JSON(400, gin.H{"error": "Invalid date format: " + err.Error()})
        return
    }

	perdin := models.Perdin{
		NoPerdin:  requestBody.NoPerdin,
		Tanggal: tanggal,
		Hotel: requestBody.Hotel,
		Transport: requestBody.Transport,
	}

	result := initializers.DB.Create(&perdin)

	if result.Error != nil {
		c.Status(400)
		return
	}

	// Return it
	c.JSON(200, gin.H{
		"Perdin": perdin,
	})

}

func PerdinIndex(c *gin.Context) {

	// Get models from DB
	var perdin []models.Perdin
	initializers.DB.Find(&perdin)

	//Respond with them
	c.JSON(200, gin.H{
		"Perdin": perdin,	
	})
}

func PerdinShow(c *gin.Context) {

	id := c.Params.ByName("id")
	// Get models from DB
	var perdin models.Perdin

	initializers.DB.First(&perdin, id)

	//Respond with them
	c.JSON(200, gin.H{
		"Perdin": perdin,
	})
}

func PerdinUpdate(c *gin.Context) {

	var requestBody perdinRequest

	if err := c.BindJSON(&requestBody); err != nil {
		c.Status(400)
		c.Error(err) // log the error
		return
	}

    id := c.Params.ByName("id")

    var perdin models.Perdin
    initializers.DB.First(&perdin, id)

    if err := initializers.DB.First(&perdin, id).Error; err != nil {
        c.JSON(404, gin.H{"error": "perdin tidak ditemukan"})
        return
    }

    if requestBody.Tanggal != "" {
        tanggal, err := time.Parse("2006-01-02", requestBody.Tanggal)
        if err != nil {
            c.JSON(400, gin.H{"error": "Format tanggal tidak valid: " + err.Error()})
            return
        }
        perdin.Tanggal = tanggal
    }

    if requestBody.NoPerdin != "" {
        perdin.NoPerdin = requestBody.NoPerdin
    } else {
		perdin.NoPerdin = perdin.NoPerdin // gunakan nilai yang ada dari database
    }

	if requestBody.Transport != "" {
		perdin.Transport = requestBody.Transport
	} else {
		perdin.Transport = perdin.Transport // gunakan nilai yang ada dari database
	}
	
	if requestBody.Hotel != "" {
		perdin.Hotel = requestBody.Hotel
	} else {
		perdin.Hotel = perdin.Hotel // gunakan nilai yang ada dari database
	}

    initializers.DB.Model(&perdin).Updates(perdin)

    c.JSON(200, gin.H{
        "perdin": perdin,
    })

}

func PerdinDelete(c *gin.Context) {

	//get id
	id := c.Params.ByName("id")

	// find the perdin
	var perdin models.Perdin

	if err := initializers.DB.First(&perdin, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Perdin not found"})
		return
	}

	/// delete it
	if err := initializers.DB.Delete(&perdin).Error; err != nil {
		c.JSON(404, gin.H{"error": "Perdin Failed to Delete"})
		return
	}

	c.JSON(200, gin.H{
		"Perdin": "Perdin deleted",
	})
}
