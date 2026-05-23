package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/mangalib/comment-service/internal/client"
	"github.com/mangalib/comment-service/internal/config"
	"github.com/mangalib/comment-service/internal/handlers"
	"github.com/mangalib/comment-service/internal/middleware"
	"github.com/mangalib/comment-service/internal/repositories"
	"github.com/mangalib/comment-service/internal/services"
)

func Setup(db *gorm.DB, cfg *config.Config) *gin.Engine {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.Use(gin.Logger())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"service": "comment", "status": "ok"})
	})

	authClient := client.NewAuthClient(cfg.AuthServiceURL)
	commentRepo := repositories.NewCommentRepository(db)
	commentSvc := services.NewCommentService(commentRepo, authClient)
	commentH := handlers.NewCommentHandler(commentSvc)

	api := r.Group("/api")

	api.GET("/comments", commentH.List)

	protected := api.Group("")
	protected.Use(middleware.Auth(cfg.JWTSecret))
	{
		protected.POST("/comments", commentH.Create)
		protected.DELETE("/comments/:id", commentH.Delete)
		protected.POST("/comments/:id/like", commentH.Like)
	}

	return r
}
