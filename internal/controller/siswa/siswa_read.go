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

// controller/siswa/siswa_read.go

// GetSiswaAktif mengambil siswa aktif (bukan alumni dan belum di-delete)
func (c *siswaController) GetSiswaAktif(req dto.FilterSiswaStatusRequest) (*dto.SiswaStatusResponse, error) {
	// Set default values
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	offset := (req.Page - 1) * req.Limit

	// ✅ FIX: Hanya cek NULL, bukan empty string
	query := c.db.Model(&model.Siswa{}).
		Where("soft_deleted = 0").
		Where("siswa_kelas_id IS NOT NULL").
		Where("siswa_kelas_id <= ?", 15).
		Where("tgl_keluar IS NULL"). // ✅ Hanya NULL
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
		return nil, fmt.Errorf("failed to count active students: %w", err)
	}

	// Get paginated data
	var siswaList []model.Siswa
	if err := query.Offset(offset).Limit(req.Limit).Find(&siswaList).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch active students: %w", err)
	}

	// Convert to response
	responses := c.convertToSiswaResponses(siswaList)

	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))

	return &dto.SiswaStatusResponse{
		Data: responses,
		Pagination: dto.PaginationResponse{
			Page:       req.Page,
			Limit:      req.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
		Total: total,
	}, nil
}

// GetSiswaAlumni mengambil alumni (kelas > 15 dan soft_deleted = 0)
func (c *siswaController) GetSiswaAlumni(req dto.FilterSiswaStatusRequest) (*dto.SiswaStatusResponse, error) {
	// Set default values
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	offset := (req.Page - 1) * req.Limit

	// ✅ FIX: Ambil alumni yang:
	// 1. Soft deleted = 0 (belum dihapus)
	// 2. SiswaKelasID > 15 (alumni)
	// ATAU siswa dengan tgl_keluar tidak NULL
	query := c.db.Model(&model.Siswa{}).
		Where("soft_deleted = 0").
		Where("siswa_kelas_id IS NOT NULL").
		Where("siswa_kelas_id > ? OR (tgl_keluar IS NOT NULL AND tgl_keluar != '')", 15).
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
		return nil, fmt.Errorf("failed to count alumni: %w", err)
	}

	// Get paginated data
	var siswaList []model.Siswa
	if err := query.Offset(offset).Limit(req.Limit).Find(&siswaList).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch alumni: %w", err)
	}

	// Convert to response
	responses := c.convertToSiswaResponses(siswaList)

	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))

	return &dto.SiswaStatusResponse{
		Data: responses,
		Pagination: dto.PaginationResponse{
			Page:       req.Page,
			Limit:      req.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
		Total: total,
	}, nil
}

// GetSiswaDeleted mengambil siswa yang sudah di-soft delete (soft_deleted = 1)
func (c *siswaController) GetSiswaDeleted(req dto.FilterSiswaStatusRequest) (*dto.SiswaStatusResponse, error) {
	// Set default values
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	offset := (req.Page - 1) * req.Limit

	// ✅ FIX: Ambil siswa yang sudah di-delete
	query := c.db.Model(&model.Siswa{}).
		Where("soft_deleted = 1").
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
		return nil, fmt.Errorf("failed to count deleted students: %w", err)
	}

	// Get paginated data
	var siswaList []model.Siswa
	if err := query.Offset(offset).Limit(req.Limit).Find(&siswaList).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch deleted students: %w", err)
	}

	// Convert to response
	responses := c.convertToSiswaResponses(siswaList)

	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))

	return &dto.SiswaStatusResponse{
		Data: responses,
		Pagination: dto.PaginationResponse{
			Page:       req.Page,
			Limit:      req.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
		Total: total,
	}, nil
}

// GetActiveStudentsForEnrollment - khusus untuk frontend (ambil semua siswa aktif tanpa pagination)
func (c *siswaController) GetActiveStudentsForEnrollment() ([]dto.SiswaResponse, error) {
	var siswaList []model.Siswa

	// ✅ Hapus kondisi tgl_keluar = ''
	err := c.db.Model(&model.Siswa{}).
		Where("soft_deleted = 0").
		Where("siswa_kelas_id IS NOT NULL").
		Where("siswa_kelas_id <= ?", 15).
		Where("tgl_keluar IS NULL"). // Cuma cek NULL
		Preload("Kelas").
		Preload("Satelit").
		Preload("Orangtua").
		Preload("Lampiran").
		Order("siswa_nama ASC").
		Find(&siswaList).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch active students: %w", err)
	}

	return c.convertToSiswaResponses(siswaList), nil
}

// ============ HELPER FUNCTIONS ============

// convertToSiswaResponses helper untuk convert list siswa ke response
func (c *siswaController) convertToSiswaResponses(siswaList []model.Siswa) []dto.SiswaResponse {
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
	return responses
}
