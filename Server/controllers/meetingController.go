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

type MeetingRequest struct {
	ID               uint      `gorm:"primaryKey"`
	Task             string    `json:"task"`
	TindakLanjut     string    `json:"tindak_lanjut"`
	Status           string    `json:"status"`
	UpdatePengerjaan *string    `json:"update_pengerjaan"`
	Pic              string    `json:"pic"`
	TanggalTarget    string    `json:"tanggal_target"`
	TanggalActual    string    `json:"tanggal_actual"`
	CreateBy         string    `json:"create_by"`
	Info             string    `json:"info"`
}
func MeetingIndex(c *gin.Context) {

	var meeting []models.Meeting

	initializers.DB.Find(&meeting)

	c.JSON(200, gin.H{
		"meeting": meeting,
	})

}

func MeetingCreate(c *gin.Context) {

	var requestBody MeetingRequest

	if err := c.BindJSON(&requestBody); err != nil {
		c.Status(400)
		c.Error(err) // log the error
		return
	}

	// Add some logging to see what's being received
	log.Println("Received request body:", requestBody)

	// Parse the date string
	tanggal_target, err := time.Parse("2006-01-02", requestBody.TanggalTarget)
	if err != nil {
		log.Printf("Error parsing date: %v", err)
		c.Status(400)
		c.JSON(400, gin.H{"error": "Invalid date format: " + err.Error()})
		return
	}

	tanggal_actual, err := time.Parse("2006-01-02", requestBody.TanggalActual)
	if err != nil {
		log.Printf("Error parsing date: %v", err)
		c.Status(400)
		c.JSON(400, gin.H{"error": "Invalid date format: " + err.Error()})
		return
	}

	requestBody.CreateBy = c.MustGet("username").(string)

	meeting := models.Meeting{
		Task:             requestBody.Task,
		TindakLanjut:     requestBody.TindakLanjut,
		Status:           requestBody.Status,
		UpdatePengerjaan: requestBody.UpdatePengerjaan,
		Pic:              requestBody.Pic,
		TanggalTarget:    tanggal_target,
		TanggalActual:    tanggal_actual,
		CreateBy:         requestBody.CreateBy,
		Info:             requestBody.CreateBy,
	}

	result := initializers.DB.Create(&meeting)

	if result.Error != nil {
		c.Status(400)
		c.JSON(400, gin.H{"error": "Failed to create Meeting: " + result.Error.Error()})
		return
	}

	c.JSON(201, gin.H{
		"meeting": meeting,
	})

}

func MeetingShow(c *gin.Context) {

	id := c.Params.ByName("id")

	var meeting models.Meeting

	initializers.DB.First(&meeting, id)

	c.JSON(200, gin.H{
		"meeting": meeting,
	})

}

func MeetingUpdate(c *gin.Context) {

	var requestBody MeetingRequest

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	id := c.Params.ByName("id")

	var meeting models.Meeting

	if err := initializers.DB.First(&meeting, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "meeting not found"})
		return
	}

	requestBody.CreateBy = c.MustGet("username").(string)

	if requestBody.TanggalTarget != "" {
		tanggal_target, err := time.Parse("2006-01-02", requestBody.TanggalTarget)
		if err != nil {
			c.JSON(400, gin.H{"error": "Format tanggal tidak valid: " + err.Error()})
			return
		}
		meeting.TanggalTarget = tanggal_target
	}

	if requestBody.TanggalActual != "" {
		tanggal_actual, err := time.Parse("2006-01-02", requestBody.TanggalActual)
		if err != nil {
			c.JSON(400, gin.H{"error": "Format tanggal tidak valid: " + err.Error()})
			return
		}
		meeting.TanggalActual = tanggal_actual
	}

	if requestBody.Task != "" {
		meeting.Task = requestBody.Task
	} else {
		meeting.Task = meeting.Task
	}

	if requestBody.TindakLanjut != "" {
		meeting.TindakLanjut = requestBody.TindakLanjut
	} else {
		meeting.TindakLanjut = meeting.TindakLanjut
	}

	if requestBody.Status != "" {
		meeting.Status = requestBody.Status
	} else {
		meeting.Status = meeting.Status
	}

	if requestBody.UpdatePengerjaan != nil {
		meeting.UpdatePengerjaan = requestBody.UpdatePengerjaan
	} else {
		meeting.UpdatePengerjaan = meeting.UpdatePengerjaan
	}

	if requestBody.Pic != "" {
		meeting.Pic = requestBody.Pic
	} else {
		meeting.Pic = meeting.Pic
	}

	initializers.DB.Save(&meeting)

	c.JSON(200, gin.H{
		"meeting": meeting,
	})
}

func MeetingDelete(c *gin.Context) {

	id := c.Params.ByName("id")

	var meeting models.Meeting

	if err := initializers.DB.First(&meeting, id); err.Error != nil {
		c.JSON(404, gin.H{"error": "meeting not found"})
		return
	}

	if err := initializers.DB.Delete(&meeting).Error; err != nil {
		c.JSON(400, gin.H{"error": "Failed to delete Memo: " + err.Error()})
		return
	}

	c.Status(204)

}

func CreateExcelMeeting(c *gin.Context) {
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
		if sheetName == "MEMO" {
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
	var memos []models.Memo
	initializers.DB.Find(&memos)

	// Write initial data to the "SAG" sheet
	memoSheetName := "MEMO"
	for i, memo := range memos {
		tanggalString := memo.Tanggal.Format("2006-01-02")
		rowNum := i + 2 // Start from the second row (first row is header)
		f.SetCellValue(memoSheetName, fmt.Sprintf("A%d", rowNum), tanggalString)
		f.SetCellValue(memoSheetName, fmt.Sprintf("B%d", rowNum), memo.NoMemo)
		f.SetCellValue(memoSheetName, fmt.Sprintf("C%d", rowNum), memo.Perihal)
		f.SetCellValue(memoSheetName, fmt.Sprintf("D%d", rowNum), memo.Pic)
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

func UpdateSheetMeeting(c *gin.Context) {
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
	sheetName := "MEMO"

	// Check if the sheet exists
	sheetIndex, err := f.GetSheetIndex(sheetName)
	if err != nil || sheetIndex == -1 {
		c.String(http.StatusBadRequest, "Lembar kerja MEMO tidak ditemukan")
		return
	}

	// Clear existing data except the header by deleting rows
	rows, err := f.GetRows(sheetName)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error getting rows: %v", err)
		return
	}
	for i := 1; i < len(rows); i++ { // Start from 1 to keep the header
		f.RemoveRow(sheetName, 2) // Always remove the second row since the sheet compresses up
	}

	// Fetch updated data from the database
	var memos []models.Memo
	initializers.DB.Find(&memos)

	// Write data rows
	for i, memo := range memos {
		rowNum := i + 2 // Start from the second row (first row is header)
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), memo.Tanggal.Format("2006-01-02"))
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), memo.NoMemo)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), memo.Perihal)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), memo.Pic)
	}

	// Save the file with updated data
	if err := f.SaveAs(filePath); err != nil {
		c.String(http.StatusInternalServerError, "Error menyimpan file: %v", err)
		return
	}

	c.Redirect(http.StatusFound, "http://localhost:8000/memo")
}

func ImportExcelMeeting(c *gin.Context) {
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
	sheetName := "MEMO"
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

		memo := models.Memo{
			Tanggal: tanggal,
			NoMemo:  noMemo,
			Perihal: perihal,
			Pic:     pic,
		}

		// Simpan ke database
		if err := initializers.DB.Create(&memo).Error; err != nil {
			log.Printf("Error saving record from row %d: %v", i+1, err)
			c.String(http.StatusInternalServerError, "Error saving record from row %d: %v", i+1, err)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data imported successfully."})
}
