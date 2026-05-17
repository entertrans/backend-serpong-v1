package siswa

import (
	"fmt"

	"github.com/entertrans/backend-bogor.git/internal/dto"
	"github.com/entertrans/backend-bogor.git/internal/model"
)

func (c *siswaController) GetAllSiswa() ([]dto.SiswaResponse, error) {
	var siswaList []model.Siswa

	err := c.db.Preload("Kelas").
		Preload("Satelit").
		Preload("Orangtua").
		Preload("Lampiran").
		Find(&siswaList).Error

	if err != nil {
		return nil, err
	}

	var responses []dto.SiswaResponse
	for _, siswa := range siswaList {
		// 👇 CONVERT pake MapLampiranToResponse yang sudah ada
		lampiranResponses := make([]dto.LampiranResponse, 0, len(siswa.Lampiran))
		for _, l := range siswa.Lampiran {
			lampiranResponses = append(lampiranResponses, dto.MapLampiranToResponse(l))
		}

		responses = append(responses, dto.SiswaResponse{
			SiswaID:       siswa.SiswaID,
			SiswaNIS:      siswa.SiswaNIS,
			SiswaNISN:     siswa.SiswaNISN,
			SiswaNama:     siswa.SiswaNama,
			NoIjazah:      siswa.NoIjazah,
			SiswaJenkel:   siswa.SiswaJenkel,
			SiswaTempat:   siswa.SiswaTempat,
			SiswaTglLahir: siswa.SiswaTglLahir,
			SiswaAlamat:   siswa.SiswaAlamat,
			SiswaEmail:    siswa.SiswaEmail,
			SiswaNoTelp:   siswa.SiswaNoTelp,
			SiswaDokumen:  siswa.SiswaDokumen,
			TglKeluar:     siswa.TglKeluar,
			TglLulus:      siswa.TglLulus,
			AnakKe:        siswa.AnakKe,
			SiswaKelasID:  siswa.SiswaKelasID,
			SoftDeleted:   siswa.SoftDeleted,
			Kelas:         siswa.Kelas,
			Satelit:       siswa.Satelit,
			Orangtua:      siswa.Orangtua,
			Lampiran:      lampiranResponses, // 👈 PAKAI YANG SUDAH DI-CONVERT
		})
	}

	return responses, nil
}

// Sama untuk GetSiswaByNIS
func (c *siswaController) GetSiswaByNIS(nis string) (*dto.SiswaResponse, error) {
	var siswa model.Siswa

	err := c.db.Preload("Kelas").
		Preload("Satelit").
		Preload("Orangtua").
		Preload("Lampiran").
		Where("siswa_nis = ?", nis).
		First(&siswa).Error

	if err != nil {
		return nil, fmt.Errorf("siswa dengan NIS %s tidak ditemukan", nis)
	}

	// 👇 CONVERT pake MapLampiranToResponse
	lampiranResponses := make([]dto.LampiranResponse, 0, len(siswa.Lampiran))
	for _, l := range siswa.Lampiran {
		lampiranResponses = append(lampiranResponses, dto.MapLampiranToResponse(l))
	}

	return &dto.SiswaResponse{
		SiswaID:       siswa.SiswaID,
		SiswaNIS:      siswa.SiswaNIS,
		SiswaNISN:     siswa.SiswaNISN,
		SiswaNama:     siswa.SiswaNama,
		NoIjazah:      siswa.NoIjazah,
		SiswaJenkel:   siswa.SiswaJenkel,
		SiswaTempat:   siswa.SiswaTempat,
		SiswaTglLahir: siswa.SiswaTglLahir,
		SiswaAlamat:   siswa.SiswaAlamat,
		SiswaEmail:    siswa.SiswaEmail,
		SiswaNoTelp:   siswa.SiswaNoTelp,
		SiswaDokumen:  siswa.SiswaDokumen,
		TglKeluar:     siswa.TglKeluar,
		TglLulus:      siswa.TglLulus,
		AnakKe:        siswa.AnakKe,
		SiswaKelasID:  siswa.SiswaKelasID,
		SoftDeleted:   siswa.SoftDeleted,
		Kelas:         siswa.Kelas,
		Satelit:       siswa.Satelit,
		Orangtua:      siswa.Orangtua,
		Lampiran:      lampiranResponses, // 👈 PAKAI YANG SUDAH DI-CONVERT
	}, nil
}

// GetSiswaByKelasID mengambil siswa berdasarkan kelas ID dengan filter search dan pagination
func (c *siswaController) GetSiswaByKelasID(req dto.FilterSiswaRequest) (*dto.SiswaListResponse, error) {
	// Set default values
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100 // Max limit 100
	}

	offset := (req.Page - 1) * req.Limit

	// Build query
	query := c.db.Model(&model.Siswa{}).
		Where("siswa_kelas_id = ? AND soft_deleted = 0", req.KelasID).
		Preload("Kelas").
		Preload("Satelit").
		Preload("Orangtua").
		Preload("Lampiran")

	// Add search filter if provided
	if req.Search != "" {
		searchPattern := "%" + req.Search + "%"
		query = query.Where("siswa_nis LIKE ? OR siswa_nama LIKE ?", searchPattern, searchPattern)
	}

	// Count total records
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count students: %w", err)
	}

	// Get paginated data
	var siswaList []model.Siswa
	if err := query.Offset(offset).Limit(req.Limit).Find(&siswaList).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch students: %w", err)
	}

	// Calculate total pages
	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))

	// Convert to response DTO
	responses := make([]dto.SiswaResponse, 0, len(siswaList))
	for _, siswa := range siswaList {
		// Convert lampiran
		lampiranResponses := make([]dto.LampiranResponse, 0, len(siswa.Lampiran))
		for _, l := range siswa.Lampiran {
			lampiranResponses = append(lampiranResponses, dto.MapLampiranToResponse(l))
		}

		responses = append(responses, dto.SiswaResponse{
			SiswaID:       siswa.SiswaID,
			SiswaNIS:      siswa.SiswaNIS,
			SiswaNISN:     siswa.SiswaNISN,
			SiswaNama:     siswa.SiswaNama,
			NoIjazah:      siswa.NoIjazah,
			SiswaJenkel:   siswa.SiswaJenkel,
			SiswaTempat:   siswa.SiswaTempat,
			SiswaTglLahir: siswa.SiswaTglLahir,
			SiswaAlamat:   siswa.SiswaAlamat,
			SiswaEmail:    siswa.SiswaEmail,
			SiswaNoTelp:   siswa.SiswaNoTelp,
			SiswaDokumen:  siswa.SiswaDokumen,
			TglKeluar:     siswa.TglKeluar,
			TglLulus:      siswa.TglLulus,
			AnakKe:        siswa.AnakKe,
			SiswaKelasID:  siswa.SiswaKelasID,
			SoftDeleted:   siswa.SoftDeleted,
			Kelas:         siswa.Kelas,
			Satelit:       siswa.Satelit,
			Orangtua:      siswa.Orangtua,
			Lampiran:      lampiranResponses,
		})
	}

	return &dto.SiswaListResponse{
		Data: responses,
		Pagination: dto.PaginationResponse{
			Page:       req.Page,
			Limit:      req.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}
