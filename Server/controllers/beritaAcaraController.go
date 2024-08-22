package controllers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"project-its/initializers"
	"project-its/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
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

func CreateExcelBerita(c *gin.Context) {
	dir := "D:\\excel"
	baseFileName := "its_report"
	filePath := filepath.Join(dir, baseFileName+".xlsx")

	// Check if the file already exists
	if _, err := os.Stat(filePath); err == nil {
		// File exists, append "_new" to the file name
		baseFileName += "_new"
	}

	fileName := baseFileName + ".xlsx"

	// File does not exist, create a new file
	f := excelize.NewFile()

	// Define sheet names
	sheetNames := []string{"SAG", "MEMO", "ISO", "SURAT", "BERITA ACARA", "SK", "PROJECT", "PERDIN", "SURAT MASUK", "SURAT KELUAR"}

	// Create sheets and set headers for "SAG" only
	for _, sheetName := range sheetNames {
		if sheetName == "BERITA ACARA" {
			f.NewSheet(sheetName)
			f.SetCellValue(sheetName, "A1", "Tanggal")
			f.SetCellValue(sheetName, "B1", "No Surat")
			f.SetCellValue(sheetName, "C1", "Perihal")
			f.SetCellValue(sheetName, "D1", "Pic")
		} else {
			f.NewSheet(sheetName)
		}
	}

	// Fetch initial data from the database
	var beritas []models.BeritaAcara
	initializers.DB.Find(&beritas)

	// Write initial data to the "SAG" sheet
	beritasSheetName := "BERITA ACARA"
	for i, berita := range beritas {
		tanggalString := berita.Tanggal.Format("2006-01-02")
		rowNum := i + 2 // Start from the second row (first row is header)
		f.SetCellValue(beritasSheetName, fmt.Sprintf("A%d", rowNum), tanggalString)
		f.SetCellValue(beritasSheetName, fmt.Sprintf("B%d", rowNum), berita.NoSurat)
		f.SetCellValue(beritasSheetName, fmt.Sprintf("C%d", rowNum), berita.Perihal)
		f.SetCellValue(beritasSheetName, fmt.Sprintf("D%d", rowNum), berita.Pic)
	}

	// Delete the default "Sheet1" sheet
	if err := f.DeleteSheet("Sheet1"); err != nil {
		panic(err) // Handle error jika bukan error "sheet tidak ditemukan"
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

func UpdateSheetBerita(c *gin.Context) {
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
	sheetName := "BERITA ACARA"

	// Check if sheet exists and delete it if it does
	if _, err := f.GetSheetIndex(sheetName); err == nil {
		f.DeleteSheet(sheetName)
	}
	f.NewSheet(sheetName)

	// Write header row
	f.SetCellValue(sheetName, "A1", "Tanggal")
	f.SetCellValue(sheetName, "B1", "No Surat")
	f.SetCellValue(sheetName, "C1", "Perihal")
	f.SetCellValue(sheetName, "D1", "PIC")

	// Fetch updated data from the database
	var beritas []models.BeritaAcara
	initializers.DB.Find(&beritas)

	// Write data rows
	for i, berita := range beritas {
		rowNum := i + 2 // Start from the second row (first row is header)
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), berita.Tanggal.Format("2006-01-02"))
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), berita.NoSurat)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), berita.Perihal)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), berita.Pic)
	}

	// Save the file with updated data
	if err := f.SaveAs(filePath); err != nil {
		c.String(http.StatusInternalServerError, "Error menyimpan file: %v", err)
		return
	}
}

func ImportExcelBerita(c *gin.Context) {
	// Mengambil file dari form upload
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "Error retrieving the file: %v", err)
		return
	}
	defer file.Close()

	// Simpan file sementara jika perlu
	tempFile, err := os.CreateTemp("", "*.xlsx")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error creating temporary file: %v", err)
		return
	}
	defer os.Remove(tempFile.Name()) // Hapus file sementara setelah selesai

	// Salin file dari request ke file sementara
	if _, err := file.Seek(0, 0); err != nil {
		c.String(http.StatusInternalServerError, "Error seeking file: %v", err)
		return
	}
	if _, err := io.Copy(tempFile, file); err != nil {
		c.String(http.StatusInternalServerError, "Error copying file: %v", err)
		return
	}

	// Buka file Excel dari file sementara
	tempFile.Seek(0, 0) // Reset pointer ke awal file
	f, err := excelize.OpenFile(tempFile.Name())
	if err != nil {
		c.String(http.StatusInternalServerError, "Error opening file: %v", err)
		return
	}
	defer f.Close()

	// Pilih sheet
	sheetName := "BERITA ACARA"
	rows, err := f.GetRows(sheetName)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error getting rows: %v", err)
		return
	}

	// Loop melalui baris dan simpan ke database
	for i, row := range rows {
		if i == 0 {
			// Lewati header baris jika ada
			continue
		}
		if len(row) < 4 {
			// Pastikan ada cukup kolom
			continue
		}
		tanggalString := row[0]
		noSurat := row[1]
		perihal := row[2]
		pic := row[3]

		// Parse tanggal
		tanggal, err := time.Parse("2006-01-02", tanggalString)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid date format in row %d: %v", i+1, err)
			return
		}

		berita := models.BeritaAcara{
			Tanggal: tanggal,
			NoSurat: noSurat,
			Perihal: perihal,
			Pic:     pic,
		}

		// Simpan ke database
		if err := initializers.DB.Create(&berita).Error; err != nil {
			log.Printf("Error saving record from row %d: %v", i+1, err)
			c.String(http.StatusInternalServerError, "Error saving record from row %d: %v", i+1, err)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data imported successfully."})
}