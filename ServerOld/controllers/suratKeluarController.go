package controllers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"project-gin/initializers"
	"project-gin/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type SuratKeluarRequest struct {
	NoSurat string `json:"no_surat"`
	Title   string `json:"title"`
	From    string `json:"from"`
	Pic     string `json:"pic"`
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
		Title:   requestBody.Title,
		From:    requestBody.From,
		Pic:     requestBody.Pic,
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

func CreateExcelSuratKeluar(c *gin.Context) {
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
		if sheetName == "SURAT KELUAR" {
			f.NewSheet(sheetName)
			f.SetCellValue(sheetName, "A1", "No Surat")
			f.SetCellValue(sheetName, "B1", "Title")
			f.SetCellValue(sheetName, "C1", "From")
			f.SetCellValue(sheetName, "D1", "Pic")
			f.SetCellValue(sheetName, "E1", "Tanggal")
		} else {
			f.NewSheet(sheetName)
		}
	}

	// Fetch initial data from the database
	var surat_keluars []models.SuratKeluar
	initializers.DB.Find(&surat_keluars)

	// Write initial data to the "SAG" sheet
	surat_keluarSheetName := "SURAT KELUAR"
	for i, surat_keluar := range surat_keluars {
		tanggalString := surat_keluar.Tanggal.Format("2006-01-02")
		rowNum := i + 2 // Start from the second row (first row is header)
		f.SetCellValue(surat_keluarSheetName, fmt.Sprintf("A%d", rowNum), surat_keluar.NoSurat)
		f.SetCellValue(surat_keluarSheetName, fmt.Sprintf("B%d", rowNum), surat_keluar.Title)
		f.SetCellValue(surat_keluarSheetName, fmt.Sprintf("C%d", rowNum), surat_keluar.From)
		f.SetCellValue(surat_keluarSheetName, fmt.Sprintf("D%d", rowNum), surat_keluar.Pic)
		f.SetCellValue(surat_keluarSheetName, fmt.Sprintf("E%d", rowNum), tanggalString)
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

func UpdateSheetSuratKeluar(c *gin.Context) {
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
	sheetName := "SURAT KELUAR"

	// Check if sheet exists and delete it if it does
	if _, err := f.GetSheetIndex(sheetName); err == nil {
		f.DeleteSheet(sheetName)
	}
	f.NewSheet(sheetName)

	// Write header row
	f.SetCellValue(sheetName, "A1", "No Surat")
	f.SetCellValue(sheetName, "B1", "Title")
	f.SetCellValue(sheetName, "C1", "From")
	f.SetCellValue(sheetName, "D1", "Pic")
	f.SetCellValue(sheetName, "E1", "Tanggal")

	// Fetch updated data from the database
	var surat_keluars []models.SuratKeluar
	initializers.DB.Find(&surat_keluars)

	// Write data rows
	for i, surat_keluar := range surat_keluars {
		rowNum := i + 2 // Start from the second row (first row is header)
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), surat_keluar.NoSurat)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), surat_keluar.Title)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), surat_keluar.From)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), surat_keluar.Pic)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), surat_keluar.Tanggal.Format("2006-01-02"))
	}

	if err := f.SaveAs(filePath); err != nil {
		c.String(http.StatusInternalServerError, "Error saving file: %v", err)
		return
	}

}

func ImportExcelSuratKeluar(c *gin.Context) {
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
	sheetName := "SURAT KELUAR"
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
		noSurat := row[0]
		title := row[1]
		from := row[2]
		pic := row[3]
		tanggalString := row[4]

		// Parse tanggal
		tanggal, err := time.Parse("2006-01-02", tanggalString)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid date format in row %d: %v", i+1, err)
			return
		}

		surat_keluar := models.SuratKeluar{
			NoSurat: noSurat,
			Title:   title,
			From:    from,
			Pic:     pic,
			Tanggal: tanggal,
		}

		// Simpan ke database
		if err := initializers.DB.Create(&surat_keluar).Error; err != nil {
			log.Printf("Error saving record from row %d: %v", i+1, err)
			c.String(http.StatusInternalServerError, "Error saving record from row %d: %v", i+1, err)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data imported successfully."})
}
