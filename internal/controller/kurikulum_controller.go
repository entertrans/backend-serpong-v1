// internal/modules/rapor/controller/kurikulum_controller.go
package controller

import (
	"errors"

	"github.com/entertrans/go-base-project.git/internal/dto"
	"github.com/entertrans/go-base-project.git/internal/model"
	"gorm.io/gorm"
)

type KurikulumController interface {
	// Untuk dropdown
	// GetKelasAktif() ([]dto.KelasAktifResponse, error)
	GetGuruAktif() ([]dto.GuruAktifResponse, error)
	GetMapelAktif() ([]dto.MapelAktifResponse, error)

	// Untuk setup kurikulum
	GetKurikulumByKelas(taID, kelasID uint) (dto.KurikulumByKelasResponse, error)
	SaveKurikulum(req dto.SaveKurikulumRequest) (dto.SaveKurikulumResponse, error)
	CopyKurikulumFromPrevious(req dto.CopyKurikulumRequest, taID, kelasID uint) (dto.KurikulumByKelasResponse, error)
	CheckKurikulumStatus(taID uint) (dto.CheckKurikulumResponse, error)
}

type kurikulumController struct {
	db *gorm.DB
}

func NewKurikulumController(db *gorm.DB) KurikulumController {
	return &kurikulumController{db: db}
}

// ==================== GET DATA UNTUK DROPDOWN ====================

func (c *kurikulumController) GetGuruAktif() ([]dto.GuruAktifResponse, error) {
	var guruList []model.Guru

	err := c.db.Where("status_aktif = ?", true).
		Order("guru_nama ASC").
		Find(&guruList).Error
	if err != nil {
		return nil, err
	}

	result := make([]dto.GuruAktifResponse, len(guruList))
	for i, g := range guruList {
		guruNIP := ""
		if g.GuruNIP != nil {
			guruNIP = *g.GuruNIP
		}
		result[i] = dto.GuruAktifResponse{
			GuruID:   g.GuruID,
			GuruNama: g.GuruNama,
			GuruNIP:  guruNIP,
		}
	}
	return result, nil
}

func (c *kurikulumController) GetMapelAktif() ([]dto.MapelAktifResponse, error) {
	var mapelList []model.Mapel

	err := c.db.Order("kd_mapel ASC").Find(&mapelList).Error
	if err != nil {
		return nil, err
	}

	// Ambil juga data kelompok dan jenjang dari tabel mapel
	// Asumsi: tabel mapel sudah punya field kelompok dan jenjang
	// Jika belum, bisa di-skip dulu atau join dengan tabel lain
	result := make([]dto.MapelAktifResponse, len(mapelList))
	for i, m := range mapelList {
		result[i] = dto.MapelAktifResponse{
			KdMapel: m.KdMapel,
			NmMapel: m.NmMapel,
			// Kelompok: m.Kelompok,
			// Jenjang:  m.Jenjang,
		}
	}
	return result, nil
}

// ==================== KURIKULUM SETUP ====================

// GetKurikulumByKelas - mengambil data kurikulum per TA per Kelas
func (c *kurikulumController) GetKurikulumByKelas(taID, kelasID uint) (dto.KurikulumByKelasResponse, error) {
	// Ambil data mapel
	var taKelasMapelList []model.TaKelasMapel
	err := c.db.
		Preload("Mapel").
		Preload("Guru").
		Where("ta_id = ? AND kelas_id = ?", taID, kelasID).
		Order("urutan ASC").
		Find(&taKelasMapelList).Error
	if err != nil {
		return dto.KurikulumByKelasResponse{}, err
	}

	// Ambil data wali kelas (tanpa pointer, cek primary key)
	var waliKelas model.TaKelasWali
	var waliKelasID *uint
	var waliKelasNama string

	err = c.db.
		Preload("Guru").
		Where("ta_id = ? AND kelas_id = ?", taID, kelasID).
		First(&waliKelas).Error

	if err == nil {
		// Cek apakah Guru terisi (primary key tidak 0)
		if waliKelas.Guru.GuruID != 0 {
			waliKelasID = &waliKelas.GuruID
			waliKelasNama = waliKelas.Guru.GuruNama
		}
	}

	// Ambil info kelas
	var kelas model.Kelas
	err = c.db.First(&kelas, kelasID).Error
	if err != nil {
		return dto.KurikulumByKelasResponse{}, err
	}

	// Map mapel list
	mapelList := make([]dto.TaKelasMapelItem, len(taKelasMapelList))
	for i, item := range taKelasMapelList {
		guruNama := ""
		// Cek apakah Guru terisi (primary key tidak 0)
		if item.Guru.GuruID != 0 {
			guruNama = item.Guru.GuruNama
		}

		mapelList[i] = dto.TaKelasMapelItem{
			TaKelasMapelID: item.TaKelasMapelID,
			KdMapel:        item.KdMapel,
			NmMapel:        item.Mapel.NmMapel,
			Kkm:            75,
			GuruID:         item.GuruID,
			GuruNama:       guruNama,
			Urutan:         item.Urutan,
		}
	}

	return dto.KurikulumByKelasResponse{
		KelasID:       kelasID,
		KelasNama:     kelas.KelasNama,
		WaliKelasID:   waliKelasID,
		WaliKelasNama: waliKelasNama,
		MapelList:     mapelList,
	}, nil
}

// SaveKurikulum - menyimpan kurikulum untuk satu kelas
func (c *kurikulumController) SaveKurikulum(req dto.SaveKurikulumRequest) (dto.SaveKurikulumResponse, error) {
	tx := c.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Hapus data mapel lama
	err := tx.Where("ta_id = ? AND kelas_id = ?", req.TaID, req.KelasID).
		Delete(&model.TaKelasMapel{}).Error
	if err != nil {
		tx.Rollback()
		return dto.SaveKurikulumResponse{}, err
	}

	// 2. Insert data mapel baru
	for _, item := range req.MapelList {
		taKelasMapel := model.TaKelasMapel{
			TaID:    req.TaID,
			KelasID: req.KelasID,
			KdMapel: item.KdMapel,
			GuruID:  item.GuruID, // ← item.GuruID sudah *uint
			Urutan:  item.Urutan,
		}
		err = tx.Create(&taKelasMapel).Error
		if err != nil {
			tx.Rollback()
			return dto.SaveKurikulumResponse{}, err
		}
	}

	// 3. Handle wali kelas (UPSERT)
	// Perhatikan: WaliKelasID bisa nil
	if req.WaliKelasID != nil && *req.WaliKelasID > 0 {
		// Cek apakah sudah ada
		var existing model.TaKelasWali
		err = tx.Where("ta_id = ? AND kelas_id = ?", req.TaID, req.KelasID).
			First(&existing).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Insert baru
			wali := model.TaKelasWali{
				TaID:    req.TaID,
				KelasID: req.KelasID,
				GuruID:  *req.WaliKelasID, // ← *req.WaliKelasID
			}
			err = tx.Create(&wali).Error
		} else if err == nil {
			// Update existing
			existing.GuruID = *req.WaliKelasID
			err = tx.Save(&existing).Error
		}

		if err != nil {
			tx.Rollback()
			return dto.SaveKurikulumResponse{}, err
		}
	} else {
		// Jika wali_kelas_id kosong/null, hapus data wali jika ada
		err = tx.Where("ta_id = ? AND kelas_id = ?", req.TaID, req.KelasID).
			Delete(&model.TaKelasWali{}).Error
		if err != nil {
			tx.Rollback()
			return dto.SaveKurikulumResponse{}, err
		}
	}

	err = tx.Commit().Error
	if err != nil {
		return dto.SaveKurikulumResponse{}, err
	}

	return dto.SaveKurikulumResponse{
		Message:    "Kurikulum berhasil disimpan",
		TaID:       req.TaID,
		KelasID:    req.KelasID,
		TotalMapel: len(req.MapelList),
	}, nil
}

// CopyKurikulumFromPrevious - copy data kurikulum dari TA sebelumnya
func (c *kurikulumController) CopyKurikulumFromPrevious(req dto.CopyKurikulumRequest, taID, kelasID uint) (dto.KurikulumByKelasResponse, error) {
	var previousMapel []model.TaKelasMapel

	err := c.db.
		Preload("Mapel").
		Preload("Guru").
		Where("ta_id = ? AND kelas_id = ?", req.FromTaID, req.ToKelasID).
		Order("urutan ASC").
		Find(&previousMapel).Error
	if err != nil {
		return dto.KurikulumByKelasResponse{}, err
	}

	// Konversi ke response
	mapelList := make([]dto.TaKelasMapelItem, len(previousMapel))
	for i, item := range previousMapel {
		guruNama := ""
		if item.Guru != nil {
			guruNama = item.Guru.GuruNama
		}

		mapelList[i] = dto.TaKelasMapelItem{
			TaKelasMapelID: 0, // ID baru nanti
			KdMapel:        item.KdMapel,
			NmMapel:        item.Mapel.NmMapel,
			Kkm:            75,
			GuruID:         item.GuruID,
			GuruNama:       guruNama,
			Urutan:         item.Urutan,
		}
	}

	// Ambil info kelas
	var kelas model.Kelas
	err = c.db.First(&kelas, kelasID).Error
	if err != nil {
		return dto.KurikulumByKelasResponse{}, err
	}

	return dto.KurikulumByKelasResponse{
		KelasID:       kelasID,
		KelasNama:     kelas.KelasNama,
		WaliKelasID:   nil,
		WaliKelasNama: "",
		MapelList:     mapelList,
	}, nil
}

// CheckKurikulumStatus - cek kelas mana saja yang belum di setup untuk TA tertentu
func (c *kurikulumController) CheckKurikulumStatus(taID uint) (dto.CheckKurikulumResponse, error) {
	// Ambil semua kelas aktif (1-12)
	var allKelas []model.Kelas
	err := c.db.Where("kelas_id BETWEEN ? AND ?", 1, 12).
		Find(&allKelas).Error
	if err != nil {
		return dto.CheckKurikulumResponse{}, err
	}

	// Ambil kelas yang sudah punya kurikulum di TA ini
	var setupKelas []struct {
		KelasID uint
	}
	err = c.db.Model(&model.TaKelasMapel{}).
		Where("ta_id = ?", taID).
		Group("kelas_id").
		Select("kelas_id").
		Find(&setupKelas).Error
	if err != nil {
		return dto.CheckKurikulumResponse{}, err
	}

	// Buat map untuk pengecekan cepat
	setupMap := make(map[uint]bool)
	for _, sk := range setupKelas {
		setupMap[sk.KelasID] = true
	}

	// Cari kelas yang belum setup
	var belumSetup []string
	for _, k := range allKelas {
		if !setupMap[k.KelasId] {
			belumSetup = append(belumSetup, k.KelasNama)
		}
	}

	return dto.CheckKurikulumResponse{
		BelumSetup: belumSetup,
		TotalKelas: len(allKelas),
		SudahSetup: len(setupKelas),
	}, nil
}
