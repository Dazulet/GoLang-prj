package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("[USER-SERVICE] Входящий запрос: %s %s", c.Request.Method, c.Request.URL.Path)
		c.Next()
	}
}

func initDB() {
	dsn := "host=localhost user=postgres password=11223344 dbname=GO-Manga port=5432 sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Ошибка подключения к БД в User-Service:", err)
	}
	log.Println("User-Service успешно подключен к базе данных")
	DB.AutoMigrate(&User{})

}

func main() {
	initDB()

	r := gin.Default()
	r.Use(LoggingMiddleware())

	r.POST("/register", Register)
	r.POST("/login", Login)
	r.GET("/users/:id", GetUser)
	r.DELETE("/users/:id", DeleteUser)

	log.Println("User-Service запущен на порту :8081")
	r.Run(":8081")
}
