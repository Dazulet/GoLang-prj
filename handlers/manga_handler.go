package handlers

import (
	"MangaLIb/config"
	"MangaLIb/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetMangas(c *gin.Context) {
	var mangas []models.Manga
	config.DB.Find(&mangas)
	c.JSON(http.StatusOK, mangas)
}

func GetManga(c *gin.Context) {
	var manga models.Manga
	id := c.Param("id")
	if err := config.DB.Preload("Chapters").First(&manga, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Манга не найдена"})
		return
	}
	c.JSON(http.StatusOK, manga)
}

func CreateManga(c *gin.Context) {
	var manga models.Manga
	if err := c.ShouldBindJSON(&manga); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Create(&manga)
	c.JSON(http.StatusCreated, manga)
}

func UpdateManga(c *gin.Context) {
	var manga models.Manga
	id := c.Param("id")
	if err := config.DB.First(&manga, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Не найдена"})
		return
	}
	c.ShouldBindJSON(&manga)
	config.DB.Save(&manga)
	c.JSON(http.StatusOK, manga)
}

func DeleteManga(c *gin.Context) {
	id := c.Param("id")
	config.DB.Delete(&models.Manga{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "Удалено"})
}

func CreateChapter(c *gin.Context) {
	var chapter models.Chapter
	if err := c.ShouldBindJSON(&chapter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Create(&chapter)
	c.JSON(http.StatusCreated, chapter)
}

func GetChapters(c *gin.Context) {
	var chapters []models.Chapter
	mangaID := c.Param("id")
	config.DB.Where("manga_id = ?", mangaID).Find(&chapters)
	c.JSON(http.StatusOK, chapters)
}

func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Create(&user)
	c.JSON(http.StatusCreated, user)
}

func AddBookmark(c *gin.Context) {
	var bookmark models.Bookmark
	c.ShouldBindJSON(&bookmark)
	config.DB.Create(&bookmark)
	c.JSON(http.StatusCreated, gin.H{"status": "success"})
}

func GetUserBookmarks(c *gin.Context) {
	userID := c.Param("id")
	var mangas []models.Manga
	config.DB.Table("mangas").
		Joins("JOIN bookmarks ON bookmarks.manga_id = mangas.id").
		Where("bookmarks.user_id = ?", userID).
		Scan(&mangas)
	c.JSON(http.StatusOK, mangas)
}
