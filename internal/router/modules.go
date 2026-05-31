// internal/router/router.go
package router

import (
	"github.com/entertrans/backend-bogor.git/internal/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ModuleRegisterFunc func(rg *gin.RouterGroup, db *gorm.DB, cfg *config.Config)

func RegisterModules(r *gin.Engine, db *gorm.DB, cfg *config.Config, modules ...ModuleRegisterFunc) {
	api := r.Group("/api/v1")

	for _, register := range modules {
		register(api, db, cfg)
	}
}
