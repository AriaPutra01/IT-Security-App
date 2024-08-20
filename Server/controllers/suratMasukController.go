package controllers


import (
	"project-gin/initializers"
	"project-gin/models"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type SuratMasukRequest struct {
	NoSurat string `json:"no_surat"`
	Title string `json:"title"`
	RelatedDiv string `json:"related_div"`
	DestinyDiv string `json:"destiny_div"`
	Tanggal string `json:"tanggal"`
}

func SuratMasukCreate(c *gin.Context) {

	// Get data off req body
	var requestBody SuratMasukRequest

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

	surat_masuk := models.SuratMasuk{
		NoSurat: requestBody.NoSurat,
		Title: requestBody.Title,
		RelatedDiv: requestBody.RelatedDiv,
		DestinyDiv: requestBody.DestinyDiv,
		Tanggal: tanggal,
	}

	result := initializers.DB.Create(&surat_masuk)

	if result.Error != nil {
		c.Status(400)
		return
	}

	// Return it
	c.JSON(200, gin.H{
		"SuratMasuk": surat_masuk,
	})

}

func SuratMasukIndex(c *gin.Context) {

	// Get models from DB
	var surat_masuk []models.SuratMasuk
	initializers.DB.Find(&surat_masuk)

	//Respond with them
	c.JSON(200, gin.H{
		"SuratMasuk": surat_masuk,
	})
}

func SuratMasukShow(c *gin.Context) {

	id := c.Params.ByName("id")
	// Get models from DB
	var surat_masuk models.SuratMasuk

	initializers.DB.First(&surat_masuk, id)

	//Respond with them
	c.JSON(200, gin.H{
		"SuratMasuk": surat_masuk,
	})
}

func SuratMasukUpdate(c *gin.Context) {

	var requestBody SuratMasukRequest

	if err := c.BindJSON(&requestBody); err != nil {
		c.Status(400)
		c.Error(err) // log the error
		return
	}
    id := c.Params.ByName("id")

    var surat_masuk models.SuratMasuk
    initializers.DB.First(&surat_masuk, id)

    if err := initializers.DB.First(&surat_masuk, id).Error; err != nil {
        c.JSON(404, gin.H{"error": "surat_masuk tidak ditemukan"})
        return
    }

    if requestBody.Tanggal != "" {
        tanggal, err := time.Parse("2006-01-02", requestBody.Tanggal)
        if err != nil {
            c.JSON(400, gin.H{"error": "Format tanggal tidak valid: " + err.Error()})
            return
        }
        surat_masuk.Tanggal = tanggal
    }

    if requestBody.NoSurat != "" {
        surat_masuk.NoSurat = requestBody.NoSurat
    } else {
        surat_masuk.NoSurat = surat_masuk.NoSurat // gunakan nilai yang ada dari database
    }

	if requestBody.Title != "" {
		surat_masuk.Title = requestBody.Title
	} else {
		surat_masuk.Title = surat_masuk.Title // gunakan nilai yang ada dari database
	}
	
	if requestBody.RelatedDiv != "" {
		surat_masuk.RelatedDiv = requestBody.RelatedDiv
	} else {
		surat_masuk.RelatedDiv = surat_masuk.RelatedDiv // gunakan nilai yang ada dari database
	}

	if requestBody.DestinyDiv != "" {
		surat_masuk.DestinyDiv = requestBody.DestinyDiv
	} else {
		surat_masuk.DestinyDiv = surat_masuk.DestinyDiv // gunakan nilai yang ada dari database
	}

	

    initializers.DB.Model(&surat_masuk).Updates(surat_masuk)

    c.JSON(200, gin.H{
        "surat_masuk": surat_masuk,
    })

}

func SuratMasukDelete(c *gin.Context) {

	//get id
	id := c.Params.ByName("id")

	// find the SuratMasuk
	var surat_masuk models.SuratMasuk

	if err := initializers.DB.First(&surat_masuk, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "surat masuk not found"})
		return
	}

	/// delete it
	if err := initializers.DB.Delete(&surat_masuk).Error; err != nil {
		c.JSON(404, gin.H{"error": "Surat Masuk Failed to Delete"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Deleted",
	})
}


