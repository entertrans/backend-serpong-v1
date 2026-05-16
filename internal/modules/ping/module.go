package pingmodule

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/entertrans/go-base-project.git/internal/config"
	"github.com/entertrans/go-base-project.git/internal/handler"
)

func Register(rg *gin.RouterGroup, _ *gorm.DB, _ *config.Config) {
	pingHandler := handler.NewPingHandler()

	rg.GET("/ping", pingHandler.Ping)

	//     rg.GET("/__whoami", func(c *gin.Context) {
	//     c.JSON(200, gin.H{"ok": true})
	// })

	// rg.GET("/__whoami/", func(c *gin.Context) {
	//     c.JSON(200, gin.H{"ok": true})
	// })
}

// internal/routes/routes.go atau di module lampiran
