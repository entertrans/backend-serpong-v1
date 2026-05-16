package siswa

import (
	"gorm.io/gorm"

	"github.com/entertrans/go-base-project.git/internal/dto"

)

type SiswaController interface {
	// Read
	GetAllSiswa() ([]dto.SiswaResponse, error)
	GetSiswaByNIS(nis string) (*dto.SiswaResponse, error)
	
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