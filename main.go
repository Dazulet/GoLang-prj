package main

import (
	"MangaLIb/config"
	"MangaLIb/handlers"
	"MangaLIb/models"

	"github.com/gin-gonic/gin"
)

func main() {
	config.Connect()

	config.DB.AutoMigrate(&models.Manga{}, &models.Chapter{}, &models.User{}, &models.Bookmark{})

	r := gin.Default()

	r.POST("/users", handlers.CreateUser)
	r.POST("/login", handlers.Login)

	protected := r.Group("/")
	protected.Use(handlers.AuthMiddleware())
	{
		protected.GET("/manga", handlers.GetMangas)
		protected.GET("/manga/:id", handlers.GetManga)
		protected.POST("/manga", handlers.CreateManga)
		protected.PUT("/manga/:id", handlers.UpdateManga)
		protected.DELETE("/manga/:id", handlers.DeleteManga)

		protected.POST("/chapters", handlers.CreateChapter)
		protected.GET("/manga/:id/chapters", handlers.GetChapters)

		protected.POST("/bookmarks", handlers.AddBookmark)
		protected.GET("/users/:id/bookmarks", handlers.GetUserBookmarks)
	}
	r.Run(":8080")

}
