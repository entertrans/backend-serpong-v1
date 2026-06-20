// internal/modules/rapor/controller/penilaian_controller.go
package controller

import (
	"errors"
	"time"

	"github.com/entertrans/backend-bogor.git/internal/dto"
	"github.com/entertrans/backend-bogor.git/internal/model"
	"gorm.io/gorm"
)

// internal/modules/rapor/controller/penilaian_controller.go

type PenilaianController interface {
	GetSiswaByKelas(taID, kelasID uint) ([]dto.SiswaByKelasResponse, error)
	GetRaportBySiswa(taID, kelasID uint, siswaNIS string) (dto.RaportResponse, error)

	GetNilaiMapel(taID, taKelasMapelID uint, kelasID uint) ([]dto.GetNilaiMapelResponse, error)

	SaveNilaiMapel(
		req dto.NilaiMapelRequest,
		userID uint,
	) (dto.NilaiMapelResponse, error)

	GetAbsensi(taID, kelasID uint) ([]dto.AbsensiItem, error)
	SaveAbsensi(req dto.AbsensiRequest) (dto.AbsensiResponse, error)

	GetEkskul(taID, kelasID uint, siswaNIS string) ([]dto.EkskulItem, error)
	SaveEkskul(req dto.EkskulRequest) (dto.EkskulResponse, error)

	UpdateStatusPublish(taID, kelasID uint, siswaNIS string, statusPublish int) error
	GetNilaiHistory(raportNilaiID uint) ([]dto.NilaiHistoryResponse, error)
	EditNilaiPerSiswa(req dto.EditNilaiPerSiswaRequest) error
}

type penilaianController struct {
	db *gorm.DB
}

func NewPenilaianController(db *gorm.DB) PenilaianController {
	return &penilaianController{db: db}
}

// ==================== GET SISWA BY KELAS ====================

func (c *penilaianController) GetSiswaByKelas(taID, kelasID uint) ([]dto.SiswaByKelasResponse, error) {
	var siswaList []model.Siswa

	// Ambil siswa berdasarkan kelas_id dan tidak soft deleted
	err := c.db.Where("siswa_kelas_id = ? AND (soft_deleted IS NULL OR soft_deleted = ?)", kelasID, 0).
		Order("siswa_nama ASC").
		Find(&siswaList).Error
	if err != nil {
		return nil, err
	}

	// Ambil semua raport untuk TA dan kelas ini
	var raportList []model.Raport
	err = c.db.Where("ta_id = ? AND kelas_id = ?", taID, kelasID).
		Find(&raportList).Error
	if err != nil {
		return nil, err
	}

	// Buat map status_publish berdasarkan siswa_nis
	statusMap := make(map[string]bool)
	for _, r := range raportList {
		statusMap[r.SiswaNIS] = r.StatusPublish
	}

	result := make([]dto.SiswaByKelasResponse, 0, len(siswaList))
	for _, s := range siswaList {
		siswaNama := ""
		if s.SiswaNama != nil {
			siswaNama = *s.SiswaNama
		}
		siswaNIS := ""
		if s.SiswaNIS != nil {
			siswaNIS = *s.SiswaNIS
		}
		siswaNISN := ""
		if s.SiswaNISN != nil {
			siswaNISN = *s.SiswaNISN
		}

		result = append(result, dto.SiswaByKelasResponse{
			SiswaNIS:      siswaNIS,
			SiswaNama:     siswaNama,
			SiswaNISN:     siswaNISN,
			StatusPublish: statusMap[siswaNIS], // ← TAMBAHKAN
		})
	}

	return result, nil
}

// ==================== NILAI MAPEL ====================

func (c *penilaianController) GetNilaiMapel(taID, taKelasMapelID uint, kelasID uint) ([]dto.GetNilaiMapelResponse, error) {
	// Ambil semua siswa di kelas
	var siswaList []model.Siswa
	err := c.db.Where("siswa_kelas_id = ? AND (soft_deleted IS NULL OR soft_deleted = ?)", kelasID, 0).
		Find(&siswaList).Error
	if err != nil {
		return nil, err
	}

	// Ambil raport_id untuk setiap siswa di TA ini
	var raportList []model.Raport
	err = c.db.Where("ta_id = ? AND kelas_id = ?", taID, kelasID).
		Find(&raportList).Error
	if err != nil {
		return nil, err
	}

	// Buat map raport_id berdasarkan siswa_nis
	raportMap := make(map[string]uint)
	for _, r := range raportList {
		raportMap[r.SiswaNIS] = r.RaportID
	}

	// Ambil nilai yang sudah ada
	var nilaiList []model.RaportNilai
	err = c.db.Where("ta_kelas_mapel_id = ?", taKelasMapelID).
		Find(&nilaiList).Error
	if err != nil {
		return nil, err
	}

	// Buat map nilai berdasarkan raport_id
	nilaiMap := make(map[uint]model.RaportNilai)
	for _, n := range nilaiList {
		nilaiMap[n.RaportID] = n
	}

	// Gabungkan data
	result := make([]dto.GetNilaiMapelResponse, 0)
	for _, siswa := range siswaList {
		siswaNIS := ""
		if siswa.SiswaNIS != nil {
			siswaNIS = *siswa.SiswaNIS
		}
		siswaNama := ""
		if siswa.SiswaNama != nil {
			siswaNama = *siswa.SiswaNama
		}

		raportID, exists := raportMap[siswaNIS]
		if !exists {
			// Jika raport tidak ada, tetap tampilkan siswa dengan nilai kosong
			result = append(result, dto.GetNilaiMapelResponse{
				RaportNilaiID: 0, // 0 berarti belum ada nilai
				SiswaNIS:      siswaNIS,
				SiswaNama:     siswaNama,
				NilaiAngka:    0,
				Deskripsi:     "",
			})
			continue
		}

		nilai, hasNilai := nilaiMap[raportID]
		if hasNilai {
			// Sudah ada nilai
			result = append(result, dto.GetNilaiMapelResponse{
				RaportNilaiID: nilai.RaportNilaiID,
				SiswaNIS:      siswaNIS,
				SiswaNama:     siswaNama,
				NilaiAngka:    nilai.NilaiAngka,
				Deskripsi:     nilai.Deskripsi,
			})
		} else {
			// Belum ada nilai, tapi raport sudah ada
			result = append(result, dto.GetNilaiMapelResponse{
				RaportNilaiID: 0, // 0 berarti belum ada nilai
				SiswaNIS:      siswaNIS,
				SiswaNama:     siswaNama,
				NilaiAngka:    0,
				Deskripsi:     "",
			})
		}
	}

	return result, nil
}

func (c *penilaianController) SaveNilaiMapel(
	req dto.NilaiMapelRequest,
	userID uint,
) (dto.NilaiMapelResponse, error) {
	tx := c.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, item := range req.NilaiList {
		// Cek apakah raport sudah ada
		var raport model.Raport
		err := tx.Where("siswa_nis = ? AND ta_id = ? AND kelas_id = ?",
			item.SiswaNIS, req.TaID, req.KelasID).
			First(&raport).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Buat raport baru
			raport = model.Raport{
				SiswaNIS: item.SiswaNIS,
				KelasID:  req.KelasID,
				TaID:     req.TaID,
				Sakit:    0,
				Izin:     0,
				Alpha:    0,
			}
			err = tx.Create(&raport).Error
			if err != nil {
				tx.Rollback()
				return dto.NilaiMapelResponse{}, err
			}
		} else if err != nil {
			tx.Rollback()
			return dto.NilaiMapelResponse{}, err
		}

		// Hitung nilai huruf dan predikat
		nilaiHuruf := konversiNilaiHuruf(item.NilaiAngka)
		predikat := konversiPredikat(item.NilaiAngka)

		if err != nil {
			tx.Rollback()
			return dto.NilaiMapelResponse{}, err
		}

		// Cek apakah nilai sudah ada
		var existing model.RaportNilai
		err = tx.Where("raport_id = ? AND ta_kelas_mapel_id = ?",
			raport.RaportID, req.TaKelasMapelID).
			First(&existing).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Insert baru
			nilai := model.RaportNilai{
				RaportID:       raport.RaportID,
				TaKelasMapelID: req.TaKelasMapelID,
				NilaiAngka:     item.NilaiAngka,
				NilaiHuruf:     nilaiHuruf,
				Predikat:       predikat,
				Deskripsi:      item.Deskripsi,

				CreatedBy: &userID,
			}
			err = tx.Create(&nilai).Error
		} else if err == nil {
			// existing.LogUpdate = datatypes.JSON(logJSON)
			nilaiChanged := existing.NilaiAngka != item.NilaiAngka

			shouldAudit :=
				nilaiChanged &&
					!(existing.NilaiAngka == 0 && item.NilaiAngka > 0)

			if shouldAudit {

				audit := model.RaportNilaiAudit{
					RaportNilaiID: existing.RaportNilaiID,
					OldValue:      existing.NilaiAngka,
					NewValue:      item.NilaiAngka,
					ChangedBy:     userID,
					ChangedAt:     time.Now(),
				}

				if err := tx.Create(&audit).Error; err != nil {
					tx.Rollback()
					return dto.NilaiMapelResponse{}, err
				}
			}

			existing.UpdatedBy = &userID
			existing.NilaiAngka = item.NilaiAngka
			existing.NilaiHuruf = nilaiHuruf
			existing.Predikat = predikat
			existing.Deskripsi = item.Deskripsi

			err = tx.Save(&existing).Error
		}

		if err != nil {
			tx.Rollback()
			return dto.NilaiMapelResponse{}, err
		}
	}

	err := tx.Commit().Error
	if err != nil {
		return dto.NilaiMapelResponse{}, err
	}

	return dto.NilaiMapelResponse{
		Message:        "Nilai berhasil disimpan",
		TotalData:      len(req.NilaiList),
		TaKelasMapelID: req.TaKelasMapelID,
	}, nil
}

// ==================== ABSENSI & CATATAN ====================

func (c *penilaianController) GetAbsensi(taID, kelasID uint) ([]dto.AbsensiItem, error) {
	var raportList []model.Raport

	err := c.db.Where("ta_id = ? AND kelas_id = ?", taID, kelasID).
		Find(&raportList).Error
	if err != nil {
		return nil, err
	}

	result := make([]dto.AbsensiItem, len(raportList))
	for i, r := range raportList {
		result[i] = dto.AbsensiItem{
			SiswaNIS:         r.SiswaNIS,
			Sakit:            r.Sakit,
			Izin:             r.Izin,
			Alpha:            r.Alpha,
			CatatanWaliKelas: r.CatatanWaliKelas,
		}
	}

	return result, nil
}

func (c *penilaianController) SaveAbsensi(req dto.AbsensiRequest) (dto.AbsensiResponse, error) {
	tx := c.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, item := range req.AbsensiList {
		var raport model.Raport
		err := tx.Where("siswa_nis = ? AND ta_id = ? AND kelas_id = ?",
			item.SiswaNIS, req.TaID, req.KelasID).
			First(&raport).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			raport = model.Raport{
				SiswaNIS:         item.SiswaNIS,
				KelasID:          req.KelasID,
				TaID:             req.TaID,
				Sakit:            item.Sakit,
				Izin:             item.Izin,
				Alpha:            item.Alpha,
				CatatanWaliKelas: item.CatatanWaliKelas,
			}
			err = tx.Create(&raport).Error
		} else if err == nil {
			raport.Sakit = item.Sakit
			raport.Izin = item.Izin
			raport.Alpha = item.Alpha
			raport.CatatanWaliKelas = item.CatatanWaliKelas
			err = tx.Save(&raport).Error
		}

		if err != nil {
			tx.Rollback()
			return dto.AbsensiResponse{}, err
		}
	}

	err := tx.Commit().Error
	if err != nil {
		return dto.AbsensiResponse{}, err
	}

	return dto.AbsensiResponse{
		Message:   "Absensi berhasil disimpan",
		TotalData: len(req.AbsensiList),
	}, nil
}

// ==================== EKSKUL ====================

func (c *penilaianController) GetEkskul(taID, kelasID uint, siswaNIS string) ([]dto.EkskulItem, error) {
	var raport model.Raport
	err := c.db.Where("siswa_nis = ? AND ta_id = ? AND kelas_id = ?",
		siswaNIS, taID, kelasID).
		First(&raport).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []dto.EkskulItem{}, nil
		}
		return nil, err
	}

	var ekskulList []model.RaportEkskul
	err = c.db.Where("raport_id = ?", raport.RaportID).
		Find(&ekskulList).Error
	if err != nil {
		return nil, err
	}

	result := make([]dto.EkskulItem, len(ekskulList))
	for i, e := range ekskulList {
		result[i] = dto.EkskulItem{
			SiswaNIS:   siswaNIS,
			NamaEkskul: e.NamaEkskul,
			Nilai:      e.Nilai,
			Deskripsi:  e.Deskripsi,
		}
	}

	return result, nil
}

func (c *penilaianController) SaveEkskul(req dto.EkskulRequest) (dto.EkskulResponse, error) {
	tx := c.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Kelompokkan berdasarkan siswa
	siswaMap := make(map[string][]dto.EkskulItem)
	for _, item := range req.EkskulList {
		siswaMap[item.SiswaNIS] = append(siswaMap[item.SiswaNIS], item)
	}

	for siswaNIS, ekskulItems := range siswaMap {
		var raport model.Raport
		err := tx.Where("siswa_nis = ? AND ta_id = ? AND kelas_id = ?",
			siswaNIS, req.TaID, req.KelasID).
			First(&raport).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			raport = model.Raport{
				SiswaNIS: siswaNIS,
				KelasID:  req.KelasID,
				TaID:     req.TaID,
			}
			err = tx.Create(&raport).Error
			if err != nil {
				tx.Rollback()
				return dto.EkskulResponse{}, err
			}
		} else if err != nil {
			tx.Rollback()
			return dto.EkskulResponse{}, err
		}

		// Hapus data ekskul lama
		err = tx.Where("raport_id = ?", raport.RaportID).
			Delete(&model.RaportEkskul{}).Error
		if err != nil {
			tx.Rollback()
			return dto.EkskulResponse{}, err
		}

		// Insert data baru
		for _, item := range ekskulItems {
			ekskul := model.RaportEkskul{
				RaportID:   raport.RaportID,
				NamaEkskul: item.NamaEkskul,
				Nilai:      item.Nilai,
				Deskripsi:  item.Deskripsi,
			}
			err = tx.Create(&ekskul).Error
			if err != nil {
				tx.Rollback()
				return dto.EkskulResponse{}, err
			}
		}
	}

	err := tx.Commit().Error
	if err != nil {
		return dto.EkskulResponse{}, err
	}

	return dto.EkskulResponse{
		Message:   "Ekstrakurikuler berhasil disimpan",
		TotalData: len(req.EkskulList),
	}, nil
}

// ==================== GET RAPORT BY SISWA ====================

func (c *penilaianController) GetRaportBySiswa(taID, kelasID uint, siswaNIS string) (dto.RaportResponse, error) {
	var raport model.Raport
	err := c.db.Where("siswa_nis = ? AND ta_id = ? AND kelas_id = ?",
		siswaNIS, taID, kelasID).
		First(&raport).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.RaportResponse{}, errors.New("raport tidak ditemukan")
		}
		return dto.RaportResponse{}, err
	}

	// Ambil data siswa
	var siswa model.Siswa
	err = c.db.Preload("Kelas").First(&siswa, "siswa_nis = ?", siswaNIS).Error
	if err != nil {
		return dto.RaportResponse{}, err
	}

	// Ambil data TA
	var ta model.TahunAjaran
	err = c.db.First(&ta, taID).Error
	if err != nil {
		return dto.RaportResponse{}, err
	}
	// Ambil wali kelas
	var waliKelas model.TaKelasWali
	var waliKelasNama string
	err = c.db.Preload("Guru").
		Where("ta_id = ? AND kelas_id = ?", taID, kelasID).
		First(&waliKelas).Error

	if err == nil && waliKelas.Guru.GuruID != 0 {
		waliKelasNama = waliKelas.Guru.GuruNama
	}

	// Ambil nilai mapel dengan preload relasi (tambah preload Mapel agar dapet kelompok)
	var nilaiMapel []model.RaportNilai
	err = c.db.Preload("TaKelasMapel.Mapel").
		Where("raport_id = ?", raport.RaportID).
		Find(&nilaiMapel).Error
	if err != nil {
		return dto.RaportResponse{}, err
	}
	// Konversi nilai mapel (dengan pengecekan yang benar)
	nilaiMapelResp := make([]dto.NilaiMapelDetail, 0, len(nilaiMapel))
	for _, n := range nilaiMapel {
		nmMapel := ""
		kelompok := ""

		// Cek apakah Mapel terisi (KdMapel tidak 0 berarti ada data)
		if n.TaKelasMapel.Mapel.KdMapel != 0 {
			nmMapel = n.TaKelasMapel.Mapel.NmMapel
			kelompok = n.TaKelasMapel.Mapel.Kelompok // ambil kelompok dari mapel
		}

		nilaiMapelResp = append(nilaiMapelResp, dto.NilaiMapelDetail{
			KdMapel:    n.TaKelasMapelID,
			NmMapel:    nmMapel,
			Kelompok:   kelompok, // tambahkan ini
			NilaiAngka: n.NilaiAngka,
			NilaiHuruf: n.NilaiHuruf,
			Predikat:   n.Predikat,
			Deskripsi:  n.Deskripsi,
		})
	}

	// Ambil nilai ekskul
	var ekskulList []model.RaportEkskul
	err = c.db.Where("raport_id = ?", raport.RaportID).
		Find(&ekskulList).Error
	if err != nil {
		return dto.RaportResponse{}, err
	}

	ekskulResp := make([]dto.EkskulDetail, len(ekskulList))
	for i, e := range ekskulList {
		ekskulResp[i] = dto.EkskulDetail{
			NamaEkskul: e.NamaEkskul,
			Nilai:      e.Nilai,
			Deskripsi:  e.Deskripsi,
		}
	}
	// Format publish_date (dari TA)
	var publishDateStr *string
	if ta.PublishDate != nil {
		formatted := ta.PublishDate.Format("2006-01-02") // Format ISO untuk frontend
		publishDateStr = &formatted
	}

	// Ambil nama kelas
	kelasNama := ""
	if siswa.Kelas.KelasId != 0 {
		kelasNama = siswa.Kelas.KelasNama
	}

	// Ambil ID kelas
	kelasIDResp := kelasID
	if siswa.Kelas.KelasId != 0 {
		kelasIDResp = siswa.Kelas.KelasId
	}

	// Ambil nama siswa dari pointer
	siswaNama := ""
	if siswa.SiswaNama != nil {
		siswaNama = *siswa.SiswaNama
	}
	siswaNISNVal := ""
	if siswa.SiswaNISN != nil {
		siswaNISNVal = *siswa.SiswaNISN
	}

	return dto.RaportResponse{
		RaportID:         raport.RaportID,
		KelasID:          kelasIDResp,
		SiswaNIS:         raport.SiswaNIS,
		SiswaNama:        siswaNama,
		SiswaNISN:        siswaNISNVal,
		KelasNama:        kelasNama,
		TahunAjaran:      ta.TahunAjaran,
		Semester:         ta.Semester,
		PublishDate:      publishDateStr,
		CatatanWaliKelas: raport.CatatanWaliKelas,
		Sakit:            raport.Sakit,
		Izin:             raport.Izin,
		Alpha:            raport.Alpha,
		NilaiMapel:       nilaiMapelResp,
		NilaiEkskul:      ekskulResp,
		WaliKelas:        waliKelasNama,
	}, nil
}

// ==================== GetNilaiHistory ====================
func (c *penilaianController) GetNilaiHistory(raportNilaiID uint) ([]dto.NilaiHistoryResponse, error) {
	var audits []model.RaportNilaiAudit

	// Ambil data audit join dengan users
	err := c.db.Table("tbl_raport_nilai_audit").
		Select("tbl_raport_nilai_audit.*, users.name, users.role").
		Joins("LEFT JOIN users ON users.id = tbl_raport_nilai_audit.changed_by").
		Where("tbl_raport_nilai_audit.raport_nilai_id = ?", raportNilaiID).
		Order("tbl_raport_nilai_audit.changed_at DESC").
		Find(&audits).Error

	if err != nil {
		return nil, err
	}

	// Konversi ke DTO
	result := make([]dto.NilaiHistoryResponse, 0, len(audits))
	for _, audit := range audits {
		// Ambil user info dari hasil join
		var userName string
		var userRole string

		// Query manual untuk ambil user info karena GORM tidak auto-populate
		var user model.User
		err := c.db.Where("id = ?", audit.ChangedBy).First(&user).Error
		if err == nil {
			userName = user.Name
			userRole = user.Role
		}

		// Handle old_value yang mungkin NULL
		var oldValue *float64
		if audit.OldValue != 0 {
			oldValue = &audit.OldValue
		}

		result = append(result, dto.NilaiHistoryResponse{
			AuditID:   audit.AuditID,
			OldValue:  oldValue,
			NewValue:  audit.NewValue,
			ChangedBy: audit.ChangedBy,
			ChangedByUser: dto.UserInfo{
				ID:   audit.ChangedBy,
				Name: userName,
				Role: userRole,
			},
			ChangedAt: audit.ChangedAt,
		})
	}

	return result, nil
}

// ==================== STATUS PUBLISH EDIT ETC ====================
// UpdateStatusPublish - update status_publish per siswa
func (c *penilaianController) UpdateStatusPublish(taID, kelasID uint, siswaNIS string, statusPublish int) error {
	var raport model.Raport
	err := c.db.Where("siswa_nis = ? AND ta_id = ? AND kelas_id = ?",
		siswaNIS, taID, kelasID).
		First(&raport).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Jika belum ada raport, buat baru
			raport = model.Raport{
				SiswaNIS:      siswaNIS,
				KelasID:       kelasID,
				TaID:          taID,
				StatusPublish: statusPublish == 1,
				Sakit:         0,
				Izin:          0,
				Alpha:         0,
			}
			return c.db.Create(&raport).Error
		}
		return err
	}

	raport.StatusPublish = statusPublish == 1
	return c.db.Save(&raport).Error
}

// EditNilaiPerSiswa - edit nilai untuk satu siswa (force majeur)
func (c *penilaianController) EditNilaiPerSiswa(req dto.EditNilaiPerSiswaRequest) error {
	tx := c.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, item := range req.NilaiList {
		var nilai model.RaportNilai

		err := tx.Where(
			"raport_id = ? AND ta_kelas_mapel_id = ?",
			req.RaportID,
			item.TaKelasMapelID,
		).First(&nilai).Error

		if err != nil {
			tx.Rollback()
			return err
		}

		nilai.NilaiAngka = item.NilaiAngka
		nilai.NilaiHuruf = konversiNilaiHuruf(item.NilaiAngka)
		nilai.Predikat = konversiPredikat(item.NilaiAngka)
		nilai.Deskripsi = item.Deskripsi

		if err := tx.Save(&nilai).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// ==================== HELPER FUNCTIONS ====================

func konversiNilaiHuruf(nilai float64) string {
	if nilai >= 90 {
		return "A"
	} else if nilai >= 80 {
		return "B"
	} else if nilai >= 70 {
		return "C"
	} else if nilai >= 60 {
		return "D"
	}
	return "E"
}

func konversiPredikat(nilai float64) string {
	if nilai >= 90 {
		return "Sangat Baik"
	} else if nilai >= 80 {
		return "Baik"
	} else if nilai >= 70 {
		return "Cukup"
	} else if nilai >= 60 {
		return "Kurang"
	}
	return "Sangat Kurang"
}
