package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"project-gin/initializers"
	"project-gin/models"

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
	sheetNames := []string{"SAG", "MEMO", "ISO", "SURAT"}

	var sags []models.Sag
	var memos []models.Memo
	var isos []models.Iso
	initializers.DB.Find(&sags)
	initializers.DB.Find(&memos)
	initializers.DB.Find(&isos)

	// Update data in each sheet
	for _, sheetName := range sheetNames {
		// Clear the sheet
		f.DeleteSheet(sheetName)
		f.NewSheet(sheetName)

		// Write header row
		f.SetCellValue(sheetName, "A1", "Tanggal")
		f.SetCellValue(sheetName, "B1", "No Memo")
		f.SetCellValue(sheetName, "C1", "Perihal")
		f.SetCellValue(sheetName, "D1", "PIC")

		var dataRows []interface{}
		switch sheetName {
		case "SAG":
			for _, sag := range sags {
				dataRows = append(dataRows, sag)
			}
		case "MEMO":
			for _, memo := range memos {
				dataRows = append(dataRows, memo)
			}
		case "ISO":
			for _, iso := range isos {
				dataRows = append(dataRows, iso)
			}
		}

		for i, dataRow := range dataRows {
			rowNum := i + 2 // Start from the second row (first row is header)
			switch sheetName {
			case "SAG":
				sag := dataRow.(models.Sag)
				f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), sag.Tanggal.Format("2006-01-02"))
				f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), sag.NoMemo)
				f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), sag.Perihal)
				f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), sag.Pic)
			case "MEMO":
				memo := dataRow.(models.Memo)
				f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), memo.Tanggal.Format("2006-01-02"))
				f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), memo.NoMemo)
				f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), memo.Perihal)
				f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), memo.Pic)
			case "ISO":
				iso := dataRow.(models.Iso)
				f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), iso.Tanggal.Format("2006-01-02"))
				f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), iso.NoMemo)
				f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), iso.Perihal)
				f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), iso.Pic)
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

func UpdateAllExcelSheets(c *gin.Context) {
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

	// Define sheet names
	sheetNames := []string{"SAG", "MEMO", "ISO"}

	// Fetch updated data from the database
	var sags []models.Sag
	var memos []models.Memo
	var isos []models.Iso
	initializers.DB.Find(&sags)
	initializers.DB.Find(&memos)
	initializers.DB.Find(&isos)

	// Update data in each sheet
	for _, sheetName := range sheetNames {
		// Clear the sheet
		f.DeleteSheet(sheetName)
		f.NewSheet(sheetName)

		// Write header row
		f.SetCellValue(sheetName, "A1", "Tanggal")
		f.SetCellValue(sheetName, "B1", "No Memo")
		f.SetCellValue(sheetName, "C1", "Perihal")
		f.SetCellValue(sheetName, "D1", "PIC")

		var dataRows []interface{}
		switch sheetName {
		case "SAG":
			for _, sag := range sags {
				dataRows = append(dataRows, sag)
			}
		case "MEMO":
			for _, memo := range memos {
				dataRows = append(dataRows, memo)
			}
		case "ISO":
			for _, iso := range isos {
				dataRows = append(dataRows, iso)
			}
		}

		for i, dataRow := range dataRows {
			rowNum := i + 2 // Start from the second row (first row is header)
			switch sheetName {
			case "SAG":
				sag := dataRow.(models.Sag)
				f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), sag.Tanggal.Format("2006-01-02"))
				f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), sag.NoMemo)
				f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), sag.Perihal)
				f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), sag.Pic)
			case "MEMO":
				memo := dataRow.(models.Memo)
				f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), memo.Tanggal.Format("2006-01-02"))
				f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), memo.NoMemo)
				f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), memo.Perihal)
				f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), memo.Pic)
			case "ISO":
				iso := dataRow.(models.Iso)
				f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), iso.Tanggal.Format("2006-01-02"))
				f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), iso.NoMemo)
				f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), iso.Perihal)
				f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), iso.Pic)
			}
		}
	}

	// Save the file with updated data
	file, err := os.OpenFile(filePath, os.O_RDWR, 0755)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error membuka file: %v", err)
		return
	}
	defer file.Close()

	if _, err := f.WriteTo(file); err != nil {
		c.String(http.StatusInternalServerError, "Error menyimpan file: %v", err)
		return
	}

	c.Redirect(http.StatusFound, "/iso")
}
