// internal/modules/rapor/controller/kurikulum_controller.go
package controller

import (
	"errors"

	"github.com/entertrans/backend-bogor.git/internal/dto"
	"github.com/entertrans/backend-bogor.git/internal/model"
	"gorm.io/gorm"
)

type KurikulumController interface {

	// Untuk setup kurikulum
	GetKurikulumByKelas(taID, kelasID uint) (dto.KurikulumByKelasResponse, error)
	SaveKurikulum(req dto.SaveKurikulumRequest) (dto.SaveKurikulumResponse, error)
	CopyKurikulumFromPrevious(req dto.CopyKurikulumRequest, taID, kelasID uint) (dto.KurikulumByKelasResponse, error)
	CheckKurikulumStatus(taID uint) (dto.CheckKurikulumResponse, error)
	GetKurikulumByGuru(taID, kelasID, guruID uint) (dto.KurikulumByGuruResponse, error)
	GetKelasWali(taID, guruID uint) (dto.GetKelasWaliResponse, error)
}

type kurikulumController struct {
	db *gorm.DB
}

func NewKurikulumController(db *gorm.DB) KurikulumController {
	return &kurikulumController{db: db}
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

// GetKurikulumByGuru - mengambil data mapel yang diajar oleh guru tertentu di suatu kelas
func (c *kurikulumController) GetKurikulumByGuru(taID, kelasID, guruID uint) (dto.KurikulumByGuruResponse, error) {
	var mapelList []model.TaKelasMapel

	// Ambil data mapel berdasarkan ta_id, kelas_id, dan guru_id
	err := c.db.
		Preload("Mapel").
		Preload("Guru").
		Where("ta_id = ? AND kelas_id = ? AND guru_id = ?", taID, kelasID, guruID).
		Order("urutan ASC").
		Find(&mapelList).Error
	if err != nil {
		return dto.KurikulumByGuruResponse{}, err
	}

	// Jika tidak ada data
	if len(mapelList) == 0 {
		return dto.KurikulumByGuruResponse{
			Message:   "Tidak ada mapel yang diajar oleh guru ini di kelas tersebut",
			MapelList: []dto.MapelByGuruItem{},
		}, nil
	}

	// Ambil info guru dari data pertama
	var guruNama string
	var guruNIP string // string, bukan pointer untuk response
	if len(mapelList) > 0 && mapelList[0].Guru.GuruID != 0 {
		guruNama = mapelList[0].Guru.GuruNama
		if mapelList[0].Guru.GuruNIP != nil {
			guruNIP = *mapelList[0].Guru.GuruNIP
		}
	}

	// Ambil info kelas
	var kelas model.Kelas
	err = c.db.First(&kelas, kelasID).Error
	if err != nil {
		return dto.KurikulumByGuruResponse{}, err
	}

	// Ambil info tahun ajaran
	var tahunAjaran model.TahunAjaran
	err = c.db.First(&tahunAjaran, taID).Error
	if err != nil {
		return dto.KurikulumByGuruResponse{}, err
	}

	// Format nama tahun ajaran (contoh: "2024/2025 - Semester 1")
	taNama := tahunAjaran.TahunAjaran
	if tahunAjaran.Semester == "1" {
		taNama = taNama + " - Semester Ganjil"
	} else {
		taNama = taNama + " - Semester Genap"
	}

	// Konversi ke response
	mapelItems := make([]dto.MapelByGuruItem, len(mapelList))
	for i, item := range mapelList {
		mapelItems[i] = dto.MapelByGuruItem{
			TaKelasMapelID: item.TaKelasMapelID,
			KdMapel:        item.KdMapel,
			NmMapel:        item.Mapel.NmMapel,
			KKM:            75, // Sesuaikan dengan logic KKM jika ada di model lain
			Urutan:         item.Urutan,
		}
	}

	return dto.KurikulumByGuruResponse{
		TaID:       taID,
		TaNama:     taNama,
		KelasID:    kelasID,
		KelasNama:  kelas.KelasNama,
		GuruID:     guruID,
		GuruNama:   guruNama,
		GuruNIP:    guruNIP,
		TotalMapel: len(mapelItems),
		MapelList:  mapelItems,
	}, nil
}

// Implementasi method GetKelasWali
func (c *kurikulumController) GetKelasWali(taID, guruID uint) (dto.GetKelasWaliResponse, error) {
	var kelasWaliList []model.TaKelasWali

	// Ambil data kelas wali berdasarkan ta_id dan guru_id
	err := c.db.
		Preload("Kelas").
		Where("ta_id = ? AND guru_id = ?", taID, guruID).
		Find(&kelasWaliList).Error

	if err != nil {
		return dto.GetKelasWaliResponse{}, err
	}

	// Ambil info tahun ajaran
	var tahunAjaran model.TahunAjaran
	err = c.db.First(&tahunAjaran, taID).Error
	if err != nil {
		return dto.GetKelasWaliResponse{}, err
	}

	// Format nama tahun ajaran
	taNama := tahunAjaran.TahunAjaran
	if tahunAjaran.Semester == "1" {
		taNama = taNama + " - Semester Ganjil"
	} else {
		taNama = taNama + " - Semester Genap"
	}

	// Ambil info guru
	var guru model.Guru
	err = c.db.First(&guru, guruID).Error
	if err != nil {
		return dto.GetKelasWaliResponse{}, err
	}

	guruNama := guru.GuruNama
	guruNIP := ""
	if guru.GuruNIP != nil {
		guruNIP = *guru.GuruNIP
	}

	// Konversi ke response
	kelasList := make([]dto.KelasWaliItem, len(kelasWaliList))
	for i, item := range kelasWaliList {
		kelasNama := ""
		if item.Kelas != nil {
			kelasNama = item.Kelas.KelasNama
		}
		kelasList[i] = dto.KelasWaliItem{
			KelasID:   item.KelasID,
			KelasNama: kelasNama,
		}
	}

	return dto.GetKelasWaliResponse{
		TaID:       taID,
		TaNama:     taNama,
		Semester:   tahunAjaran.Semester,
		Status:     tahunAjaran.Status,
		GuruID:     guruID,
		GuruNama:   guruNama,
		GuruNIP:    guruNIP,
		TotalKelas: len(kelasList),
		KelasList:  kelasList,
	}, nil
}
