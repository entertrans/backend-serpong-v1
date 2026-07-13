// internal/controller/kelasonline/meeting_crud.go
package kelasonline

import (
	"net/http"
	"strconv"
	"time"

	"github.com/entertrans/backend-bogor.git/internal/dto"
	"github.com/entertrans/backend-bogor.git/internal/model"
	"github.com/gin-gonic/gin"
)

// ============================================
// MEETING CRUD
// ============================================

func (ctrl *kelasonlineController) BukaKelas(c *gin.Context) {
	var req dto.BukaKelasRequest
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

	// Cek status - hanya draft yang bisa dibuka
	if pertemuan.Status != "draft" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pertemuan tidak dalam status draft"})
		return
	}

	// Cek apakah materi sudah ada
	var materi model.PertemuanMateri
	if err := ctrl.db.Where("pertemuan_id = ?", req.PertemuanID).First(&materi).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mohon buat materi terlebih dahulu sebelum membuka kelas"})
		return
	}

	// Buat Google Meet
	meetingURL, calendarEventID, err := ctrl.createGoogleMeet(pertemuan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat Google Meet: " + err.Error()})
		return
	}

	now := time.Now()

	// Simpan meeting
	meeting := &model.PertemuanMeeting{
		PertemuanID:     req.PertemuanID,
		Provider:        "google_meet",
		CalendarEventID: &calendarEventID,
		MeetingURL:      &meetingURL,
		StartedAt:       &now,
	}

	if err := ctrl.db.Create(meeting).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update status pertemuan
	if err := ctrl.db.Model(&pertemuan).Updates(map[string]interface{}{
		"status": "berlangsung",
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctrl.db.Preload("TaKelasMapel.Kelas").Preload("TaKelasMapel.Mapel").First(&pertemuan, pertemuan.PertemuanID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Kelas berhasil dibuka",
		"data":    ctrl.toPertemuanResponse(&pertemuan),
	})
}

func (ctrl *kelasonlineController) TutupKelas(c *gin.Context) {
	var req dto.TutupKelasRequest
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

	// Cek status - hanya berlangsung yang bisa ditutup
	if pertemuan.Status != "berlangsung" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pertemuan tidak dalam status berlangsung"})
		return
	}

	now := time.Now()

	// Update meeting ended_at
	if err := ctrl.db.Model(&model.PertemuanMeeting{}).Where("pertemuan_id = ?", req.PertemuanID).Update("ended_at", now).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update status pertemuan
	if err := ctrl.db.Model(&pertemuan).Updates(map[string]interface{}{
		"status": "selesai",
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctrl.db.Preload("TaKelasMapel.Kelas").Preload("TaKelasMapel.Mapel").First(&pertemuan, pertemuan.PertemuanID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Kelas berhasil ditutup",
		"data":    ctrl.toPertemuanResponse(&pertemuan),
	})
}

func (ctrl *kelasonlineController) GetMeetingByPertemuan(c *gin.Context) {
	pertemuanID := c.Param("pertemuan_id")
	id, err := strconv.ParseUint(pertemuanID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID pertemuan tidak valid"})
		return
	}

	var meeting model.PertemuanMeeting
	if err := ctrl.db.Where("pertemuan_id = ?", id).First(&meeting).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Meeting tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": dto.MeetingResponse{
			MeetingID:       meeting.MeetingID,
			PertemuanID:     meeting.PertemuanID,
			Provider:        meeting.Provider,
			CalendarEventID: meeting.CalendarEventID,
			MeetingURL:      meeting.MeetingURL,
			StartedAt:       meeting.StartedAt,
			EndedAt:         meeting.EndedAt,
			CreatedAt:       meeting.CreatedAt,
		},
	})
}

// ============================================
// GOOGLE MEET INTEGRATION (Stub)
// ============================================

func (ctrl *kelasonlineController) createGoogleMeet(pertemuan model.Pertemuan) (string, string, error) {
	// TODO: Implementasi integrasi Google Meet API
	// Ini adalah stub, nanti diganti dengan implementasi sebenarnya

	// Simulasi pembuatan Google Meet
	meetingURL := "https://meet.google.com/xxx-xxxx-xxx"
	calendarEventID := "event_123456789"

	return meetingURL, calendarEventID, nil
}
