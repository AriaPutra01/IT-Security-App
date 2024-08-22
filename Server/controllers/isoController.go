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

type isoRequest struct {
	Tanggal string `json:"tanggal"`
	NoMemo  string `json:"no_memo"`
	Perihal string `json:"perihal"`
	Pic     string `json:"pic"`
}

func IsoCreate(c *gin.Context) {
	var requestBody isoRequest

	if err := c.BindJSON(&requestBody); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.Status(400)
		c.Error(err) // log the error
		return
	}

	log.Println("Received request body:", requestBody)

	tanggal, err := time.Parse("2006-01-02", requestBody.Tanggal)
	if err != nil {
		log.Printf("Error parsing date: %v", err)
		c.JSON(400, gin.H{"error": "Invalid date format: " + err.Error()})
		return
	}

	iso := models.Iso{
		Tanggal: tanggal,
		NoMemo:  requestBody.NoMemo,
		Perihal: requestBody.Perihal,
		Pic:     requestBody.Pic,
	}

	result := initializers.DB.Create(&iso)
	if result.Error != nil {
		log.Printf("Error saving to database: %v", result.Error)
		c.Status(400)
		return
	}

	log.Println("Iso created successfully:", iso)
	c.JSON(200, gin.H{"Iso": iso})
}

func IsoIndex(c *gin.Context) {

	// Get models from DB
	var iso []models.Iso
	initializers.DB.Find(&iso)

	//Respond with them
	c.JSON(200, gin.H{
		"Iso": iso,
	})
}

func IsoShow(c *gin.Context) {

	id := c.Params.ByName("id")
	// Get models from DB
	var iso models.Iso

	initializers.DB.First(&iso, id)

	//Respond with them
	c.JSON(200, gin.H{
		"Iso": &iso,
	})
}

func IsoUpdate(c *gin.Context) {
	var requestBody isoRequest

	if err := c.BindJSON(&requestBody); err != nil {
		c.Status(400)
		c.Error(err) // log the error
		return
	}

	id := c.Params.ByName("id")

	var iso models.Iso
	initializers.DB.First(&iso, id)

	if err := initializers.DB.First(&iso, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Iso tidak ditemukan"})
		return
	}

	if requestBody.Tanggal != "" {
		tanggal, err := time.Parse("2006-01-02", requestBody.Tanggal)
		if err != nil {
			c.JSON(400, gin.H{"error": "Format tanggal tidak valid: " + err.Error()})
			return
		}
		iso.Tanggal = tanggal
	}

	if requestBody.NoMemo != "" {
		iso.NoMemo = requestBody.NoMemo
	} else {
		iso.NoMemo = iso.NoMemo // gunakan nilai yang ada dari database
	}

	if requestBody.Perihal != "" {
		iso.Perihal = requestBody.Perihal
	} else {
		iso.Perihal = iso.Perihal // gunakan nilai yang ada dari database
	}

	if requestBody.Pic != "" {
		iso.Pic = requestBody.Pic
	} else {
		iso.Pic = iso.Pic // gunakan nilai yang ada dari database
	}

	initializers.DB.Model(&iso).Updates(iso)

	c.JSON(200, gin.H{
		"Iso": &iso,
	})
}

func IsoDelete(c *gin.Context) {

	//get id
	id := c.Params.ByName("id")

	// find the iso
	var iso models.Iso

	if err := initializers.DB.First(&iso, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Iso not found"})
		return
	}

	/// delete it
	if err := initializers.DB.Delete(&iso).Error; err != nil {
		c.JSON(404, gin.H{"error": "Iso Failed to Delete"})
		return
	}

	c.JSON(200, gin.H{
		"Iso": "Iso deleted",
	})
}

func CreateExcelIso(c *gin.Context) {
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
		if sheetName == "ISO" {
			f.NewSheet(sheetName)
			f.SetCellValue(sheetName, "A1", "Tanggal")
			f.SetCellValue(sheetName, "B1", "NoMemo")
			f.SetCellValue(sheetName, "C1", "Perihal")
			f.SetCellValue(sheetName, "D1", "Pic")
		} else {
			f.NewSheet(sheetName)
		}
	}

	// Fetch initial data from the database
	var isos []models.Iso
	initializers.DB.Find(&isos)

	// Write initial data to the "SAG" sheet
	isoSheetName := "ISO"
	for i, iso := range isos {
		tanggalString := iso.Tanggal.Format("2006-01-02")
		rowNum := i + 2 // Start from the second row (first row is header)
		f.SetCellValue(isoSheetName, fmt.Sprintf("A%d", rowNum), tanggalString)
		f.SetCellValue(isoSheetName, fmt.Sprintf("B%d", rowNum), iso.NoMemo)
		f.SetCellValue(isoSheetName, fmt.Sprintf("C%d", rowNum), iso.Perihal)
		f.SetCellValue(isoSheetName, fmt.Sprintf("D%d", rowNum), iso.Pic)
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

func UpdateSheetIso(c *gin.Context) {
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
	sheetName := "ISO"

	// Check if sheet exists and delete it if it does
	if _, err := f.GetSheetIndex(sheetName); err == nil {
		f.DeleteSheet(sheetName)
	}
	f.NewSheet(sheetName)

	// Write header row
	f.SetCellValue(sheetName, "A1", "Tanggal")
	f.SetCellValue(sheetName, "B1", "No Memo")
	f.SetCellValue(sheetName, "C1", "Perihal")
	f.SetCellValue(sheetName, "D1", "PIC")

	// Fetch updated data from the database
	var isos []models.Iso
	initializers.DB.Find(&isos)

	// Write data rows
	for i, iso := range isos {
		rowNum := i + 2 // Start from the second row (first row is header)
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), iso.Tanggal.Format("2006-01-02"))
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), iso.NoMemo)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), iso.Perihal)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), iso.Pic)
	}

	// Save the file with updated data
	if err := f.SaveAs(filePath); err != nil {
		c.String(http.StatusInternalServerError, "Error menyimpan file: %v", err)
		return
	}

	c.Redirect(http.StatusFound, "/iso")
}

func ImportExcelIso(c *gin.Context) {
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
	sheetName := "ISO"
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
		noMemo := row[1]
		perihal := row[2]
		pic := row[3]

		// Parse tanggal
		tanggal, err := time.Parse("2006-01-02", tanggalString)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid date format in row %d: %v", i+1, err)
			return
		}

		iso := models.Iso{
			Tanggal: tanggal,
			NoMemo:  noMemo,
			Perihal: perihal,
			Pic:     pic,
		}

		// Simpan ke database
		if err := initializers.DB.Create(&iso).Error; err != nil {
			log.Printf("Error saving record from row %d: %v", i+1, err)
			c.String(http.StatusInternalServerError, "Error saving record from row %d: %v", i+1, err)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data imported successfully."})
}
