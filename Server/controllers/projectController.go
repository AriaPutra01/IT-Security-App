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
	"strings"
	"time"
	"unicode"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/xuri/excelize/v2"
)

type ProjectRequest struct {
	ID              uint    `gorm:"primaryKey"`
	KodeProject     *string `json:"kode_project"`
	JenisPengadaan  *string `json:"jenis_pengadaan"`
	NamaPengadaan   *string `json:"nama_pengadaan"`
	DivInisiasi     *string `json:"div_inisiasi"`
	Bulan           *string `json:"bulan"`
	SumberPendanaan *string `json:"sumber_pendanaan"`
	Anggaran        *string `json:"anggaran"`
	NoIzin          *string `json:"no_izin"`
	TanggalIzin     *string `json:"tanggal_izin"`
	TanggalTor      *string `json:"tanggal_tor"`
	Pic             *string `json:"pic"`
	CreateBy        string  `json:"create_by"`
}

func init() {
	err := godotenv.Load() // Memuat file .env
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	accountNameProject = os.Getenv("ACCOUNT_NAME")             // Mengambil nilai dari .env
	accountKeyProject = os.Getenv("ACCOUNT_KEY")               // Mengambil nilai dari .env
	containerNameProject = os.Getenv("CONTAINER_NAME_PROJECT") // Mengambil nilai dari .env
}

// Tambahkan variabel global untuk menyimpan kredensial
var (
	accountNameProject   string
	accountKeyProject    string
	containerNameProject string
)

func getBlobServiceClientProject() azblob.ServiceURL {
	creds, err := azblob.NewSharedKeyCredential(accountNameProject, accountKeyProject)
	if err != nil {
		panic("Failed to create shared key credential: " + err.Error())
	}

	pipeline := azblob.NewPipeline(creds, azblob.PipelineOptions{})

	// Build the URL for the Azure Blob Storage account
	URL, err := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/", accountNameProject))
	if err != nil {
		log.Fatal("Invalid URL format")
	}

	// Create a ServiceURL object that wraps the URL and the pipeline
	serviceURL := azblob.NewServiceURL(*URL, pipeline)

	return serviceURL
}

func UploadHandlerProject(c *gin.Context) {
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
	containerURL := getBlobServiceClient().NewContainerURL(containerNameProject)
	blobURL := containerURL.NewBlockBlobURL(filename)

	_, err = azblob.UploadStreamToBlockBlob(context.TODO(), src, blobURL, azblob.UploadStreamToBlockBlobOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengunggah file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File berhasil diunggah"})
}

func GetFilesByIDProject(c *gin.Context) {
	id := c.Param("id") // Mendapatkan ID dari URL

	containerURL := getBlobServiceClient().NewContainerURL(containerNameProject)
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
func DeleteFileHandlerProject(c *gin.Context) {
	filename := c.Param("filename")
	id := c.Param("id")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Filename is required"})
		return
	}

	// Membuat path lengkap berdasarkan ID dan nama file
	fullPath := fmt.Sprintf("%s/%s", id, filename)

	containerURL := getBlobServiceClient().NewContainerURL(containerNameProject)
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
func DownloadFileHandlerProject(c *gin.Context) {
	id := c.Param("id") // Mendapatkan ID dari URL
	filename := c.Param("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Filename is required"})
		return
	}

	// Membuat path lengkap berdasarkan ID dan nama file
	fullPath := fmt.Sprintf("%s/%s", id, filename)

	containerURL := getBlobServiceClient().NewContainerURL(containerNameProject)
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

func ProjectCreate(c *gin.Context) {
	// Get data off req body
	var requestBody ProjectRequest

	if err := c.BindJSON(&requestBody); err != nil {
		c.Status(400)
		c.Error(err) // log the error
		return
	}

	// Add some logging to see what's being received
	log.Println("Received request body:", requestBody)

	requestBody.CreateBy = c.MustGet("username").(string)

	var bulan, tanggal_izin, tanggal_tor *time.Time // Deklarasi variabel tanggal sebagai pointer ke time.Time

	// Parse the date string only if it's not nil and not empty
	if requestBody.Bulan != nil && *requestBody.Bulan != "" {
		parsedBulan, err := time.Parse("2006-01-02", *requestBody.Bulan)
		if err != nil {
			log.Printf("Error parsing date: %v", err)
			c.JSON(400, gin.H{"error": "Invalid date format: " + err.Error()})
			return
		}
		bulan = &parsedBulan
	}

	if requestBody.TanggalIzin != nil && *requestBody.TanggalIzin != "" {
		parsedTanggalIzin, err := time.Parse("2006-01-02", *requestBody.TanggalIzin)
		if err != nil {
			log.Printf("Error parsing date: %v", err)
			c.JSON(400, gin.H{"error": "Invalid date format: " + err.Error()})
			return
		}
		tanggal_izin = &parsedTanggalIzin
	}

	if requestBody.TanggalTor != nil && *requestBody.TanggalTor != "" {
		parsedTanggalTor, err := time.Parse("2006-01-02", *requestBody.TanggalTor)
		if err != nil {
			log.Printf("Error parsing date: %v", err)
			c.JSON(400, gin.H{"error": "Invalid date format: " + err.Error()})
			return
		}
		tanggal_tor = &parsedTanggalTor
	}

	project := models.Project{
		KodeProject:     requestBody.KodeProject,
		JenisPengadaan:  requestBody.JenisPengadaan,
		NamaPengadaan:   requestBody.NamaPengadaan,
		DivInisiasi:     requestBody.DivInisiasi,
		Bulan:           bulan,
		SumberPendanaan: requestBody.SumberPendanaan,
		Anggaran:        requestBody.Anggaran,
		NoIzin:          requestBody.NoIzin,
		TanggalIzin:     tanggal_izin,
		TanggalTor:      tanggal_tor,
		Pic:             requestBody.Pic,
		CreateBy:        requestBody.CreateBy,
	}

	result := initializers.DB.Create(&project)

	if result.Error != nil {
		c.Status(400)
		return
	}

	// Return it
	c.JSON(200, gin.H{
		"project": project,
	})
}

func ProjectIndex(c *gin.Context) {

	// Get models from DB
	var project []models.Project
	initializers.DB.Find(&project)

	//Respond with them
	c.JSON(200, gin.H{
		"project": project,
	})
}

func ProjectShow(c *gin.Context) {

	//get id
	id := c.Params.ByName("id")
	// Get models from DB
	var project models.Project

	initializers.DB.First(&project, id)

	//Respond with them
	c.JSON(200, gin.H{
		"project": project,
	})
}

func ProjectUpdate(c *gin.Context) {

	var requestBody ProjectRequest

	if err := c.BindJSON(&requestBody); err != nil {
		c.Status(400)
		c.Error(err) // log the error
		return
	}

	id := c.Params.ByName("id")

	var project models.Project
	initializers.DB.First(&project, id)

	if err := initializers.DB.First(&project, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "project tidak ditemukan"})
		return
	}

	requestBody.CreateBy = c.MustGet("username").(string)
	project.CreateBy = requestBody.CreateBy

	if requestBody.Bulan != nil {
		bulan, err := time.Parse("2006-01-02", *requestBody.Bulan)
		if err != nil {
			c.JSON(400, gin.H{"error": "Format tanggal tidak valid: " + err.Error()})
			return
		}
		project.TanggalIzin = &bulan
	}

	if requestBody.TanggalIzin != nil {
		tanggal_izin, err := time.Parse("2006-01-02", *requestBody.TanggalIzin)
		if err != nil {
			c.JSON(400, gin.H{"error": "Format tanggal tidak valid: " + err.Error()})
			return
		}
		project.TanggalIzin = &tanggal_izin
	}

	if requestBody.TanggalTor != nil {
		tanggal_tor, err := time.Parse("2006-01-02", *requestBody.TanggalTor)
		if err != nil {
			c.JSON(400, gin.H{"error": "Format tanggal tidak valid: " + err.Error()})
			return
		}
		project.TanggalTor = &tanggal_tor
	}

	if requestBody.KodeProject != nil {
		project.KodeProject = requestBody.KodeProject
	} else {
		project.KodeProject = project.KodeProject // gunakan nilai yang ada dari database
	}

	if requestBody.JenisPengadaan != nil {
		project.JenisPengadaan = requestBody.JenisPengadaan
	} else {
		project.JenisPengadaan = project.JenisPengadaan // gunakan nilai yang ada dari database
	}

	if requestBody.NamaPengadaan != nil {
		project.NamaPengadaan = requestBody.NamaPengadaan
	} else {
		project.NamaPengadaan = project.NamaPengadaan // gunakan nilai yang ada dari database
	}

	if requestBody.DivInisiasi != nil {
		project.DivInisiasi = requestBody.DivInisiasi
	} else {
		project.DivInisiasi = project.DivInisiasi // gunakan nilai yang ada dari database
	}

	if requestBody.SumberPendanaan != nil {
		project.SumberPendanaan = requestBody.SumberPendanaan
	} else {
		project.SumberPendanaan = project.SumberPendanaan // gunakan nilai yang ada dari database
	}

	if requestBody.Anggaran != nil {
		project.Anggaran = requestBody.Anggaran
	} else {
		project.Anggaran = project.Anggaran // gunakan nilai yang ada dari database
	}

	if requestBody.NoIzin != nil {
		project.NoIzin = requestBody.NoIzin
	} else {
		project.NoIzin = project.NoIzin // gunakan nilai yang ada dari database
	}

	if requestBody.Pic != nil {
		project.Pic = requestBody.Pic
	} else {
		project.Pic = project.Pic // gunakan nilai yang ada dari database
	}

	if requestBody.CreateBy != "" {
		project.CreateBy = requestBody.CreateBy
	} else {
		project.CreateBy = project.CreateBy // gunakan nilai yang ada dari database
	}

	initializers.DB.Model(&project).Updates(project)

	c.JSON(200, gin.H{
		"project": project,
	})

}

func ProjectDelete(c *gin.Context) {

	//get id
	id := c.Params.ByName("id")

	// find the project
	var project models.Project

	if err := initializers.DB.First(&project, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Project not found"})
		return
	}

	/// delete it
	if err := initializers.DB.Delete(&project).Error; err != nil {
		c.JSON(404, gin.H{"error": "Project Failed to Delete"})
		return
	}

	c.JSON(200, gin.H{
		"project": "Deleted",
	})
}

func CreateExcelProject(c *gin.Context) {
	dir := "D:\\excel"
	baseFileName := "its_report"
	filePath := filepath.Join(dir, baseFileName+".xlsx")

	// Check if the file already exists
	if _, err := os.Stat(filePath); err == nil {
		// File exists, append "_new" to the file name
		baseFileName += "_new"
	}

	fileName := baseFileName + ".xlsx"

	// Create a new Excel file
	f := excelize.NewFile()

	// Define sheet names
	sheetNames := []string{"MEMO", "PROJECT", "PERDIN", "SURAT MASUK", "SURAT KELUAR", "ARSIP", "MEETING", "MEETING SCHEDULE"}

	// Create sheets and set headers
	for _, sheetName := range sheetNames {
		f.NewSheet(sheetName)
		if sheetName == "PROJECT" {
			f.SetCellValue(sheetName, "A1", "Kode Project")
			f.SetCellValue(sheetName, "B1", "Jenis Pengadaan")
			f.SetCellValue(sheetName, "C1", "Nama Pengadaan")
			f.SetCellValue(sheetName, "D1", "Divisi Inisiasi")
			f.SetCellValue(sheetName, "E1", "Bulan")
			f.SetCellValue(sheetName, "F1", "Sumber Pendanaan")
			f.SetCellValue(sheetName, "G1", "Anggaran")
			f.SetCellValue(sheetName, "H1", "No Izin")
			f.SetCellValue(sheetName, "I1", "Tgl Izin")
			f.SetCellValue(sheetName, "J1", "Tgl TOR")
			f.SetCellValue(sheetName, "K1", "Pic")

			f.SetColWidth(sheetName, "A", "K", 20)
		}
	}

	// Fetch initial data from the database
	var projects []models.Project
	initializers.DB.Find(&projects)

	// Write initial data to the "PROJECT" sheet
	projectSheetName := "PROJECT"
	for i, project := range projects {
		izinString := project.TanggalIzin.Format("02-01-2006")
		torString := project.TanggalTor.Format("02-01-2006")
		bulanString := project.Bulan.Format("02-01-2006")
		rowNum := i + 2 // Start from the second row (first row is header)

		// Ensure data is correctly written to cells
		f.SetCellValue(projectSheetName, fmt.Sprintf("A%d", rowNum), project.KodeProject)
		f.SetCellValue(projectSheetName, fmt.Sprintf("B%d", rowNum), project.JenisPengadaan)
		f.SetCellValue(projectSheetName, fmt.Sprintf("C%d", rowNum), project.NamaPengadaan)
		f.SetCellValue(projectSheetName, fmt.Sprintf("D%d", rowNum), project.DivInisiasi)
		f.SetCellValue(projectSheetName, fmt.Sprintf("E%d", rowNum), bulanString) // Ensure this is the correct format
		f.SetCellValue(projectSheetName, fmt.Sprintf("F%d", rowNum), project.SumberPendanaan)
		f.SetCellValue(projectSheetName, fmt.Sprintf("G%d", rowNum), project.Anggaran)
		f.SetCellValue(projectSheetName, fmt.Sprintf("H%d", rowNum), project.NoIzin)
		f.SetCellValue(projectSheetName, fmt.Sprintf("I%d", rowNum), izinString)
		f.SetCellValue(projectSheetName, fmt.Sprintf("J%d", rowNum), torString)
		f.SetCellValue(projectSheetName, fmt.Sprintf("K%d", rowNum), project.Pic)
	}

	// Delete the default "Sheet1" sheet if it exists
	if err := f.DeleteSheet("Sheet1"); err != nil {
		c.String(http.StatusInternalServerError, "Error deleting default sheet: %v", err)
		return
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

func UpdateSheetProject(c *gin.Context) {
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
	sheetName := "PROJECT"

	// Check if sheet exists and delete it if it does
	if _, err := f.GetSheetIndex(sheetName); err == nil {
		f.DeleteSheet(sheetName)
	}
	f.NewSheet(sheetName)

	// Write header row
	f.SetCellValue(sheetName, "A1", "Kode Project")
	f.SetCellValue(sheetName, "B1", "Jenis Pengadaan")
	f.SetCellValue(sheetName, "C1", "Nama Pengadaan")
	f.SetCellValue(sheetName, "D1", "Divisi Inisiasi")
	f.SetCellValue(sheetName, "E1", "Bulan")
	f.SetCellValue(sheetName, "F1", "Sumber Pendanaan")
	f.SetCellValue(sheetName, "G1", "Anggaran")
	f.SetCellValue(sheetName, "H1", "No Izin")
	f.SetCellValue(sheetName, "I1", "Tgl Izin")
	f.SetCellValue(sheetName, "J1", "Tgl TOR")
	f.SetCellValue(sheetName, "K1", "Pic")

	// Fetch updated data from the database
	var projects []models.Project
	initializers.DB.Find(&projects)

	// Write data rows
	for i, project := range projects {
		rowNum := i + 2 // Start from the second row (first row is header)

		// Convert date to string with specific format
		bulanString := project.Bulan.Format("02-01-2006")

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), project.KodeProject)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), project.JenisPengadaan)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), project.NamaPengadaan)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), project.DivInisiasi)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), bulanString) // Write month as text
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", rowNum), project.SumberPendanaan)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", rowNum), project.Anggaran)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", rowNum), project.NoIzin)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", rowNum), project.TanggalIzin.Format("02-01-2006"))
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", rowNum), project.TanggalTor.Format("02-01-2006"))
		f.SetCellValue(sheetName, fmt.Sprintf("K%d", rowNum), project.Pic)

		f.SetColWidth(sheetName, "A", "K", 20)
	}

	// Save the file with updated data
	if err := f.SaveAs(filePath); err != nil {
		c.String(http.StatusInternalServerError, "Error menyimpan file: %v", err)
		return
	}

	c.Redirect(http.StatusFound, "/Project")
}

func ImportExcelProject(c *gin.Context) {
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

	sheetName := "PROJECT"
	rows, err := f.GetRows(sheetName)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error getting rows: %v", err)
		return
	}

	for i, row := range rows {
		if i < 7 { // Skip header or initial rows if necessary
			continue
		}

		// Count non-empty columns
		nonEmptyCount := 0
		for _, cell := range row {
			if cell != "" {
				nonEmptyCount++
			}
		}

		// Skip rows with less than 3 non-empty columns
		if nonEmptyCount < 3 {
			log.Printf("Row %d skipped: less than 3 columns filled, only %d filled", i+1, nonEmptyCount)
			continue
		}

		// Membersihkan string anggaran dari karakter non-numerik
		rawAnggaran := getStringOrNil(getColumn(row, 7))
		var anggaranCleaned *string
		if rawAnggaran != nil {
			cleanedAnggaran := cleanNumericString(*rawAnggaran)
			anggaranCleaned = &cleanedAnggaran
		}

		project := models.Project{
			KodeProject:     getStringOrNil(getColumn(row, 1)),
			JenisPengadaan:  getStringOrNil(getColumn(row, 2)),
			NamaPengadaan:   getStringOrNil(getColumn(row, 3)),
			DivInisiasi:     getStringOrNil(getColumn(row, 4)),
			Bulan:           parseDateOrNil(getStringOrNil(getColumn(row, 5))),
			SumberPendanaan: getStringOrNil(getColumn(row, 6)),
			Anggaran:        anggaranCleaned,
			NoIzin:          getStringOrNil(getColumn(row, 8)),
			TanggalIzin:     parseDateOrNil(getStringOrNil(getColumn(row, 9))),
			TanggalTor:      parseDateOrNil(getStringOrNil(getColumn(row, 10))),
			Pic:             getStringOrNil(getColumn(row, 11)),
			CreateBy:        c.MustGet("username").(string),
		}

		if err := initializers.DB.Create(&project).Error; err != nil {
			log.Printf("Error saving record from row %d: %v", i+1, err)
			continue
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data imported successfully, check logs for any skipped rows."})
}

// Helper function to safely get column data or return empty if index is out of range
func getColumn(row []string, index int) string {
	if index >= len(row) {
		return ""
	}
	return row[index]
}

// Helper function to return nil if the string is empty
func getStringOrNil(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

// Helper function to parse date from various formats
func parseDate(dateStr string) (time.Time, error) {
	dateFormats := []string{
		"2 January 2006",
		"2006-01-02",
		"02-01-2006",
		"01-02-2006",
		"01/02/2006",
		"02/01/2006",
		"2006.01.02",
		"Jan-06",
		"02-January-2006",
		"02-Jan-06",
	}

	for _, format := range dateFormats {
		parsedDate, err := time.Parse(format, dateStr)
		if err == nil {
			return parsedDate, nil
		}
	}
	return time.Time{}, fmt.Errorf("no valid date format found")
}

// Fungsi untuk membersihkan string dari karakter non-numerik
func cleanNumericString(input string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsDigit(r) {
			return r
		}
		return -1
	}, input)
}

// Helper function to parse date or return nil if input is nil
func parseDateOrNil(dateStr *string) *time.Time {
	if dateStr == nil {
		return nil
	}
	parsedDate, err := parseDate(*dateStr)
	if err != nil {
		return nil
	}
	return &parsedDate
}

func exportProjectToExcel(projects []models.Project) (*excelize.File, error) {
	// Buat file Excel baru
	f := excelize.NewFile()

	sheetNames := []string{"MEMO", "BERITA ACARA", "SK", "SURAT", "PROJECT", "PERDIN", "SURAT MASUK", "SURAT KELUAR", "ARSIP", "MEETING", "MEETING SCHEDULE"}

	for _, sheetName := range sheetNames {
		f.NewSheet(sheetName)
		if sheetName == "PROJECT" {
			// Header untuk SAG (kolom kiri)
			f.SetCellValue(sheetName, "A1", "Kode Project")
			f.SetCellValue(sheetName, "B1", "Jenis Pengadaan")
			f.SetCellValue(sheetName, "C1", "Nama Pengadaan")
			f.SetCellValue(sheetName, "D1", "Divisi Inisiasi")
			f.SetCellValue(sheetName, "E1", "Bulan")
			f.SetCellValue(sheetName, "F1", "Sumber Pendanaan")
			f.SetCellValue(sheetName, "G1", "Anggaran")
			f.SetCellValue(sheetName, "H1", "No Izin")
			f.SetCellValue(sheetName, "I1", "Tgl Izin")
			f.SetCellValue(sheetName, "J1", "Tgl TOR")
			f.SetCellValue(sheetName, "K1", "Pic")

			// Header untuk ISO (kolom kanan)
			f.SetCellValue(sheetName, "A1", "Kode Project")
			f.SetCellValue(sheetName, "B1", "Jenis Pengadaan")
			f.SetCellValue(sheetName, "C1", "Nama Pengadaan")
			f.SetCellValue(sheetName, "D1", "Divisi Inisiasi")
			f.SetCellValue(sheetName, "E1", "Bulan")
			f.SetCellValue(sheetName, "F1", "Sumber Pendanaan")
			f.SetCellValue(sheetName, "G1", "Anggaran")
			f.SetCellValue(sheetName, "H1", "No Izin")
			f.SetCellValue(sheetName, "I1", "Tgl Izin")
			f.SetCellValue(sheetName, "J1", "Tgl TOR")
			f.SetCellValue(sheetName, "K1", "Pic")
		}
	}
	f.DeleteSheet("Sheet1")

	// Inisialisasi baris awal
	rowSAG := 2
	rowISO := 2

	// Loop melalui data memo
	for _, project := range projects {
		// Pastikan untuk dereferensikan pointer jika tidak nil
		var kodeProject, jenisPengadaan, namaPengadaan, divInisiasi, bulan, sumberPendanaan, anggaran, noIzin, tanggalIzin, tanggalTor, pic string
		if project.KodeProject != nil {
			kodeProject = *project.KodeProject
		}
		if project.JenisPengadaan != nil {
			jenisPengadaan = *project.JenisPengadaan
		}
		if project.NamaPengadaan != nil {
			namaPengadaan = *project.NamaPengadaan
		}
		if project.DivInisiasi != nil {
			divInisiasi = *project.DivInisiasi
		}
		if project.Bulan != nil {
			bulan = project.Bulan.Format("Jan-06") // Format tanggal sesuai kebutuhan
		}
		if project.SumberPendanaan != nil {
			sumberPendanaan = *project.SumberPendanaan
		}
		if project.Anggaran != nil {
			anggaran = *project.Anggaran
		}
		if project.NoIzin != nil {
			noIzin = *project.NoIzin
		}
		if project.TanggalIzin != nil {
			tanggalIzin = project.TanggalIzin.Format("2006-01-02")
		}
		if project.TanggalTor != nil {
			tanggalTor = project.TanggalTor.Format("2006-01-02")
		}
		if project.Pic != nil {
			pic = *project.Pic
		}

		// Pisahkan NoMemo untuk mendapatkan tipe memo
		parts := strings.Split(*project.KodeProject, "/")
		if len(parts) > 1 && parts[1] == "ITS-SAG" {
			// Isi kolom SAG di sebelah kiri
			f.SetCellValue("PROJECT", fmt.Sprintf("A%d", rowSAG), kodeProject)
			f.SetCellValue("PROJECT", fmt.Sprintf("B%d", rowSAG), jenisPengadaan)
			f.SetCellValue("PROJECT", fmt.Sprintf("C%d", rowSAG), namaPengadaan)
			f.SetCellValue("PROJECT", fmt.Sprintf("D%d", rowSAG), divInisiasi)
			f.SetCellValue("PROJECT", fmt.Sprintf("E%d", rowSAG), bulan)
			f.SetCellValue("PROJECT", fmt.Sprintf("F%d", rowSAG), sumberPendanaan)
			f.SetCellValue("PROJECT", fmt.Sprintf("G%d", rowSAG), anggaran)
			f.SetCellValue("PROJECT", fmt.Sprintf("H%d", rowSAG), noIzin)
			f.SetCellValue("PROJECT", fmt.Sprintf("I%d", rowSAG), tanggalIzin)
			f.SetCellValue("PROJECT", fmt.Sprintf("J%d", rowSAG), tanggalTor)
			f.SetCellValue("PROJECT", fmt.Sprintf("K%d", rowSAG), pic)
			rowSAG++
		} else if len(parts) > 1 && parts[1] == "ITS-ISO" {
			// Isi kolom ISO di sebelah kanan
			f.SetCellValue("PROJECT", fmt.Sprintf("L%d", rowSAG), kodeProject)
			f.SetCellValue("PROJECT", fmt.Sprintf("M%d", rowSAG), jenisPengadaan)
			f.SetCellValue("PROJECT", fmt.Sprintf("N%d", rowSAG), namaPengadaan)
			f.SetCellValue("PROJECT", fmt.Sprintf("O%d", rowSAG), divInisiasi)
			f.SetCellValue("PROJECT", fmt.Sprintf("P%d", rowSAG), bulan)
			f.SetCellValue("PROJECT", fmt.Sprintf("Q%d", rowSAG), sumberPendanaan)
			f.SetCellValue("PROJECT", fmt.Sprintf("R%d", rowSAG), anggaran)
			f.SetCellValue("PROJECT", fmt.Sprintf("S%d", rowSAG), noIzin)
			f.SetCellValue("PROJECT", fmt.Sprintf("T%d", rowSAG), tanggalIzin)
			f.SetCellValue("PROJECT", fmt.Sprintf("U%d", rowSAG), tanggalTor)
			f.SetCellValue("PROJECT", fmt.Sprintf("V%d", rowSAG), pic)
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
	f.SetColWidth("PROJECT", "A", "D", 20)
	f.SetColWidth("PROJECT", "F", "I", 20)
	f.SetColWidth("PROJECT", "E", "E", 2)
	for i := 2; i <= lastRow; i++ {
		f.SetRowHeight("PROJECT", i, 30)
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
	err = f.SetCellStyle("PROJECT", "E1", fmt.Sprintf("E%d", lastRow), styleLine)

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
	err = f.SetCellStyle("PROJECT", "A1", fmt.Sprintf("D%d", lastRow), styleBorder)
	err = f.SetCellStyle("PROJECT", "F1", fmt.Sprintf("I%d", lastRow), styleBorder)

	return f, nil
}

// Handler untuk melakukan export Excel dengan Gin
func ExportProjectHandler(c *gin.Context) {
	// Data memo contoh
	var projects []models.Project
	initializers.DB.Find(&projects)

	// Buat file Excel
	f, err := exportProjectToExcel(projects)
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
