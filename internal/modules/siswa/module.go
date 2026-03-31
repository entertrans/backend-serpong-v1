package siswamodule

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/entertrans/go-base-project.git/internal/config"
	"github.com/entertrans/go-base-project.git/internal/controller"
	"github.com/entertrans/go-base-project.git/internal/handler"
	"github.com/entertrans/go-base-project.git/internal/middleware"

)

func Register(rg *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	siswaController := controller.NewSiswaController(db)
	siswaHandler := handler.NewSiswaHandler(siswaController)

	siswa := rg.Group("/siswa")
	siswa.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		siswa.GET("/", siswaHandler.GetAllSiswa)
		siswa.POST("/", siswaHandler.CreateSiswa)
		// Contoh kalau nanti mau tambah route:
		// siswa.GET("/:nis", siswaHandler.GetSiswaByNIS)
	}
}
