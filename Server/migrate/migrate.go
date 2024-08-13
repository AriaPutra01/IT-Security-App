package main

import (
	"project-gin/initializers"
	"project-gin/models"
)

func init() {

	initializers.LoadEnvVariables()
	initializers.ConnectToDB()

}

func main() {

	initializers.DB.AutoMigrate(&models.Sag{}, &models.Memo{})

}
