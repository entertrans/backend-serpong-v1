// internal/controller/kelasonline/absensi_crud.go
package kelasonline

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/entertrans/backend-bogor.git/internal/dto"
	"github.com/entertrans/backend-bogor.git/internal/model"
)

// ============================================
// ABSENSI CRUD
// ============================================

func (ctrl *kelasonlineController) JoinKelas(c *gin.Context) {
	var req dto.JoinKelasRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cek pertemuan
	var pertemuan model.Pertemuan
	if err := ctrl.db.First(&pertemuan, req.PertemuanID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pertemuan tidak ditemukan"})
		return
	}

	// Cek status - hanya dibuka atau berlangsung
	if pertemuan.Status != "dibuka" && pertemuan.Status != "berlangsung" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kelas belum dibuka atau sudah selesai"})
		return
	}

	// Cek siswa
	siswaID, err := ctrl.getSiswaIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Cek apakah siswa terdaftar di kelas ini
	var taKelasMapel model.TaKelasMapel
	ctrl.db.First(&taKelasMapel, pertemuan.TaKelasMapelID)
	if !ctrl.isSiswaInKelas(siswaID, taKelasMapel.KelasID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak terdaftar di kelas ini"})
		return
	}

	// Dapatkan meeting URL
	var meeting model.PertemuanMeeting
	if err := ctrl.db.Where("pertemuan_id = ?", req.PertemuanID).First(&meeting).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Meeting tidak ditemukan"})
		return
	}

	if meeting.MeetingURL == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Meeting URL tidak ditemukan"})
		return
	}

	now := time.Now()

	// Catat join
	absensi := &model.PertemuanAbsensi{
		PertemuanID: req.PertemuanID,
		SiswaID:     siswaID,
		JoinTime:    &now,
		ViaLMS:      true,
		StatusFinal: "belum_divalidasi",
	}

	// UPSERT - jika sudah ada, update join_time
	if err := ctrl.db.Where("pertemuan_id = ? AND siswa_id = ?", req.PertemuanID, siswaID).Assign(absensi).FirstOrCreate(absensi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Berhasil join kelas",
		"meeting_url": *meeting.MeetingURL,
	})
}

func (ctrl *kelasonlineController) CreateAbsensi(c *gin.Context) {
	var req dto.CreateAbsensiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cek pertemuan
	var pertemuan model.Pertemuan
	if err := ctrl.db.First(&pertemuan, req.PertemuanID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pertemuan tidak ditemukan"})
		return
	}

	// Cek akses guru
	guruID, err := ctrl.getGuruIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if pertemuan.CreatedBy != guruID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses"})
		return
	}

	// Cek status - hanya berlangsung yang bisa diabsensi
	if pertemuan.Status != "berlangsung" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Hanya kelas yang sedang berlangsung yang bisa diabsensi"})
		return
	}

	// Cek apakah siswa ada
	var siswa model.Siswa
	if err := ctrl.db.First(&siswa, req.SiswaID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Siswa tidak ditemukan"})
		return
	}

	// Buat absensi
	absensi := &model.PertemuanAbsensi{
		PertemuanID: req.PertemuanID,
		SiswaID:     req.SiswaID,
		Catatan:     &req.Catatan,
		StatusFinal: "belum_divalidasi",
	}

	if err := ctrl.db.Create(absensi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctrl.db.Preload("Siswa").First(absensi, absensi.AbsensiID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Absensi berhasil dibuat",
		"data":    ctrl.toAbsensiDetailResponse(absensi),
	})
}

func (ctrl *kelasonlineController) GetAbsensiByPertemuan(c *gin.Context) {
	pertemuanID := c.Param("pertemuan_id")
	id, err := strconv.ParseUint(pertemuanID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID pertemuan tidak valid"})
		return
	}

	// Cek akses
	userRole, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Role tidak ditemukan"})
		return
	}

	var pertemuan model.Pertemuan
	if err := ctrl.db.First(&pertemuan, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pertemuan tidak ditemukan"})
		return
	}

	if userRole == "guru" {
		guruID, err := ctrl.getGuruIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		if pertemuan.CreatedBy != guruID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses"})
			return
		}
	} else if userRole == "siswa" {
		siswaID, err := ctrl.getSiswaIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		var taKelasMapel model.TaKelasMapel
		ctrl.db.First(&taKelasMapel, pertemuan.TaKelasMapelID)
		if !ctrl.isSiswaInKelas(siswaID, taKelasMapel.KelasID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses"})
			return
		}
	}

	var absensi []model.PertemuanAbsensi
	if err := ctrl.db.Where("pertemuan_id = ?", id).Preload("Siswa").Find(&absensi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var responses []dto.AbsensiDetailResponse
	for _, a := range absensi {
		responses = append(responses, ctrl.toAbsensiDetailResponse(&a))
	}

	c.JSON(http.StatusOK, gin.H{"data": responses})
}

func (ctrl *kelasonlineController) GetAbsensiBySiswa(c *gin.Context) {
	siswaID, err := ctrl.getSiswaIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var absensi []model.PertemuanAbsensi
	if err := ctrl.db.Where("siswa_id = ?", siswaID).
		Preload("Pertemuan").
		Preload("Pertemuan.TaKelasMapel").
		Preload("Pertemuan.TaKelasMapel.Kelas").
		Preload("Pertemuan.TaKelasMapel.Mapel").
		Order("created_at DESC").
		Find(&absensi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var responses []dto.AbsensiDetailResponse
	for _, a := range absensi {
		responses = append(responses, ctrl.toAbsensiDetailResponse(&a))
	}

	c.JSON(http.StatusOK, gin.H{"data": responses})
}

func (ctrl *kelasonlineController) UpdateAbsensi(c *gin.Context) {
	id := c.Param("id")
	absensiID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var req dto.UpdateAbsensiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var absensi model.PertemuanAbsensi
	if err := ctrl.db.First(&absensi, absensiID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Absensi tidak ditemukan"})
		return
	}

	// Cek akses guru
	var pertemuan model.Pertemuan
	ctrl.db.First(&pertemuan, absensi.PertemuanID)
	guruID, err := ctrl.getGuruIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if pertemuan.CreatedBy != guruID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses"})
		return
	}

	// Cek status pertemuan - hanya berlangsung yang bisa diupdate
	if pertemuan.Status != "berlangsung" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Hanya kelas yang sedang berlangsung yang bisa diubah absensinya"})
		return
	}

	updates := map[string]interface{}{
		"status_final": req.StatusFinal,
	}
	if req.Catatan != "" {
		updates["catatan"] = req.Catatan
	}

	if err := ctrl.db.Model(&absensi).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctrl.db.Preload("Siswa").First(&absensi, absensiID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Absensi berhasil diupdate",
		"data":    ctrl.toAbsensiDetailResponse(&absensi),
	})
}

func (ctrl *kelasonlineController) ValidasiAbsensi(c *gin.Context) {
	id := c.Param("id")
	absensiID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	// Cek absensi
	var absensi model.PertemuanAbsensi
	if err := ctrl.db.First(&absensi, absensiID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Absensi tidak ditemukan"})
		return
	}

	// Cek akses guru
	var pertemuan model.Pertemuan
	ctrl.db.First(&pertemuan, absensi.PertemuanID)
	guruID, err := ctrl.getGuruIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if pertemuan.CreatedBy != guruID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses"})
		return
	}

	// Cek status pertemuan - selesai atau berlangsung
	if pertemuan.Status == "draft" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kelas belum dibuka"})
		return
	}

	now := time.Now()
	if err := ctrl.db.Model(&absensi).Updates(map[string]interface{}{
		"validated_by": guruID,
		"validated_at": now,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctrl.db.Preload("Siswa").First(&absensi, absensiID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Absensi berhasil divalidasi",
		"data":    ctrl.toAbsensiDetailResponse(&absensi),
	})
}
