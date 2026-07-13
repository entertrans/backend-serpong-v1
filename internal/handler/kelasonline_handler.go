// internal/handler/kelasonline_handler.go (update semua method)

package handler

import (
	"github.com/entertrans/backend-bogor.git/internal/controller/kelasonline"
	"github.com/gin-gonic/gin"
)

type KelasOnlineHandler struct {
	controller kelasonline.KelasOnlineController
}

func NewKelasOnlineHandler(controller kelasonline.KelasOnlineController) *KelasOnlineHandler {
	return &KelasOnlineHandler{controller: controller}
}

// Pertemuan
func (h *KelasOnlineHandler) CreatePertemuan(c *gin.Context) {
	h.controller.CreatePertemuan(c)
}

func (h *KelasOnlineHandler) GetPertemuanByID(c *gin.Context) {
	h.controller.GetPertemuanByID(c)
}

func (h *KelasOnlineHandler) GetPertemuanByGuru(c *gin.Context) {
	h.controller.GetPertemuanByGuru(c)
}

func (h *KelasOnlineHandler) GetPertemuanBySiswa(c *gin.Context) {
	h.controller.GetPertemuanBySiswa(c)
}

func (h *KelasOnlineHandler) UpdatePertemuan(c *gin.Context) {
	h.controller.UpdatePertemuan(c)
}

func (h *KelasOnlineHandler) DeletePertemuan(c *gin.Context) {
	h.controller.DeletePertemuan(c)
}

// Materi
func (h *KelasOnlineHandler) CreateMateri(c *gin.Context) {
	h.controller.CreateMateri(c)
}

func (h *KelasOnlineHandler) GetMateriByPertemuan(c *gin.Context) {
	h.controller.GetMateriByPertemuan(c)
}

func (h *KelasOnlineHandler) UpdateMateri(c *gin.Context) {
	h.controller.UpdateMateri(c)
}

func (h *KelasOnlineHandler) DeleteMateri(c *gin.Context) {
	h.controller.DeleteMateri(c)
}

// Meeting
func (h *KelasOnlineHandler) BukaKelas(c *gin.Context) {
	h.controller.BukaKelas(c)
}

func (h *KelasOnlineHandler) TutupKelas(c *gin.Context) {
	h.controller.TutupKelas(c)
}

func (h *KelasOnlineHandler) GetMeetingByPertemuan(c *gin.Context) {
	h.controller.GetMeetingByPertemuan(c)
}

// Absensi
func (h *KelasOnlineHandler) JoinKelas(c *gin.Context) {
	h.controller.JoinKelas(c)
}

func (h *KelasOnlineHandler) CreateAbsensi(c *gin.Context) {
	h.controller.CreateAbsensi(c)
}

func (h *KelasOnlineHandler) GetAbsensiByPertemuan(c *gin.Context) {
	h.controller.GetAbsensiByPertemuan(c)
}

func (h *KelasOnlineHandler) GetAbsensiBySiswa(c *gin.Context) {
	h.controller.GetAbsensiBySiswa(c)
}

func (h *KelasOnlineHandler) UpdateAbsensi(c *gin.Context) {
	h.controller.UpdateAbsensi(c)
}

func (h *KelasOnlineHandler) ValidasiAbsensi(c *gin.Context) {
	h.controller.ValidasiAbsensi(c)
}
