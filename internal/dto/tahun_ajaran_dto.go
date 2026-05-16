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
