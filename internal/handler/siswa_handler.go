package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/entertrans/backend-bogor.git/internal/controller/siswa"
	"github.com/entertrans/backend-bogor.git/internal/dto"
)

type SiswaHandler interface {
	GetAllSiswa(c *gin.Context)
	GetSiswaByKelas(c *gin.Context)                // Tambahkan ini
	GetSiswaAktif(c *gin.Context)                  // NEW
	GetSiswaAlumni(c *gin.Context)                 // NEW
	GetSiswaDeleted(c *gin.Context)                // NEW
	GetActiveStudentsForEnrollment(c *gin.Context) // ✅ NEW
	CreateSiswa(c *gin.Context)
	UpdateSiswa(c *gin.Context)
	UpdateOrangtua(c *gin.Context)
}

type siswaHandler struct {
	siswaController siswa.SiswaController // ← ubah tipe
}

func NewSiswaHandler(siswaController siswa.SiswaController) SiswaHandler {
	return &siswaHandler{siswaController: siswaController}
}

func (h *siswaHandler) GetAllSiswa(c *gin.Context) {
	siswaList, err := h.siswaController.GetAllSiswa()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data siswa berhasil diambil",
		"data":    siswaList,
	})
}

func (h *siswaHandler) CreateSiswa(c *gin.Context) {
	var req dto.CreateSiswaRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		// fmt.Println("❌ BIND ERROR:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("✅ REQUEST MASUK: NIS=%s, Nama=%s, KelasID=%d\n",
		req.SiswaNIS, req.SiswaNama, req.SiswaKelasID)

	err := h.siswaController.CreateSiswa(req)
	if err != nil {
		// fmt.Println("❌ CONTROLLER ERROR:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// fmt.Println("✅ SUCCESS INSERT")

	// Response dengan informasi login
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Siswa berhasil ditambahkan",
		"data": gin.H{
			"nis":              req.SiswaNIS,
			"email":            req.SiswaNIS + "@siswa.sch.id",
			"default_password": req.SiswaNIS,
		},
	})
}

func (h *siswaHandler) UpdateSiswa(c *gin.Context) {
	nis := c.Param("nis") // ambil NIS dari URL: /api/v1/siswa/:nis
	if nis == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "NIS tidak boleh kosong"})
		return
	}

	var req dto.UpdateSiswaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// fmt.Println("❌ BIND ERROR:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("✅ UPDATE SISWA: NIS=%s, Payload=%+v\n", nis, req)

	err := h.siswaController.UpdateSiswa(nis, req)
	if err != nil {
		// fmt.Println("❌ CONTROLLER ERROR:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data siswa berhasil diupdate",
	})
}

// UpdateOrangtua handler untuk update data orangtua (PATCH)
func (h *siswaHandler) UpdateOrangtua(c *gin.Context) {
	nis := c.Param("nis") // ambil NIS dari URL: /api/v1/siswa/:nis/orangtua
	if nis == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "NIS tidak boleh kosong"})
		return
	}

	var req dto.UpdateOrangtuaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// fmt.Println("❌ BIND ERROR:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("✅ UPDATE ORANGTUA: NIS=%s, Payload=%+v\n", nis, req)

	err := h.siswaController.UpdateOrangtua(nis, req)
	if err != nil {
		// fmt.Println("❌ CONTROLLER ERROR:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data orangtua berhasil diupdate",
	})
}

// GetSiswaByKelas handler untuk mendapatkan siswa berdasarkan kelas dengan filter dan pagination
func (h *siswaHandler) GetSiswaByKelas(c *gin.Context) {
	var req dto.FilterSiswaRequest

	// Bind query parameters
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	// Validate kelas_id is provided
	if req.KelasID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "kelas_id is required",
		})
		return
	}

	// Call controller
	result, err := h.siswaController.GetSiswaByKelasID(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"message":    "Data siswa berhasil diambil",
		"data":       result.Data,
		"pagination": result.Pagination,
	})
}

// GetSiswaAktif handler untuk mendapatkan siswa aktif
func (h *siswaHandler) GetSiswaAktif(c *gin.Context) {
	var req dto.FilterSiswaStatusRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	result, err := h.siswaController.GetSiswaAktif(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"message":    "Data siswa aktif berhasil diambil",
		"data":       result.Data,
		"pagination": result.Pagination,
		"total":      result.Total,
	})
}

// GetSiswaAlumni handler untuk mendapatkan alumni
func (h *siswaHandler) GetSiswaAlumni(c *gin.Context) {
	var req dto.FilterSiswaStatusRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	result, err := h.siswaController.GetSiswaAlumni(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"message":    "Data alumni berhasil diambil",
		"data":       result.Data,
		"pagination": result.Pagination,
		"total":      result.Total,
	})
}

// GetSiswaDeleted handler untuk mendapatkan siswa yang di-soft delete
func (h *siswaHandler) GetSiswaDeleted(c *gin.Context) {
	var req dto.FilterSiswaStatusRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	result, err := h.siswaController.GetSiswaDeleted(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"message":    "Data siswa deleted berhasil diambil",
		"data":       result.Data,
		"pagination": result.Pagination,
		"total":      result.Total,
	})
}

// GetActiveStudentsForEnrollment handler untuk mengambil semua siswa aktif (tanpa pagination)
func (h *siswaHandler) GetActiveStudentsForEnrollment(c *gin.Context) {
	result, err := h.siswaController.GetActiveStudentsForEnrollment()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data siswa aktif berhasil diambil",
		"data":    result,
		"total":   len(result),
	})
}
