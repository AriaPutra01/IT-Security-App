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
	"strconv"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/xuri/excelize/v2"
)

type SuratKeluarRequest struct {
	ID       uint   `gorm:"primaryKey"`
	NoSurat  *string `json:"no_surat"`
	Title    *string `json:"title"`
	From     *string `json:"from"`
	Pic      *string `json:"pic"`
	Tanggal  *string `json:"tanggal"`
	CreateBy string `json:"create_by"`
}

func init() {
	err := godotenv.Load() // Memuat file .env
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	accountNameSuratKeluar = os.Getenv("ACCOUNT_NAME")                 // Mengambil nilai dari .env
	accountKeySuratKeluar = os.Getenv("ACCOUNT_KEY")                   // Mengambil nilai dari .env
	containerNameSuratKeluar = os.Getenv("CONTAINER_NAME_SURATKELUAR") // Mengambil nilai dari .env
}

// Tambahkan variabel global untuk menyimpan kredensial
var (
	accountNameSuratKeluar   string
	accountKeySuratKeluar    string
	containerNameSuratKeluar string
)

func getBlobServiceClientKeluar() azblob.ServiceURL {
	creds, err := azblob.NewSharedKeyCredential(accountNameSuratKeluar, accountKeySuratKeluar)
	if err != nil {
		panic("Failed to create shared key credential: " + err.Error())
	}

	pipeline := azblob.NewPipeline(creds, azblob.PipelineOptions{})

	// Build the URL for the Azure Blob Storage account
	URL, err := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/", accountNameSuratKeluar))
	if err != nil {
		log.Fatal("Invalid URL format")
	}

	// Create a ServiceURL object that wraps the URL and the pipeline
	serviceURL := azblob.NewServiceURL(*URL, pipeline)

	return serviceURL
}

func UploadHandlerSuratKeluar(c *gin.Context) {
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
	containerURL := getBlobServiceClient().NewContainerURL(containerNameSuratKeluar)
	blobURL := containerURL.NewBlockBlobURL(filename)

	_, err = azblob.UploadStreamToBlockBlob(context.TODO(), src, blobURL, azblob.UploadStreamToBlockBlobOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengunggah file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File berhasil diunggah"})
}

func GetFilesByIDSuratKeluar(c *gin.Context) {
	id := c.Param("id") // Mendapatkan ID dari URL

	containerURL := getBlobServiceClient().NewContainerURL(containerNameSuratKeluar)
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
func DeleteFileHandlerSuratKeluar(c *gin.Context) {
	filename := c.Param("filename")
	id := c.Param("id")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Filename is required"})
		return
	}

	// Membuat path lengkap berdasarkan ID dan nama file
	fullPath := fmt.Sprintf("%s/%s", id, filename)

	containerURL := getBlobServiceClient().NewContainerURL(containerNameSuratKeluar)
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
func DownloadFileHandlerSuratKeluar(c *gin.Context) {
	id := c.Param("id") // Mendapatkan ID dari URL
	filename := c.Param("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Filename is required"})
		return
	}

	// Membuat path lengkap berdasarkan ID dan nama file
	fullPath := fmt.Sprintf("%s/%s", id, filename)

	containerURL := getBlobServiceClient().NewContainerURL(containerNameSuratKeluar)
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

	requestBody.CreateBy = c.MustGet("username").(string)

	var tanggal *time.Time // Deklarasi variabel tanggal sebagai pointer ke time.Time
	if requestBody.Tanggal != nil && *requestBody.Tanggal != "" {
		// Parse the date string only if it's not nil and not empty
		parsedTanggal, err := time.Parse("2006-01-02", *requestBody.Tanggal)
		if err != nil {
			log.Printf("Error parsing date: %v", err)
			c.JSON(400, gin.H{"error": "Invalid date format: " + err.Error()})
			return
		}
		tanggal = &parsedTanggal
	}

	surat_keluar := models.SuratKeluar{
		NoSurat:  requestBody.NoSurat,
		Title:    requestBody.Title,
		From:     requestBody.From,
		Pic:      requestBody.Pic,
		Tanggal:  tanggal, // Gunakan tanggal yang telah diparsing, bisa jadi nil jika input kosong
		CreateBy: requestBody.CreateBy,
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

	requestBody.CreateBy = c.MustGet("username").(string)
	surat_keluar.CreateBy = requestBody.CreateBy

	if requestBody.Tanggal != nil {
		tanggal, err := time.Parse("2006-01-02", *requestBody.Tanggal)
		if err != nil {
			c.JSON(400, gin.H{"error": "Format tanggal tidak valid: " + err.Error()})
			return
		}
		surat_keluar.Tanggal = &tanggal
	}

	if requestBody.NoSurat != nil {
		surat_keluar.NoSurat = requestBody.NoSurat
	} else {
		surat_keluar.NoSurat = surat_keluar.NoSurat // gunakan nilai yang ada dari database
	}

	if requestBody.Title != nil {
		surat_keluar.Title = requestBody.Title
	} else {
		surat_keluar.Title = surat_keluar.Title // gunakan nilai yang ada dari database
	}

	if requestBody.Pic != nil {
		surat_keluar.Pic = requestBody.Pic
	} else {
		surat_keluar.Pic = surat_keluar.Pic // gunakan nilai yang ada dari database
	}

		if requestBody.CreateBy != "" {
		surat_keluar.CreateBy = requestBody.CreateBy
	} else {
		surat_keluar.CreateBy = surat_keluar.CreateBy // gunakan nilai yang ada dari database
	}

	initializers.DB.Model(&surat_keluar).Updates(surat_keluar)

	c.JSON(200, gin.H{
		"SuratKeluar": surat_keluar,
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
		"SuratKeluar": "Deleted",
	})
}

func CreateExcelSuratKeluar(c *gin.Context) {
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
		if sheetName == "SURAT KELUAR" {
			f.NewSheet(sheetName)
			f.SetCellValue(sheetName, "A1", "No Surat")
			f.SetCellValue(sheetName, "B1", "Title Of Letter")
			f.SetCellValue(sheetName, "C1", "From")
			f.SetCellValue(sheetName, "D1", "Pic")
			f.SetCellValue(sheetName, "E1", "Date Issue")

			f.SetColWidth(sheetName, "A", "E", 20)
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

func excelDateToTimeSuratKeluar(excelDate int) (time.Time, error) {
	baseDate := time.Date(1899, time.December, 30, 0, 0, 0, 0, time.UTC)
	days := time.Duration(excelDate) * 24 * time.Hour
	return baseDate.Add(days), nil
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

	// Definisikan semua format tanggal yang mungkin
	dateFormats := []string{
		"2 January 2006",
		"2006-01-02",
		"02-01-2006",
		"01/02/2006",
		"2006.01.02",
		"02/01/2006",
	}

	// Loop melalui baris dan simpan ke database
	for i, row := range rows {
		if i < 5 { // Lewati baris pertama yang merupakan header
			continue
		}
		if len(row) < 5 { // Pastikan ada cukup kolom
			continue
		}
		noSurat := row[0]
		title := row[1]
		from := row[2]
		pic := row[3]
		tanggalString := row[4]

		var tanggal time.Time
		var parseErr error

		// Coba konversi dari serial Excel jika tanggalString adalah angka
		if serial, err := strconv.Atoi(tanggalString); err == nil {
			tanggal, parseErr = excelDateToTimeSuratKeluar(serial)
		} else {
			// Coba parse menggunakan format tanggal yang sudah ada
			for _, format := range dateFormats {
				tanggal, parseErr = time.Parse(format, tanggalString)
				if parseErr == nil {
					break // Keluar dari loop jika parsing berhasil
				}
			}
		}

		surat_keluar := models.SuratKeluar{
			NoSurat:  &noSurat,
			Title:    &title,
			From:     &from,
			Pic:      &pic,
			Tanggal:  &tanggal,
			CreateBy: c.MustGet("username").(string),
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
