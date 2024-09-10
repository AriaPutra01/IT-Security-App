package main

import (
	"project-its/controllers"
	"project-its/initializers"

	// "project-its/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func init() {

	initializers.LoadEnvVariables()
	initializers.ConnectToDB()

}

func main() {

	r := gin.Default()

	// Enable CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:8000"
		},
		MaxAge: 12 * time.Hour,
	}))

	// Route yang tidak memerlukan autentikasi
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	// Terapkan middleware autentikasi ke semua route selanjutnya
	// r.Use(middleware.TokenAuthMiddleware())

	// Routes for User
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.GET("/exportAll", controllers.ExportAllSheets)

	// Setup session store
	store = cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	//logout must be after middleware
	r.POST("/logout", controllers.Logout)

	r.GET("/user", controllers.UserIndex)
	r.PUT("/user/:id", controllers.UserUpdate)
	r.DELETE("/user/:id", controllers.UserDelete)

	// Routes for Memo
	r.GET("/memos", controllers.MemoIndex)
	r.POST("/memos", controllers.MemoCreate)
	r.GET("/memos/:id", controllers.MemoShow)
	r.PUT("/memos/:id", controllers.MemoUpdate)
	r.DELETE("/memos/:id", controllers.MemoDelete)
	r.GET("/exportMemo", controllers.CreateExcelMemo)
	r.POST("/uploadMemo", controllers.ImportExcelMemo)

	//Project routes
	r.POST("/Project", controllers.ProjectCreate)
	r.PUT("/Project/:id", controllers.ProjectUpdate)
	r.GET("/Project", controllers.ProjectIndex)
	r.DELETE("/Project/:id", controllers.ProjectDelete)
	r.GET("/exportProject", controllers.CreateExcelProject)
	r.POST("/uploadProject", controllers.ImportExcelProject)

	//Timeline routes
	r.GET("/timelines", controllers.GetEventsTimeline)
	r.POST("/timelines", controllers.CreateEventTimeline)
	r.DELETE("/timelines/:id", controllers.DeleteEventTimeline)
	r.GET("/resources", controllers.GetResourcesTimeline)
	r.POST("/resources", controllers.CreateResourceTimeline)
	r.DELETE("/resources/:id", controllers.DeleteResourceTimeline)

	//Booking Rapat routes
	r.GET("/booking-rapat", controllers.GetEventsBookingRapat)
	r.POST("/booking-rapat", controllers.CreateEventBookingRapat)
	r.DELETE("/booking-rapat/:id", controllers.DeleteEventBookingRapat)

	// Ruang Rapat routes
	r.GET("/ruang-rapat", controllers.GetEvents)
	r.POST("/ruang-rapat", controllers.CreateEvent)
	r.DELETE("/ruang-rapat/:id", controllers.DeleteEvent)
	r.GET("/notifications", controllers.GetNotifications)
	r.DELETE("/notifications/:id", controllers.DeleteNotification)

	// Jadwal Cuti routes
	r.GET("/jadwal-cuti", controllers.GetEventsCuti)
	r.POST("/jadwal-cuti", controllers.CreateEventCuti)
	r.DELETE("/jadwal-cuti/:id", controllers.DeleteEventCuti)

	//Perdin routes
	r.POST("/Perdin", controllers.PerdinCreate)
	r.PUT("/Perdin/:id", controllers.PerdinUpdate)
	r.GET("/Perdin", controllers.PerdinIndex)
	r.DELETE("/Perdin/:id", controllers.PerdinDelete)
	r.GET("/Perdin/:id", controllers.PerdinShow)
	r.GET("/exportPerdin", controllers.CreateExcelPerdin)
	r.POST("/uploadPerdin", controllers.ImportExcelPerdin)

	//Surat  Masuk routes
	r.POST("/SuratMasuk", controllers.SuratMasukCreate)
	r.PUT("/SuratMasuk/:id", controllers.SuratMasukUpdate)
	r.GET("/SuratMasuk", controllers.SuratMasukIndex)
	r.DELETE("/SuratMasuk/:id", controllers.SuratMasukDelete)
	r.GET("/SuratMasuk/:id", controllers.SuratMasukShow)
	r.GET("/exportSuratMasuk", controllers.CreateExcelSuratMasuk)
	r.POST("/uploadSuratMasuk", controllers.ImportExcelSuratMasuk)

	//Surat  Keluar routes
	r.POST("/SuratKeluar", controllers.SuratKeluarCreate)
	r.PUT("/SuratKeluar/:id", controllers.SuratKeluarUpdate)
	r.GET("/SuratKeluar", controllers.SuratKeluarIndex)
	r.DELETE("/SuratKeluar/:id", controllers.SuratKeluarDelete)
	r.GET("/SuratKeluar/:id", controllers.SuratKeluarShow)
	r.GET("/exportSuratKeluar", controllers.CreateExcelSuratKeluar)
	r.POST("/uploadSuratKeluar", controllers.ImportExcelSuratKeluar)

	// Rute untuk upload file
	r.POST("/upload", controllers.UploadHandler)
	r.GET("/download/:id/:filename", controllers.DownloadFileHandler) // Ubah endpoint

	// Rute untuk hapus file
	r.DELETE("/delete/:id/:filename", controllers.DeleteFileHandler)

	// Routes for Arsip
	r.GET("/Arsip", controllers.ArsipIndex) // Tambahkan rute untuk membuat arsip
	r.POST("/Arsip", controllers.CreateArsip)
	r.PUT("/Arsip/:id", controllers.UpdateArsip)    // Tambahkan rute untuk memperbarui arsip
	r.DELETE("/Arsip/:id", controllers.DeleteArsip) // Tambahkan rute untuk menghapus arsip

	r.GET("/files/:id", controllers.GetFilesByID)

	r.Run()
}
