package main

import (
	"bjb-crud/controllers"
	perdincontrollers "bjb-crud/controllers/perdincontroller"
	projectcontrollers "bjb-crud/controllers/projectcontroller"
	suratmasukcontrollers "bjb-crud/controllers/suratMasukcontroller"
	suratkeluarcontrollers "bjb-crud/controllers/suratKeluarcontroller"
	"bjb-crud/initializers"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariable()
	initializers.ConnectToDB()
}

func main() {
	r := gin.Default()

	//Iso routes
	r.POST("/Iso", controllers.IsoCreate)
	r.PUT("/Iso/:id", controllers.IsoUpdate)
	r.DELETE("/Iso/:id", controllers.IsoDelete)
	r.GET("/Iso", controllers.IsoIndex)
	r.GET("/Iso/:id", controllers.IsoShow)

	//Surat routes
	r.POST("/Surat", controllers.SuratCreate)
	r.PUT("/Surat/:id", controllers.SuratUpdate)
	r.GET("/Surat", controllers.SuratIndex)
	r.DELETE("/Surat/:id", controllers.SuratDelete)
	r.GET("/Surat/:id", controllers.SuratShow)

	//BeritaAcara routes
	r.POST("/BeritaAcara", controllers.BeritaAcaraCreate)
	r.PUT("/BeritaAcara/:id", controllers.BeritaAcaraUpdate)
	r.GET("/BeritaAcara", controllers.BeritaAcaraIndex)
	r.DELETE("/BeritaAcara/:id", controllers.BeritaAcaraDelete)
	r.GET("/BeritaAcara/:id", controllers.BeritaAcaraShow)

	//Sk routes
	r.POST("/Sk", controllers.SkCreate)
	r.PUT("/Sk/:id", controllers.SkUpdate)
	r.GET("/Sk", controllers.SkIndex)
	r.DELETE("/Sk/:id", controllers.SkDelete)
	r.GET("/Sk/:id", controllers.SkShow)

	//Project routes
	r.POST("/Project", projectcontrollers.ProjectCreate)
	r.PUT("/Project/:id", projectcontrollers.ProjectUpdate)
	r.GET("/Project", projectcontrollers.ProjectIndex)
	r.DELETE("/Project/:id", projectcontrollers.ProjectDelete)
	r.GET("/Project/:id", projectcontrollers.ProjectShow)

	//Perdin routes
	r.POST("/Perdin", perdincontrollers.PerdinCreate)
	r.PUT("/Perdin/:id", perdincontrollers.PerdinUpdate)
	r.GET("/Perdin", perdincontrollers.PerdinIndex)
	r.DELETE("/Perdin/:id", perdincontrollers.PerdinDelete)
	r.GET("/Perdin/:id", perdincontrollers.PerdinShow)

	//Surat  Masuk routes
	r.POST("/SuratMasuk", suratmasukcontrollers.SuratMasukCreate)
	r.PUT("/SuratMasuk/:id", suratmasukcontrollers.SuratMasukUpdate)
	r.GET("/SuratMasuk", suratmasukcontrollers.SuratMasukIndex)
	r.DELETE("/SuratMasuk/:id", suratmasukcontrollers.SuratMasukDelete)
	r.GET("/SuratMasuk/:id", suratmasukcontrollers.SuratMasukShow)

	//Surat  Keluar routes
	r.POST("/SuratKeluar", suratkeluarcontrollers.SuratKeluarCreate)
	r.PUT("/SuratKeluar/:id", suratkeluarcontrollers.SuratKeluarUpdate)
	r.GET("/SuratKeluar", suratkeluarcontrollers.SuratKeluarIndex)
	r.DELETE("/SuratKeluar/:id", suratkeluarcontrollers.SuratKeluarDelete)
	r.GET("/SuratKeluar/:id", suratkeluarcontrollers.SuratKeluarShow)

	r.Run() // listen and serve on 0.0.0.0:8080
}
