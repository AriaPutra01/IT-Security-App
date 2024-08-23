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

	initializers.DB.AutoMigrate(
		&models.Sag{},
		&models.Memo{},
		&models.Iso{},
		&models.Project{},
		&models.Surat{},
		&models.BeritaAcara{},
		&models.SuratMasuk{},
		&models.SuratKeluar{},
		&models.Sk{},
		&models.Perdin{},
		&models.User{},
	)

}
