// internal/controller/kelasonline/materi_crud.go
package kelasonline

import (
	"net/http"
	"strconv"

	"github.com/entertrans/backend-bogor.git/internal/dto"
	"github.com/entertrans/backend-bogor.git/internal/model"
	"github.com/gin-gonic/gin"
)

// ============================================
// MATERI CRUD
// ============================================

func (ctrl *kelasonlineController) CreateMateri(c *gin.Context) {
	var req dto.CreateMateriRequest
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

	// Cek apakah materi sudah ada
	var existingMateri model.PertemuanMateri
	if err := ctrl.db.Where("pertemuan_id = ?", req.PertemuanID).First(&existingMateri).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Materi untuk pertemuan ini sudah ada. Gunakan update"})
		return
	}

	materi := &model.PertemuanMateri{
		PertemuanID: req.PertemuanID,
		Judul:       req.Judul,
		Konten:      req.Konten,
	}

	if err := ctrl.db.Create(materi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Materi berhasil dibuat",
		"data": dto.MateriResponse{
			MateriID:    materi.MateriID,
			PertemuanID: materi.PertemuanID,
			Judul:       materi.Judul,
			Konten:      materi.Konten,
			CreatedAt:   materi.CreatedAt,
			UpdatedAt:   materi.UpdatedAt,
		},
	})
}

func (ctrl *kelasonlineController) GetMateriByPertemuan(c *gin.Context) {
	pertemuanID := c.Param("pertemuan_id")
	id, err := strconv.ParseUint(pertemuanID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID pertemuan tidak valid"})
		return
	}

	// Cek akses
	userRole, _ := c.Get("role")

	var pertemuan model.Pertemuan
	if err := ctrl.db.First(&pertemuan, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pertemuan tidak ditemukan"})
		return
	}

	if userRole == "guru" {
		guruID, _ := ctrl.getGuruIDFromContext(c)
		if pertemuan.CreatedBy != guruID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses"})
			return
		}
	} else if userRole == "siswa" {
		siswaID, _ := ctrl.getSiswaIDFromContext(c)
		var taKelasMapel model.TaKelasMapel
		ctrl.db.First(&taKelasMapel, pertemuan.TaKelasMapelID)
		if !ctrl.isSiswaInKelas(siswaID, taKelasMapel.KelasID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses"})
			return
		}
	}

	var materi model.PertemuanMateri
	if err := ctrl.db.Where("pertemuan_id = ?", id).First(&materi).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Materi tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": dto.MateriResponse{
			MateriID:    materi.MateriID,
			PertemuanID: materi.PertemuanID,
			Judul:       materi.Judul,
			Konten:      materi.Konten,
			CreatedAt:   materi.CreatedAt,
			UpdatedAt:   materi.UpdatedAt,
		},
	})
}

func (ctrl *kelasonlineController) UpdateMateri(c *gin.Context) {
	id := c.Param("id")
	materiID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var req dto.UpdateMateriRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var materi model.PertemuanMateri
	if err := ctrl.db.First(&materi, materiID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Materi tidak ditemukan"})
		return
	}

	// Cek akses
	var pertemuan model.Pertemuan
	ctrl.db.First(&pertemuan, materi.PertemuanID)
	guruID, _ := ctrl.getGuruIDFromContext(c)
	if pertemuan.CreatedBy != guruID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses"})
		return
	}

	// Update
	updates := make(map[string]interface{})
	if req.Judul != "" {
		updates["judul"] = req.Judul
	}
	if req.Konten != "" {
		updates["konten"] = req.Konten
	}

	if err := ctrl.db.Model(&materi).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctrl.db.First(&materi, materiID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Materi berhasil diupdate",
		"data": dto.MateriResponse{
			MateriID:    materi.MateriID,
			PertemuanID: materi.PertemuanID,
			Judul:       materi.Judul,
			Konten:      materi.Konten,
			CreatedAt:   materi.CreatedAt,
			UpdatedAt:   materi.UpdatedAt,
		},
	})
}

func (ctrl *kelasonlineController) DeleteMateri(c *gin.Context) {
	id := c.Param("id")
	materiID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var materi model.PertemuanMateri
	if err := ctrl.db.First(&materi, materiID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Materi tidak ditemukan"})
		return
	}

	// Cek akses
	var pertemuan model.Pertemuan
	ctrl.db.First(&pertemuan, materi.PertemuanID)
	guruID, _ := ctrl.getGuruIDFromContext(c)
	if pertemuan.CreatedBy != guruID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses"})
		return
	}

	if err := ctrl.db.Delete(&materi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Materi berhasil dihapus"})
}
