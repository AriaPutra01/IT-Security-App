package main

import (
	"project-gin/controllers"
	"project-gin/initializers"

	"github.com/gin-gonic/gin"
)

func init() {

	initializers.LoadEnvVariables()
	initializers.ConnectToDB()

}

func main() {

	r := gin.Default()

	// Routes for SAG
	r.GET("/sag", controllers.SagIndex)
	r.POST("/sag", controllers.CreateSag)
	r.GET("/sag/:id", controllers.PostsShow)
	r.PUT("/sag/:id", controllers.PostsUpdate)
	r.DELETE("/sag/:id", controllers.PostsDelete)

	// Routes for Memo
	r.GET("/memos", controllers.MemoIndex)
	r.POST("/memos", controllers.MemoCreate)
	r.GET("/memos/:id", controllers.MemoShow)
	r.PUT("/memos/:id", controllers.MemoUpdate)
	r.DELETE("/memos/:id", controllers.MemoDelete)
	r.GET("/export", controllers.ExportExcel)

	r.Run()
}
