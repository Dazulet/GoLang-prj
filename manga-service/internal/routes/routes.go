package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/mangalib/manga-service/internal/client"
	"github.com/mangalib/manga-service/internal/config"
	"github.com/mangalib/manga-service/internal/handlers"
	"github.com/mangalib/manga-service/internal/middleware"
	"github.com/mangalib/manga-service/internal/repositories"
	"github.com/mangalib/manga-service/internal/services"
)

func Setup(db *gorm.DB, cfg *config.Config) *gin.Engine {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.Use(gin.Logger())

	r.Static("/uploads", cfg.UploadDir)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"service": "manga", "status": "ok"})
	})

	authClient := client.NewAuthClient(cfg.AuthServiceURL)

	mangaRepo := repositories.NewMangaRepository(db)
	chapterRepo := repositories.NewChapterRepository(db)
	pageRepo := repositories.NewPageRepository(db)
	genreRepo := repositories.NewGenreRepository(db)
	tagRepo := repositories.NewTagRepository(db)
	bookmarkRepo := repositories.NewBookmarkRepository(db)
	ratingRepo := repositories.NewRatingRepository(db)
	progressRepo := repositories.NewProgressRepository(db)

	mangaSvc := services.NewMangaService(mangaRepo, genreRepo, tagRepo)
	chapterSvc := services.NewChapterService(chapterRepo, pageRepo, mangaRepo)
	genreSvc := services.NewGenreService(genreRepo)
	tagSvc := services.NewTagService(tagRepo)
	bookmarkSvc := services.NewBookmarkService(bookmarkRepo)
	ratingSvc := services.NewRatingService(ratingRepo)
	progressSvc := services.NewProgressService(progressRepo)

	mangaH := handlers.NewMangaHandler(mangaSvc, cfg)
	chapterH := handlers.NewChapterHandler(chapterSvc, cfg)
	genreTagH := handlers.NewGenreTagHandler(genreSvc, tagSvc)
	bookmarkH := handlers.NewBookmarkHandler(bookmarkSvc)
	ratingH := handlers.NewRatingHandler(ratingSvc)
	progressH := handlers.NewProgressHandler(progressSvc)
	userProxyH := handlers.NewUserProxyHandler(authClient) // Resty v2 proxy

	api := r.Group("/api")

	manga := api.Group("/manga")
	{
		manga.GET("", mangaH.List)
		manga.GET("/:id", mangaH.Get)
		manga.GET("/slug/:slug", mangaH.GetBySlug)
		manga.GET("/:id/chapters", chapterH.ListByManga)
	}

	api.GET("/chapters/:id", chapterH.Get)
	api.GET("/genres", genreTagH.ListGenres)
	api.GET("/tags", genreTagH.ListTags)

	api.GET("/users/:id/profile", userProxyH.GetUserProfile)

	protected := api.Group("")
	protected.Use(middleware.Auth(cfg.JWTSecret))
	{
		protected.GET("/bookmarks", bookmarkH.List)
		protected.POST("/bookmarks", bookmarkH.Upsert)
		protected.DELETE("/bookmarks/:id", bookmarkH.Remove)

		protected.POST("/ratings", ratingH.Upsert)
		protected.GET("/ratings/manga/:id", ratingH.GetMine)

		protected.POST("/progress", progressH.Update)
		protected.GET("/progress", progressH.History)

		admin := protected.Group("")
		admin.Use(middleware.AdminOnly())
		{
			admin.POST("/manga", mangaH.Create)
			admin.PUT("/manga/:id", mangaH.Update)
			admin.DELETE("/manga/:id", mangaH.Delete)
			admin.POST("/manga/:id/cover", mangaH.UploadCover)

			admin.POST("/manga/:id/chapters", chapterH.Create)
			admin.PUT("/chapters/:id", chapterH.Update)
			admin.DELETE("/chapters/:id", chapterH.Delete)
			admin.POST("/chapters/:id/pages", chapterH.UploadPages)

			admin.POST("/genres", genreTagH.CreateGenre)
			admin.POST("/tags", genreTagH.CreateTag)
		}
	}

	return r
}
