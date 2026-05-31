package masterdatamodule

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/entertrans/backend-bogor.git/internal/config"
	"github.com/entertrans/backend-bogor.git/internal/controller"
	"github.com/entertrans/backend-bogor.git/internal/handler"
	"github.com/entertrans/backend-bogor.git/internal/middleware"
)

func Register(rg *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	// log.Println("=== MASTERDATA MODULE REGISTERING ===") // <-- TAMBAHKAN
	masterController := controller.NewMasterdataController(db)
	masterHandler := handler.NewMasterdataHandler(masterController)

	master := rg.Group("/master")
	master.Use(middleware.AuthMiddleware(cfg))
	{
		// log.Println("Registering route: /kelas/all11111111111111111")        // <-- TAMBAHKAN
		master.GET("/kelas/all", masterHandler.GetKelasAll)       // 1 - alumni
		master.GET("/kelas/aktif", masterHandler.GetKelasAktif)   // 1 - 12
		master.GET("/kelas/alumni", masterHandler.GetKelasAlumni) // alumni sd-sma
		master.GET("/satelit", masterHandler.GetSatelit)          // tbl_satelit
		master.GET("/guru/aktif", masterHandler.GetGuruAktifHandler)
		master.GET("/mapel/aktif", masterHandler.GetMapelAktifHandler)
	}
	// log.Println("=== MASTERDATA MODULE REGISTERED ===") // <-- TAMBAHKAN
}
