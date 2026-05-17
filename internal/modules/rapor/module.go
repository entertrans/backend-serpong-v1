// internal/modules/rapor/module.go
package rapor

import (
	"github.com/entertrans/backend-bogor.git/internal/config"
	"github.com/entertrans/backend-bogor.git/internal/controller"
	"github.com/entertrans/backend-bogor.git/internal/handler"
	"github.com/entertrans/backend-bogor.git/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(rg *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	// Inisialisasi controller
	tahunAjaranController := controller.NewTahunAjaranController(db)
	kurikulumController := controller.NewKurikulumController(db)
	penilaianController := controller.NewPenilaianController(db)

	// Inisialisasi handler
	tahunAjaranHandler := handler.NewTahunAjaranHandler(tahunAjaranController)
	kurikulumHandler := handler.NewKurikulumHandler(kurikulumController)
	penilaianHandler := handler.NewPenilaianHandler(penilaianController)

	// ==================== ROUTES TAHUN AJARAN ====================
	tahunAjaran := rg.Group("/tahun-ajaran")
	tahunAjaran.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		tahunAjaran.GET("/all", tahunAjaranHandler.GetAllTahunAjaranHandler)
		tahunAjaran.POST("/create", tahunAjaranHandler.CreateTahunAjaranHandler)
		tahunAjaran.PUT("/:ta_id/activate", tahunAjaranHandler.ActivateTahunAjaranHandler)
		tahunAjaran.PUT("/:ta_id/publish", tahunAjaranHandler.PublishTahunAjaranHandler)       // NEW
		tahunAjaran.PUT("/:ta_id/reactivate", tahunAjaranHandler.ReactivateTahunAjaranHandler) // NEW
	}

	// ==================== ROUTES MASTER DATA (untuk dropdown) ====================
	master := rg.Group("/master")
	master.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		// master.GET("/kelas/aktif", kurikulumHandler.GetKelasAktifHandler)
		master.GET("/guru/aktif", kurikulumHandler.GetGuruAktifHandler)
		master.GET("/mapel/aktif", kurikulumHandler.GetMapelAktifHandler)
	}

	// ==================== ROUTES KURIKULUM ====================
	taKelasMapel := rg.Group("/ta-kelas-mapel")
	taKelasMapel.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		taKelasMapel.GET("/:ta_id/:kelas_id", kurikulumHandler.GetKurikulumByKelasHandler)
		taKelasMapel.POST("/save", kurikulumHandler.SaveKurikulumHandler)
		taKelasMapel.POST("/copy/:ta_id/:kelas_id", kurikulumHandler.CopyKurikulumHandler)
		taKelasMapel.GET("/check/:ta_id", kurikulumHandler.CheckKurikulumStatusHandler)
	}

	// ==================== ROUTES PENILAIAN ====================
	penilaian := rg.Group("/penilaian")
	penilaian.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		// Siswa
		penilaian.GET("/siswa/:ta_id/:kelas_id", penilaianHandler.GetSiswaByKelasHandler)

		// Nilai Mapel
		penilaian.GET("/nilai/:ta_id/:ta_kelas_mapel_id/:kelas_id", penilaianHandler.GetNilaiMapelHandler)
		penilaian.POST("/nilai/save", penilaianHandler.SaveNilaiMapelHandler)

		// Absensi
		penilaian.GET("/absensi/:ta_id/:kelas_id", penilaianHandler.GetAbsensiHandler)
		penilaian.POST("/absensi/save", penilaianHandler.SaveAbsensiHandler)

		// Ekskul
		penilaian.GET("/ekskul/:ta_id/:kelas_id/:siswa_nis", penilaianHandler.GetEkskulHandler)
		penilaian.POST("/ekskul/save", penilaianHandler.SaveEkskulHandler)

		// Raport
		penilaian.GET("/raport/:ta_id/:kelas_id/:siswa_nis", penilaianHandler.GetRaportBySiswaHandler)
		// Update status publish per siswa
		penilaian.PUT("/status-publish/:ta_id/:kelas_id/:siswa_nis", penilaianHandler.UpdateStatusPublishHandler)

		// Edit nilai per siswa (force majeur)
		penilaian.PUT("/edit-siswa", penilaianHandler.EditNilaiPerSiswaHandler)
	}
}
