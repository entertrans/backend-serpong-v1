// internal/controller/kelasonline/helper.go
package kelasonline

import (
	"errors"

	"github.com/entertrans/backend-bogor.git/internal/dto"
	"github.com/entertrans/backend-bogor.git/internal/model"
	"github.com/gin-gonic/gin"
)

// Helper functions untuk mengambil ID dari context
func (ctrl *kelasonlineController) getGuruIDFromContext(c *gin.Context) (uint, error) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		return 0, errors.New("user tidak ditemukan dalam context")
	}

	// Konversi userID ke uint
	var userID uint
	switch v := userIDValue.(type) {
	case float64:
		userID = uint(v)
	case int:
		userID = uint(v)
	case uint:
		userID = v
	case int64:
		userID = uint(v)
	default:
		return 0, errors.New("invalid user id type")
	}

	// Ambil guru_id dari user_id
	var guru model.Guru
	if err := ctrl.db.Where("user_id = ?", userID).First(&guru).Error; err != nil {
		return 0, errors.New("user bukan guru")
	}

	return guru.GuruID, nil
}

func (ctrl *kelasonlineController) getSiswaIDFromContext(c *gin.Context) (uint, error) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		return 0, errors.New("user tidak ditemukan dalam context")
	}

	// Konversi userID ke uint
	var userID uint
	switch v := userIDValue.(type) {
	case float64:
		userID = uint(v)
	case int:
		userID = uint(v)
	case uint:
		userID = v
	case int64:
		userID = uint(v)
	default:
		return 0, errors.New("invalid user id type")
	}

	var siswa model.Siswa
	if err := ctrl.db.Where("user_id = ?", userID).First(&siswa).Error; err != nil {
		return 0, errors.New("user bukan siswa")
	}

	return siswa.SiswaID, nil
}

func (ctrl *kelasonlineController) isSiswaInKelas(siswaID, kelasID uint) bool {
	var siswa model.Siswa
	if err := ctrl.db.First(&siswa, siswaID).Error; err != nil {
		return false
	}
	return siswa.SiswaKelasID != nil && *siswa.SiswaKelasID == int(kelasID)
}

// Helper untuk convert ke response
func (ctrl *kelasonlineController) toPertemuanAbsensiResponse(a *model.PertemuanAbsensi) dto.PertemuanAbsensiResponse {
	resp := dto.PertemuanAbsensiResponse{
		AbsensiID:   a.AbsensiID,
		SiswaID:     a.SiswaID,
		JoinTime:    a.JoinTime,
		ViaLMS:      a.ViaLMS,
		StatusFinal: a.StatusFinal,
		Catatan:     a.Catatan,
		ValidatedBy: a.ValidatedBy,
		ValidatedAt: a.ValidatedAt,
	}

	if a.Siswa != nil {
		resp.SiswaNama = a.Siswa.SiswaNama
		resp.SiswaNIS = a.Siswa.SiswaNIS
	}

	return resp
}

func (ctrl *kelasonlineController) toAbsensiDetailResponse(a *model.PertemuanAbsensi) dto.AbsensiDetailResponse {
	resp := dto.AbsensiDetailResponse{
		AbsensiID:   a.AbsensiID,
		PertemuanID: a.PertemuanID,
		SiswaID:     a.SiswaID,
		JoinTime:    a.JoinTime,
		ViaLMS:      a.ViaLMS,
		StatusFinal: a.StatusFinal,
		Catatan:     a.Catatan,
		ValidatedBy: a.ValidatedBy,
		ValidatedAt: a.ValidatedAt,
		CreatedAt:   a.CreatedAt,
	}

	if a.Siswa != nil {
		resp.SiswaNama = a.Siswa.SiswaNama
		resp.SiswaNIS = a.Siswa.SiswaNIS
	}

	if a.Pertemuan != nil {
		resp.Pertemuan = &dto.PertemuanSimpleResponse{
			PertemuanID: a.Pertemuan.PertemuanID,
			PertemuanKe: a.Pertemuan.PertemuanKe,
			Tema:        a.Pertemuan.Tema,
			Tanggal:     a.Pertemuan.Tanggal,
			Status:      a.Pertemuan.Status,
		}
	}

	return resp
}

func (ctrl *kelasonlineController) toPertemuanResponse(p *model.Pertemuan) dto.PertemuanResponse {
	resp := dto.PertemuanResponse{
		PertemuanID:    p.PertemuanID,
		TaKelasMapelID: p.TaKelasMapelID,
		PertemuanKe:    p.PertemuanKe,
		Tema:           p.Tema,
		Deskripsi:      p.Deskripsi,
		Tanggal:        p.Tanggal,
		Status:         p.Status,
		CreatedBy:      p.CreatedBy,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}

	if p.Guru != nil {
		resp.Guru = &dto.GuruResponse{
			GuruID:   p.Guru.GuruID,
			GuruNama: p.Guru.GuruNama,
		}
	}

	if p.TaKelasMapel != nil && p.TaKelasMapel.Kelas != nil && p.TaKelasMapel.Mapel != nil {
		resp.TaKelasMapel = &dto.TaKelasMapelResponse{
			TaKelasMapelID: p.TaKelasMapel.TaKelasMapelID,
			TaID:           p.TaKelasMapel.TaID,
			KelasID:        p.TaKelasMapel.KelasID,
			KdMapel:        p.TaKelasMapel.KdMapel,
			GuruID:         p.TaKelasMapel.GuruID,
			Urutan:         p.TaKelasMapel.Urutan,
			Kelas: &dto.KelasResponse{
				KelasId:   p.TaKelasMapel.Kelas.KelasId,
				KelasNama: p.TaKelasMapel.Kelas.KelasNama,
			},
			Mapel: &dto.MapelResponse{
				KdMapel: p.TaKelasMapel.Mapel.KdMapel,
				NmMapel: p.TaKelasMapel.Mapel.NmMapel,
			},
		}
	}

	if p.Materi != nil {
		resp.Materi = &dto.MateriResponse{
			MateriID:    p.Materi.MateriID,
			PertemuanID: p.Materi.PertemuanID,
			Judul:       p.Materi.Judul,
			Konten:      p.Materi.Konten,
			CreatedAt:   p.Materi.CreatedAt,
			UpdatedAt:   p.Materi.UpdatedAt,
		}
	}

	if p.Meeting != nil {
		resp.Meeting = &dto.MeetingResponse{
			MeetingID:       p.Meeting.MeetingID,
			PertemuanID:     p.Meeting.PertemuanID,
			Provider:        p.Meeting.Provider,
			CalendarEventID: p.Meeting.CalendarEventID,
			MeetingURL:      p.Meeting.MeetingURL,
			StartedAt:       p.Meeting.StartedAt,
			EndedAt:         p.Meeting.EndedAt,
			CreatedAt:       p.Meeting.CreatedAt,
		}
	}

	if len(p.Absensi) > 0 {
		var absensiResponses []dto.PertemuanAbsensiResponse
		for _, a := range p.Absensi {
			absensiResponses = append(absensiResponses, ctrl.toPertemuanAbsensiResponse(&a))
		}
		resp.Absensi = absensiResponses
	}

	return resp
}
