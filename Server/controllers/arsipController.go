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
	"github.com/xuri/excelize/v2"
)

type arsipRequest struct {
	ID                uint    `gorm:"primaryKey"`
	NoArsip           *string `json:"no_arsip"`
	JenisDokumen      *string `json:"jenis_dokumen"`
	NoDokumen         *string `json:"no_dokumen"`
	Perihal           *string `json:"perihal"`
	NoBox             *string `json:"no_box"`
	TanggalDokumen    *string `json:"tanggal_dokumen"`
	TanggalPenyerahan *string `json:"tanggal_penyerahan"`
	Keterangan        *string `json:"keterangan"`
	CreateBy          string  `json:"create_by"`
}

var (
	accountName   = "itsproject"
	accountKey    = "EnrPkwbyOBKlj57MliEipaIyhiYopF8RxlJL3htHGCLXg2vlTfIwiGQedB+GS9XiN95azazsLANb+ASt72N5xQ=="
	containerName = "projectits"
)

func getBlobServiceClient() azblob.ServiceURL {
	creds, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		panic("Failed to create shared key credential: " + err.Error())
	}

	pipeline := azblob.NewPipeline(creds, azblob.PipelineOptions{})

	// Build the URL for the Azure Blob Storage account
	URL, err := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/", accountName))
	if err != nil {
		log.Fatal("Invalid URL format")
	}

	// Create a ServiceURL object that wraps the URL and the pipeline
	serviceURL := azblob.NewServiceURL(*URL, pipeline)

	return serviceURL
}

func UploadHandler(c *gin.Context) {
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
	containerURL := getBlobServiceClient().NewContainerURL(containerName)
	blobURL := containerURL.NewBlockBlobURL(filename)

	_, err = azblob.UploadStreamToBlockBlob(context.TODO(), src, blobURL, azblob.UploadStreamToBlockBlobOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengunggah file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File berhasil diunggah"})
}

func GetFilesByID(c *gin.Context) {
	id := c.Param("id") // Mendapatkan ID dari URL

	containerURL := getBlobServiceClient().NewContainerURL(containerName)
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
func DeleteFileHandler(c *gin.Context) {
	filename := c.Param("filename")
	id := c.Param("id")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Filename is required"})
		return
	}

	// Membuat path lengkap berdasarkan ID dan nama file
	fullPath := fmt.Sprintf("%s/%s", id, filename)

	containerURL := getBlobServiceClient().NewContainerURL(containerName)
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
func DownloadFileHandler(c *gin.Context) {
	id := c.Param("id") // Mendapatkan ID dari URL
	filename := c.Param("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Filename is required"})
		return
	}

	// Membuat path lengkap berdasarkan ID dan nama file
	fullPath := fmt.Sprintf("%s/%s", id, filename)

	containerURL := getBlobServiceClient().NewContainerURL(containerName)
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

func ArsipIndex(c *gin.Context) {
	var arsip []models.Arsip
	initializers.DB.Find(&arsip)
	c.JSON(200, gin.H{
		"arsip": arsip,
	})
}

// Fungsi untuk membuat arsip baru
func ArsipCreate(c *gin.Context) {
	var requestBody arsipRequest
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}
	requestBody.CreateBy = c.MustGet("username").(string)

	var tanggal *time.Time
	if requestBody.TanggalDokumen != nil && *requestBody.TanggalDokumen != "" {
		parsedTanggal, err := time.Parse("2006-01-02", *requestBody.TanggalDokumen)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
			return
		}
		tanggal = &parsedTanggal
	}

	arsip := models.Arsip{
		NoArsip:           requestBody.NoArsip,
		JenisDokumen:      requestBody.JenisDokumen,
		NoDokumen:         requestBody.NoDokumen,
		Perihal:           requestBody.Perihal,
		NoBox:             requestBody.NoBox,
		Keterangan:        requestBody.Keterangan,
		TanggalDokumen:    tanggal,
		TanggalPenyerahan: tanggal, // Assuming same date handling for TanggalPenyerahan
		CreateBy:          requestBody.CreateBy,
	}

	if err := initializers.DB.Create(&arsip).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create arsip"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"arsip": arsip})
}

func ArsipShow(c *gin.Context) {
	id := c.Param("id")
	var arsip models.Arsip
	if err := initializers.DB.First(&arsip, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Arsip not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"arsip": arsip})
}

func ArsipUpdate(c *gin.Context) {
	id := c.Param("id")
	var requestBody arsipRequest
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	var arsip models.Arsip
	if err := initializers.DB.First(&arsip, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Arsip not found"})
		return
	}

	if requestBody.TanggalDokumen != nil {
		tanggal, err := time.Parse("2006-01-02", *requestBody.TanggalDokumen)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
			return
		}
		arsip.TanggalDokumen = &tanggal
	}

	// Update fields if provided in request
	if requestBody.NoArsip != nil {
		arsip.NoArsip = requestBody.NoArsip
	}
	if requestBody.JenisDokumen != nil {
		arsip.JenisDokumen = requestBody.JenisDokumen
	}
	if requestBody.NoDokumen != nil {
		arsip.NoDokumen = requestBody.NoDokumen
	}
	if requestBody.Perihal != nil {
		arsip.Perihal = requestBody.Perihal
	}
	if requestBody.NoBox != nil {
		arsip.NoBox = requestBody.NoBox
	}
	if requestBody.Keterangan != nil {
		arsip.Keterangan = requestBody.Keterangan
	}
	if requestBody.CreateBy != "" {
		arsip.CreateBy = requestBody.CreateBy
	}

	if err := initializers.DB.Save(&arsip).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update arsip"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"arsip": arsip})
}

func ArsipDelete(c *gin.Context) {
	id := c.Param("id")
	var arsip models.Arsip
	if err := initializers.DB.Where("id = ?", id).Delete(&arsip).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete arsip"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Arsip deleted successfully"})
}

func CreateExcelArsip(c *gin.Context) {
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

	// Create sheets and set headers for "ARSIP" only
	for _, sheetName := range sheetNames {
		if sheetName == "ARSIP" {
			f.NewSheet(sheetName)
			f.SetCellValue(sheetName, "A1", "No Arsip")
			f.SetCellValue(sheetName, "B1", "Jenis Dokumen")
			f.SetCellValue(sheetName, "C1", "No Dokumen")
			f.SetCellValue(sheetName, "D1", "Perihal")
			f.SetCellValue(sheetName, "E1", "No Box")
			f.SetCellValue(sheetName, "F1", "Keterangan")
			f.SetCellValue(sheetName, "G1", "Tanggal Dokumen")
			f.SetCellValue(sheetName, "H1", "Tanggal Penyerahan")

			// Set column widths for better readability
			f.SetColWidth(sheetName, "A", "H", 20)
		} else {
			f.NewSheet(sheetName)
		}
	}

	// Fetch initial data from the database
	var arsips []models.Arsip
	initializers.DB.Find(&arsips)

	// Write initial data to the "ARSIP" sheet
	// ... existing code ...

	// Write initial data to the "ARSIP" sheet
	arsipSheetName := "ARSIP"
	for i, arsip := range arsips {
		rowNum := i + 2 // Start from the second row (first row is header)
		f.SetCellValue(arsipSheetName, fmt.Sprintf("A%d", rowNum), derefString(arsip.NoArsip))
		f.SetCellValue(arsipSheetName, fmt.Sprintf("B%d", rowNum), derefString(arsip.JenisDokumen))
		f.SetCellValue(arsipSheetName, fmt.Sprintf("C%d", rowNum), derefString(arsip.NoDokumen))
		f.SetCellValue(arsipSheetName, fmt.Sprintf("D%d", rowNum), derefString(arsip.Perihal))
		f.SetCellValue(arsipSheetName, fmt.Sprintf("E%d", rowNum), derefString(arsip.NoBox))
		f.SetCellValue(arsipSheetName, fmt.Sprintf("F%d", rowNum), derefString(arsip.Keterangan))
		f.SetCellValue(arsipSheetName, fmt.Sprintf("G%d", rowNum), arsip.TanggalDokumen.Format("2006-01-02"))
		f.SetCellValue(arsipSheetName, fmt.Sprintf("H%d", rowNum), arsip.TanggalPenyerahan.Format("2006-01-02"))
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

func derefString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func UpdateSheetArsip(c *gin.Context) {
	dir := "C:\\excel"
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

	sheetName := "ARSIP"

	if _, err := f.GetSheetIndex(sheetName); err == nil {
		f.DeleteSheet(sheetName)
	}
	f.NewSheet(sheetName)

	f.SetCellValue(sheetName, "A1", "No Arsip")
	f.SetCellValue(sheetName, "B1", "Jenis Dokumen")
	f.SetCellValue(sheetName, "C1", "No Dokumen")
	f.SetCellValue(sheetName, "D1", "Perihal")
	f.SetCellValue(sheetName, "E1", "No Box")
	f.SetCellValue(sheetName, "F1", "Keterangan")
	f.SetCellValue(sheetName, "G1", "Tanggal Dokumen")
	f.SetCellValue(sheetName, "H1", "Tanggal Penyerahan")

	var arsips []models.Arsip
	initializers.DB.Find(&arsips)

	for i, arsip := range arsips {
		rowNum := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), arsip.NoArsip)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), arsip.JenisDokumen)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), arsip.NoDokumen)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), arsip.Perihal)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), arsip.NoBox)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", rowNum), arsip.Keterangan)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", rowNum), arsip.TanggalDokumen.Format("2006-01-02"))
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", rowNum), arsip.TanggalPenyerahan.Format("2006-01-02"))
	}

	if err := f.SaveAs(filePath); err != nil {
		c.String(http.StatusInternalServerError, "Error menyimpan file: %v", err)
		return
	}
}

func ImportExcelArsip(c *gin.Context) {
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
	sheetName := "ARSIP"
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

		noArsip := row[0]
		jenisDokumen := row[1]
		noDokumen := row[2]
		perihal := row[3]
		noBox := row[4]
		keterangan := row[5]
		tanggalDokumen := row[6]
		tanggalPenyerahan := row[7]

		// Parse tanggal
		tanggalDokumenString, err := time.Parse("2006-01-02", tanggalDokumen)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid date format in row %d: %v", i+1, err)
			return
		}
		tanggalPenyerahanString, err := time.Parse("2006-01-02", tanggalPenyerahan)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid date format in row %d: %v", i+1, err)
			return
		}

		arsip := models.Arsip{
			NoArsip:           &noArsip,
			JenisDokumen:      &jenisDokumen,
			NoDokumen:         &noDokumen,
			Perihal:           &perihal,
			NoBox:             &noBox,
			Keterangan:        &keterangan,
			TanggalDokumen:    &tanggalDokumenString,
			TanggalPenyerahan: &tanggalPenyerahanString,
			// CreateBy:          c.MustGet("username").(string),
		}

		// Simpan ke database
		if err := initializers.DB.Create(&arsip).Error; err != nil {
			log.Printf("Error saving record from row %d: %v", i+1, err)
			c.String(http.StatusInternalServerError, "Error saving record from row %d: %v", i+1, err)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data imported successfully."})
}
