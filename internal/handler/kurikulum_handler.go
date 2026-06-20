// internal/modules/rapor/handler/kurikulum_handler.go
package handler

import (
	"net/http"
	"strconv"

	"github.com/entertrans/backend-bogor.git/internal/controller"
	"github.com/entertrans/backend-bogor.git/internal/dto"
	"github.com/gin-gonic/gin"
)

type KurikulumHandler struct {
	kurikulumController controller.KurikulumController
}

func NewKurikulumHandler(kurikulumController controller.KurikulumController) *KurikulumHandler {
	return &KurikulumHandler{
		kurikulumController: kurikulumController,
	}
}

// ==================== KURIKULUM SETUP ====================

// GetKurikulumByKelasHandler - GET /ta-kelas-mapel/:ta_id/:kelas_id
func (h *KurikulumHandler) GetKurikulumByKelasHandler(c *gin.Context) {
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

	result, err := h.kurikulumController.GetKurikulumByKelas(uint(taID), uint(kelasID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// SaveKurikulumHandler - POST /ta-kelas-mapel/save
func (h *KurikulumHandler) SaveKurikulumHandler(c *gin.Context) {
	var req dto.SaveKurikulumRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.kurikulumController.SaveKurikulum(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// CopyKurikulumHandler - POST /ta-kelas-mapel/copy/:ta_id/:kelas_id
func (h *KurikulumHandler) CopyKurikulumHandler(c *gin.Context) {
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

	var req dto.CopyKurikulumRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.kurikulumController.CopyKurikulumFromPrevious(req, uint(taID), uint(kelasID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// CheckKurikulumStatusHandler - GET /ta-kelas-mapel/check/:ta_id
func (h *KurikulumHandler) CheckKurikulumStatusHandler(c *gin.Context) {
	taID, err := strconv.ParseUint(c.Param("ta_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ta_id"})
		return
	}

	result, err := h.kurikulumController.CheckKurikulumStatus(uint(taID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetKurikulumByGuruHandler - GET /ta-kelas-mapel/:ta_id/:kelas_id/:guru_id
func (h *KurikulumHandler) GetKurikulumByGuruHandler(c *gin.Context) {
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

	guruID, err := strconv.ParseUint(c.Param("guru_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid guru_id"})
		return
	}

	result, err := h.kurikulumController.GetKurikulumByGuru(uint(taID), uint(kelasID), uint(guruID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetKelasWaliHandler - GET /ta-kelas-wali/:ta_id/:guru_id
func (h *KurikulumHandler) GetKelasWaliHandler(c *gin.Context) {
	taID, err := strconv.ParseUint(c.Param("ta_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ta_id"})
		return
	}

	guruID, err := strconv.ParseUint(c.Param("guru_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid guru_id"})
		return
	}

	result, err := h.kurikulumController.GetKelasWali(uint(taID), uint(guruID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}
