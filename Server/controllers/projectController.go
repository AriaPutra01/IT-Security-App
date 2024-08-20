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
	Anggaran        string  `json:"anggaran"`
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
	bulan, err := time.Parse("2006-01-02", requestBody.Bulan)
    if err != nil {
        log.Printf("Error parsing date: %v", err)
        c.JSON(400, gin.H{"error": "Invalid date format: " + err.Error()})
        return
    }

	tanggal_izin, err := time.Parse("2006-01-02", requestBody.TanggalIzin)
    if err != nil {
        log.Printf("Error parsing date: %v", err)
        c.JSON(400, gin.H{"error": "Invalid date format: " + err.Error()})
        return
    }

	tanggal_tor, err := time.Parse("2006-01-02", requestBody.TanggalTor)
    if err != nil {
        log.Printf("Error parsing date: %v", err)
        c.JSON(400, gin.H{"error": "Invalid date format: " + err.Error()})
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

	id := c.Params.ByName("id")

	var project models.Project
	initializers.DB.First(&project, id)

	if err := initializers.DB.First(&project, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "project tidak ditemukan"})
		return
	}

	if requestBody.Bulan != "" {
		bulan, err := time.Parse("2006-01-02", requestBody.Bulan)
		if err != nil {
			c.JSON(400, gin.H{"error": "Format tanggal tidak valid: " + err.Error()})
			return
		}
		project.TanggalIzin = bulan
	}

	if requestBody.TanggalIzin != "" {
		tanggal_izin, err := time.Parse("2006-01-02", requestBody.TanggalIzin)
		if err != nil {
			c.JSON(400, gin.H{"error": "Format tanggal tidak valid: " + err.Error()})
			return
		}
		project.TanggalIzin = tanggal_izin
	}

	if requestBody.TanggalTor != "" {
		tanggal_tor, err := time.Parse("2006-01-02", requestBody.TanggalTor)
		if err != nil {
			c.JSON(400, gin.H{"error": "Format tanggal tidak valid: " + err.Error()})
			return
		}
		project.TanggalIzin = tanggal_tor
	}

	if requestBody.KodeProject != "" {
		project.KodeProject = requestBody.KodeProject
	} else {
		project.KodeProject = project.KodeProject // gunakan nilai yang ada dari database
	}

	if requestBody.JenisPengadaan != "" {
		project.JenisPengadaan = requestBody.JenisPengadaan
	} else {
		project.JenisPengadaan = project.JenisPengadaan // gunakan nilai yang ada dari database
	}

	if requestBody.NamaPengadaan != "" {
		project.NamaPengadaan = requestBody.NamaPengadaan
	} else {
		project.NamaPengadaan = project.NamaPengadaan // gunakan nilai yang ada dari database
	}

	if requestBody.DivInisiasi != "" {
		project.DivInisiasi = requestBody.DivInisiasi
	} else {
		project.DivInisiasi = project.DivInisiasi // gunakan nilai yang ada dari database
	}

	if requestBody.SumberPendanaan != "" {
		project.SumberPendanaan = requestBody.SumberPendanaan
	} else {
		project.SumberPendanaan = project.SumberPendanaan // gunakan nilai yang ada dari database
	}

	if requestBody.Anggaran != "" {
		project.Anggaran = requestBody.Anggaran
	} else {
		project.Anggaran = project.Anggaran // gunakan nilai yang ada dari database
	}

	if requestBody.NoIzin != "" {
		project.NoIzin = requestBody.NoIzin
	} else {
		project.NoIzin = project.NoIzin // gunakan nilai yang ada dari database
	}

	if requestBody.Pic != "" {
		project.Pic = requestBody.Pic
	} else {
		project.Pic = project.Pic // gunakan nilai yang ada dari database
	}

	initializers.DB.Model(&project).Updates(project)

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
