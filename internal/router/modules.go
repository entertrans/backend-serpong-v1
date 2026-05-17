package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/entertrans/backend-bogor.git/internal/config"
)

type Module func(rg *gin.RouterGroup, db *gorm.DB, cfg *config.Config)

func RegisterModules(engine *gin.Engine, db *gorm.DB, cfg *config.Config, modules ...Module) {
	api := engine.Group("/api/v1")

	for _, m := range modules {
		m(api, db, cfg)
	}
}
