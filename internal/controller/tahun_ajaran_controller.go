// internal/modules/rapor/controller/tahun_ajaran_controller.go
package controller

import (
	"errors"
	"time"

	"github.com/entertrans/backend-bogor.git/internal/dto"
	"github.com/entertrans/backend-bogor.git/internal/model"
	"gorm.io/gorm"
)

type TahunAjaranController interface {
	GetAllTahunAjaran() (dto.TahunAjaranListResponse, error)
	CreateTahunAjaran(req dto.CreateTahunAjaranRequest) (dto.TahunAjaranResponse, error)
	ActivateTahunAjaran(taID uint) (dto.ActivateTahunAjaranResponse, error)
	PublishTahunAjaran(taID uint, publishDate string) (dto.PublishTahunAjaranResponse, error) // NEW
	ReactivateTahunAjaran(taID uint) (dto.ReactivateTahunAjaranResponse, error)               // NEW
}

type tahunAjaranController struct {
	db *gorm.DB
}

func NewTahunAjaranController(db *gorm.DB) TahunAjaranController {
	return &tahunAjaranController{db: db}
}

func (c *tahunAjaranController) GetAllTahunAjaran() (dto.TahunAjaranListResponse, error) {
	var tahunAjaranList []model.TahunAjaran

	err := c.db.Order("ta_id DESC").Find(&tahunAjaranList).Error
	if err != nil {
		return dto.TahunAjaranListResponse{}, err
	}

	data := make([]dto.TahunAjaranResponse, len(tahunAjaranList))
	for i, ta := range tahunAjaranList {
		data[i] = dto.TahunAjaranResponse{
			TaID:        ta.TaID,
			TahunAjaran: ta.TahunAjaran,
			Semester:    ta.Semester,
			Status:      ta.Status,
			PublishDate: ta.PublishDate,
			IsActive:    ta.IsActive,
			CreatedAt:   ta.CreatedAt,
		}
	}

	return dto.TahunAjaranListResponse{Data: data}, nil
}

func (c *tahunAjaranController) CreateTahunAjaran(req dto.CreateTahunAjaranRequest) (dto.TahunAjaranResponse, error) {
	// Cek duplikasi
	var existing model.TahunAjaran
	err := c.db.Where("tahun_ajaran = ? AND semester = ?", req.TahunAjaran, req.Semester).
		First(&existing).Error

	if err == nil {
		return dto.TahunAjaranResponse{}, errors.New("tahun ajaran dengan semester ini sudah ada")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.TahunAjaranResponse{}, err
	}

	tahunAjaran := model.TahunAjaran{
		TahunAjaran: req.TahunAjaran,
		Semester:    req.Semester,
		Status:      "draft",
		IsActive:    false,
	}

	err = c.db.Create(&tahunAjaran).Error
	if err != nil {
		return dto.TahunAjaranResponse{}, err
	}

	return dto.TahunAjaranResponse{
		TaID:        tahunAjaran.TaID,
		TahunAjaran: tahunAjaran.TahunAjaran,
		Semester:    tahunAjaran.Semester,
		Status:      tahunAjaran.Status,
		PublishDate: tahunAjaran.PublishDate,
		IsActive:    tahunAjaran.IsActive,
		CreatedAt:   tahunAjaran.CreatedAt,
	}, nil
}

func (c *tahunAjaranController) ActivateTahunAjaran(taID uint) (dto.ActivateTahunAjaranResponse, error) {
	tx := c.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Cek apakah TA ada
	var tahunAjaran model.TahunAjaran
	err := tx.First(&tahunAjaran, taID).Error
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ActivateTahunAjaranResponse{}, errors.New("tahun ajaran tidak ditemukan")
		}
		return dto.ActivateTahunAjaranResponse{}, err
	}

	// Cek status
	if tahunAjaran.Status != "draft" {
		tx.Rollback()
		return dto.ActivateTahunAjaranResponse{}, errors.New("hanya tahun ajaran dengan status DRAFT yang bisa diaktifkan")
	}

	// Non-aktifkan semua TA yang aktif
	err = tx.Model(&model.TahunAjaran{}).
		Where("is_active = ?", true).
		Update("is_active", false).Error
	if err != nil {
		tx.Rollback()
		return dto.ActivateTahunAjaranResponse{}, err
	}

	// Aktifkan TA ini
	tahunAjaran.Status = "aktif"
	tahunAjaran.IsActive = true

	err = tx.Save(&tahunAjaran).Error
	if err != nil {
		tx.Rollback()
		return dto.ActivateTahunAjaranResponse{}, err
	}

	err = tx.Commit().Error
	if err != nil {
		return dto.ActivateTahunAjaranResponse{}, err
	}

	return dto.ActivateTahunAjaranResponse{
		TaID:        tahunAjaran.TaID,
		TahunAjaran: tahunAjaran.TahunAjaran,
		Semester:    tahunAjaran.Semester,
		Status:      tahunAjaran.Status,
		IsActive:    tahunAjaran.IsActive,
		Message:     "Tahun ajaran berhasil diaktifkan",
	}, nil
}

// NEW: Publish tahun ajaran dengan tanggal
func (c *tahunAjaranController) PublishTahunAjaran(taID uint, publishDateStr string) (dto.PublishTahunAjaranResponse, error) {
	// Parse tanggal
	publishDate, err := time.Parse("2006-01-02", publishDateStr)
	if err != nil {
		return dto.PublishTahunAjaranResponse{}, errors.New("format tanggal tidak valid, gunakan YYYY-MM-DD")
	}

	tx := c.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Cek apakah TA ada
	var tahunAjaran model.TahunAjaran
	err = tx.First(&tahunAjaran, taID).Error
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.PublishTahunAjaranResponse{}, errors.New("tahun ajaran tidak ditemukan")
		}
		return dto.PublishTahunAjaranResponse{}, err
	}

	// Cek status (hanya yang aktif bisa dipublish)
	if tahunAjaran.Status != "aktif" {
		tx.Rollback()
		return dto.PublishTahunAjaranResponse{}, errors.New("hanya tahun ajaran dengan status AKTIF yang bisa dipublish")
	}

	// Update status dan publish date
	tahunAjaran.Status = "selesai"
	tahunAjaran.PublishDate = &publishDate

	err = tx.Save(&tahunAjaran).Error
	if err != nil {
		tx.Rollback()
		return dto.PublishTahunAjaranResponse{}, err
	}

	err = tx.Commit().Error
	if err != nil {
		return dto.PublishTahunAjaranResponse{}, err
	}

	return dto.PublishTahunAjaranResponse{
		TaID:        tahunAjaran.TaID,
		TahunAjaran: tahunAjaran.TahunAjaran,
		Semester:    tahunAjaran.Semester,
		Status:      tahunAjaran.Status,
		PublishDate: *tahunAjaran.PublishDate,
		Message:     "Rapor berhasil dipublish",
	}, nil
}

// NEW: Kembalikan ke status aktif (untuk perbaikan nilai)
func (c *tahunAjaranController) ReactivateTahunAjaran(taID uint) (dto.ReactivateTahunAjaranResponse, error) {
	tx := c.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Cek apakah TA ada
	var tahunAjaran model.TahunAjaran
	err := tx.First(&tahunAjaran, taID).Error
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ReactivateTahunAjaranResponse{}, errors.New("tahun ajaran tidak ditemukan")
		}
		return dto.ReactivateTahunAjaranResponse{}, err
	}

	// Cek status (hanya yang selesai bisa direactivate)
	if tahunAjaran.Status != "selesai" {
		tx.Rollback()
		return dto.ReactivateTahunAjaranResponse{}, errors.New("hanya tahun ajaran dengan status SELESAI yang bisa dikembalikan ke AKTIF")
	}

	// Update status dan hapus publish_date
	tahunAjaran.Status = "aktif"
	tahunAjaran.PublishDate = nil

	err = tx.Save(&tahunAjaran).Error
	if err != nil {
		tx.Rollback()
		return dto.ReactivateTahunAjaranResponse{}, err
	}

	err = tx.Commit().Error
	if err != nil {
		return dto.ReactivateTahunAjaranResponse{}, err
	}

	return dto.ReactivateTahunAjaranResponse{
		TaID:        tahunAjaran.TaID,
		TahunAjaran: tahunAjaran.TahunAjaran,
		Semester:    tahunAjaran.Semester,
		Status:      tahunAjaran.Status,
		Message:     "Tahun ajaran berhasil dikembalikan ke status AKTIF. Anda dapat memperbaiki nilai.",
	}, nil
}
