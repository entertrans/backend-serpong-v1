// internal/controller/kelasonline/pertemuan_crud.go
package kelasonline

import (
	"net/http"
	"strconv"

	"github.com/entertrans/backend-bogor.git/internal/dto"
	"github.com/entertrans/backend-bogor.git/internal/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type kelasonlineController struct {
	db *gorm.DB
}

func NewKelasOnlineController(db *gorm.DB) KelasOnlineController {
	return &kelasonlineController{db: db}
}

// ============================================
// PERTEMUAN CRUD
// ============================================

func (ctrl *kelasonlineController) CreatePertemuan(c *gin.Context) {
	var req dto.CreatePertemuanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ambil userID dari context
	guruID, err := ctrl.getGuruIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Cek apakah guru memiliki akses ke penugasan ini
	var taKelasMapel model.TaKelasMapel
	if err := ctrl.db.Where("ta_kelas_mapel_id = ? AND guru_id = ?", req.TaKelasMapelID, guruID).First(&taKelasMapel).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses ke penugasan ini"})
		return
	}

	// Hitung pertemuan ke berapa
	var count int64
	ctrl.db.Model(&model.Pertemuan{}).Where("ta_kelas_mapel_id = ?", req.TaKelasMapelID).Count(&count)
	pertemuanKe := int(count) + 1

	pertemuan := &model.Pertemuan{
		TaKelasMapelID: req.TaKelasMapelID,
		PertemuanKe:    pertemuanKe,
		Tema:           req.Tema,
		Deskripsi:      req.Deskripsi,
		Tanggal:        req.Tanggal,
		Status:         "draft",
		CreatedBy:      guruID,
	}

	if err := ctrl.db.Create(pertemuan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Load relasi untuk response
	ctrl.db.Preload("TaKelasMapel.Kelas").Preload("TaKelasMapel.Mapel").Preload("Guru").First(pertemuan, pertemuan.PertemuanID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pertemuan berhasil dibuat",
		"data":    ctrl.toPertemuanResponse(pertemuan),
	})
}

func (ctrl *kelasonlineController) GetPertemuanByID(c *gin.Context) {
	id := c.Param("id")
	pertemuanID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var pertemuan model.Pertemuan
	if err := ctrl.db.Preload("TaKelasMapel.Kelas").Preload("TaKelasMapel.Mapel").Preload("Guru").Preload("Materi").Preload("Meeting").Preload("Absensi.Siswa").First(&pertemuan, pertemuanID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pertemuan tidak ditemukan"})
		return
	}

	// Cek akses berdasarkan role
	userRole, _ := c.Get("role")
	guruID, _ := ctrl.getGuruIDFromContext(c)

	if userRole == "guru" {
		if pertemuan.CreatedBy != guruID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses ke pertemuan ini"})
			return
		}
	} else if userRole == "siswa" {
		// Cek apakah siswa terdaftar di kelas ini
		siswaID, _ := ctrl.getSiswaIDFromContext(c)
		if !ctrl.isSiswaInKelas(siswaID, pertemuan.TaKelasMapel.KelasID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses ke pertemuan ini"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": ctrl.toPertemuanResponse(&pertemuan),
	})
}

func (ctrl *kelasonlineController) GetPertemuanByGuru(c *gin.Context) {
	guruID, err := ctrl.getGuruIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Ambil semua penugasan guru
	var taKelasMapelIDs []uint
	ctrl.db.Model(&model.TaKelasMapel{}).Where("guru_id = ?", guruID).Pluck("ta_kelas_mapel_id", &taKelasMapelIDs)

	var pertemuan []model.Pertemuan
	if err := ctrl.db.Where("ta_kelas_mapel_id IN ?", taKelasMapelIDs).Preload("TaKelasMapel.Kelas").Preload("TaKelasMapel.Mapel").Order("tanggal DESC, pertemuan_ke DESC").Find(&pertemuan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var responses []dto.ListPertemuanResponse
	for _, p := range pertemuan {
		responses = append(responses, dto.ListPertemuanResponse{
			PertemuanID: p.PertemuanID,
			PertemuanKe: p.PertemuanKe,
			Tema:        p.Tema,
			Tanggal:     p.Tanggal,
			Status:      p.Status,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": responses})
}

func (ctrl *kelasonlineController) GetPertemuanBySiswa(c *gin.Context) {
	siswaID, err := ctrl.getSiswaIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Ambil kelas siswa
	var siswa model.Siswa
	if err := ctrl.db.First(&siswa, siswaID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Siswa tidak ditemukan"})
		return
	}

	// Ambil semua penugasan untuk kelas siswa
	var taKelasMapelIDs []uint
	ctrl.db.Model(&model.TaKelasMapel{}).Where("kelas_id = ?", siswa.SiswaKelasID).Pluck("ta_kelas_mapel_id", &taKelasMapelIDs)

	var pertemuan []model.Pertemuan
	if err := ctrl.db.Where("ta_kelas_mapel_id IN ?", taKelasMapelIDs).Preload("TaKelasMapel.Kelas").Preload("TaKelasMapel.Mapel").Order("tanggal DESC, pertemuan_ke DESC").Find(&pertemuan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var responses []dto.ListPertemuanResponse
	for _, p := range pertemuan {
		responses = append(responses, dto.ListPertemuanResponse{
			PertemuanID: p.PertemuanID,
			PertemuanKe: p.PertemuanKe,
			Tema:        p.Tema,
			Tanggal:     p.Tanggal,
			Status:      p.Status,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": responses})
}

func (ctrl *kelasonlineController) UpdatePertemuan(c *gin.Context) {
	id := c.Param("id")
	pertemuanID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var req dto.UpdatePertemuanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var pertemuan model.Pertemuan
	if err := ctrl.db.First(&pertemuan, pertemuanID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pertemuan tidak ditemukan"})
		return
	}

	// Cek akses
	guruID, _ := ctrl.getGuruIDFromContext(c)
	if pertemuan.CreatedBy != guruID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses"})
		return
	}

	// Cek status - hanya draft yang bisa diubah
	if pertemuan.Status != "draft" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pertemuan dengan status selain draft tidak dapat diubah"})
		return
	}

	// Update data
	updates := make(map[string]interface{})
	if req.Tema != "" {
		updates["tema"] = req.Tema
	}
	if req.Deskripsi != "" {
		updates["deskripsi"] = req.Deskripsi
	}
	if !req.Tanggal.IsZero() {
		updates["tanggal"] = req.Tanggal
	}

	if err := ctrl.db.Model(&pertemuan).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctrl.db.Preload("TaKelasMapel.Kelas").Preload("TaKelasMapel.Mapel").First(&pertemuan, pertemuan.PertemuanID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Pertemuan berhasil diupdate",
		"data":    ctrl.toPertemuanResponse(&pertemuan),
	})
}

func (ctrl *kelasonlineController) DeletePertemuan(c *gin.Context) {
	id := c.Param("id")
	pertemuanID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var pertemuan model.Pertemuan
	if err := ctrl.db.First(&pertemuan, pertemuanID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pertemuan tidak ditemukan"})
		return
	}

	// Cek akses
	guruID, _ := ctrl.getGuruIDFromContext(c)
	if pertemuan.CreatedBy != guruID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses"})
		return
	}

	// Hanya draft yang bisa dihapus
	if pertemuan.Status != "draft" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pertemuan dengan status selain draft tidak dapat dihapus"})
		return
	}

	if err := ctrl.db.Delete(&pertemuan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pertemuan berhasil dihapus"})
}
