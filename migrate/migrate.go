package main

import (
	"bjb-crud/initializers"
	"bjb-crud/models"
)

func init() {
	initializers.LoadEnvVariable()
	initializers.ConnectToDB()
}

func main() {
	initializers.DB.AutoMigrate(
		&models.Iso{}, 
		&models.Surat{}, 
		&models.BeritaAcara{}, 
		&models.Sk{},
		&models.Project{},
		&models.Perdin{},
		&models.SuratMasuk{},
		&models.SuratKeluar{},
		&models.User{},
	)
}
