package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var jwtSecret = []byte("mangalib_secret_key")

// Middleware для проверки токена
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Токен отсутствует"})
			c.Abort()
			return
		}

		// Убираем "Bearer " из строки
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Невалидный токен"})
			c.Abort()
			return
		}

		c.Next()
	}
}

var DB *gorm.DB

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("[MANGA-SERVICE] Входящий запрос: %s %s", c.Request.Method, c.Request.URL.Path)
		c.Next()
	}
}

func initDB() {
	dsn := "host=localhost user=postgres password=11223344 dbname=GO-Manga port=5432 sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	log.Println("Manga-Service успешно подключен к базе данных")
}

func main() {
	initDB()
	r := gin.Default()
	r.Use(LoggingMiddleware())
	r.GET("/manga", GetMangas)
	r.GET("/manga/:id", GetManga)

	protected := r.Group("/")
	protected.Use(AuthMiddleware())
	{
		protected.POST("/manga", CreateManga)
		protected.PUT("/manga/:id", UpdateManga)
		protected.DELETE("/manga/:id", DeleteManga)
	}

	r.Run(":8082")
}
