// internal/modules/rapor/handler/penilaian_handler.go
package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/entertrans/backend-bogor.git/internal/controller"
	"github.com/entertrans/backend-bogor.git/internal/dto"
	"github.com/gin-gonic/gin"
)

type PenilaianHandler struct {
	penilaianController controller.PenilaianController
}

func NewPenilaianHandler(penilaianController controller.PenilaianController) *PenilaianHandler {
	return &PenilaianHandler{
		penilaianController: penilaianController,
	}
}

// ==================== GET SISWA BY KELAS ====================

// GetSiswaByKelasHandler - GET /penilaian/siswa/:ta_id/:kelas_id
func (h *PenilaianHandler) GetSiswaByKelasHandler(c *gin.Context) {
	taID, err := strconv.ParseUint(c.Param("ta_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ta_id"})
		return
	}

	kelasID, err := strconv.ParseUint(c.Param("kelas_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid kelas_id"})
		return
	}

	result, err := h.penilaianController.GetSiswaByKelas(uint(taID), uint(kelasID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

// ==================== NILAI MAPEL ====================

// GetNilaiMapelHandler - GET /penilaian/nilai/:ta_id/:ta_kelas_mapel_id/:kelas_id
func (h *PenilaianHandler) GetNilaiMapelHandler(c *gin.Context) {
	taID, err := strconv.ParseUint(c.Param("ta_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ta_id"})
		return
	}

	taKelasMapelID, err := strconv.ParseUint(c.Param("ta_kelas_mapel_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ta_kelas_mapel_id"})
		return
	}

	kelasID, err := strconv.ParseUint(c.Param("kelas_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid kelas_id"})
		return
	}

	result, err := h.penilaianController.GetNilaiMapel(uint(taID), uint(taKelasMapelID), uint(kelasID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

// SaveNilaiMapelHandler - POST /penilaian/nilai/save
func (h *PenilaianHandler) SaveNilaiMapelHandler(c *gin.Context) {
	var req dto.NilaiMapelRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	userIDFloat, ok := userIDRaw.(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
		return
	}

	userID := uint(userIDFloat)

	result, err := h.penilaianController.SaveNilaiMapel(
		req,
		userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ==================== ABSENSI ====================

// GetAbsensiHandler - GET /penilaian/absensi/:ta_id/:kelas_id
func (h *PenilaianHandler) GetAbsensiHandler(c *gin.Context) {
	taID, err := strconv.ParseUint(c.Param("ta_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ta_id"})
		return
	}

	kelasID, err := strconv.ParseUint(c.Param("kelas_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid kelas_id"})
		return
	}

	result, err := h.penilaianController.GetAbsensi(uint(taID), uint(kelasID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

// SaveAbsensiHandler - POST /penilaian/absensi/save
func (h *PenilaianHandler) SaveAbsensiHandler(c *gin.Context) {
	var req dto.AbsensiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.penilaianController.SaveAbsensi(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// ==================== EKSKUL ====================

// GetEkskulHandler - GET /penilaian/ekskul/:ta_id/:kelas_id/:siswa_nis
func (h *PenilaianHandler) GetEkskulHandler(c *gin.Context) {
	taID, err := strconv.ParseUint(c.Param("ta_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ta_id"})
		return
	}

	kelasID, err := strconv.ParseUint(c.Param("kelas_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid kelas_id"})
		return
	}

	siswaNIS := c.Param("siswa_nis")
	if siswaNIS == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "siswa_nis is required"})
		return
	}

	result, err := h.penilaianController.GetEkskul(uint(taID), uint(kelasID), siswaNIS)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

// SaveEkskulHandler - POST /penilaian/ekskul/save
func (h *PenilaianHandler) SaveEkskulHandler(c *gin.Context) {
	var req dto.EkskulRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.penilaianController.SaveEkskul(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// ==================== RAPORT ====================

// GetRaportBySiswaHandler - GET /penilaian/raport/:ta_id/:kelas_id/:siswa_nis
func (h *PenilaianHandler) GetRaportBySiswaHandler(c *gin.Context) {
	taID, err := strconv.ParseUint(c.Param("ta_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ta_id"})
		return
	}

	kelasID, err := strconv.ParseUint(c.Param("kelas_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid kelas_id"})
		return
	}

	siswaNIS := c.Param("siswa_nis")
	if siswaNIS == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "siswa_nis is required"})
		return
	}

	result, err := h.penilaianController.GetRaportBySiswa(uint(taID), uint(kelasID), siswaNIS)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// ==================== STATUS ====================
// internal/modules/rapor/handler/penilaian_handler.go

func (h *PenilaianHandler) UpdateStatusPublishHandler(c *gin.Context) {
	taID, err := strconv.ParseUint(c.Param("ta_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ta_id"})
		return
	}

	kelasID, err := strconv.ParseUint(c.Param("kelas_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid kelas_id"})
		return
	}

	siswaNIS := c.Param("siswa_nis")
	if siswaNIS == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "siswa_nis is required"})
		return
	}

	var req dto.UpdateStatusPublishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.penilaianController.UpdateStatusPublish(uint(taID), uint(kelasID), siswaNIS, req.StatusPublish)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status publish berhasil diupdate"})
}

// EditNilaiPerSiswaHandler - PUT /penilaian/edit-siswa
func (h *PenilaianHandler) EditNilaiPerSiswaHandler(c *gin.Context) {
	// Baca raw body untuk debugging
	body, _ := io.ReadAll(c.Request.Body)
	fmt.Printf("Raw Body: %s\n", string(body))

	// Restore body
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	var req dto.EditNilaiPerSiswaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("Binding error: %v\n", err) // ← lihat error binding
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("Parsed request: raport_id=%d, nilai_list_count=%d\n", req.RaportID, len(req.NilaiList))

	err := h.penilaianController.EditNilaiPerSiswa(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Nilai berhasil diupdate"})
}

// GetNilaiHistoryHandler - get /nilai/history/:raport_nilai_id
func (h *PenilaianHandler) GetNilaiHistoryHandler(c *gin.Context) {
	raportNilaiID, err := strconv.ParseUint(c.Param("raport_nilai_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid raport_nilai_id"})
		return
	}

	result, err := h.penilaianController.GetNilaiHistory(uint(raportNilaiID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}
