// internal/modules/kelasonline/kelasonlinemodule.go
package kelasonline

import (
	"github.com/entertrans/backend-bogor.git/internal/config"
	"github.com/entertrans/backend-bogor.git/internal/controller/kelasonline"
	"github.com/entertrans/backend-bogor.git/internal/handler"
	"github.com/entertrans/backend-bogor.git/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(rg *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	controller := kelasonline.NewKelasOnlineController(db)
	handler := handler.NewKelasOnlineHandler(controller)

	kelasGroup := rg.Group("/kelasonline")
	kelasGroup.Use(middleware.AuthMiddleware(cfg))
	{
		// ============================================
		// PERTEMUAN
		// ============================================
		pertemuanGroup := kelasGroup.Group("/pertemuan")
		{
			// Guru
			pertemuanGroup.POST("", handler.CreatePertemuan)
			pertemuanGroup.PUT("/:id", handler.UpdatePertemuan)
			pertemuanGroup.DELETE("/:id", handler.DeletePertemuan)
			pertemuanGroup.GET("/guru", handler.GetPertemuanByGuru)

			// Siswa
			pertemuanGroup.GET("/siswa", handler.GetPertemuanBySiswa)

			// Public (dengan akses check di controller)
			pertemuanGroup.GET("/:id", handler.GetPertemuanByID)
		}

		// ============================================
		// MATERI
		// ============================================
		materiGroup := kelasGroup.Group("/materi")
		{
			materiGroup.POST("", handler.CreateMateri)
			materiGroup.GET("/pertemuan/:pertemuan_id", handler.GetMateriByPertemuan)
			materiGroup.PUT("/:id", handler.UpdateMateri)
			materiGroup.DELETE("/:id", handler.DeleteMateri)
		}

		// ============================================
		// MEETING
		// ============================================
		meetingGroup := kelasGroup.Group("/meeting")
		{
			// Guru
			meetingGroup.POST("/buka", handler.BukaKelas)
			meetingGroup.POST("/tutup", handler.TutupKelas)

			// Public
			meetingGroup.GET("/pertemuan/:pertemuan_id", handler.GetMeetingByPertemuan)
		}

		// ============================================
		// ABSENSI
		// ============================================
		absensiGroup := kelasGroup.Group("/absensi")
		{
			// Siswa
			absensiGroup.POST("/join", handler.JoinKelas)
			absensiGroup.GET("/siswa", handler.GetAbsensiBySiswa)

			// Guru
			absensiGroup.POST("", handler.CreateAbsensi)
			absensiGroup.GET("/pertemuan/:pertemuan_id", handler.GetAbsensiByPertemuan)
			absensiGroup.PUT("/:id", handler.UpdateAbsensi)
			absensiGroup.PUT("/:id/validasi", handler.ValidasiAbsensi)
		}
	}
}
