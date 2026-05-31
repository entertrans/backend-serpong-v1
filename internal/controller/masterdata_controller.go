package controller

import (
	"gorm.io/gorm"

	"github.com/entertrans/backend-bogor.git/internal/dto"
	"github.com/entertrans/backend-bogor.git/internal/model"
)

type MasterdataController interface {
	GetKelasAll() ([]dto.KelasListItem, error)
	GetKelasAktif() ([]dto.KelasListItem, error)
	GetKelasAlumni() ([]dto.KelasListItem, error)
	GetSatelit() ([]dto.SatelitListItem, error)
	GetGuruAktif() ([]dto.GuruAktifResponse, error)
	GetMapelAktif() ([]dto.MapelAktifResponse, error)
}

type masterdataController struct {
	db *gorm.DB
}

func NewMasterdataController(db *gorm.DB) MasterdataController {
	return &masterdataController{db: db}
}

func (c *masterdataController) GetKelasAll() ([]dto.KelasListItem, error) {
	var kelasList []model.Kelas

	err := c.db.
		Order("kelas_id asc").
		Find(&kelasList).Error
	if err != nil {
		return nil, err
	}

	return mapKelasList(kelasList), nil
}

// sesuai permintaan: kelas aktif = 1 - 12
func (c *masterdataController) GetKelasAktif() ([]dto.KelasListItem, error) {
	var kelasList []model.Kelas

	err := c.db.
		Where("kelas_id BETWEEN ? AND ?", 1, 12).
		Order("kelas_id asc").
		Find(&kelasList).Error
	if err != nil {
		return nil, err
	}

	return mapKelasList(kelasList), nil
}

// alumni = kelas_id > 15 (atau > 12 kalau kamu ingin 13-15 masuk alumni, tapi dari data kamu alumni mulai 16)
func (c *masterdataController) GetKelasAlumni() ([]dto.KelasListItem, error) {
	var kelasList []model.Kelas

	err := c.db.
		Where("kelas_id > ?", 15).
		Order("kelas_id asc").
		Find(&kelasList).Error
	if err != nil {
		return nil, err
	}

	return mapKelasList(kelasList), nil
}

func (c *masterdataController) GetSatelit() ([]dto.SatelitListItem, error) {
	var satelitList []model.DtSatelit

	err := c.db.
		Order("satelit_id asc").
		Find(&satelitList).Error
	if err != nil {
		return nil, err
	}

	result := make([]dto.SatelitListItem, 0, len(satelitList))
	for _, s := range satelitList {
		result = append(result, dto.SatelitListItem{
			SatelitId:   s.SatelitId,
			SatelitNama: s.SatelitNama,
		})
	}

	return result, nil
}

func mapKelasList(kelasList []model.Kelas) []dto.KelasListItem {
	result := make([]dto.KelasListItem, 0, len(kelasList))
	for _, k := range kelasList {
		result = append(result, dto.KelasListItem{
			KelasId:   k.KelasId,
			KelasNama: k.KelasNama,
			IsAlumni:  k.KelasId > 15,
		})
	}
	return result
}

//new

func (c *masterdataController) GetGuruAktif() ([]dto.GuruAktifResponse, error) {
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

func (c *masterdataController) GetMapelAktif() ([]dto.MapelAktifResponse, error) {
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
