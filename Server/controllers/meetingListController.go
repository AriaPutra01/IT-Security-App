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

type MeetingListRequest struct {
	ID       uint    `gorm:"primaryKey"`
	Hari     *string `json:"hari"`
	Tanggal  *string `json:"tanggal"`
	Perihal  *string `json:"perihal"`
	Waktu    *string `json:"waktu"`
	Selesai  *string `json:"selesai"`
	Tempat   *string `json:"tempat"`
	Pic      *string  `json:"pic"`
	Status   *string  `json:"status"`
	CreateBy string  `json:"create_by"`
	Info     string  `json:"info"`
	Color    string  `json:"color"`
}

func init() {
	err := godotenv.Load() // Memuat file .env
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	accountNameMeetingList = os.Getenv("ACCOUNT_NAME")                  // Mengambil nilai dari .env
	accountKeyMeetingList = os.Getenv("ACCOUNT_KEY")                    // Mengambil nilai dari .env
	containerNameMeetingList = os.Getenv("CONTAINER_NAME_MEETING_LIST") // Mengambil nilai dari .env
}

// Tambahkan variabel global untuk menyimpan kredensial
var (
	accountNameMeetingList   string
	accountKeyMeetingList    string
	containerNameMeetingList string
)

func getBlobServiceClientMeetingList() azblob.ServiceURL {
	creds, err := azblob.NewSharedKeyCredential(accountNameMeetingList, accountKeyMeetingList)
	if err != nil {
		panic("Failed to create shared key credential: " + err.Error())
	}

	pipeline := azblob.NewPipeline(creds, azblob.PipelineOptions{})

	// Build the URL for the Azure Blob Storage account
	URL, err := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/", accountNameMeetingList))
	if err != nil {
		log.Fatal("Invalid URL format")
	}

	// Create a ServiceURL object that wraps the URL and the pipeline
	serviceURL := azblob.NewServiceURL(*URL, pipeline)

	return serviceURL
}

func UploadHandlerMeetingList(c *gin.Context) {
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
	containerURL := getBlobServiceClient().NewContainerURL(containerNameMeetingList)
	blobURL := containerURL.NewBlockBlobURL(filename)

	_, err = azblob.UploadStreamToBlockBlob(context.TODO(), src, blobURL, azblob.UploadStreamToBlockBlobOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengunggah file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File berhasil diunggah"})
}

func GetFilesByIDMeetingList(c *gin.Context) {
	id := c.Param("id") // Mendapatkan ID dari URL

	containerURL := getBlobServiceClient().NewContainerURL(containerNameMeetingList)
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
func DeleteFileHandlerMeetingList(c *gin.Context) {
	filename := c.Param("filename")
	id := c.Param("id")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Filename is required"})
		return
	}

	// Membuat path lengkap berdasarkan ID dan nama file
	fullPath := fmt.Sprintf("%s/%s", id, filename)

	containerURL := getBlobServiceClient().NewContainerURL(containerNameMeetingList)
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
func DownloadFileHandlerMeetingList(c *gin.Context) {
	id := c.Param("id") // Mendapatkan ID dari URL
	filename := c.Param("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Filename is required"})
		return
	}

	// Membuat path lengkap berdasarkan ID dan nama file
	fullPath := fmt.Sprintf("%s/%s", id, filename)

	containerURL := getBlobServiceClient().NewContainerURL(containerNameMeetingList)
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

func MeetingListIndex(c *gin.Context) {

	var meetingList []models.MeetingSchedule

	initializers.DB.Find(&meetingList)

	c.JSON(200, gin.H{
		"meetingschedule": meetingList,
	})

}

func MeetingListCreate(c *gin.Context) {

	var requestBody MeetingListRequest

	if err := c.BindJSON(&requestBody); err != nil {
		c.Status(400)
		c.Error(err) // log the error
		return
	}

	// Add some logging to see what's being received
	log.Println("Received request body:", requestBody)

	// Parse the date string
	tanggal, err := time.Parse("2006-01-02", *requestBody.Tanggal)
	if err != nil {
		log.Printf("Error parsing date: %v", err)
		c.Status(400)
		c.JSON(400, gin.H{"error": "Invalid date format: " + err.Error()})
		return
	}

	requestBody.CreateBy = c.MustGet("username").(string)

	meetingList := models.MeetingSchedule{
		Hari:     requestBody.Hari,
		Tanggal:  &tanggal,
		Perihal:  requestBody.Perihal,
		Waktu:    requestBody.Waktu,
		Selesai:  requestBody.Selesai,
		Tempat:   requestBody.Tempat,
		Pic:      requestBody.Pic,
		Status:   requestBody.Status,
		CreateBy: requestBody.CreateBy,
		Color:    requestBody.Color,
	}

	result := initializers.DB.Create(&meetingList)

	if result.Error != nil {
		c.Status(400)
		c.JSON(400, gin.H{"error": "Failed to create Meeting: " + result.Error.Error()})
		return
	}

	c.JSON(201, gin.H{
		"meetingschedule": meetingList,
	})

}

func MeetingListShow(c *gin.Context) {

	id := c.Params.ByName("id")

	var meetingList models.MeetingSchedule

	initializers.DB.First(&meetingList, id)

	c.JSON(200, gin.H{
		"meetinglist": meetingList,
	})

}

func MeetingListUpdate(c *gin.Context) {

	var requestBody MeetingListRequest

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	id := c.Params.ByName("id")

	var meetingList models.MeetingSchedule

	if err := initializers.DB.First(&meetingList, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "meeting not found"})
		return
	}

	requestBody.CreateBy = c.MustGet("username").(string)
	meetingList.CreateBy = requestBody.CreateBy

	if requestBody.Tanggal != nil {
		tanggal, err := time.Parse("2006-01-02", *requestBody.Tanggal)
		if err != nil {
			c.JSON(400, gin.H{"error": "Format tanggal tidak valid: " + err.Error()})
			return
		}
		meetingList.Tanggal = &tanggal
	}

	if requestBody.Hari != nil {
		meetingList.Hari = requestBody.Hari
	} else {
		meetingList.Hari = meetingList.Hari
	}

	if requestBody.Perihal != nil {
		meetingList.Perihal = requestBody.Perihal
	} else {
		meetingList.Perihal = meetingList.Perihal

	}

	if requestBody.Waktu != nil {
		meetingList.Waktu = requestBody.Waktu
	} else {
		meetingList.Waktu = meetingList.Waktu
	}

	if requestBody.Selesai != nil {
		meetingList.Selesai = requestBody.Selesai
	} else {
		meetingList.Selesai = meetingList.Selesai
	}

	if requestBody.Tempat != nil {
		meetingList.Tempat = requestBody.Tempat
	} else {
		meetingList.Tempat = meetingList.Tempat
	}

	if requestBody.Pic != nil {
		meetingList.Pic = requestBody.Pic
	} else {
		meetingList.Pic = meetingList.Pic
	}

	if requestBody.Status != nil {
		meetingList.Status = requestBody.Status
	} else {
		meetingList.Status = meetingList.Status
	}

	if requestBody.Color != "" {
		meetingList.Color = requestBody.Color
	} else {
		meetingList.Color = meetingList.Color
	}

	initializers.DB.Save(&meetingList)

	c.JSON(200, gin.H{
		"meetinglist": meetingList,
	})
}

func MeetingListDelete(c *gin.Context) {

	id := c.Params.ByName("id")

	var meetingList models.MeetingSchedule

	if err := initializers.DB.First(&meetingList, id); err.Error != nil {
		c.JSON(404, gin.H{"error": "meeting not found"})
		return
	}

	if err := initializers.DB.Delete(&meetingList).Error; err != nil {
		c.JSON(400, gin.H{"error": "Failed to delete Memo: " + err.Error()})
		return
	}

	c.Status(204)

}

func CreateExcelMeetingList(c *gin.Context) {
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
		if sheetName == "MEETING SCHEDULE" {
			f.NewSheet(sheetName)
			f.SetCellValue(sheetName, "A1", "Hari")
			f.SetCellValue(sheetName, "B1", "Tanggal")
			f.SetCellValue(sheetName, "C1", "Perihal")
			f.SetCellValue(sheetName, "D1", "Waktu")
			f.SetCellValue(sheetName, "E1", "Selesai")
			f.SetCellValue(sheetName, "F1", "Tempat")
			f.SetCellValue(sheetName, "G1", "Status")
			f.SetCellValue(sheetName, "H1", "PIC")

			f.SetColWidth(sheetName, "A", "H", 20)
		} else {
			f.NewSheet(sheetName)
		}
	}

	// Fetch initial data from the database
	var meetingList []models.MeetingSchedule
	initializers.DB.Find(&meetingList)

	// Write initial data to the "SAG" sheet
	meetingListSheetName := "MEETING SCHEDULE"
	for i, meetingList := range meetingList {
		tanggalString := meetingList.Tanggal.Format("2006-01-02")
		rowNum := i + 2 // Start from the second row (first row is header)
		f.SetCellValue(meetingListSheetName, fmt.Sprintf("A%d", rowNum), meetingList.Hari)
		f.SetCellValue(meetingListSheetName, fmt.Sprintf("B%d", rowNum), tanggalString)
		f.SetCellValue(meetingListSheetName, fmt.Sprintf("C%d", rowNum), meetingList.Perihal)

		// Handle Waktu
		if meetingList.Waktu != nil {
			f.SetCellValue(meetingListSheetName, fmt.Sprintf("D%d", rowNum), *meetingList.Waktu)
		} else {
			f.SetCellValue(meetingListSheetName, fmt.Sprintf("D%d", rowNum), "")
		}

		// Handle Selesai
		if meetingList.Selesai != nil {
			f.SetCellValue(meetingListSheetName, fmt.Sprintf("E%d", rowNum), *meetingList.Selesai)
		} else {
			f.SetCellValue(meetingListSheetName, fmt.Sprintf("E%d", rowNum), "")
		}

		if meetingList.Tempat != nil {
			f.SetCellValue(meetingListSheetName, fmt.Sprintf("F%d", rowNum), *meetingList.Tempat)
		} else {
			f.SetCellValue(meetingListSheetName, fmt.Sprintf("F%d", rowNum), "")
		}

		f.SetCellValue(meetingListSheetName, fmt.Sprintf("G%d", rowNum), meetingList.Status)
		f.SetCellValue(meetingListSheetName, fmt.Sprintf("H%d", rowNum), meetingList.Pic)
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

func UpdateSheetMeetingList(c *gin.Context) {
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
	sheetName := "MEETING SCHEDULE"

	// Check if the sheet exists
	sheetIndex, err := f.GetSheetIndex(sheetName)
	if err != nil || sheetIndex == -1 {
		c.String(http.StatusBadRequest, "Lembar kerja MEETING SCHEDULE tidak ditemukan")
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
	var meetingList []models.MeetingSchedule
	initializers.DB.Find(&meetingList)

	// Write data rows
	for i, meetingList := range meetingList {
		rowNum := i + 2 // Start from the second row (first row is header)
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), meetingList.Hari)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), meetingList.Tanggal.Format("2006-01-02"))
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), meetingList.Perihal)

		// Handle Waktu
		if meetingList.Waktu != nil {
			f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), *meetingList.Waktu)
		} else {
			f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), "")
		}

		// Handle Selesai
		if meetingList.Selesai != nil {
			f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), *meetingList.Selesai)
		} else {
			f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), "")
		}

		if meetingList.Tempat != nil {
			f.SetCellValue(sheetName, fmt.Sprintf("F%d", rowNum), *meetingList.Tempat)
		} else {
			f.SetCellValue(sheetName, fmt.Sprintf("F%d", rowNum), "")
		}

		f.SetCellValue(sheetName, fmt.Sprintf("G%d", rowNum), meetingList.Status)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", rowNum), meetingList.Pic)
	}

	// Save the file with updated data
	if err := f.SaveAs(filePath); err != nil {
		c.String(http.StatusInternalServerError, "Error menyimpan file: %v", err)
		return
	}

}

func ImportExcelMeetingList(c *gin.Context) {
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
	sheetName := "MEETING SCHEDULE"
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
		hari := row[0]
		tanggal := row[1]
		perihal := row[2]
		waktu := row[3]
		selesai := row[4]
		tempat := row[5]
		status := row[6]
		pic := row[7]

		// Parse tanggal
		tanggalString, err := time.Parse("2006-01-02", tanggal)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid date format in row %d: %v", i+1, err)
			return
		}

		meetingList := models.MeetingSchedule{
			Hari:     &hari,
			Tanggal:  &tanggalString,
			Perihal:  &perihal,
			Waktu:    &waktu,
			Selesai:  &selesai,
			Tempat:   &tempat,
			Status:   &status,
			Pic:      &pic,
			CreateBy: c.MustGet("username").(string),
		}

		// Simpan ke database
		if err := initializers.DB.Create(&meetingList).Error; err != nil {
			log.Printf("Error saving record from row %d: %v", i+1, err)
			c.String(http.StatusInternalServerError, "Error saving record from row %d: %v", i+1, err)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data imported successfully."})
}
