package controllers

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"project-its/initializers"
	"project-its/models"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/xuri/excelize/v2"
)

type MeetingRequest struct {
	ID               uint    `gorm:"primaryKey"`
	Task             *string `json:"task"`
	TindakLanjut     *string `json:"tindak_lanjut"`
	Status           *string `json:"status"`
	UpdatePengerjaan *string `json:"update_pengerjaan"`
	Pic              *string `json:"pic"`
	TanggalTarget    *string `json:"tanggal_target"`
	TanggalActual    *string `json:"tanggal_actual"`
	CreateBy         string  `json:"create_by"`
}

func init() {
	err := godotenv.Load() // Memuat file .env
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	accountNameMeeting = os.Getenv("ACCOUNT_NAME")                 // Mengambil nilai dari .env
	accountKeyMeeting = os.Getenv("ACCOUNT_KEY")                   // Mengambil nilai dari .env
	containerNameMeeting = os.Getenv("CONTAINER_NAME_MEETING") // Mengambil nilai dari .env
}

// Tambahkan variabel global untuk menyimpan kredensial
var (
	accountNameMeeting   string
	accountKeyMeeting    string
	containerNameMeeting string
)

func getBlobServiceClientMeeting() azblob.ServiceURL {
	creds, err := azblob.NewSharedKeyCredential(accountNameMeeting, accountKeyMeeting)
	if err != nil {
		panic("Failed to create shared key credential: " + err.Error())
	}

	pipeline := azblob.NewPipeline(creds, azblob.PipelineOptions{})

	// Build the URL for the Azure Blob Storage account
	URL, err := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/", accountNameMeeting))
	if err != nil {
		log.Fatal("Invalid URL format")
	}

	// Create a ServiceURL object that wraps the URL and the pipeline
	serviceURL := azblob.NewServiceURL(*URL, pipeline)

	return serviceURL
}

func UploadHandlerMeeting(c *gin.Context) {
	id := c.PostForm("id") // Mendapatkan ID dari form data
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File diperlukan"})
		return
	}

	// Membuat path berdasarkan ID
	filename := fmt.Sprintf("%s/%s", id, file.Filename)

	// Membuka file
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuka file"})
		return
	}
	defer src.Close()

	// Mengunggah file ke Azure Blob Storage
	containerURL := getBlobServiceClient().NewContainerURL(containerNameMeeting)
	blobURL := containerURL.NewBlockBlobURL(filename)

	_, err = azblob.UploadStreamToBlockBlob(context.TODO(), src, blobURL, azblob.UploadStreamToBlockBlobOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengunggah file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File berhasil diunggah"})
}

func GetFilesByIDMeeting(c *gin.Context) {
	id := c.Param("id") // Mendapatkan ID dari URL

	containerURL := getBlobServiceClient().NewContainerURL(containerNameMeeting)
	prefix := fmt.Sprintf("%s/", id) // Prefix untuk daftar blob di folder tertentu (ID)

	var files []string
	for marker := (azblob.Marker{}); marker.NotDone(); {
		listBlob, err := containerURL.ListBlobsFlatSegment(context.TODO(), marker, azblob.ListBlobsSegmentOptions{
			Prefix: prefix, // Hanya daftar blob dengan prefix yang ditentukan (folder)
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat daftar file"})
			return
		}

		for _, blobInfo := range listBlob.Segment.BlobItems {
			files = append(files, blobInfo.Name)
		}

		marker = listBlob.NextMarker
	}

	c.JSON(http.StatusOK, gin.H{"files": files}) // Pastikan mengembalikan array files
}

// Fungsi untuk menghapus file dari Azure Blob Storage
func DeleteFileHandlerMeeting(c *gin.Context) {
	filename := c.Param("filename")
	id := c.Param("id")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Filename is required"})
		return
	}

	// Membuat path lengkap berdasarkan ID dan nama file
	fullPath := fmt.Sprintf("%s/%s", id, filename)

	containerURL := getBlobServiceClient().NewContainerURL(containerNameMeeting)
	blobURL := containerURL.NewBlockBlobURL(fullPath)

	// Menghapus blob
	_, err := blobURL.Delete(context.TODO(), azblob.DeleteSnapshotsOptionNone, azblob.BlobAccessConditions{})
	if err != nil {
		log.Printf("Error deleting file: %v", err) // Log kesalahan
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"}) // Pastikan ini ada
}

// Fungsi untuk mendownload file dari Azure Blob Storage
func DownloadFileHandlerMeeting(c *gin.Context) {
	id := c.Param("id") // Mendapatkan ID dari URL
	filename := c.Param("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Filename is required"})
		return
	}

	// Membuat path lengkap berdasarkan ID dan nama file
	fullPath := fmt.Sprintf("%s/%s", id, filename)

	containerURL := getBlobServiceClient().NewContainerURL(containerNameMeeting)
	blobURL := containerURL.NewBlockBlobURL(fullPath)

	downloadResponse, err := blobURL.Download(context.TODO(), 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false, azblob.ClientProvidedKeyOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download file"})
		return
	}

	bodyStream := downloadResponse.Body(azblob.RetryReaderOptions{})
	defer bodyStream.Close()

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "application/octet-stream")

	// Mengirimkan data file ke client
	io.Copy(c.Writer, bodyStream)
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
	tanggal_target, err := time.Parse("2006-01-02", *requestBody.TanggalTarget)
	if err != nil {
		log.Printf("Error parsing date: %v", err)
		c.Status(400)
		c.JSON(400, gin.H{"error": "Invalid date format: " + err.Error()})
		return
	}

	tanggal_actual, err := time.Parse("2006-01-02", *requestBody.TanggalActual)
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
		TanggalTarget:    &tanggal_target,
		TanggalActual:    &tanggal_actual,
		CreateBy:         requestBody.CreateBy,
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
	meeting.CreateBy = requestBody.CreateBy

	if requestBody.TanggalTarget != nil {
		tanggal_target, err := time.Parse("2006-01-02", *requestBody.TanggalTarget)
		if err != nil {
			c.JSON(400, gin.H{"error": "Format tanggal tidak valid: " + err.Error()})
			return
		}
		meeting.TanggalTarget = &tanggal_target
	}

	if requestBody.TanggalActual != nil {
		tanggal_actual, err := time.Parse("2006-01-02", *requestBody.TanggalActual)
		if err != nil {
			c.JSON(400, gin.H{"error": "Format tanggal tidak valid: " + err.Error()})
			return
		}
		meeting.TanggalActual = &tanggal_actual
	}

	if requestBody.Task != nil {
		meeting.Task = requestBody.Task
	} else {
		meeting.Task = meeting.Task
	}

	if requestBody.TindakLanjut != nil {
		meeting.TindakLanjut = requestBody.TindakLanjut
	} else {
		meeting.TindakLanjut = meeting.TindakLanjut
	}

	if requestBody.Status != nil {
		meeting.Status = requestBody.Status
	} else {
		meeting.Status = meeting.Status
	}

	if requestBody.UpdatePengerjaan != nil {
		meeting.UpdatePengerjaan = requestBody.UpdatePengerjaan
	} else {
		meeting.UpdatePengerjaan = meeting.UpdatePengerjaan
	}

	if requestBody.Pic != nil {
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
	dir := "C:\\excel"
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
	sheetNames := []string{"MEMO", "PROJECT", "PERDIN", "SURAT MASUK", "SURAT KELUAR", "ARSIP", "MEETING", "MEETING SCHEDULE"}

	// Create sheets and set headers for "SAG" only
	for _, sheetName := range sheetNames {
		if sheetName == "MEETING" {
			f.NewSheet(sheetName)
			f.SetCellValue(sheetName, "A1", "TASK")
			f.SetCellValue(sheetName, "B1", "TINDAK LANJUT")
			f.SetCellValue(sheetName, "C1", "STATUS")
			f.SetCellValue(sheetName, "D1", "UPDATE PENGERJAAN")
			f.SetCellValue(sheetName, "E1", "PIC")
			f.SetCellValue(sheetName, "F1", "TANGGAL TARGET")
			f.SetCellValue(sheetName, "G1", "TANGGAL ACTUAL")

			f.SetColWidth(sheetName, "A", "G", 20)
		} else {
			f.NewSheet(sheetName)
		}
	}

	// Fetch initial data from the database
	var meetings []models.Meeting
	initializers.DB.Find(&meetings)

	// Write initial data to the "SAG" sheet
	meetingSheetName := "MEETING"
	for i, meeting := range meetings {
		tanggalTargetString := meeting.TanggalTarget.Format("2006-01-02")
		tanggalActualString := meeting.TanggalActual.Format("2006-01-02")
		rowNum := i + 2 // Start from the second row (first row is header)
		f.SetCellValue(meetingSheetName, fmt.Sprintf("A%d", rowNum), meeting.Task)
		f.SetCellValue(meetingSheetName, fmt.Sprintf("B%d", rowNum), meeting.TindakLanjut)
		f.SetCellValue(meetingSheetName, fmt.Sprintf("C%d", rowNum), meeting.Status)
		if meeting.UpdatePengerjaan != nil {
			f.SetCellValue(meetingSheetName, fmt.Sprintf("D%d", rowNum), *meeting.UpdatePengerjaan)
		} else {
			f.SetCellValue(meetingSheetName, fmt.Sprintf("D%d", rowNum), "")
		}
		f.SetCellValue(meetingSheetName, fmt.Sprintf("E%d", rowNum), meeting.Pic)
		f.SetCellValue(meetingSheetName, fmt.Sprintf("F%d", rowNum), tanggalTargetString)
		f.SetCellValue(meetingSheetName, fmt.Sprintf("G%d", rowNum), tanggalActualString)
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
	dir := "C:\\excel"
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
	sheetName := "MEETING"

	// Check if the sheet exists
	sheetIndex, err := f.GetSheetIndex(sheetName)
	if err != nil || sheetIndex == -1 {
		c.String(http.StatusBadRequest, "Lembar kerja MEETING tidak ditemukan")
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
	var meetings []models.Meeting
	initializers.DB.Find(&meetings)

	// Write data rows
	for i, meeting := range meetings {
		rowNum := i + 2 // Start from the second row (first row is header)
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), meeting.Task)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), meeting.TindakLanjut)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), meeting.Status)
		if meeting.UpdatePengerjaan != nil {
			f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), *meeting.UpdatePengerjaan)
		} else {
			f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), "")
		}
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), meeting.Pic)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", rowNum), meeting.TanggalTarget.Format("2006-01-02"))
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", rowNum), meeting.TanggalActual.Format("2006-01-02"))
	}

	// Save the file with updated data
	if err := f.SaveAs(filePath); err != nil {
		c.String(http.StatusInternalServerError, "Error menyimpan file: %v", err)
		return
	}


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
	sheetName := "MEETING"
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
		if len(row) < 7 {
			// Pastikan ada cukup kolom
			continue
		}
		task := row[0]
		tindakLanjut := row[1]
		status := row[2]
		updatePengerjaan := row[3]
		pic := row[4]
		tanggalTargetString := row[5]
		tanggalActualString := row[6]

		// Parse tanggal
		tanggalTarget, err := time.Parse("2006-01-02", tanggalTargetString)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid date format in row %d: %v", i+1, err)
			return
		}
		tanggalActual, err := time.Parse("2006-01-02", tanggalActualString)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid date format in row %d: %v", i+1, err)
			return
		}

		meeting := models.Meeting{
			Task:             &task,
			TindakLanjut:     &tindakLanjut,
			Status:           &status,
			UpdatePengerjaan: &updatePengerjaan,
			Pic:              &pic,
			TanggalTarget:    &tanggalTarget,
			TanggalActual:    &tanggalActual,
			CreateBy:         c.MustGet("username").(string),
		}

		// Simpan ke database
		if err := initializers.DB.Create(&meeting).Error; err != nil {
			log.Printf("Error saving record from row %d: %v", i+1, err)
			c.String(http.StatusInternalServerError, "Error saving record from row %d: %v", i+1, err)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data imported successfully."})
}