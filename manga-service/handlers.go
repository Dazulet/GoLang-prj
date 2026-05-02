package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetMangas(c *gin.Context) {
	var mangas []Manga
	DB.Find(&mangas)
	c.JSON(http.StatusOK, mangas)
}

func GetManga(c *gin.Context) {
	var manga Manga
	id := c.Param("id")
	if err := DB.Preload("Chapters").First(&manga, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Манга не найдена"})
		return
	}
	c.JSON(http.StatusOK, manga)
}

func CreateManga(c *gin.Context) {
	var manga Manga
	if err := c.ShouldBindJSON(&manga); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := DB.Create(&manga)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Ошибка базы данных",
			"details": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, manga)
}
func UpdateManga(c *gin.Context) {
	var manga Manga
	id := c.Param("id")
	if err := DB.First(&manga, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Не найдена"})
		return
	}
	c.ShouldBindJSON(&manga)
	DB.Save(&manga)
	c.JSON(200, manga)
}

func DeleteManga(c *gin.Context) {
	id := c.Param("id")
	DB.Delete(&Manga{}, id)
	c.JSON(200, gin.H{"message": "Манга удалена"})
}
func (Manga) TableName() string {
	return "mangas"
}

func (Chapter) TableName() string {
	return "chapters"
}
