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

type skRequest struct {
	Tanggal string `json:"tanggal"`
	NoSurat string `json:"no_surat"`
	Perihal string `json:"perihal"`
	Pic     string `json:"pic"`
}

func SkCreate(c *gin.Context) {

	// Get data off req body
	var requestBody skRequest

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

	sk := models.Sk{
		Tanggal: tanggal,
		NoSurat: requestBody.NoSurat,
		Perihal: requestBody.Perihal,
		Pic:     requestBody.Pic,
	}

	result := initializers.DB.Create(&sk)

	if result.Error != nil {
		c.Status(400)
		return
	}

	// Return it
	c.JSON(200, gin.H{
		"SK": sk,
	})

}

func SkIndex(c *gin.Context) {

	// Get models from DB
	var sk []models.Sk
	initializers.DB.Find(&sk)

	//Respond with them
	c.JSON(200, gin.H{
		"SK": sk,
	})
}

func SkShow(c *gin.Context) {

	id := c.Params.ByName("id")
	// Get models from DB
	var sk models.Sk

	initializers.DB.First(&sk, id)

	//Respond with them
	c.JSON(200, gin.H{
		"SK": sk,
	})
}

func SkUpdate(c *gin.Context) {

	var requestBody skRequest

	if err := c.BindJSON(&requestBody); err != nil {
		c.Status(400)
		c.Error(err) // log the error
		return
	}

	id := c.Params.ByName("id")

	var sk models.Sk
	initializers.DB.First(&sk, id)

	if err := initializers.DB.First(&sk, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "SK tidak ditemukan"})
		return
	}

	if requestBody.Tanggal != "" {
		tanggal, err := time.Parse("2006-01-02", requestBody.Tanggal)
		if err != nil {
			c.JSON(400, gin.H{"error": "Format tanggal tidak valid: " + err.Error()})
			return
		}
		sk.Tanggal = tanggal
	}

	if requestBody.NoSurat != "" {
		sk.NoSurat = requestBody.NoSurat
	} else {
		sk.NoSurat = sk.NoSurat // gunakan nilai yang ada dari database
	}

	if requestBody.Perihal != "" {
		sk.Perihal = requestBody.Perihal
	} else {
		sk.Perihal = sk.Perihal // gunakan nilai yang ada dari database
	}

	if requestBody.Pic != "" {
		sk.Pic = requestBody.Pic
	} else {
		sk.Pic = sk.Pic // gunakan nilai yang ada dari database
	}

	initializers.DB.Model(&sk).Updates(sk)

	c.JSON(200, gin.H{
		"sk": sk,
	})

}

func SkDelete(c *gin.Context) {

	//get id
	id := c.Params.ByName("id")

	// find the sk
	var sk models.Sk

	if err := initializers.DB.First(&sk, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "sk not found"})
		return
	}

	/// delete it
	if err := initializers.DB.Delete(&sk).Error; err != nil {
		c.JSON(404, gin.H{"error": "sk Failed to Delete"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Deleted",
	})
}

func CreateExcelSk(c *gin.Context) {
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
		if sheetName == "SK" {
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
	var sks []models.Sk
	initializers.DB.Find(&sks)

	// Write initial data to the "SAG" sheet
	skSheetName := "SK"
	for i, sk := range sks {
		tanggalString := sk.Tanggal.Format("2006-01-02")
		rowNum := i + 2 // Start from the second row (first row is header)
		f.SetCellValue(skSheetName, fmt.Sprintf("A%d", rowNum), tanggalString)
		f.SetCellValue(skSheetName, fmt.Sprintf("B%d", rowNum), sk.NoSurat)
		f.SetCellValue(skSheetName, fmt.Sprintf("C%d", rowNum), sk.Perihal)
		f.SetCellValue(skSheetName, fmt.Sprintf("D%d", rowNum), sk.Pic)
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

func UpdateSheetSk(c *gin.Context) {
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
	sheetName := "SK"

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
	var sks []models.Sk
	initializers.DB.Find(&sks)

	// Write data rows
	for i, sk := range sks {
		rowNum := i + 2 // Start from the second row (first row is header)
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), sk.Tanggal.Format("2006-01-02"))
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), sk.NoSurat)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), sk.Perihal)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), sk.Pic)
	}

	if err := f.SaveAs(filePath); err != nil {
		c.String(http.StatusInternalServerError, "Error menyimpan file: %v", err)
		return
	}

}

func ImportExcelSk(c *gin.Context) {
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
	sheetName := "SK"
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

		sk := models.Sk{
			Tanggal: tanggal,
			NoSurat: noSurat,
			Perihal: perihal,
			Pic:     pic,
		}

		// Simpan ke database
		if err := initializers.DB.Create(&sk).Error; err != nil {
			log.Printf("Error saving record from row %d: %v", i+1, err)
			c.String(http.StatusInternalServerError, "Error saving record from row %d: %v", i+1, err)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data imported successfully."})
}
