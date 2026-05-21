package main

import (
	"log"

	"github.com/mangalib/manga-service/internal/config"
	"github.com/mangalib/manga-service/internal/database"
	"github.com/mangalib/manga-service/internal/routes"
	"github.com/mangalib/manga-service/internal/storage"
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

	if err := storage.EnsureDirectories(cfg); err != nil {
		log.Fatalf("storage setup failed: %v", err)
	}

	router := routes.Setup(db, cfg)

	log.Printf("📚 Manga Service running on :%s", cfg.AppPort)
	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
