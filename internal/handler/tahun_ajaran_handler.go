// internal/modules/rapor/handler/tahun_ajaran_handler.go
package handler

import (
	"net/http"
	"strconv"

	"github.com/entertrans/go-base-project.git/internal/controller"
	"github.com/entertrans/go-base-project.git/internal/dto"
	"github.com/gin-gonic/gin"
)

type TahunAjaranHandler struct {
	tahunAjaranController controller.TahunAjaranController
}

func NewTahunAjaranHandler(tahunAjaranController controller.TahunAjaranController) *TahunAjaranHandler {
	return &TahunAjaranHandler{
		tahunAjaranController: tahunAjaranController,
	}
}

func (h *TahunAjaranHandler) GetAllTahunAjaranHandler(c *gin.Context) {
	result, err := h.tahunAjaranController.GetAllTahunAjaran()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *TahunAjaranHandler) CreateTahunAjaranHandler(c *gin.Context) {
	var req dto.CreateTahunAjaranRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.tahunAjaranController.CreateTahunAjaran(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *TahunAjaranHandler) ActivateTahunAjaranHandler(c *gin.Context) {
	taID, err := strconv.ParseUint(c.Param("ta_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ta_id"})
		return
	}

	result, err := h.tahunAjaranController.ActivateTahunAjaran(uint(taID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// NEW: Publish tahun ajaran dengan tanggal
func (h *TahunAjaranHandler) PublishTahunAjaranHandler(c *gin.Context) {
	taID, err := strconv.ParseUint(c.Param("ta_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ta_id"})
		return
	}

	var req dto.PublishTahunAjaranRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.tahunAjaranController.PublishTahunAjaran(uint(taID), req.PublishDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// NEW: Reactivate tahun ajaran (kembalikan ke status aktif)
func (h *TahunAjaranHandler) ReactivateTahunAjaranHandler(c *gin.Context) {
	taID, err := strconv.ParseUint(c.Param("ta_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ta_id"})
		return
	}

	result, err := h.tahunAjaranController.ReactivateTahunAjaran(uint(taID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
