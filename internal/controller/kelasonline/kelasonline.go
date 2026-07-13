// internal/controller/kelasonline/kelasonline.go
package kelasonline

import (
	"github.com/gin-gonic/gin"
)

type KelasOnlineController interface {
	// Pertemuan
	CreatePertemuan(c *gin.Context)
	GetPertemuanByID(c *gin.Context)
	GetPertemuanByGuru(c *gin.Context)
	GetPertemuanBySiswa(c *gin.Context)
	UpdatePertemuan(c *gin.Context)
	DeletePertemuan(c *gin.Context)

	// Materi
	CreateMateri(c *gin.Context)
	GetMateriByPertemuan(c *gin.Context)
	UpdateMateri(c *gin.Context)
	DeleteMateri(c *gin.Context)

	// Meeting
	BukaKelas(c *gin.Context)
	TutupKelas(c *gin.Context)
	GetMeetingByPertemuan(c *gin.Context)

	// Absensi
	JoinKelas(c *gin.Context)
	CreateAbsensi(c *gin.Context)
	GetAbsensiByPertemuan(c *gin.Context)
	GetAbsensiBySiswa(c *gin.Context)
	UpdateAbsensi(c *gin.Context)
	ValidasiAbsensi(c *gin.Context)
}
