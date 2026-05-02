package main

import (
	"fmt" // Добавь это, если подчеркивает fmt ниже
	"net/http"
	"order-service/clients"

	"github.com/gin-gonic/gin"
)

func CreateOrder(c *gin.Context) {
	var order Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
		return
	}

	userExists, _ := clients.CheckUserExists(fmt.Sprintf("%d", order.UserID))
	if !userExists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Такого пользователя не существует"})
		return
	}

	mangaExists, _ := clients.CheckMangaExists(fmt.Sprintf("%d", order.MangaID))
	if !mangaExists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Такой манги не существует"})
		return
	}

	order.Status = "created"
	DB.Create(&order)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Заказ успешно оформлен",
		"order":   order,
	})
}
