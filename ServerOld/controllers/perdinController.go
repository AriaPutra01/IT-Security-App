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

type perdinRequest struct {
	NoPerdin  string `json:"no_perdin"`
	Tanggal   string `json:"tanggal"`
	Hotel     string `json:"hotel"`
	Transport string `json:"transport"`
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
		Tanggal:   tanggal,
		Hotel:     requestBody.Hotel,
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

func CreateExcelPerdin(c *gin.Context) {
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
		if sheetName == "PERDIN" {
			f.NewSheet(sheetName)
			f.SetCellValue(sheetName, "A1", "No Perdin")
			f.SetCellValue(sheetName, "B1", "Tanggal")
			f.SetCellValue(sheetName, "C1", "Hotel")
			f.SetCellValue(sheetName, "D1", "Transport")
		} else {
			f.NewSheet(sheetName)
		}
	}

	// Fetch initial data from the database
	var perdins []models.Perdin
	initializers.DB.Find(&perdins)

	// Write initial data to the "SAG" sheet
	perdinSheetName := "PERDIN"
	for i, perdin := range perdins {
		tanggalString := perdin.Tanggal.Format("2006-01-02")
		rowNum := i + 2 // Start from the second row (first row is header)
		f.SetCellValue(perdinSheetName, fmt.Sprintf("A%d", rowNum), perdin.NoPerdin)
		f.SetCellValue(perdinSheetName, fmt.Sprintf("B%d", rowNum), tanggalString)
		f.SetCellValue(perdinSheetName, fmt.Sprintf("C%d", rowNum), perdin.Hotel)
		f.SetCellValue(perdinSheetName, fmt.Sprintf("D%d", rowNum), perdin.Transport)
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

func UpdateSheetPerdin(c *gin.Context) {
	dir := "D:\\excel"
	fileName := "its_report.xlsx"
	filePath := filepath.Join(dir, fileName)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.String(http.StatusBadRequest, "File tidak ada")
		return
	}

	f, err := excelize.OpenFile(filePath)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error membuka file: %v", err)
		return
	}
	defer f.Close()

	sheetName := "PERDIN"

	if _, err := f.GetSheetIndex(sheetName); err == nil {
		f.DeleteSheet(sheetName)
	}
	f.NewSheet(sheetName)

	f.SetCellValue(sheetName, "A1", "No Perdin")
	f.SetCellValue(sheetName, "B1", "Tanggal")
	f.SetCellValue(sheetName, "C1", "Hotel")
	f.SetCellValue(sheetName, "D1", "Transport")

	var perdins []models.Perdin
	initializers.DB.Find(&perdins)

	for i, perdin := range perdins {
		rowNum := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), perdin.NoPerdin)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), perdin.Tanggal.Format("2006-01-02"))
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), perdin.Hotel)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), perdin.Transport)
	}

	if err := f.SaveAs(filePath); err != nil {
		c.String(http.StatusInternalServerError, "Error menyimpan file: %v", err)
		return
	}

	c.Redirect(http.StatusFound, "/Perdin")
}

func ImportExcelPerdin(c *gin.Context) {
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
	sheetName := "PERDIN"
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
		noPerdin := row[0]
		tanggalString := row[1]
		hotel := row[2]
		transport := row[3]

		// Parse tanggal
		tanggal, err := time.Parse("2006-01-02", tanggalString)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid date format in row %d: %v", i+1, err)
			return
		}

		perdin := models.Perdin{
			Tanggal:   tanggal,
			NoPerdin:  noPerdin,
			Hotel:     hotel,
			Transport: transport,
		}

		// Simpan ke database
		if err := initializers.DB.Create(&perdin).Error; err != nil {
			log.Printf("Error saving record from row %d: %v", i+1, err)
			c.String(http.StatusInternalServerError, "Error saving record from row %d: %v", i+1, err)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data imported successfully."})
}
