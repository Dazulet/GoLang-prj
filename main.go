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

	r.GET("/manga", handlers.GetMangas)
	r.GET("/manga/:id", handlers.GetManga)
	r.POST("/manga", handlers.CreateManga)
	r.PUT("/manga/:id", handlers.UpdateManga)
	r.DELETE("/manga/:id", handlers.DeleteManga)

	r.POST("/chapters", handlers.CreateChapter)
	r.GET("/manga/:id/chapters", handlers.GetChapters)

	r.POST("/users", handlers.CreateUser)
	r.POST("/bookmarks", handlers.AddBookmark)
	r.GET("/users/:id/bookmarks", handlers.GetUserBookmarks)

	// 4. Запуск
	r.Run(":8080")
}
