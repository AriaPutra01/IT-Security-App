package controllers

import (
	"context"
	"errors"
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
	"strings"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type MemoRequest struct {
	ID       uint    `gorm:"primaryKey"`
	Tanggal  *string `json:"tanggal"`
	NoMemo   string  `json:"no_memo"`
	Perihal  *string `json:"perihal"`
	Pic      *string `json:"pic"`
	Kategori *string `json:"kategori"`
	CreateBy string  `json:"create_by"`
}

func init() {
	err := godotenv.Load() // Memuat file .env
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	accountNameMemo = os.Getenv("ACCOUNT_NAME") // Mengambil nilai dari .env
	accountKeyMemo = os.Getenv("ACCOUNT_KEY")   // Mengambil nilai dari .env
	containerNameMemo = "memoits"               // Mengambil nilai dari .env
}

// Tambahkan variabel global untuk menyimpan kredensial
var (
	accountNameMemo   string
	accountKeyMemo    string
	containerNameMemo string
)

func getBlobServiceClientMemo() azblob.ServiceURL {
	creds, err := azblob.NewSharedKeyCredential(accountNameMemo, accountKeyMemo)
	if err != nil {
		panic("Failed to create shared key credential: " + err.Error())
	}

	pipeline := azblob.NewPipeline(creds, azblob.PipelineOptions{})

	// Build the URL for the Azure Blob Storage account
	URL, err := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/", accountNameMemo))
	if err != nil {
		log.Fatal("Invalid URL format")
	}

	// Create a ServiceURL object that wraps the URL and the pipeline
	serviceURL := azblob.NewServiceURL(*URL, pipeline)

	return serviceURL
}

func UploadHandlerMemo(c *gin.Context) {
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
	containerURL := getBlobServiceClientMemo().NewContainerURL(containerNameMemo)
	blobURL := containerURL.NewBlockBlobURL(filename)

	_, err = azblob.UploadStreamToBlockBlob(context.TODO(), src, blobURL, azblob.UploadStreamToBlockBlobOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengunggah file"})
		return
	}

	// Menambahkan log untuk menunjukkan ke kontainer mana file diunggah
	log.Printf("File %s berhasil diunggah ke kontainer %s", filename, containerNameMemo)

	c.JSON(http.StatusOK, gin.H{"message": "File berhasil diunggah"})
}

func GetFilesByIDMemo(c *gin.Context) {
	id := c.Param("id") // Mendapatkan ID dari URL

	containerURL := getBlobServiceClient().NewContainerURL(containerNameMemo)
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
func DeleteFileHandlerMemo(c *gin.Context) {
	filename := c.Param("filename")
	id := c.Param("id")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Filename is required"})
		return
	}

	// Membuat path lengkap berdasarkan ID dan nama file
	fullPath := fmt.Sprintf("%s/%s", id, filename)

	containerURL := getBlobServiceClient().NewContainerURL(containerNameMemo)
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
func DownloadFileHandlerMemo(c *gin.Context) {
	id := c.Param("id") // Mendapatkan ID dari URL
	filename := c.Param("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Filename is required"})
		return
	}

	// Membuat path lengkap berdasarkan ID dan nama file
	fullPath := fmt.Sprintf("%s/%s", id, filename)

	containerURL := getBlobServiceClient().NewContainerURL(containerNameMemo)
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

func GetLatestMemoNumber(NoMemo string) (string, error) {
	var lastMemo models.Memo
	// Ubah pencarian untuk menggunakan format yang benar
	searchPattern := fmt.Sprintf("%%/%s/M/%%", NoMemo) // Ini akan mencari format seperti '%/ITS-SAG/M/%'
	if err := initializers.DB.Where("no_memo LIKE ?", searchPattern).Order("id desc").First(&lastMemo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "00001", nil // Jika tidak ada catatan, kembalikan 00001
		}
		return "", err
	}

	// Ambil nomor memo terakhir, pisahkan, dan tambahkan 1
	parts := strings.Split(*lastMemo.NoMemo, "/")
	if len(parts) > 0 {
		number, err := strconv.Atoi(parts[0])
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%05d", number+1), nil // Tambahkan 1 ke nomor terakhir
	}

	return "00001", nil
}

func MemoIndex(c *gin.Context) {

	var memosag []models.Memo

	initializers.DB.Find(&memosag)

	c.JSON(200, gin.H{
		"memo": memosag,
	})

}

func MemoCreate(c *gin.Context) {
	var requestBody MemoRequest

	if err := c.BindJSON(&requestBody); err != nil {
		c.Status(400)
		c.Error(err) // log the error
		return
	}

	log.Println("Received request body:", requestBody)

	var tanggal *time.Time
	if requestBody.Tanggal != nil && *requestBody.Tanggal != "" {
		parsedTanggal, err := time.Parse("2006-01-02", *requestBody.Tanggal)
		if err != nil {
			log.Printf("Error parsing date: %v", err)
			c.JSON(400, gin.H{"error": "Invalid date format: " + err.Error()})
			return
		}
		tanggal = &parsedTanggal
	}

	log.Printf("Parsed date: %v", tanggal) // Tambahkan log ini untuk melihat tanggal yang diparsing

	nomor, err := GetLatestMemoNumber(requestBody.NoMemo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get latest memo number"})
		return
	}

	// Cek apakah nomor yang diterima adalah "00001"
	if nomor == "00001" {
		// Jika "00001", berarti ini adalah entri pertama
		log.Println("This is the first memo entry.")
	}

	tahun := time.Now().Year()

	// Menentukan format NoMemo berdasarkan kategori
	if requestBody.NoMemo == "ITS-SAG" {
		noMemo := fmt.Sprintf("%s/ITS-SAG/M/%d", nomor, tahun)
		requestBody.NoMemo = noMemo
		log.Printf("Generated NoMemo for ITS-SAG: %s", requestBody.NoMemo) // Log nomor memo
	} else if requestBody.NoMemo == "ITS-ISO" {
		noMemo := fmt.Sprintf("%s/ITS-ISO/M/%d", nomor, tahun)
		requestBody.NoMemo = noMemo
		log.Printf("Generated NoMemo for ITS-ISO: %s", requestBody.NoMemo) // Log nomor memo
	}

	requestBody.CreateBy = c.MustGet("username").(string)

	memosag := models.Memo{
		Tanggal: tanggal,             // Gunakan tanggal yang telah diparsing, bisa jadi nil jika input kosong
		NoMemo:  &requestBody.NoMemo, // Menggunakan NoMemo yang sudah diformat
		Perihal: requestBody.Perihal,
		Pic:     requestBody.Pic,
		// Kategori: requestBody.Kategori,
		CreateBy: requestBody.CreateBy,
	}

	result := initializers.DB.Create(&memosag)
	if result.Error != nil {
		log.Printf("Error saving memo: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Memo Sag"})
		return
	}
	log.Printf("Memo created successfully: %v", memosag)

	c.JSON(201, gin.H{
		"memo": memosag,
	})
}

func MemoShow(c *gin.Context) {

	id := c.Params.ByName("id")

	var memosag models.Memo

	initializers.DB.First(&memosag, id)

	c.JSON(200, gin.H{
		"memo": memosag,
	})

}

func MemoUpdate(c *gin.Context) {
	var requestBody MemoRequest

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	id := c.Param("id")
	var memo models.Memo

	// Cari memo berdasarkan ID
	if err := initializers.DB.First(&memo, id).Error; err != nil {
		log.Printf("Memo with ID %s not found: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Memo not found"})
		return
	}

	// Update tanggal jika diberikan dan tidak kosong
	if requestBody.Tanggal != nil && *requestBody.Tanggal != "" {
		parsedTanggal, err := time.Parse("2006-01-02", *requestBody.Tanggal)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
			return
		}
		memo.Tanggal = &parsedTanggal
	}

	// Update fields lainnya
	if requestBody.Perihal != nil {
		memo.Perihal = requestBody.Perihal
	}
	if requestBody.Pic != nil {
		memo.Pic = requestBody.Pic
	}

	// Simpan perubahan
	if err := initializers.DB.Save(&memo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update memo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Memo updated successfully", "memo": memo})
}

func MemoDelete(c *gin.Context) {

	id := c.Params.ByName("id")

	var memosag models.Memo

	if err := initializers.DB.First(&memosag, id); err.Error != nil {
		c.JSON(404, gin.H{"error": "Memo not found"})
		return
	}

	if err := initializers.DB.Delete(&memosag).Error; err != nil {
		c.JSON(400, gin.H{"error": "Failed to delete Memo: " + err.Error()})
		return
	}

	c.Status(204)

}

// func CreateExcelMemo(c *gin.Context) {
// 	dir := "C:\\excel"
// 	baseFileName := "its_report"
// 	filePath := filepath.Join(dir, baseFileName+".xlsx")

// 	// Check if the file already exists
// 	if _, err := os.Stat(filePath); err == nil {
// 		// File exists, append "_new" to the file name
// 		baseFileName += "_new"
// 	}

// 	fileName := baseFileName + ".xlsx"

// 	// File does not exist, create a new file
// 	f := excelize.NewFile()

// 	// Define sheet names
// 	sheetNames := []string{"MEMO", "PROJECT", "PERDIN", "SURAT MASUK", "SURAT KELUAR", "ARSIP", "MEETING", "MEETING SCHEDULE"}

// 	// Create sheets and set headers for "SAG" only
// 	for _, sheetName := range sheetNames {
// 		if sheetName == "MEMO" {
// 			f.NewSheet(sheetName)
// 			f.SetCellValue(sheetName, "A1", "Tanggal")
// 			f.SetCellValue(sheetName, "B1", "NoMemo")
// 			f.SetCellValue(sheetName, "C1", "Perihal")
// 			f.SetCellValue(sheetName, "D1", "Kategori")
// 			f.SetCellValue(sheetName, "E1", "Pic")

// 			f.SetColWidth(sheetName, "A", "E", 20)
// 		} else {
// 			f.NewSheet(sheetName)
// 		}
// 	}

// 	// Fetch initial data from the database
// 	var memos []models.Memo
// 	initializers.DB.Find(&memos)

// 	// Write initial data to the "SAG" sheet
// 	memoSheetName := "MEMO"
// 	for i, memo := range memos {
// 		tanggalString := memo.Tanggal.Format("2006-01-02")
// 		rowNum := i + 2 // Start from the second row (first row is header)
// 		f.SetCellValue(memoSheetName, fmt.Sprintf("A%d", rowNum), tanggalString)
// 		f.SetCellValue(memoSheetName, fmt.Sprintf("B%d", rowNum), memo.NoMemo)
// 		f.SetCellValue(memoSheetName, fmt.Sprintf("C%d", rowNum), memo.Perihal)
// 		// f.SetCellValue(memoSheetName, fmt.Sprintf("D%d", rowNum), memo.Kategori)
// 		f.SetCellValue(memoSheetName, fmt.Sprintf("E%d", rowNum), memo.Pic)
// 	}

// 	// Delete the default "Sheet1" sheet
// 	if err := f.DeleteSheet("Sheet1"); err != nil {
// 		panic(err) // Handle error jika bukan error "sheet tidak ditemukan"
// 	}

// 	// Save the newly created file
// 	buf, err := f.WriteToBuffer()
// 	if err != nil {
// 		c.String(http.StatusInternalServerError, "Error saving file: %v", err)
// 		return
// 	}

// 	// Serve the file to the client
// 	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
// 	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
// 	c.Writer.Write(buf.Bytes())
// }

func exportMemoToExcel(memos []models.Memo) (*excelize.File, error) {
	// Buat file Excel baru
	f := excelize.NewFile()

	sheetNames := []string{"MEMO", "BERITA ACARA", "SK", "SURAT", "PROJECT", "PERDIN", "SURAT MASUK", "SURAT KELUAR", "ARSIP", "MEETING", "MEETING SCHEDULE"}

	for _, sheetName := range sheetNames {
		f.NewSheet(sheetName)
		if sheetName == "MEMO" {
			// Header untuk SAG (kolom kiri)
			f.SetCellValue(sheetName, "A1", "Tanggal")
			f.SetCellValue(sheetName, "B1", "No Surat")
			f.SetCellValue(sheetName, "C1", "Perihal")
			f.SetCellValue(sheetName, "D1", "PIC")

			// Header untuk ISO (kolom kanan)
			f.SetCellValue(sheetName, "F1", "Tanggal")
			f.SetCellValue(sheetName, "G1", "No Surat")
			f.SetCellValue(sheetName, "H1", "Perihal")
			f.SetCellValue(sheetName, "I1", "PIC")
		}
	}
	f.DeleteSheet("Sheet1")

	// Inisialisasi baris awal
	rowSAG := 2
	rowISO := 2

	// Loop melalui data memo
	for _, memo := range memos {
		// Pastikan untuk dereferensikan pointer jika tidak nil
		var tanggal, noMemo, perihal, pic string
		if memo.Tanggal != nil {
			tanggal = memo.Tanggal.Format("2006-01-02") // Format tanggal sesuai kebutuhan
		}
		if memo.NoMemo != nil {
			noMemo = *memo.NoMemo
		}
		if memo.Perihal != nil {
			perihal = *memo.Perihal
		}
		if memo.Pic != nil {
			pic = *memo.Pic
		}

		// Pisahkan NoMemo untuk mendapatkan tipe memo
		parts := strings.Split(*memo.NoMemo, "/")
		if len(parts) > 1 && parts[1] == "ITS-SAG" {
			// Isi kolom SAG di sebelah kiri
			f.SetCellValue("MEMO", fmt.Sprintf("A%d", rowSAG), tanggal)
			f.SetCellValue("MEMO", fmt.Sprintf("B%d", rowSAG), noMemo)
			f.SetCellValue("MEMO", fmt.Sprintf("C%d", rowSAG), perihal)
			f.SetCellValue("MEMO", fmt.Sprintf("D%d", rowSAG), pic)
			rowSAG++
		} else if len(parts) > 1 && parts[1] == "ITS-ISO" {
			// Isi kolom ISO di sebelah kanan
			f.SetCellValue("MEMO", fmt.Sprintf("F%d", rowISO), tanggal)
			f.SetCellValue("MEMO", fmt.Sprintf("G%d", rowISO), noMemo)
			f.SetCellValue("MEMO", fmt.Sprintf("H%d", rowISO), perihal)
			f.SetCellValue("MEMO", fmt.Sprintf("I%d", rowISO), pic)
			rowISO++
		}
	}

	// style Line
	lastRowSAG := rowSAG - 1
	lastRowISO := rowISO - 1
	lastRow := lastRowSAG
	if lastRowISO > lastRowSAG {
		lastRow = lastRowISO
	}

	// Set lebar kolom agar rapi
	f.SetColWidth("MEMO", "A", "D", 20)
	f.SetColWidth("MEMO", "F", "I", 20)
	f.SetColWidth("MEMO", "E", "E", 2)
	for i := 2; i <= lastRow; i++ {
		f.SetRowHeight("MEMO", i, 30)
	}

	// style Line
	styleLine, err := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"000000"}, Pattern: 1},
		Border: []excelize.Border{
			{Type: "bottom", Color: "FFFFFF", Style: 2},
		},
	})
	if err != nil {
		fmt.Println(err)
	}
	err = f.SetCellStyle("MEMO", "E1", fmt.Sprintf("E%d", lastRow), styleLine)

	// style Border
	styleBorder, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "8E8E8E", Style: 2},
			{Type: "top", Color: "8E8E8E", Style: 2},
			{Type: "bottom", Color: "8E8E8E", Style: 2},
			{Type: "right", Color: "8E8E8E", Style: 2},
		},
	})
	if err != nil {
		fmt.Println(err)
	}
	err = f.SetCellStyle("MEMO", "A1", fmt.Sprintf("D%d", lastRow), styleBorder)
	err = f.SetCellStyle("MEMO", "F1", fmt.Sprintf("I%d", lastRow), styleBorder)

	return f, nil
}

// Handler untuk melakukan export Excel dengan Gin
func ExportMemoHandler(c *gin.Context) {
	// Data memo contoh
	var memos []models.Memo
	initializers.DB.Find(&memos)

	// Buat file Excel
	f, err := exportMemoToExcel(memos)
	if err != nil {
		c.String(http.StatusInternalServerError, "Gagal mengekspor data ke Excel")
		return
	}

	// Set nama file dan header untuk download
	fileName := fmt.Sprintf("MemoData_%s.xlsx", time.Now().Format("20060102"))
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Type", "application/octet-stream")

	// Simpan file Excel ke dalam buffer
	if err := f.Write(c.Writer); err != nil {
		c.String(http.StatusInternalServerError, "Gagal menyimpan file Excel")
	}
}

func UpdateSheetMemo(c *gin.Context) {
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
		// f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), memo.Kategori)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), memo.Pic)
	}

	// Save the file with updated data
	if err := f.SaveAs(filePath); err != nil {
		c.String(http.StatusInternalServerError, "Error menyimpan file: %v", err)
		return
	}
}

func excelDateToTimeMemo(excelDate int) (time.Time, error) {
	baseDate := time.Date(1899, time.December, 30, 0, 0, 0, 0, time.UTC)
	days := time.Duration(excelDate) * 24 * time.Hour
	return baseDate.Add(days), nil
}

func ImportExcelMemo(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "Error retrieving the file: %v", err)
		return
	}
	defer file.Close()

	tempFile, err := os.CreateTemp("", "*.xlsx")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error creating temporary file: %v", err)
		return
	}
	defer os.Remove(tempFile.Name())

	if _, err := io.Copy(tempFile, file); err != nil {
		c.String(http.StatusInternalServerError, "Error copying file: %v", err)
		return
	}

	tempFile.Seek(0, 0)
	f, err := excelize.OpenFile(tempFile.Name())
	if err != nil {
		c.String(http.StatusInternalServerError, "Error opening file: %v", err)
		return
	}
	defer f.Close()

	sheetName := "MEMO"
	rows, err := f.GetRows(sheetName)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error getting rows: %v", err)
		return
	}

	log.Println("Processing rows...") // Log untuk memulai proses baris
	for i, row := range rows {
		if i == 0 { // Skip header or initial rows if necessary
			continue
		}

		// Count non-empty columns
		nonEmptyCount := 0
		for _, cell := range row {
			if cell != "" {
				nonEmptyCount++
			}
		}

		// Skip rows with less than 5 non-empty columns
		if nonEmptyCount < 4 {
			log.Printf("Row %d skipped: less than 5 columns filled, only %d filled", i+1, nonEmptyCount)
			continue
		}

		// Assign values from row, using nil if the cell is empty or column does not exist
		var (
			tanggal = getStringOrNil(getColumn(row, 0))
			noMemo  = getStringOrNil(getColumn(row, 1))
			perihal = getStringOrNil(getColumn(row, 2))
			pic     = getStringOrNil(getColumn(row, 3))
			// kategori    = getStringOrNil(getColumn(row, 4))
		)

		// Convert string dates to time.Time pointers if not nil
		var tanggalTime *time.Time
		if tanggal != nil {
			parsed, err := parseDate(*tanggal)
			if err != nil {
				log.Printf("Error parsing date from row %d: %v", i+1, err)
				continue
			}
			tanggalTime = &parsed
		}

		memo := models.Memo{
			Tanggal: tanggalTime,
			NoMemo:  noMemo,
			Perihal: perihal,
			Pic:     pic,
			// Kategori: kategori,
			CreateBy: c.MustGet("username").(string),
		}

		if err := initializers.DB.Create(&memo).Error; err != nil {
			log.Printf("Error saving record from row %d: %v", i+1, err)
			continue
		}
		log.Printf("Row %d imported successfully", i+1) // Log untuk setiap baris yang berhasil diimpor
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data imported successfully, check logs for any skipped rows."})
}
