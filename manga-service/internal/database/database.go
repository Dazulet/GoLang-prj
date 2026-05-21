package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/mangalib/manga-service/internal/config"
	"github.com/mangalib/manga-service/internal/models"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	lvl := logger.Info
	if cfg.AppEnv == "production" {
		lvl = logger.Error
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(lvl),
	})
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	log.Println("✅ Manga DB connected")
	return db, nil
}

func Migrate(db *gorm.DB) error {
	log.Println("🔄 Manga migrations running...")
	err := db.AutoMigrate(
		&models.Genre{},
		&models.Tag{},
		&models.Manga{},
		&models.Chapter{},
		&models.Page{},
		&models.Bookmark{},
		&models.Rating{},
		&models.ReadingProgress{},
	)
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	log.Println("✅ Manga migrations done")
	return nil
}
