package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"project-its/initializers"
	"project-its/models"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func ExportAllSheets(c *gin.Context) {
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

	// Create all sheets
	for _, sheetName := range sheetNames {
		f.NewSheet(sheetName)
	}

	// Fetch data from database
	var memos []models.Memo
	var projects []models.Project
	var perdins []models.Perdin
	var suratMasuks []models.SuratMasuk
	var suratKeluars []models.SuratKeluar
	initializers.DB.Find(&memos)
	initializers.DB.Find(&projects)
	initializers.DB.Find(&perdins)
	initializers.DB.Find(&suratMasuks)
	initializers.DB.Find(&suratKeluars)

	// Update data in each sheet
	for _, sheetName := range sheetNames {
		// Write header row
		switch sheetName {
		case "SAG":
			f.SetCellValue(sheetName, "A1", "Tanggal")
			f.SetCellValue(sheetName, "B1", "No Memo")
			f.SetCellValue(sheetName, "C1", "Perihal")
			f.SetCellValue(sheetName, "D1", "PIC")
		case "MEMO":
			f.SetCellValue(sheetName, "A1", "Tanggal")
			f.SetCellValue(sheetName, "B1", "No Memo")
			f.SetCellValue(sheetName, "C1", "Perihal")
			f.SetCellValue(sheetName, "D1", "PIC")
		case "ISO":
			f.SetCellValue(sheetName, "A1", "Tanggal")
			f.SetCellValue(sheetName, "B1", "No Memo")
			f.SetCellValue(sheetName, "C1", "Perihal")
			f.SetCellValue(sheetName, "D1", "PIC")
		case "SURAT":
			f.SetCellValue(sheetName, "A1", "Tanggal")
			f.SetCellValue(sheetName, "B1", "No Surat")
			f.SetCellValue(sheetName, "C1", "Perihal")
			f.SetCellValue(sheetName, "D1", "PIC")
		case "BERITA ACARA":
			f.SetCellValue(sheetName, "A1", "Tanggal")
			f.SetCellValue(sheetName, "B1", "No Surat")
			f.SetCellValue(sheetName, "C1", "Perihal")
			f.SetCellValue(sheetName, "D1", "PIC")
		case "SK":
			f.SetCellValue(sheetName, "A1", "Tanggal")
			f.SetCellValue(sheetName, "B1", "No Surat")
			f.SetCellValue(sheetName, "C1", "Perihal")
			f.SetCellValue(sheetName, "D1", "PIC")
		case "PROJECT":
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
		case "PERDIN":
			f.SetCellValue(sheetName, "A1", "No Perdin")
			f.SetCellValue(sheetName, "B1", "Tanggal")
			f.SetCellValue(sheetName, "C1", "Hotel")
			f.SetCellValue(sheetName, "D1", "Transport")
		case "SURAT MASUK":
			f.SetCellValue(sheetName, "A1", "No Surat")
			f.SetCellValue(sheetName, "B1", "Title")
			f.SetCellValue(sheetName, "C1", "Related Divisi")
			f.SetCellValue(sheetName, "D1", "Destiny Divisi")
			f.SetCellValue(sheetName, "E1", "Tanggal")
		case "SURAT KELUAR":
			f.SetCellValue(sheetName, "A1", "No Surat")
			f.SetCellValue(sheetName, "B1", "Title")
			f.SetCellValue(sheetName, "C1", "From")
			f.SetCellValue(sheetName, "D1", "Pic")
			f.SetCellValue(sheetName, "E1", "Tanggal")
		}

		// Write data rows
		var dataRows []interface{}
		switch sheetName {
		case "MEMO":
			for _, memo := range memos {
				dataRows = append(dataRows, memo)
			}
		case "PROJECT":
			for _, project := range projects {
				dataRows = append(dataRows, project)
			}
		case "PERDIN":
			for _, perdin := range perdins {
				dataRows = append(dataRows, perdin)
			}
		case "SURAT MASUK":
			for _, suratMasuk := range suratMasuks {
				dataRows = append(dataRows, suratMasuk)
			}
		case "SURAT KELUAR":
			for _, suratKeluar := range suratKeluars {
				dataRows = append(dataRows, suratKeluar)
			}
		}

		for i, dataRow := range dataRows {
			rowNum := i + 2 // Start from the second row (first row is header)
			switch sheetName {
			case "MEMO":
				memo := dataRow.(models.Memo)
				f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), memo.Tanggal.Format("2006-01-02"))
				f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), memo.NoMemo)
				f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), memo.Perihal)
				f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), memo.Pic)
			case "PROJECT":
				project := dataRow.(models.Project)
				f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), project.KodeProject)
				f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), project.JenisPengadaan)
				f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), project.NamaPengadaan)
				f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), project.DivInisiasi)
				f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), project.Bulan)
				f.SetCellValue(sheetName, fmt.Sprintf("F%d", rowNum), project.SumberPendanaan)
				f.SetCellValue(sheetName, fmt.Sprintf("G%d", rowNum), project.Anggaran)
				f.SetCellValue(sheetName, fmt.Sprintf("H%d", rowNum), project.NoIzin)
				f.SetCellValue(sheetName, fmt.Sprintf("I%d", rowNum), project.TanggalIzin.Format("2006-01-02"))
				f.SetCellValue(sheetName, fmt.Sprintf("J%d", rowNum), project.TanggalTor.Format("2006-01-02"))
				f.SetCellValue(sheetName, fmt.Sprintf("K%d", rowNum), project.Pic)
			case "PERDIN":
				perdin := dataRow.(models.Perdin)
				f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), perdin.NoPerdin)
				f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), perdin.Tanggal.Format("2006-01-02"))
				f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), perdin.Hotel)
				f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), perdin.Transport)
			case "SURAT MASUK":
				suratMasuk := dataRow.(models.SuratMasuk)
				f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), suratMasuk.NoSurat)
				f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), suratMasuk.Title)
				f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), suratMasuk.RelatedDiv)
				f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), suratMasuk.DestinyDiv)
				f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), suratMasuk.Tanggal.Format("2006-01-02"))
			case "SURAT KELUAR":
				suratKeluar := dataRow.(models.SuratKeluar)
				f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), suratKeluar.NoSurat)
				f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), suratKeluar.Title)
				f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), suratKeluar.From)
				f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), suratKeluar.Pic)
				f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), suratKeluar.Tanggal.Format("2006-01-02"))
			}
		}
	}

	// Delete the default "Sheet1" sheet
	if err := f.DeleteSheet("Sheet1"); err != nil {
		panic(err) // Handle error jika bukan error "sheet tidak ditemukan"
	}

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