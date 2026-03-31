package controller

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/entertrans/go-base-project.git/internal/dto"
	"github.com/entertrans/go-base-project.git/internal/model"

)

type SiswaController interface {
	GetAllSiswa() ([]dto.SiswaResponse, error)
	CreateSiswa(req dto.CreateSiswaRequest) error // <-- TAMBAH INI
}

type siswaController struct {
	db *gorm.DB
}

func NewSiswaController(db *gorm.DB) SiswaController {
	return &siswaController{
		db: db,
	}
}
func (c *siswaController) GetAllSiswa() ([]dto.SiswaResponse, error) {
	var siswaList []model.Siswa

	// Query dengan preload kelas
	err := c.db.Preload("Kelas").
		Preload("Satelit").
		Preload("Orangtua").
		Preload("Lampiran").
		Find(&siswaList).Error

	if err != nil {
		return nil, err
	}

	// Convert ke response DTO
	var responses []dto.SiswaResponse
	for _, siswa := range siswaList {
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
			// SiswaPhoto:    siswa.SiswaPhoto,
			TglKeluar:    siswa.TglKeluar,
			TglLulus:     siswa.TglLulus,
			AnakKe:       siswa.AnakKe,
			SiswaKelasID: siswa.SiswaKelasID,
			SoftDeleted:  siswa.SoftDeleted,
			Kelas:        siswa.Kelas,
			Satelit:      siswa.Satelit,
			Orangtua:     siswa.Orangtua,
			Lampiran:     siswa.Lampiran,
		})
	}

	return responses, nil
}


//create siswa + orangtua dalam 1 transaksi
func (c *siswaController) CreateSiswa(req dto.CreateSiswaRequest) error {
	fmt.Println("🚀 MASUK CONTROLLER")

	tx := c.db.Begin()

	// ======================
	// INSERT SISWA
	// ======================
	siswa := model.Siswa{
		SiswaNIS:      req.SiswaNIS,
		SiswaNISN:     req.SiswaNISN,
		SiswaNama:     req.SiswaNama,
		SiswaJenkel:   req.SiswaJenkel,
		SiswaTempat:   req.SiswaTempat,
		SiswaTglLahir: req.SiswaTglLahir,
		SiswaAlamat:   req.SiswaAlamat,
		SiswaEmail:    req.SiswaEmail,
		SiswaNoTelp:   req.SiswaNoTelp,
		SiswaKelasID:  req.SiswaKelasID,
		AnakKe:        req.AnakKe,
		NoIjazah:      req.NoIjazah,
	}

	if err := tx.Create(&siswa).Error; err != nil {
		tx.Rollback()
		return err
	}

	fmt.Println("✅ SISWA KE-INSERT")

	// ======================
	// INSERT ORANGTUA (OPTIONAL)
	// ======================
	ortu := req.Orangtua

	// cek apakah ada data minimal (contoh: nama atau no telp)
	if ortu.AyahNama != nil || ortu.NoTelpAyah != nil {
		fmt.Println("👨‍👩‍👧 INSERT ORTU")

		ortuModel := model.Orangtua{
			SiswaNIS:   req.SiswaNIS,
			AyahNama:   ortu.AyahNama,
			AyahNotelp: ortu.NoTelpAyah,
			AyahNik:    ortu.AyahNik,
			// tambahin sesuai kebutuhan
		}

		if err := tx.Create(&ortuModel).Error; err != nil {
			fmt.Println("❌ GAGAL INSERT ORTU:", err)
			tx.Rollback()
			return err
		}
	} else {
		fmt.Println("⏭️ SKIP INSERT ORTU (kosong)")
	}

	return tx.Commit().Error
}