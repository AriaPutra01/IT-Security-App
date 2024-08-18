package controllers

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"project-gin/initializers"
	"project-gin/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
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

	t, err := template.ParseFiles("views/project.html")
	if err != nil {
		log.Fatal(err)
	}

	err = t.Execute(c.Writer, gin.H{
		"Project": project,
	})
	if err != nil {
		log.Fatal(err)
	}

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

	if requestBody.Anggaran != 0 {
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

func CreateExcelProject(c *gin.Context) {
	dir := "D:\\excel"
	baseFileName := "its_report"
	filePath := filepath.Join(dir, baseFileName+".xlsx")

	// Check if the file already exists
	if _, err := os.Stat(filePath); err == nil {
		// File exists, append "_new" to the file name
		baseFileName += "_new"
	}

	fileName := baseFileName + ".xlsx"

	// Create a new Excel file
	f := excelize.NewFile()

	// Define sheet names
	sheetNames := []string{"SAG", "MEMO", "ISO", "SURAT", "BERITA ACARA", "SK", "PROJECT", "PERDIN", "SURAT MASUK", "SURAT KELUAR"}

	// Create sheets and set headers
	for _, sheetName := range sheetNames {
		f.NewSheet(sheetName)
		if sheetName == "PROJECT" {
			f.SetCellValue(sheetName, "A1", "Kode Project")
			f.SetCellValue(sheetName, "B1", "Jenis Pengadaan")
			f.SetCellValue(sheetName, "C1", "Nama Pengadaan")
			f.SetCellValue(sheetName, "D1", "Divisi Inisiasi")
			f.SetCellValue(sheetName, "E1", "Bulan")
			f.SetCellValue(sheetName, "F1", "Sumber Pendanaan")
			f.SetCellValue(sheetName, "G1", "Anggaran")
			f.SetCellValue(sheetName, "H1", "No Izin")
			f.SetCellValue(sheetName, "I1", "Tgl Izin")
			f.SetCellValue(sheetName, "J1", "Tgl TOR")
			f.SetCellValue(sheetName, "K1", "Pic")
		}
	}

	// Fetch initial data from the database
	var projects []models.Project
	initializers.DB.Find(&projects)

	// Write initial data to the "PROJECT" sheet
	projectSheetName := "PROJECT"
	for i, project := range projects {
		izinString := project.TanggalIzin.Format("02-01-2006")
		torString := project.TanggalTor.Format("02-01-2006")
		bulanString := project.Bulan.Format("02-01-2006")
		rowNum := i + 2 // Start from the second row (first row is header)

		// Ensure data is correctly written to cells
		f.SetCellValue(projectSheetName, fmt.Sprintf("A%d", rowNum), project.KodeProject)
		f.SetCellValue(projectSheetName, fmt.Sprintf("B%d", rowNum), project.JenisPengadaan)
		f.SetCellValue(projectSheetName, fmt.Sprintf("C%d", rowNum), project.NamaPengadaan)
		f.SetCellValue(projectSheetName, fmt.Sprintf("D%d", rowNum), project.DivInisiasi)
		f.SetCellValue(projectSheetName, fmt.Sprintf("E%d", rowNum), bulanString) // Ensure this is the correct format
		f.SetCellValue(projectSheetName, fmt.Sprintf("F%d", rowNum), project.SumberPendanaan)
		f.SetCellValue(projectSheetName, fmt.Sprintf("G%d", rowNum), project.Anggaran)
		f.SetCellValue(projectSheetName, fmt.Sprintf("H%d", rowNum), project.NoIzin)
		f.SetCellValue(projectSheetName, fmt.Sprintf("I%d", rowNum), izinString)
		f.SetCellValue(projectSheetName, fmt.Sprintf("J%d", rowNum), torString)
		f.SetCellValue(projectSheetName, fmt.Sprintf("K%d", rowNum), project.Pic)
	}

	// Delete the default "Sheet1" sheet if it exists
	if err := f.DeleteSheet("Sheet1"); err != nil {
		c.String(http.StatusInternalServerError, "Error deleting default sheet: %v", err)
		return
	}

	// Save the newly created file
	buf, err := f.WriteToBuffer()
	if err != nil {
		c.String(http.StatusInternalServerError, "Error saving file: %v", err)
		return
	}

	// Serve the file to the client
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Writer.Write(buf.Bytes())
}

func UpdateSheetProject(c *gin.Context) {
    dir := "D:\\excel"
    fileName := "its_report.xlsx"
    filePath := filepath.Join(dir, fileName)

    // Check if the file exists
    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        c.String(http.StatusBadRequest, "File tidak ada")
        return
    }

    // Open the existing Excel file
    f, err := excelize.OpenFile(filePath)
    if err != nil {
        c.String(http.StatusInternalServerError, "Error membuka file: %v", err)
        return
    }
    defer f.Close()

    // Define sheet name
    sheetName := "PROJECT"

    // Check if sheet exists and delete it if it does
    if _, err := f.GetSheetIndex(sheetName); err == nil {
        f.DeleteSheet(sheetName)
    }
    f.NewSheet(sheetName)

    // Write header row
    f.SetCellValue(sheetName, "A1", "Kode Project")
    f.SetCellValue(sheetName, "B1", "Jenis Pengadaan")
    f.SetCellValue(sheetName, "C1", "Nama Pengadaan")
    f.SetCellValue(sheetName, "D1", "Divisi Inisiasi")
    f.SetCellValue(sheetName, "E1", "Bulan")
    f.SetCellValue(sheetName, "F1", "Sumber Pendanaan")
    f.SetCellValue(sheetName, "G1", "Anggaran")
    f.SetCellValue(sheetName, "H1", "No Izin")
    f.SetCellValue(sheetName, "I1", "Tgl Izin")
    f.SetCellValue(sheetName, "J1", "Tgl TOR")
    f.SetCellValue(sheetName, "K1", "Pic")

    // Fetch updated data from the database
    var projects []models.Project
    initializers.DB.Find(&projects)

    // Write data rows
    for i, project := range projects {
        rowNum := i + 2 // Start from the second row (first row is header)
        
        // Convert date to string with specific format
        bulanString := project.Bulan.Format("02-01-2006")
        
        f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), project.KodeProject)
        f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), project.JenisPengadaan)
        f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), project.NamaPengadaan)
        f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), project.DivInisiasi)
        f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), bulanString) // Write month as text
        f.SetCellValue(sheetName, fmt.Sprintf("F%d", rowNum), project.SumberPendanaan)
        f.SetCellValue(sheetName, fmt.Sprintf("G%d", rowNum), project.Anggaran)
        f.SetCellValue(sheetName, fmt.Sprintf("H%d", rowNum), project.NoIzin)
        f.SetCellValue(sheetName, fmt.Sprintf("I%d", rowNum), project.TanggalIzin.Format("02-01-2006"))
        f.SetCellValue(sheetName, fmt.Sprintf("J%d", rowNum), project.TanggalTor.Format("02-01-2006"))
        f.SetCellValue(sheetName, fmt.Sprintf("K%d", rowNum), project.Pic)

    }

    // Save the file with updated data
    file, err := os.OpenFile(filePath, os.O_RDWR, 0755)
    if err != nil {
        c.String(http.StatusInternalServerError, "Error membuka file: %v", err)
        return
    }
    defer file.Close()

    if _, err := f.WriteTo(file); err != nil {
        c.String(http.StatusInternalServerError, "Error menyimpan file: %v", err)
        return
    }

    c.Redirect(http.StatusFound, "/Project")
}

func ImportExcelProject(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "Error retrieving the file: %v", err)
		return
	}
	defer file.Close()

	tempFile, err := os.CreateTemp("", "*.xlsx")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error creating temporary file: %v", err)
		return
	}
	defer os.Remove(tempFile.Name())

	if _, err := io.Copy(tempFile, file); err != nil {
		c.String(http.StatusInternalServerError, "Error copying file: %v", err)
		return
	}

	tempFile.Seek(0, 0)
	f, err := excelize.OpenFile(tempFile.Name())
	if err != nil {
		c.String(http.StatusInternalServerError, "Error opening file: %v", err)
		return
	}
	defer f.Close()

	sheetName := "PROJECT"
	rows, err := f.GetRows(sheetName)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error getting rows: %v", err)
		return
	}

	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 11 {
			continue
		}

		kodeProject := row[0]
		jenisPengadaan := row[1]
		namaPengadaan := row[2]
		divInisiasi := row[3]
		bulanString := row[4]
		sumberPendanaan := row[5]
		anggaran := row[6]
		noIzin := row[7]
		tanggalIzinString := row[8]
		tanggalTorString := row[9]
		pic := row[10]

		bulan, err := time.Parse("2006-01-02", bulanString)
		if err != nil {
			c.Status(400)
			c.JSON(400, gin.H{"error": "Format tanggal tidak valid: " + err.Error()})
			return
		}

		tanggalIzin, err := time.Parse("02-01-2006", tanggalIzinString)
		if err != nil {
			c.Status(400)
			c.JSON(400, gin.H{"error": "Format tanggal tidak valid: " + err.Error()})
			return
		}

		tanggalTor, err := time.Parse("02-01-2006", tanggalTorString)
		if err != nil {
			c.Status(400)
			c.JSON(400, gin.H{"error": "Format tanggal tidak valid: " + err.Error()})
			return
		}

		anggaranInt, err := strconv.Atoi(anggaran)
		if err != nil {
			c.Status(400)
			c.JSON(400, gin.H{"error": "Format anggaran tidak valid: " + err.Error()})
			return
		}

		project := models.Project{
			KodeProject:     kodeProject,
			JenisPengadaan:  jenisPengadaan,
			NamaPengadaan:   namaPengadaan,
			DivInisiasi:     divInisiasi,
			Bulan:           bulan,
			SumberPendanaan: sumberPendanaan,
			Anggaran:        int64(anggaranInt),
			NoIzin:          noIzin,
			TanggalIzin:     tanggalIzin,
			TanggalTor:      tanggalTor,
			Pic:             pic,
		}

		if err := initializers.DB.Create(&project).Error; err != nil {
			log.Printf("Error saving record from row %d: %v", i+1, err)
			c.String(http.StatusInternalServerError, "Error saving record from row %d: %v", i+1, err)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data imported successfully."})
}
