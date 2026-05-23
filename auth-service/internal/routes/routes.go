package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mangalib/auth-service/internal/config"
	"github.com/mangalib/auth-service/internal/handlers"
	"github.com/mangalib/auth-service/internal/middleware"
	"github.com/mangalib/auth-service/internal/repositories"
	"github.com/mangalib/auth-service/internal/services"
	"gorm.io/gorm"
	"net/http"
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
		c.JSON(http.StatusOK, gin.H{"service": "auth", "status": "ok"})
	})

	userRepo := repositories.NewUserRepository(db)
	authSvc := services.NewAuthService(userRepo, cfg)
	userSvc := services.NewUserService(userRepo, cfg)

	authH := handlers.NewAuthHandler(authSvc)
	userH := handlers.NewUserHandler(userSvc)

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authH.Register)
			auth.POST("/login", authH.Login)
			auth.POST("/validate", authH.Validate)
		}

		protected := api.Group("")
		protected.Use(middleware.Auth(cfg.JWTSecret))
		{
			protected.GET("/users/me", userH.Me)
			protected.PATCH("/users/me", userH.UpdateProfile)

			admin := protected.Group("")
			admin.Use(middleware.AdminOnly())
			{
				admin.GET("/users", userH.GetAllUsers)
				admin.PATCH("/users/:id", userH.UpdateUserProfile)
			}
		}

		api.GET("/users/:id", userH.GetProfile)
	}

	return r
}
