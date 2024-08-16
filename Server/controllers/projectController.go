package controllers

import (
	"project-gin/initializers"
	"project-gin/models"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type ProjectRequest struct {
	KodeProject     string `json:"kode_project"`
	JenisPengadaan  string `json:"jenis_pengadaan"`
	NamaPengadaan   string `json:"nama_pengadaan"`
	DivInisiasi     string `json:"div_inisiasi"`
	Bulan           string `json:"bulan"`
	SumberPendanaan string `json:"sumber_pendanaan"`
	Anggaran        int64  `json:"anggaran"`
	NoIzin          string `json:"no_izin"`
	TanggalIzin     string `json:"tanggal_izin"`
	TanggalTor      string `json:"tanggal_tor"`
	Pic             string `json:"pic"`
}

func ProjectCreate(c *gin.Context) {

	// Get data off req body
	var requestBody ProjectRequest

	if err := c.BindJSON(&requestBody); err != nil {
		c.Status(400)
		c.Error(err) // log the error
		return
	}

	// Add some logging to see what's being received
	log.Println("Received request body:", requestBody)

	// Parse the date string
	bulanString := requestBody.Bulan
	bulan, err := time.Parse("01-02-2006", bulanString)
	if err != nil {
		c.Status(400)
		c.JSON(400, gin.H{"error": "Format tanggal tidak valid: " + err.Error()})
		return
	}
	// Parse the date string
	izinString := requestBody.TanggalIzin
	tanggal_izin, err := time.Parse("02-01-2006", izinString)
	if err != nil {
		c.Status(400)
		c.JSON(400, gin.H{"error": "Format tanggal tidak valid: " + err.Error()})
		return
	}
	// Parse the date string
	torString := requestBody.TanggalTor
	tanggal_tor, err := time.Parse("01-02-2006", torString)
	if err != nil {
		c.Status(400)
		c.JSON(400, gin.H{"error": "Format tanggal tidak valid: " + err.Error()})
		return
	}

	project := models.Project{
		KodeProject:     requestBody.KodeProject,
		JenisPengadaan:  requestBody.JenisPengadaan,
		NamaPengadaan:   requestBody.NamaPengadaan,
		DivInisiasi:     requestBody.DivInisiasi,
		Bulan:           bulan,
		SumberPendanaan: requestBody.SumberPendanaan,
		Anggaran:        requestBody.Anggaran,
		NoIzin:          requestBody.NoIzin,
		TanggalIzin:     tanggal_izin,
		TanggalTor:      tanggal_tor,
		Pic:             requestBody.Pic,
	}

	result := initializers.DB.Create(&project)

	if result.Error != nil {
		c.Status(400)
		return
	}

	// Return it
	c.JSON(200, gin.H{
		"Project": project,
	})

}

func ProjectIndex(c *gin.Context) {

	// Get models from DB
	var project []models.Project
	initializers.DB.Find(&project)

	//Respond with them
	c.JSON(200, gin.H{
		"Project": project,
	})
}

func ProjectShow(c *gin.Context) {

	//get id
	id := c.Params.ByName("id")
	// Get models from DB
	var project models.Project

	initializers.DB.First(&project, id)

	//Respond with them
	c.JSON(200, gin.H{
		"Project": project,
	})
}

func ProjectUpdate(c *gin.Context) {

	var requestBody ProjectRequest

	if err := c.BindJSON(&requestBody); err != nil {
		c.Status(400)
		c.Error(err) // log the error
		return
	}

	// Parse the date string
	bulanString := requestBody.Bulan
	bulan, err := time.Parse("01-02-2006", bulanString)
	if err != nil {
		c.Status(400)
		c.JSON(400, gin.H{"error": "Format tanggal tidak valid: " + err.Error()})
		return
	}
	// Parse the date string
	izinString := requestBody.TanggalIzin
	tanggal_izin, err := time.Parse("02-01-2006", izinString)
	if err != nil {
		c.Status(400)
		c.JSON(400, gin.H{"error": "Format tanggal tidak valid: " + err.Error()})
		return
	}
	// Parse the date string
	torString := requestBody.TanggalTor
	tanggal_tor, err := time.Parse("01-02-2006", torString)
	if err != nil {
		c.Status(400)
		c.JSON(400, gin.H{"error": "Format tanggal tidak valid: " + err.Error()})
		return
	}

	//get id
	id := c.Params.ByName("id")

	// find the project
	var project models.Project
	initializers.DB.First(&project, id)

	if err := initializers.DB.First(&project, id); err.Error != nil {
		c.JSON(404, gin.H{"error": "project not found"})
		return
	}

	// update it
	initializers.DB.Model(&project).Updates(models.Project{
		KodeProject:     requestBody.KodeProject,
		JenisPengadaan:  requestBody.JenisPengadaan,
		NamaPengadaan:   requestBody.NamaPengadaan,
		DivInisiasi:     requestBody.DivInisiasi,
		Bulan:           bulan,
		SumberPendanaan: requestBody.SumberPendanaan,
		Anggaran:        requestBody.Anggaran,
		NoIzin:          requestBody.NoIzin,
		TanggalIzin:     tanggal_izin,
		TanggalTor:      tanggal_tor,
		Pic:             requestBody.Pic,
	})

	//Respond with them
	c.JSON(200, gin.H{
		"Project": project,
	})

}

func ProjectDelete(c *gin.Context) {

	//get id
	id := c.Params.ByName("id")

	// find the project
	var project models.Project

	if err := initializers.DB.First(&project, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Project not found"})
		return
	}

	/// delete it
	if err := initializers.DB.Delete(&project).Error; err != nil {
		c.JSON(404, gin.H{"error": "Project Failed to Delete"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Deleted",
	})
}
