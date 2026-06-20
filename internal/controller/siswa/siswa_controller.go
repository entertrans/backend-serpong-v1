package siswa

import (
	"gorm.io/gorm"

	"github.com/entertrans/backend-bogor.git/internal/dto"
)

type SiswaController interface {
	// Read
	GetAllSiswa() ([]dto.SiswaResponse, error)
	GetSiswaByNIS(nis string) (*dto.SiswaResponse, error)
	GetSiswaByKelasID(req dto.FilterSiswaRequest) (*dto.SiswaListResponse, error) // Tambahkan ini
	// NEW: Filter by status
	GetSiswaAktif(req dto.FilterSiswaStatusRequest) (*dto.SiswaStatusResponse, error)
	GetSiswaAlumni(req dto.FilterSiswaStatusRequest) (*dto.SiswaStatusResponse, error)
	GetSiswaDeleted(req dto.FilterSiswaStatusRequest) (*dto.SiswaStatusResponse, error)
	// ✅ NEW: Khusus untuk frontend (tanpa pagination)
	GetActiveStudentsForEnrollment() ([]dto.SiswaResponse, error)
	// Create
	CreateSiswa(req dto.CreateSiswaRequest) error

	// Update
	UpdateSiswa(nis string, req dto.UpdateSiswaRequest) error
	UpdateOrangtua(nis string, req dto.UpdateOrangtuaRequest) error

	// Delete (opsional)
	// DeleteSiswa(nis string) error
}

type siswaController struct {
	db *gorm.DB
}

func NewSiswaController(db *gorm.DB) SiswaController {
	return &siswaController{
		db: db,
	}
}
