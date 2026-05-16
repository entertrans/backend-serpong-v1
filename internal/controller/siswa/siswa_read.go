package siswa

import (
	"fmt"

	"github.com/entertrans/go-base-project.git/internal/dto"
	"github.com/entertrans/go-base-project.git/internal/model"
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
