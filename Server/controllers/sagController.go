package controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"project-gin/initializers"
	"project-gin/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type SagRequestBody struct {
	Tanggal string `json:"tanggal"`
	NoMemo  string `json:"no_memo"`
	Perihal string `json:"perihal"`
	Pic     string `json:"pic"`
}

func CreateSag(c *gin.Context) {
	var requestBody SagRequestBody

	if err := c.BindJSON(&requestBody); err != nil {
		c.Status(400)
		c.Error(err) // log the error
		return
	}

	// Add some logging to see what's being received
	log.Println("Received request body:", requestBody)

	// Parse the date string
	tanggalString := requestBody.Tanggal
	tanggal, err := time.Parse("2006-01-02", tanggalString)
	if err != nil {
		c.Status(400)
		c.JSON(400, gin.H{"error": "Invalid date format: " + err.Error()})
		return
	}

	sag := models.Sag{
		Tanggal: tanggal,
		NoMemo:  requestBody.NoMemo,
		Perihal: requestBody.Perihal,
		Pic:     requestBody.Pic,
	}

	result := initializers.DB.Create(&sag)

	if result.Error != nil {
		c.Status(400)
		c.JSON(400, gin.H{"error": "Failed to create Sag: " + result.Error.Error()})
		return
	}

	c.JSON(201, gin.H{
		"sag": sag,
	})
}

func SagIndex(c *gin.Context) {

	// Get the Posts
	var posts []models.Sag
	initializers.DB.Find(&posts)

	t, err := template.ParseFiles("views/index.html")
	if err != nil {
		log.Fatal(err)
	}

	err = t.Execute(c.Writer, gin.H{
		"Sags": posts,
	})
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(200, gin.H{
		"posts": posts,
	})

}

func PostsShow(c *gin.Context) {

	// Get id off url
	id := c.Params.ByName("id")

	// Get the Posts
	var post models.Sag
	initializers.DB.First(&post, id)

	// Respond with them
	c.JSON(200, gin.H{
		"post": post,
	})

}

func PostsUpdate(c *gin.Context) {
	var requestBody SagRequestBody

	// Bind the request body to the requestBody struct
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Parse the date string from the request body
	tanggalString := requestBody.Tanggal
	tanggal, err := time.Parse("2006-01-02", tanggalString)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid date format: " + err.Error()})
		return
	}

	// Get the Id from the URL parameters
	id := c.Param("id")

	// Find the post we are updating
	var post models.Sag
	if err := initializers.DB.First(&post, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Post not found"})
		return
	}

	// Update the post
	if err := initializers.DB.Model(&post).Updates(models.Sag{
		Tanggal: tanggal,
		NoMemo:  requestBody.NoMemo,
		Perihal: requestBody.Perihal,
		Pic:     requestBody.Pic,
	}).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update post: " + err.Error()})
		return
	}

	// Respond with the updated post
	c.JSON(200, gin.H{
		"post": post,
	})
}

func PostsDelete(c *gin.Context) {
	// Get id

	id := c.Param("id")

	// Delete Posts

	initializers.DB.Delete(&models.Sag{}, id)

	// Response
	c.Status(200)
}

func ExportExcel(c *gin.Context) {
	fileName := "sag_report.xlsx"

	// Check if the file exists
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		// File does not exist, create a new file
		f := excelize.NewFile()
		sheetName := "Sheet1"
		f.NewSheet(sheetName)
		f.SetCellValue(sheetName, "A1", "Tanggal")
		f.SetCellValue(sheetName, "B1", "NoMemo")
		f.SetCellValue(sheetName, "C1", "Perihal")
		f.SetCellValue(sheetName, "D1", "Pic")

		// Save the newly created file
		if err := f.SaveAs(fileName); err != nil {
			c.String(http.StatusInternalServerError, "Error saving file: %v", err)
			return
		}
	}

	// Open the existing or newly created Excel file
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error opening file: %v", err)
		return
	}
	defer f.Close()

	sheetName := "Sheet1"

	// Get the last row with data
	rows, err := f.GetRows(sheetName)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error getting rows: %v", err)
		return
	}
	lastRow := len(rows) // The number of rows in the sheet

	// Fetch updated data from the database
	var sags []models.Sag
	initializers.DB.Find(&sags)

	// Write data to the Excel file
	for i, sag := range sags {
		tanggalString := sag.Tanggal.Format("2006-01-02")
		rowNum := lastRow + i + 1 // Start from the row after the last existing row
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), tanggalString)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), sag.NoMemo)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), sag.Perihal)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), sag.Pic)
	}

	// Save the file with updated data
	if err := f.SaveAs(fileName); err != nil {
		c.String(http.StatusInternalServerError, "Error saving file: %v", err)
		return
	}

	// Serve the file to the client
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.File(fileName)
}