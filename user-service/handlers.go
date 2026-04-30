package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("mangalib_secret_key") // Секретный ключ для подписи

func Register(c *gin.Context) {
	var user User
	c.ShouldBindJSON(&user)

	result := DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Не удалось создать пользователя",
			"details": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, user)
}
func Login(c *gin.Context) {
	var input User
	c.ShouldBindJSON(&input)

	var user User
	// Ищем пользователя по email и паролю
	if err := DB.Where("email = ? AND password = ?", input.Email, input.Password).First(&user).Error; err != nil {
		c.JSON(401, gin.H{"error": "Неверный email или пароль"})
		return
	}

	// Создаем JWT токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Токен на 24 часа
	})

	tokenString, _ := token.SignedString(jwtSecret)
	c.JSON(200, gin.H{"token": tokenString})
}

func CreateUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании пользователя"})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func GetUser(c *gin.Context) {
	var user User
	id := c.Param("id")
	if err := DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}
	c.JSON(http.StatusOK, user)
}
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	DB.Delete(&User{}, id)
	c.JSON(200, gin.H{"message": "Пользователь удален"})
}
func (User) TableName() string {
	return "users"
}
