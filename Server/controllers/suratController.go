package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"project-its/initializers"
	"project-its/models"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type SuratRequest struct {
	ID       uint    `gorm:"primaryKey"`
	Tanggal  *string `json:"tanggal"`
	NoSurat  string  `json:"no_surat"`
	Perihal  *string `json:"perihal"`
	Pic      *string `json:"pic"`
	CreateBy string  `json:"create_by"`
}

// func init() {
// 	err := godotenv.Load() // Memuat file .env
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}

// 	accountNameMemo = os.Getenv("ACCOUNT_NAME") // Mengambil nilai dari .env
// 	accountKeyMemo = os.Getenv("ACCOUNT_KEY")   // Mengambil nilai dari .env
// 	containerNameMemo = "memoits"               // Mengambil nilai dari .env
// }

// // Tambahkan variabel global untuk menyimpan kredensial
// var (
// 	accountNameMemo   string
// 	accountKeyMemo    string
// 	containerNameMemo string
// )

// func getBlobServiceClientMemo() azblob.ServiceURL {
// 	creds, err := azblob.NewSharedKeyCredential(accountNameMemo, accountKeyMemo)
// 	if err != nil {
// 		panic("Failed to create shared key credential: " + err.Error())
// 	}

// 	pipeline := azblob.NewPipeline(creds, azblob.PipelineOptions{})

// 	// Build the URL for the Azure Blob Storage account
// 	URL, err := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/", accountNameMemo))
// 	if err != nil {
// 		log.Fatal("Invalid URL format")
// 	}

// 	// Create a ServiceURL object that wraps the URL and the pipeline
// 	serviceURL := azblob.NewServiceURL(*URL, pipeline)

// 	return serviceURL
// }

// func UploadHandlerMemo(c *gin.Context) {
// 	id := c.PostForm("id") // Mendapatkan ID dari form data
// 	file, err := c.FormFile("file")
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "File diperlukan"})
// 		return
// 	}

// 	// Membuat path berdasarkan ID
// 	filename := fmt.Sprintf("%s/%s", id, file.Filename)

// 	// Membuka file
// 	src, err := file.Open()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuka file"})
// 		return
// 	}
// 	defer src.Close()

// 	// Mengunggah file ke Azure Blob Storage
// 	containerURL := getBlobServiceClientMemo().NewContainerURL(containerNameMemo)
// 	blobURL := containerURL.NewBlockBlobURL(filename)

// 	_, err = azblob.UploadStreamToBlockBlob(context.TODO(), src, blobURL, azblob.UploadStreamToBlockBlobOptions{})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengunggah file"})
// 		return
// 	}

// 	// Menambahkan log untuk menunjukkan ke kontainer mana file diunggah
// 	log.Printf("File %s berhasil diunggah ke kontainer %s", filename, containerNameMemo)

// 	c.JSON(http.StatusOK, gin.H{"message": "File berhasil diunggah"})
// }

// func GetFilesByIDMemo(c *gin.Context) {
// 	id := c.Param("id") // Mendapatkan ID dari URL

// 	containerURL := getBlobServiceClient().NewContainerURL(containerNameMemo)
// 	prefix := fmt.Sprintf("%s/", id) // Prefix untuk daftar blob di folder tertentu (ID)

// 	var files []string
// 	for marker := (azblob.Marker{}); marker.NotDone(); {
// 		listBlob, err := containerURL.ListBlobsFlatSegment(context.TODO(), marker, azblob.ListBlobsSegmentOptions{
// 			Prefix: prefix, // Hanya daftar blob dengan prefix yang ditentukan (folder)
// 		})
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat daftar file"})
// 			return
// 		}

// 		for _, blobInfo := range listBlob.Segment.BlobItems {
// 			files = append(files, blobInfo.Name)
// 		}

// 		marker = listBlob.NextMarker
// 	}

// 	c.JSON(http.StatusOK, gin.H{"files": files}) // Pastikan mengembalikan array files
// }

// // Fungsi untuk menghapus file dari Azure Blob Storage
// func DeleteFileHandlerMemo(c *gin.Context) {
// 	filename := c.Param("filename")
// 	id := c.Param("id")
// 	if filename == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Filename is required"})
// 		return
// 	}

// 	// Membuat path lengkap berdasarkan ID dan nama file
// 	fullPath := fmt.Sprintf("%s/%s", id, filename)

// 	containerURL := getBlobServiceClient().NewContainerURL(containerNameMemo)
// 	blobURL := containerURL.NewBlockBlobURL(fullPath)

// 	// Menghapus blob
// 	_, err := blobURL.Delete(context.TODO(), azblob.DeleteSnapshotsOptionNone, azblob.BlobAccessConditions{})
// 	if err != nil {
// 		log.Printf("Error deleting file: %v", err) // Log kesalahan
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"}) // Pastikan ini ada
// }

// // Fungsi untuk mendownload file dari Azure Blob Storage
// func DownloadFileHandlerMemo(c *gin.Context) {
// 	id := c.Param("id") // Mendapatkan ID dari URL
// 	filename := c.Param("filename")
// 	if filename == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Filename is required"})
// 		return
// 	}

// 	// Membuat path lengkap berdasarkan ID dan nama file
// 	fullPath := fmt.Sprintf("%s/%s", id, filename)

// 	containerURL := getBlobServiceClient().NewContainerURL(containerNameMemo)
// 	blobURL := containerURL.NewBlockBlobURL(fullPath)

// 	downloadResponse, err := blobURL.Download(context.TODO(), 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false, azblob.ClientProvidedKeyOptions{})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download file"})
// 		return
// 	}

// 	bodyStream := downloadResponse.Body(azblob.RetryReaderOptions{})
// 	defer bodyStream.Close()

// 	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
// 	c.Header("Content-Type", "application/octet-stream")

// 	// Mengirimkan data file ke client
// 	io.Copy(c.Writer, bodyStream)
// }

func SuratIndex(c *gin.Context) {

	var surat []models.Surat
	if err := initializers.DB.Find(&surat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data surat: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"surat": surat})

}

func GetLatestSuratNumber(NoSurat string) (string, error) {
	var lastSurat models.Surat
	// Ubah pencarian untuk menggunakan format yang benar
	searchPattern := fmt.Sprintf("%%/%s/S/%%", NoSurat) // Ini akan mencari format seperti '%ITS-SAG/S/%'
	if err := initializers.DB.Where("no_surat LIKE ?", searchPattern).Order("id desc").First(&lastSurat).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "00001", nil // Jika tidak ada catatan, kembalikan 00001
		}
		return "", err
	}

	// Ambil nomor surat terakhir, pisahkan, dan tambahkan 1
	parts := strings.Split(*lastSurat.NoSurat, "/")
	if len(parts) > 0 {
		// Ambil bagian pertama dari parts yang seharusnya adalah nomor
		numberPart := parts[0]
		number, err := strconv.Atoi(numberPart)
		if err != nil {
			log.Printf("Error parsing number from surat: %v", err)
			return "", err
		}
		return fmt.Sprintf("%05d", number+1), nil // Tambahkan 1 ke nomor terakhir
	}

	return "00001", nil
}

func SuratCreate(c *gin.Context) {
	var requestBody SuratRequest

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

	nomor, err := GetLatestSuratNumber(requestBody.NoSurat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get latest surat number"})
		return
	}

	// Cek apakah nomor yang diterima adalah "00001"
	if nomor == "00001" {
		// Jika "00001", berarti ini adalah entri pertama
		log.Println("This is the first surat entry.")
	}

	tahun := time.Now().Year()

	// Menentukan format NoSurat berdasarkan kategori
	if requestBody.NoSurat == "ITS-SAG" {
		noSurat := fmt.Sprintf("%s/ITS-SAG/S/%d", nomor, tahun)
		requestBody.NoSurat = noSurat
		log.Printf("Generated NoSurat for ITS-SAG: %s", requestBody.NoSurat) // Log nomor surat
	} else if requestBody.NoSurat == "ITS-ISO" {
		noSurat := fmt.Sprintf("%s/ITS-ISO/S/%d", nomor, tahun)
		requestBody.NoSurat = noSurat
		log.Printf("Generated NoSurat for ITS-ISO: %s", requestBody.NoSurat) // Log nomor surat
	}

	requestBody.CreateBy = c.MustGet("username").(string)

	surat := models.Surat{
		Tanggal:  tanggal,              // Gunakan tanggal yang telah diparsing, bisa jadi nil jika input kosong
		NoSurat:  &requestBody.NoSurat, // Menggunakan NoMemo yang sudah diformat
		Perihal:  requestBody.Perihal,
		Pic:      requestBody.Pic,
		CreateBy: requestBody.CreateBy,
	}

	result := initializers.DB.Create(&surat)
	if result.Error != nil {
		log.Printf("Error saving surat: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Memo Sag"})
		return
	}
	log.Printf("Surat created successfully: %v", surat)

	c.JSON(201, gin.H{
		"surat": surat,
	})
}

func SuratShow(c *gin.Context) {

	id := c.Params.ByName("id")

	var surat models.Surat

	initializers.DB.First(&surat, id)

	c.JSON(200, gin.H{
		"surat": surat,
	})

}

func SuratUpdate(c *gin.Context) {

	var requestBody SuratRequest

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	id := c.Params.ByName("id")

	var surat models.Surat

	if err := initializers.DB.First(&surat, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Surat not found"})
		return
	}

	nomor, err := GetLatestSuratNumber(requestBody.NoSurat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get latest memo number"})
		return
	}

	// Cek apakah nomor yang diterima adalah "00001"
	if nomor == "00001" {
		// Jika "00001", berarti ini adalah entri pertama
		log.Println("This is the first surat entry.")
	}

	tahun := time.Now().Year()

	// Menentukan format NoSurat berdasarkan kategori
	if requestBody.NoSurat == "ITS-SAG" {
		noSurat := fmt.Sprintf("%s/ITS-SAG/S/%d", nomor, tahun)
		requestBody.NoSurat = noSurat
		log.Printf("Generated NoSurat for ITS-SAG: %s", requestBody.NoSurat) // Log nomor surat
	} else if requestBody.NoSurat == "ITS-ISO" {
		noSurat := fmt.Sprintf("%s/ITS-ISO/S/%d", nomor, tahun)
		requestBody.NoSurat = noSurat
		log.Printf("Generated NoSurat for ITS-ISO: %s", requestBody.NoSurat) // Log nomor surat
	}

	requestBody.CreateBy = c.MustGet("username").(string)
	surat.CreateBy = requestBody.CreateBy

	if requestBody.Tanggal != nil {
		tanggal, err := time.Parse("2006-01-02", *requestBody.Tanggal)
		if err != nil {
			c.JSON(400, gin.H{"error": "Format tanggal tidak valid: " + err.Error()})
			return
		}
		surat.Tanggal = &tanggal
	}

	if requestBody.NoSurat != "" {
		surat.NoSurat = &requestBody.NoSurat
	} else {
		surat.NoSurat = surat.NoSurat
	}

	if requestBody.Perihal != nil {
		surat.Perihal = requestBody.Perihal
	} else {
		surat.Perihal = surat.Perihal
	}

	if requestBody.Pic != nil {
		surat.Pic = requestBody.Pic
	} else {
		surat.Pic = surat.Pic
	}

	if requestBody.CreateBy != "" {
		surat.CreateBy = requestBody.CreateBy
	} else {
		surat.CreateBy = surat.CreateBy
	}

	initializers.DB.Save(&surat)

	c.JSON(200, gin.H{
		"surat": surat,
	})
}

func SuratDelete(c *gin.Context) {

	id := c.Params.ByName("id")

	var surat models.Surat

	if err := initializers.DB.First(&surat, id); err.Error != nil {
		c.JSON(404, gin.H{"error": "Surat not found"})
		return
	}

	if err := initializers.DB.Delete(&surat).Error; err != nil {
		c.JSON(400, gin.H{"error": "Failed to delete Surat: " + err.Error()})
		return
	}

	c.Status(204)

}

func exportSuratToExcel(surats []models.Surat) (*excelize.File, error) {
	// Buat file Excel baru
	f := excelize.NewFile()

	sheetNames := []string{"MEMO", "BERITA ACARA", "SK", "SURAT", "PROJECT", "PERDIN", "SURAT MASUK", "SURAT KELUAR", "ARSIP", "MEETING", "MEETING SCHEDULE"}

	for _, sheetName := range sheetNames {
		f.NewSheet(sheetName)
		if sheetName == "SURAT" {
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
	for _, surat := range surats {
		// Pastikan untuk dereferensikan pointer jika tidak nil
		var tanggal, nosurat, perihal, pic string
		if surat.Tanggal != nil {
			tanggal = surat.Tanggal.Format("2006-01-02") // Format tanggal sesuai kebutuhan
		}
		if surat.NoSurat != nil {
			nosurat = *surat.NoSurat
		}
		if surat.Perihal != nil {
			perihal = *surat.Perihal
		}
		if surat.Pic != nil {
			pic = *surat.Pic
		}

		// Pisahkan NoMemo untuk mendapatkan tipe memo
		parts := strings.Split(*surat.NoSurat, "/")
		if len(parts) > 1 && parts[1] == "ITS-SAG" {
			// Isi kolom SAG di sebelah kiri
			f.SetCellValue("SURAT", fmt.Sprintf("A%d", rowSAG), tanggal)
			f.SetCellValue("SURAT", fmt.Sprintf("B%d", rowSAG), nosurat)
			f.SetCellValue("SURAT", fmt.Sprintf("C%d", rowSAG), perihal)
			f.SetCellValue("SURAT", fmt.Sprintf("D%d", rowSAG), pic)
			rowSAG++
		} else if len(parts) > 1 && parts[1] == "ITS-ISO" {
			// Isi kolom ISO di sebelah kanan
			f.SetCellValue("SURAT", fmt.Sprintf("F%d", rowISO), tanggal)
			f.SetCellValue("SURAT", fmt.Sprintf("G%d", rowISO), nosurat)
			f.SetCellValue("SURAT", fmt.Sprintf("H%d", rowISO), perihal)
			f.SetCellValue("SURAT", fmt.Sprintf("I%d", rowISO), pic)
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
	f.SetColWidth("SURAT", "A", "D", 20)
	f.SetColWidth("SURAT", "F", "I", 20)
	f.SetColWidth("SURAT", "E", "E", 2)
	for i := 2; i <= lastRow; i++ {
		f.SetRowHeight("SURAT", i, 30)
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
	err = f.SetCellStyle("SURAT", "E1", fmt.Sprintf("E%d", lastRow), styleLine)

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
	err = f.SetCellStyle("SURAT", "A1", fmt.Sprintf("D%d", lastRow), styleBorder)
	err = f.SetCellStyle("SURAT", "F1", fmt.Sprintf("I%d", lastRow), styleBorder)

	return f, nil
}

// Handler untuk melakukan export Excel dengan Gin
func ExportSuratHandler(c *gin.Context) {
	// Data memo contoh
	var surats []models.Surat
	initializers.DB.Find(&surats)

	// Buat file Excel
	f, err := exportSuratToExcel(surats)
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
