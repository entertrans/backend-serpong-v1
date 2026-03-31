package masterdatamodule

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/entertrans/go-base-project.git/internal/config"
	"github.com/entertrans/go-base-project.git/internal/controller"
	"github.com/entertrans/go-base-project.git/internal/handler"
	"github.com/entertrans/go-base-project.git/internal/middleware"
)

func Register(rg *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	masterController := controller.NewMasterdataController(db)
	masterHandler := handler.NewMasterdataHandler(masterController)

	master := rg.Group("/master")
	master.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		master.GET("/kelas/all", masterHandler.GetKelasAll)       // 1 - alumni
		master.GET("/kelas/aktif", masterHandler.GetKelasAktif)   // 1 - 12
		master.GET("/kelas/alumni", masterHandler.GetKelasAlumni) // alumni sd-sma
		master.GET("/satelit", masterHandler.GetSatelit)          // tbl_satelit
	}
}
