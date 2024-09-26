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
		&models.BeritaAcara{},
		&models.Surat{},
		&models.Sk{},
		&models.Project{},
		&models.SuratMasuk{},
		&models.SuratKeluar{},
		&models.Perdin{},
		&models.TimelineProject{},
		&models.ResourceProject{},
		&models.TimelineDesktop{},
		&models.ResourceDesktop{},
		&models.BookingRapat{},
		&models.JadwalRapat{},
		&models.JadwalCuti{},
		&models.Notification{},
		&models.Arsip{},
		&models.Meeting{},
		&models.MeetingSchedule{},
	)

}
