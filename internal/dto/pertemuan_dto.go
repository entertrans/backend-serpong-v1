// internal/dto/pertemuan_dto.go
package dto

import (
	"time"
)

// ============================================
// PERTEMUAN DTO
// ============================================

type CreatePertemuanRequest struct {
	TaKelasMapelID uint      `json:"ta_kelas_mapel_id" binding:"required"`
	Tema           string    `json:"tema" binding:"required"`
	Deskripsi      string    `json:"deskripsi"`
	Tanggal        time.Time `json:"tanggal" binding:"required"`
}

type UpdatePertemuanRequest struct {
	Tema      string    `json:"tema"`
	Deskripsi string    `json:"deskripsi"`
	Tanggal   time.Time `json:"tanggal"`
}

type PertemuanResponse struct {
	PertemuanID    uint      `json:"pertemuan_id"`
	TaKelasMapelID uint      `json:"ta_kelas_mapel_id"`
	PertemuanKe    int       `json:"pertemuan_ke"`
	Tema           string    `json:"tema"`
	Deskripsi      string    `json:"deskripsi"`
	Tanggal        time.Time `json:"tanggal"`
	Status         string    `json:"status"`
	CreatedBy      uint      `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	TaKelasMapel *TaKelasMapelResponse      `json:"ta_kelas_mapel,omitempty"`
	Guru         *GuruResponse              `json:"guru,omitempty"`
	Materi       *MateriResponse            `json:"materi,omitempty"`
	Meeting      *MeetingResponse           `json:"meeting,omitempty"`
	Absensi      []PertemuanAbsensiResponse `json:"absensi,omitempty"`
}

type ListPertemuanResponse struct {
	PertemuanID uint      `json:"pertemuan_id"`
	PertemuanKe int       `json:"pertemuan_ke"`
	Tema        string    `json:"tema"`
	Tanggal     time.Time `json:"tanggal"`
	Status      string    `json:"status"`
}

// ============================================
// MATERI DTO
// ============================================

type CreateMateriRequest struct {
	PertemuanID uint   `json:"pertemuan_id" binding:"required"`
	Judul       string `json:"judul" binding:"required"`
	Konten      string `json:"konten" binding:"required"`
}

type UpdateMateriRequest struct {
	Judul  string `json:"judul"`
	Konten string `json:"konten"`
}

type MateriResponse struct {
	MateriID    uint      `json:"materi_id"`
	PertemuanID uint      `json:"pertemuan_id"`
	Judul       string    `json:"judul"`
	Konten      string    `json:"konten"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ============================================
// MEETING DTO
// ============================================

type CreateMeetingRequest struct {
	PertemuanID uint `json:"pertemuan_id" binding:"required"`
}

type MeetingResponse struct {
	MeetingID       uint       `json:"meeting_id"`
	PertemuanID     uint       `json:"pertemuan_id"`
	Provider        string     `json:"provider"`
	CalendarEventID *string    `json:"calendar_event_id,omitempty"`
	MeetingURL      *string    `json:"meeting_url,omitempty"`
	StartedAt       *time.Time `json:"started_at,omitempty"`
	EndedAt         *time.Time `json:"ended_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

// ============================================
// ABSENSI DTO - Ganti nama untuk menghindari konflik
// ============================================

type CreateAbsensiRequest struct {
	PertemuanID uint   `json:"pertemuan_id" binding:"required"`
	SiswaID     uint   `json:"siswa_id" binding:"required"`
	Catatan     string `json:"catatan"`
}

type UpdateAbsensiRequest struct {
	StatusFinal string `json:"status_final" binding:"required,oneof=hadir izin sakit alfa"`
	Catatan     string `json:"catatan"`
}

// PertemuanAbsensiResponse - untuk response absensi di pertemuan
type PertemuanAbsensiResponse struct {
	AbsensiID   uint       `json:"absensi_id"`
	SiswaID     uint       `json:"siswa_id"`
	JoinTime    *time.Time `json:"join_time,omitempty"`
	ViaLMS      bool       `json:"via_lms"`
	StatusFinal string     `json:"status_final"`
	Catatan     *string    `json:"catatan,omitempty"`
	ValidatedBy *uint      `json:"validated_by,omitempty"`
	ValidatedAt *time.Time `json:"validated_at,omitempty"`
	SiswaNama   *string    `json:"siswa_nama,omitempty"`
	SiswaNIS    *string    `json:"siswa_nis,omitempty"`
}

// AbsensiDetailResponse - untuk detail absensi
type AbsensiDetailResponse struct {
	AbsensiID   uint                     `json:"absensi_id"`
	PertemuanID uint                     `json:"pertemuan_id"`
	SiswaID     uint                     `json:"siswa_id"`
	JoinTime    *time.Time               `json:"join_time,omitempty"`
	ViaLMS      bool                     `json:"via_lms"`
	StatusFinal string                   `json:"status_final"`
	Catatan     *string                  `json:"catatan,omitempty"`
	ValidatedBy *uint                    `json:"validated_by,omitempty"`
	ValidatedAt *time.Time               `json:"validated_at,omitempty"`
	CreatedAt   time.Time                `json:"created_at"`
	SiswaNama   *string                  `json:"siswa_nama,omitempty"`
	SiswaNIS    *string                  `json:"siswa_nis,omitempty"`
	Pertemuan   *PertemuanSimpleResponse `json:"pertemuan,omitempty"`
}

type PertemuanSimpleResponse struct {
	PertemuanID uint      `json:"pertemuan_id"`
	PertemuanKe int       `json:"pertemuan_ke"`
	Tema        string    `json:"tema"`
	Tanggal     time.Time `json:"tanggal"`
	Status      string    `json:"status"`
}

// ============================================
// KELAS ONLINE DTO
// ============================================

type BukaKelasRequest struct {
	PertemuanID uint `json:"pertemuan_id" binding:"required"`
}

type TutupKelasRequest struct {
	PertemuanID uint `json:"pertemuan_id" binding:"required"`
}

type JoinKelasRequest struct {
	PertemuanID uint `json:"pertemuan_id" binding:"required"`
}

type JoinKelasResponse struct {
	Message    string `json:"message"`
	MeetingURL string `json:"meeting_url"`
}

// ============================================
// SUPPORTING DTO
// ============================================

type TaKelasMapelResponse struct {
	TaKelasMapelID uint  `json:"ta_kelas_mapel_id"`
	TaID           uint  `json:"ta_id"`
	KelasID        uint  `json:"kelas_id"`
	KdMapel        uint  `json:"kd_mapel"`
	GuruID         *uint `json:"guru_id,omitempty"`
	Urutan         int   `json:"urutan"`

	Kelas *KelasResponse `json:"kelas,omitempty"`
	Mapel *MapelResponse `json:"mapel,omitempty"`
}

type GuruResponse struct {
	GuruID   uint   `json:"guru_id"`
	GuruNama string `json:"guru_nama"`
}

type KelasResponse struct {
	KelasId   uint   `json:"kelas_id"`
	KelasNama string `json:"kelas_nama"`
}

type MapelResponse struct {
	KdMapel uint   `json:"kd_mapel"`
	NmMapel string `json:"nm_mapel"`
}
