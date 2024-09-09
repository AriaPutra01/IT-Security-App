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
		&models.Memo{},
		&models.Project{},
		&models.SuratMasuk{},
		&models.SuratKeluar{},
		&models.Perdin{},
		&models.RuangRapat{},
		&models.Notification{},
		&models.JadwalCuti{},
		&models.Timeline{},
		&models.ResourceTimeline{},
		&models.BookingRapat{},
		&models.Arsip{},
	)

}
