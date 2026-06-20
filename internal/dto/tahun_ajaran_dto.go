// internal/modules/rapor/dto/tahun_ajaran_dto.go
package dto

import "time"

type CreateTahunAjaranRequest struct {
	TahunAjaran string `json:"tahun_ajaran" binding:"required"`
	Semester    string `json:"semester" binding:"required,oneof=1 2"`
}

type TahunAjaranResponse struct {
	TaID        uint       `json:"ta_id"`
	TahunAjaran string     `json:"tahun_ajaran"`
	Semester    string     `json:"semester"`
	Status      string     `json:"status"` // draft, aktif, selesai
	PublishDate *time.Time `json:"publish_date"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
}

type TahunAjaranListResponse struct {
	Data []TahunAjaranResponse `json:"data"`
}

type ActivateTahunAjaranResponse struct {
	TaID        uint   `json:"ta_id"`
	TahunAjaran string `json:"tahun_ajaran"`
	Semester    string `json:"semester"`
	Status      string `json:"status"`
	IsActive    bool   `json:"is_active"`
	Message     string `json:"message"`
}

// Request untuk publish dengan tanggal
type PublishTahunAjaranRequest struct {
	PublishDate string `json:"publish_date" binding:"required"` // format: YYYY-MM-DD
}

// Response untuk publish
type PublishTahunAjaranResponse struct {
	TaID        uint      `json:"ta_id"`
	TahunAjaran string    `json:"tahun_ajaran"`
	Semester    string    `json:"semester"`
	Status      string    `json:"status"`
	PublishDate time.Time `json:"publish_date"`
	Message     string    `json:"message"`
}

// Response untuk reactivate
type ReactivateTahunAjaranResponse struct {
	TaID        uint   `json:"ta_id"`
	TahunAjaran string `json:"tahun_ajaran"`
	Semester    string `json:"semester"`
	Status      string `json:"status"`
	Message     string `json:"message"`
}

// internal/dto/rapor_dto.go
// KelasWaliItem - item kelas wali
type KelasWaliItem struct {
	KelasID   uint   `json:"kelas_id"`
	KelasNama string `json:"kelas_nama"`
}

// GetKelasWaliResponse - response untuk mengambil kelas wali
type GetKelasWaliResponse struct {
	TaID       uint            `json:"ta_id"`
	TaNama     string          `json:"ta_nama"`
	Semester   string          `json:"semester"`
	Status     string          `json:"status"`
	GuruID     uint            `json:"guru_id"`
	GuruNama   string          `json:"guru_nama"`
	GuruNIP    string          `json:"guru_nip"`
	TotalKelas int             `json:"total_kelas"`
	KelasList  []KelasWaliItem `json:"kelas_list"`
}
