package main

import (
	"log"

	"github.com/mangalib/comment-service/internal/config"
	"github.com/mangalib/comment-service/internal/database"
	"github.com/mangalib/comment-service/internal/routes"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	router := routes.Setup(db, cfg)

	log.Printf("💬 Comment Service running on :%s", cfg.AppPort)
	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
