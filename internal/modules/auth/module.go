package authmodule

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/entertrans/backend-bogor.git/internal/config"
	"github.com/entertrans/backend-bogor.git/internal/controller"
	"github.com/entertrans/backend-bogor.git/internal/handler"
	"github.com/entertrans/backend-bogor.git/internal/middleware"
	"github.com/entertrans/backend-bogor.git/internal/repository"
	"github.com/entertrans/backend-bogor.git/internal/service"
)

func Register(rg *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, cfg)
	authController := controller.NewAuthController(authService)
	authHandler := handler.NewAuthHandler(authController)

	// Public
	authGroup := rg.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
	}

	// Protected
	protectedGroup := rg.Group("/")
	protectedGroup.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		protectedGroup.GET("/profile", authHandler.Profile)
	}
}
