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
	var meetingLists []models.MeetingSchedule
	var meetings []models.Meeting
	var arsips []models.Arsip

	initializers.DB.Find(&memos)
	initializers.DB.Find(&projects)
	initializers.DB.Find(&perdins)
	initializers.DB.Find(&suratMasuks)
	initializers.DB.Find(&suratKeluars)
	initializers.DB.Find(&meetingLists)
	initializers.DB.Find(&meetings)
	initializers.DB.Find(&arsips)

	// Update data in each sheet
	for _, sheetName := range sheetNames {
		// Write header row
		switch sheetName {
		case "MEMO":
			f.SetCellValue(sheetName, "A1", "Tanggal")
			f.SetCellValue(sheetName, "B1", "No Memo")
			f.SetCellValue(sheetName, "C1", "Perihal")
			f.SetCellValue(sheetName, "D1", "Kategori")
			f.SetCellValue(sheetName, "E1", "PIC")
			f.SetColWidth(sheetName, "A", "E", 20)
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
			f.SetColWidth(sheetName, "A", "K", 20)
		case "PERDIN":
			f.SetCellValue(sheetName, "A1", "No Perdin")
			f.SetCellValue(sheetName, "B1", "Tanggal")
			f.SetCellValue(sheetName, "C1", "Deskripsi")
			f.MergeCell(sheetName, "C1", "D1")
			f.SetColWidth(sheetName, "A", "E", 20)
		case "SURAT MASUK":
			f.SetCellValue(sheetName, "A1", "No Surat")
			f.SetCellValue(sheetName, "B1", "Title")
			f.SetCellValue(sheetName, "C1", "Related Divisi")
			f.SetCellValue(sheetName, "D1", "Destiny Divisi")
			f.SetCellValue(sheetName, "E1", "Tanggal")
			f.SetColWidth(sheetName, "A", "E", 20)
		case "SURAT KELUAR":
			f.SetCellValue(sheetName, "A1", "No Surat")
			f.SetCellValue(sheetName, "B1", "Title")
			f.SetCellValue(sheetName, "C1", "From")
			f.SetCellValue(sheetName, "D1", "Pic")
			f.SetCellValue(sheetName, "E1", "Tanggal")
			f.SetColWidth(sheetName, "A", "E", 20)
		case "MEETING":
			f.SetCellValue(sheetName, "A1", "TASK")
			f.SetCellValue(sheetName, "B1", "TINDAK LANJUT")
			f.SetCellValue(sheetName, "C1", "STATUS")
			f.SetCellValue(sheetName, "D1", "UPDATE PENGERJAAN")
			f.SetCellValue(sheetName, "E1", "PIC")
			f.SetCellValue(sheetName, "F1", "TANGGAL TARGET")
			f.SetCellValue(sheetName, "G1", "TANGGAL ACTUAL")
			f.SetColWidth(sheetName, "A", "G", 20)
		case "MEETING SCHEDULE":
			f.SetCellValue(sheetName, "A1", "Hari")
			f.SetCellValue(sheetName, "B1", "Tanggal")
			f.SetCellValue(sheetName, "C1", "Perihal")
			f.SetCellValue(sheetName, "D1", "Waktu")
			f.SetCellValue(sheetName, "E1", "Selesai")
			f.SetCellValue(sheetName, "F1", "Tempat")
			f.SetCellValue(sheetName, "G1", "Status")
			f.SetCellValue(sheetName, "H1", "PIC")
			f.SetColWidth(sheetName, "A", "Z", 20)
		case "ARSIP":
			f.SetCellValue(sheetName, "A1", "No Arsip")
			f.SetCellValue(sheetName, "B1", "Jenis Dokumen")
			f.SetCellValue(sheetName, "C1", "No Dokumen")
			f.SetCellValue(sheetName, "D1", "Perihal")
			f.SetCellValue(sheetName, "E1", "No Box")
			f.SetCellValue(sheetName, "F1", "Keterangan")
			f.SetCellValue(sheetName, "G1", "Tanggal Dokumen")
			f.SetCellValue(sheetName, "H1", "Tanggal Penyerahan")
			f.SetColWidth(sheetName, "A", "Z", 20)
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
		case "MEETING":
			var meetings []models.Meeting
			initializers.DB.Find(&meetings)
			for _, meeting := range meetings {
				dataRows = append(dataRows, meeting)
			}
		case "MEETING SCHEDULE":
			var meetingLists []models.MeetingSchedule
			initializers.DB.Find(&meetingLists)
			for _, meetingList := range meetingLists {
				dataRows = append(dataRows, meetingList)
			}
		case "ARSIP":
			var arsips []models.Arsip
			initializers.DB.Find(&arsips)
			for _, arsip := range arsips {
				dataRows = append(dataRows, arsip)
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
				// f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), memo.Kategori)
				f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), memo.Pic)
			case "PROJECT":
				project := dataRow.(models.Project)
				f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), project.KodeProject)
				f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), project.JenisPengadaan)
				f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), project.NamaPengadaan)
				f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), project.DivInisiasi)
				f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), project.Bulan.Format("2006-01-02"))
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
			case "MEETING":
				meeting := dataRow.(models.Meeting)
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
			case "MEETING SCHEDULE":
				meetingList := dataRow.(models.MeetingSchedule)
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
			case "ARSIP":
				arsip := dataRow.(models.Arsip)
				f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), arsip.NoArsip)
				f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), arsip.JenisDokumen)
				f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), arsip.NoDokumen)
				f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), arsip.Perihal)
				f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), arsip.NoBox)
				f.SetCellValue(sheetName, fmt.Sprintf("F%d", rowNum), arsip.Keterangan)
				f.SetCellValue(sheetName, fmt.Sprintf("G%d", rowNum), arsip.TanggalDokumen.Format("2006-01-02"))
				f.SetCellValue(sheetName, fmt.Sprintf("H%d", rowNum), arsip.TanggalPenyerahan.Format("2006-01-02"))
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

/* INI FUNCTION COBA COCA */

/******************************** INI FUNCTION UPDATE ALL SHEET ***************************************/

func UpdateAllSheets(c *gin.Context) {
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

	// Define sheet names
	sheetNames := []string{"MEMO", "PROJECT", "PERDIN", "SURAT MASUK", "SURAT KELUAR", "ARSIP", "MEETING", "MEETING SCHEDULE"}

	// Fetch updated data from the database
	var memos []models.Memo
	var projects []models.Project
	var perdins []models.Perdin
	var suratMasuks []models.SuratMasuk
	var suratKeluars []models.SuratKeluar
	var meetingLists []models.MeetingSchedule
	var meetings []models.Meeting
	var arsips []models.Arsip

	initializers.DB.Find(&memos)
	initializers.DB.Find(&projects)
	initializers.DB.Find(&perdins)
	initializers.DB.Find(&suratMasuks)
	initializers.DB.Find(&suratKeluars)
	initializers.DB.Find(&meetingLists)
	initializers.DB.Find(&meetings)
	initializers.DB.Find(&arsips)

	for _, sheetName := range sheetNames {
		// Write header row
		switch sheetName {
		case "MEMO":
			f.SetCellValue(sheetName, "A1", "Tanggal")
			f.SetCellValue(sheetName, "B1", "No Memo")
			f.SetCellValue(sheetName, "C1", "Perihal")
			f.SetCellValue(sheetName, "D1", "Kategori")
			f.SetCellValue(sheetName, "E1", "PIC")

			f.SetColWidth(sheetName, "A", "E", 20)
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

			f.SetColWidth(sheetName, "A", "K", 20)
		case "PERDIN":
			f.SetCellValue(sheetName, "A1", "No Perdin")
			f.SetCellValue(sheetName, "B1", "Tanggal")
			f.SetCellValue(sheetName, "C1", "Deskripsi")
			f.MergeCell(sheetName, "C1", "D1")

			f.SetColWidth(sheetName, "A", "D", 20)
		case "SURAT MASUK":
			f.SetCellValue(sheetName, "A1", "No Surat")
			f.SetCellValue(sheetName, "B1", "Title")
			f.SetCellValue(sheetName, "C1", "Related Divisi")
			f.SetCellValue(sheetName, "D1", "Destiny Divisi")
			f.SetCellValue(sheetName, "E1", "Tanggal")

			f.SetColWidth(sheetName, "A", "E", 20)
		case "SURAT KELUAR":
			f.SetCellValue(sheetName, "A1", "No Surat")
			f.SetCellValue(sheetName, "B1", "Title")
			f.SetCellValue(sheetName, "C1", "From")
			f.SetCellValue(sheetName, "D1", "Pic")
			f.SetCellValue(sheetName, "E1", "Tanggal")

			f.SetColWidth(sheetName, "A", "E", 20)
		case "ARSIP":
			f.SetCellValue(sheetName, "A1", "No Arsip")
			f.SetCellValue(sheetName, "B1", "Jenis Dokumen")
			f.SetCellValue(sheetName, "C1", "No Dokumen")
			f.SetCellValue(sheetName, "D1", "Perihal")
			f.SetCellValue(sheetName, "E1", "No Box")
			f.SetCellValue(sheetName, "F1", "Keterangan")
			f.SetCellValue(sheetName, "G1", "Tanggal Dokumen")
			f.SetCellValue(sheetName, "H1", "Tanggal Penyerahan")

			f.SetColWidth(sheetName, "A", "H", 20)
		case "MEETING":
			f.SetCellValue(sheetName, "A1", "TASK")
			f.SetCellValue(sheetName, "B1", "TINDAK LANJUT")
			f.SetCellValue(sheetName, "C1", "STATUS")
			f.SetCellValue(sheetName, "D1", "UPDATE PENGERJAAN")
			f.SetCellValue(sheetName, "E1", "PIC")
			f.SetCellValue(sheetName, "F1", "TANGGAL TARGET")
			f.SetCellValue(sheetName, "G1", "TANGGAL ACTUAL")

			f.SetColWidth(sheetName, "A", "G", 20)
		case "MEETING SCHEDULE":
			f.SetCellValue(sheetName, "A1", "Hari")
			f.SetCellValue(sheetName, "B1", "Tanggal")
			f.SetCellValue(sheetName, "C1", "Perihal")
			f.SetCellValue(sheetName, "D1", "Waktu")
			f.SetCellValue(sheetName, "E1", "Selesai")
			f.SetCellValue(sheetName, "F1", "Tempat")
			f.SetCellValue(sheetName, "G1", "Status")
			f.SetCellValue(sheetName, "H1", "PIC")

			f.SetColWidth(sheetName, "A", "H", 20)
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
		case "ARSIP":
			for _, arsip := range arsips {
				dataRows = append(dataRows, arsip)
			}
		case "MEETING":
			for _, meeting := range meetings {
				dataRows = append(dataRows, meeting)
			}
		case "MEETING SCHEDULE":
			for _, meetingList := range meetingLists {
				dataRows = append(dataRows, meetingList)
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
				// f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), memo.Kategori)
				f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), memo.Pic)

				f.SetColWidth(sheetName, "A", "E", 20)
			case "PROJECT":
				project := dataRow.(models.Project)
				f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), project.KodeProject)
				f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), project.JenisPengadaan)
				f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), project.NamaPengadaan)
				f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), project.DivInisiasi)
				f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), project.Bulan.Format("02-01-2006")) // Write month as text
				f.SetCellValue(sheetName, fmt.Sprintf("F%d", rowNum), project.SumberPendanaan)
				f.SetCellValue(sheetName, fmt.Sprintf("G%d", rowNum), project.Anggaran)
				f.SetCellValue(sheetName, fmt.Sprintf("H%d", rowNum), project.NoIzin)
				f.SetCellValue(sheetName, fmt.Sprintf("I%d", rowNum), project.TanggalIzin.Format("02-01-2006"))
				f.SetCellValue(sheetName, fmt.Sprintf("J%d", rowNum), project.TanggalTor.Format("02-01-2006"))
				f.SetCellValue(sheetName, fmt.Sprintf("K%d", rowNum), project.Pic)

				f.SetColWidth(sheetName, "A", "K", 20)
			case "PERDIN":
				perdin := dataRow.(models.Perdin)
				f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), perdin.NoPerdin)
				f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), perdin.Tanggal.Format("2006-01-02"))
				f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), perdin.Hotel)
				f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), perdin.Transport)

				f.SetColWidth(sheetName, "A", "D", 20)
			case "SURAT MASUK":
				suratMasuk := dataRow.(models.SuratMasuk)
				f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), suratMasuk.NoSurat)
				f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), suratMasuk.Title)
				f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), suratMasuk.RelatedDiv)
				f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), suratMasuk.DestinyDiv)
				f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), suratMasuk.Tanggal.Format("2006-01-02"))

				f.SetColWidth(sheetName, "A", "E", 20)
			case "SURAT KELUAR":
				suratKeluar := dataRow.(models.SuratKeluar)
				f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), suratKeluar.NoSurat)
				f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), suratKeluar.Title)
				f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), suratKeluar.From)
				f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), suratKeluar.Pic)
				f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), suratKeluar.Tanggal.Format("2006-01-02"))

				f.SetColWidth(sheetName, "A", "E", 20)
			case "ARSIP":
				arsip := dataRow.(models.Arsip)
				f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), arsip.NoArsip)
				f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), arsip.JenisDokumen)
				f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), arsip.NoDokumen)
				f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), arsip.Perihal)
				f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), arsip.NoBox)
				f.SetCellValue(sheetName, fmt.Sprintf("F%d", rowNum), arsip.Keterangan)
				f.SetCellValue(sheetName, fmt.Sprintf("G%d", rowNum), arsip.TanggalDokumen.Format("2006-01-02"))
				f.SetCellValue(sheetName, fmt.Sprintf("H%d", rowNum), arsip.TanggalPenyerahan.Format("2006-01-02"))

				f.SetColWidth(sheetName, "A", "H", 20)
			case "MEETING":
				meeting := dataRow.(models.Meeting)
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

				f.SetColWidth(sheetName, "A", "G", 20)
			case "MEETING SCHEDULE":
				meetingList := dataRow.(models.MeetingSchedule)
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

				f.SetColWidth(sheetName, "A", "H", 20)
			}
		}
	}

	// Save the file with updated data
	if err := f.SaveAs(filePath); err != nil {
		c.String(http.StatusInternalServerError, "Error menyimpan file: %v", err)
		return
	}
}