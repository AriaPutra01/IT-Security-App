package controllers

import (
	"project-gin/initializers"
	"project-gin/models"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type isoRequest struct {
	Tanggal string `json:"tanggal"`
	NoMemo  string `json:"no_memo"`
	Perihal string `json:"perihal"`
	Pic     string `json:"pic"`
}


func IsoCreate(c *gin.Context) {
    var requestBody isoRequest

    if err := c.BindJSON(&requestBody); err != nil {
        log.Printf("Error binding JSON: %v", err)
        c.Status(400)
        c.Error(err) // log the error
        return
    }

    log.Println("Received request body:", requestBody)

    tanggal, err := time.Parse("2006-01-02", requestBody.Tanggal)
    if err != nil {
        log.Printf("Error parsing date: %v", err)
        c.JSON(400, gin.H{"error": "Invalid date format: " + err.Error()})
        return
    }

    iso := models.Iso{
        Tanggal: tanggal,
        NoMemo:  requestBody.NoMemo,
        Perihal: requestBody.Perihal,
        Pic:     requestBody.Pic,
    }

    result := initializers.DB.Create(&iso)
    if result.Error != nil {
        log.Printf("Error saving to database: %v", result.Error)
        c.Status(400)
        return
    }

    log.Println("Iso created successfully:", iso)
    c.JSON(200, gin.H{"Iso": iso})
}

func IsoIndex(c *gin.Context) {

	// Get models from DB
	var iso []models.Iso
	initializers.DB.Find(&iso)

	//Respond with them
	c.JSON(200, gin.H{
		"Iso": iso,
	})
}

func IsoShow(c *gin.Context) {

	id := c.Params.ByName("id")
	// Get models from DB
	var iso models.Iso

	initializers.DB.First(&iso, id)

	//Respond with them
	c.JSON(200, gin.H{
		"Iso": &iso,
	})
}

func IsoUpdate(c *gin.Context) {
    var requestBody isoRequest

    if err := c.BindJSON(&requestBody); err != nil {
        c.Status(400)
        c.Error(err) // log the error
        return
    }

    id := c.Params.ByName("id")

    var iso models.Iso
    initializers.DB.First(&iso, id)

    if err := initializers.DB.First(&iso, id).Error; err != nil {
        c.JSON(404, gin.H{"error": "Iso tidak ditemukan"})
        return
    }

    if requestBody.Tanggal != "" {
        tanggal, err := time.Parse("2006-01-02", requestBody.Tanggal)
        if err != nil {
            c.JSON(400, gin.H{"error": "Format tanggal tidak valid: " + err.Error()})
            return
        }
        iso.Tanggal = tanggal
    }

    if requestBody.NoMemo != "" {
        iso.NoMemo = requestBody.NoMemo
    } else {
        iso.NoMemo = iso.NoMemo // gunakan nilai yang ada dari database
    }

	if requestBody.Perihal != "" {
		iso.Perihal = requestBody.Perihal
	} else {
		iso.Perihal = iso.Perihal // gunakan nilai yang ada dari database
	}
	
	if requestBody.Pic != "" {
		iso.Pic = requestBody.Pic
	} else {
		iso.Pic = iso.Pic // gunakan nilai yang ada dari database
	}

    initializers.DB.Model(&iso).Updates(iso)

    c.JSON(200, gin.H{
        "Iso": &iso,
    })
}

func IsoDelete(c *gin.Context) {

	//get id
	id := c.Params.ByName("id")

	// find the iso
	var iso models.Iso

	if err := initializers.DB.First(&iso, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Iso not found"})
		return
	}

	/// delete it
	if err := initializers.DB.Delete(&iso).Error; err != nil {
		c.JSON(404, gin.H{"error": "Iso Failed to Delete"})
		return
	}

	c.JSON(200, gin.H{
		"Iso": "Iso deleted",
	})
}
