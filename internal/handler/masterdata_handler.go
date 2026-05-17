package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/entertrans/backend-bogor.git/internal/controller"
)

type MasterdataHandler interface {
	GetKelasAll(c *gin.Context)
	GetKelasAktif(c *gin.Context)
	GetKelasAlumni(c *gin.Context)
	GetSatelit(c *gin.Context)
}

type masterdataHandler struct {
	masterController controller.MasterdataController
}

func NewMasterdataHandler(masterController controller.MasterdataController) MasterdataHandler {
	return &masterdataHandler{masterController: masterController}
}

func (h *masterdataHandler) GetKelasAll(c *gin.Context) {
	result, err := h.masterController.GetKelasAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *masterdataHandler) GetKelasAktif(c *gin.Context) {
	result, err := h.masterController.GetKelasAktif()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *masterdataHandler) GetKelasAlumni(c *gin.Context) {
	result, err := h.masterController.GetKelasAlumni()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *masterdataHandler) GetSatelit(c *gin.Context) {
	result, err := h.masterController.GetSatelit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
