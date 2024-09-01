package main

import (
	"project-its/initializers"
	"project-its/models"
)

func init() {

	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

}

func main() {

	initializers.DB.AutoMigrate(
		&models.User{},
		&models.UserToken{},
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
		&models.RuangRapat{},
		&models.Notification{},
		&models.JadwalCuti{},
		&models.Timeline{},
		&models.BookingRapat{},
	)

}
