package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/entertrans/go-base-project.git/internal/controller"
	"github.com/entertrans/go-base-project.git/internal/dto"

)

type SiswaHandler interface {
	GetAllSiswa(c *gin.Context)
	CreateSiswa(c *gin.Context)
}

type siswaHandler struct {
	siswaController controller.SiswaController
}

func NewSiswaHandler(siswaController controller.SiswaController) SiswaHandler {
	return &siswaHandler{siswaController: siswaController}
}

// GetAllSiswa handler untuk mendapatkan semua siswa
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

// CreateSiswa handler untuk membuat siswa baru
func (h *siswaHandler) CreateSiswa(c *gin.Context) {
	var req dto.CreateSiswaRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("❌ BIND ERROR:", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("✅ REQUEST MASUK:", req)

	err := h.siswaController.CreateSiswa(req)
	if err != nil {
		fmt.Println("❌ CONTROLLER ERROR:", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("✅ SUCCESS INSERT")

	c.JSON(200, gin.H{"success": true})
}