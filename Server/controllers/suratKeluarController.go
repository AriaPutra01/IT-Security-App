package controllers


import (
	"project-gin/initializers"
	"project-gin/models"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type SuratKeluarRequest struct {
	NoSurat string `json:"no_surat"`
	Title string `json:"title"`
	From string `json:"from"`
	Pic string `json:"pic"`
	Tanggal string `json:"tanggal"`
}

func SuratKeluarCreate(c *gin.Context) {

	// Get data off req body
	var requestBody SuratKeluarRequest

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

	surat_keluar := models.SuratKeluar{
		NoSurat: requestBody.NoSurat,
		Title: requestBody.Title,
		From: requestBody.From,
		Pic: requestBody.Pic,
		Tanggal: tanggal,
	}

	result := initializers.DB.Create(&surat_keluar)

	if result.Error != nil {
		c.Status(400)
		return
	}

	// Return it
	c.JSON(200, gin.H{
		"SuratKeluar": surat_keluar,
	})

}

func SuratKeluarIndex(c *gin.Context) {

	// Get models from DB
	var surat_keluar []models.SuratKeluar
	initializers.DB.Find(&surat_keluar)

	//Respond with them
	c.JSON(200, gin.H{
		"SuratKeluar": surat_keluar,
	})
}

func SuratKeluarShow(c *gin.Context) {

	id := c.Params.ByName("id")
	// Get models from DB
	var surat_keluar models.SuratKeluar

	initializers.DB.First(&surat_keluar, id)

	//Respond with them
	c.JSON(200, gin.H{
		"SuratKeluar": surat_keluar,
	})
}

func SuratKeluarUpdate(c *gin.Context) {

	var requestBody SuratKeluarRequest

	if err := c.BindJSON(&requestBody); err != nil {
		c.Status(400)
		c.Error(err) // log the error
		return
	}

	id := c.Params.ByName("id")

    var surat_keluar models.SuratKeluar
    initializers.DB.First(&surat_keluar, id)

    if err := initializers.DB.First(&surat_keluar, id).Error; err != nil {
        c.JSON(404, gin.H{"error": "surat_keluar tidak ditemukan"})
        return
    }

    if requestBody.Tanggal != "" {
        tanggal, err := time.Parse("2006-01-02", requestBody.Tanggal)
        if err != nil {
            c.JSON(400, gin.H{"error": "Format tanggal tidak valid: " + err.Error()})
            return
        }
        surat_keluar.Tanggal = tanggal
    }

    if requestBody.NoSurat != "" {
        surat_keluar.NoSurat = requestBody.NoSurat
    } else {
        surat_keluar.NoSurat = surat_keluar.NoSurat // gunakan nilai yang ada dari database
    }

	if requestBody.Title != "" {
		surat_keluar.Title = requestBody.Title
	} else {
		surat_keluar.Title = surat_keluar.Title // gunakan nilai yang ada dari database
	}
	
	if requestBody.Pic != "" {
		surat_keluar.Pic = requestBody.Pic
	} else {
		surat_keluar.Pic = surat_keluar.Pic // gunakan nilai yang ada dari database
	}

    initializers.DB.Model(&surat_keluar).Updates(surat_keluar)

    c.JSON(200, gin.H{
        "surat_keluar": surat_keluar,
    })
}

func SuratKeluarDelete(c *gin.Context) {

	//get id
	id := c.Params.ByName("id")

	// find the Surat Keluar
	var surat_keluar models.SuratKeluar

	if err := initializers.DB.First(&surat_keluar, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "surat_keluar not found"})
		return
	}

	/// delete it
	if err := initializers.DB.Delete(&surat_keluar).Error; err != nil {
		c.JSON(404, gin.H{"error": "Surat Keluar Failed to Delete"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Deleted",
	})
}


