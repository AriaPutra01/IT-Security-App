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

type SuratMasukRequest struct {
	NoSurat    string `json:"no_surat"`
	Title      string `json:"title"`
	RelatedDiv string `json:"related_div"`
	DestinyDiv string `json:"destiny_div"`
	Tanggal    string `json:"tanggal"`
	CreateBy string `json:"create_by"`
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

	requestBody.CreateBy = c.MustGet("username").(string)

	// Parse the date string
	tanggal, err := time.Parse("2006-01-02", requestBody.Tanggal)
	if err != nil {
		log.Printf("Error parsing date: %v", err)
		c.JSON(400, gin.H{"error": "Invalid date format: " + err.Error()})
		return
	}

	surat_masuk := models.SuratMasuk{
		NoSurat:    requestBody.NoSurat,
		Title:      requestBody.Title,
		RelatedDiv: requestBody.RelatedDiv,
		DestinyDiv: requestBody.DestinyDiv,
		Tanggal:    tanggal,
		CreateBy: requestBody.CreateBy,
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

	requestBody.CreateBy = c.MustGet("username").(string)

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

	if requestBody.CreateBy != "" {
		surat_masuk.CreateBy = requestBody.CreateBy
	} else {
		surat_masuk.CreateBy = surat_masuk.CreateBy // gunakan nilai yang ada dari database
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

func CreateExcelSuratMasuk(c *gin.Context) {
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
		if sheetName == "SURAT MASUK" {
			f.NewSheet(sheetName)
			f.SetCellValue(sheetName, "A1", "No Surat")
			f.SetCellValue(sheetName, "B1", "Title")
			f.SetCellValue(sheetName, "C1", "Related Divisi")
			f.SetCellValue(sheetName, "D1", "Destiny Divisi")
			f.SetCellValue(sheetName, "E1", "Tanggal")
		} else {
			f.NewSheet(sheetName)
		}
	}

	// Fetch initial data from the database
	var surat_masuks []models.SuratMasuk
	initializers.DB.Find(&surat_masuks)

	// Write initial data to the "SAG" sheet
	surat_masukSheetName := "SURAT MASUK"
	for i, surat_masuk := range surat_masuks {
		tanggalString := surat_masuk.Tanggal.Format("2006-01-02")
		rowNum := i + 2 // Start from the second row (first row is header)
		f.SetCellValue(surat_masukSheetName, fmt.Sprintf("A%d", rowNum), surat_masuk.NoSurat)
		f.SetCellValue(surat_masukSheetName, fmt.Sprintf("B%d", rowNum), surat_masuk.Title)
		f.SetCellValue(surat_masukSheetName, fmt.Sprintf("C%d", rowNum), surat_masuk.RelatedDiv)
		f.SetCellValue(surat_masukSheetName, fmt.Sprintf("D%d", rowNum), surat_masuk.DestinyDiv)
		f.SetCellValue(surat_masukSheetName, fmt.Sprintf("E%d", rowNum), tanggalString)
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

func UpdateSheetSuratMasuk(c *gin.Context) {
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
	sheetName := "SURAT MASUK"

	// Check if sheet exists and delete it if it does
	if _, err := f.GetSheetIndex(sheetName); err == nil {
		f.DeleteSheet(sheetName)
	}
	f.NewSheet(sheetName)

	// Write header row
	f.SetCellValue(sheetName, "A1", "No Surat")
	f.SetCellValue(sheetName, "B1", "Title")
	f.SetCellValue(sheetName, "C1", "Related Divisi")
	f.SetCellValue(sheetName, "D1", "Destiny Divisi")
	f.SetCellValue(sheetName, "E1", "Tanggal")

	// Fetch updated data from the database
	var surat_masuks []models.SuratMasuk
	initializers.DB.Find(&surat_masuks)

	// Write data rows
	for i, surat_masuk := range surat_masuks {
		rowNum := i + 2 // Start from the second row (first row is header)
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), surat_masuk.NoSurat)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), surat_masuk.Title)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), surat_masuk.RelatedDiv)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), surat_masuk.DestinyDiv)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), surat_masuk.Tanggal.Format("02-01-2006"))
	}

	if err := f.SaveAs(filePath); err != nil {
		c.String(http.StatusInternalServerError, "Error saving file: %v", err)
		return
	}

}

func ImportExcelSuratMasuk(c *gin.Context) {
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
	sheetName := "SURAT MASUK"
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
		related_div := row[2]
		destiny_div := row[3]
		tanggalString := row[4]

		// Parse tanggal
		tanggal, err := time.Parse("2006-01-02", tanggalString)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid date format in row %d: %v", i+1, err)
			return
		}

		surat_masuk := models.SuratMasuk{
			NoSurat:    noSurat,
			Title:      title,
			RelatedDiv: related_div,
			DestinyDiv: destiny_div,
			Tanggal:    tanggal,
		}

		// Simpan ke database
		if err := initializers.DB.Create(&surat_masuk).Error; err != nil {
			log.Printf("Error saving record from row %d: %v", i+1, err)
			c.String(http.StatusInternalServerError, "Error saving record from row %d: %v", i+1, err)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data imported successfully."})
}
