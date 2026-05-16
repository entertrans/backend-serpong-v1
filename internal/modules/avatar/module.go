// internal/modules/avatar/module.go
package avatarmodule

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/entertrans/go-base-project.git/internal/config"
	"github.com/entertrans/go-base-project.git/internal/handler"
)

func Register(rg *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	avatarHandler := handler.NewAvatarHandler(cfg)

	// Public endpoint - no auth needed
	rg.GET("/avatar/:nis", avatarHandler.ServeAvatar)
}
