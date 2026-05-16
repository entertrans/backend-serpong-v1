package siswamodule

import (
	// ❌ HAPUS import ini
	// ✅ TAMBAHKAN import ini

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/entertrans/go-base-project.git/internal/config"
	"github.com/entertrans/go-base-project.git/internal/controller/siswa"
	"github.com/entertrans/go-base-project.git/internal/handler"
	"github.com/entertrans/go-base-project.git/internal/middleware"

)

func Register(rg *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	siswaController := siswa.NewSiswaController(db)
	siswaHandler := handler.NewSiswaHandler(siswaController)

	siswa := rg.Group("/siswa")
	siswa.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		siswa.GET("/", siswaHandler.GetAllSiswa)
		siswa.POST("/", siswaHandler.CreateSiswa)
		siswa.PATCH("/:nis", siswaHandler.UpdateSiswa)           
		siswa.PATCH("/:nis/orangtua", siswaHandler.UpdateOrangtua) 
		// Contoh kalau nanti mau tambah route:
		// siswa.GET("/:nis", siswaHandler.GetSiswaByNIS)
	}
}