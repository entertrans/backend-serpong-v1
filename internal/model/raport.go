// internal/model/raport.go
package model

import (
	"time"
)

type Raport struct {
	RaportID         uint      `gorm:"column:raport_id;primaryKey;autoIncrement" json:"raport_id"`
	SiswaNIS         string    `gorm:"column:siswa_nis;size:30;not null" json:"siswa_nis"`
	KelasID          uint      `gorm:"column:kelas_id;not null" json:"kelas_id"`
	TaID             uint      `gorm:"column:ta_id;not null" json:"ta_id"`
	CatatanWaliKelas string    `gorm:"column:catatan_wali_kelas;type:text" json:"catatan_wali_kelas"`
	Sakit            int       `gorm:"column:sakit;default:0" json:"sakit"`
	Izin             int       `gorm:"column:izin;default:0" json:"izin"`
	Alpha            int       `gorm:"column:alpha;default:0" json:"alpha"`
	StatusPublish    bool      `gorm:"column:status_publish;default:false" json:"status_publish"`
	CreatedAt        time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	// Relasi
	Siswa         Siswa          `gorm:"foreignKey:SiswaNIS;references:SiswaNIS" json:"siswa,omitempty"`
	Kelas         Kelas          `gorm:"foreignKey:KelasID;references:KelasId" json:"kelas,omitempty"`
	TahunAjaran   TahunAjaran    `gorm:"foreignKey:TaID;references:TaID" json:"tahun_ajaran,omitempty"`
	RaportNilais  []RaportNilai  `gorm:"foreignKey:RaportID;references:RaportID" json:"raport_nilais,omitempty"`
	RaportEkskuls []RaportEkskul `gorm:"foreignKey:RaportID;references:RaportID" json:"raport_ekskuls,omitempty"`
}

func (Raport) TableName() string {
	return "tbl_raport"
}

type RaportNilai struct {
	RaportNilaiID  uint `gorm:"column:raport_nilai_id;primaryKey;autoIncrement" json:"raport_nilai_id"`
	RaportID       uint `gorm:"column:raport_id;not null" json:"raport_id"`
	TaKelasMapelID uint `gorm:"column:ta_kelas_mapel_id;not null" json:"ta_kelas_mapel_id"`

	NilaiAngka float64 `gorm:"column:nilai_angka;type:decimal(5,2);default:0.00" json:"nilai_angka"`
	NilaiHuruf string  `gorm:"column:nilai_huruf;size:5" json:"nilai_huruf"`
	Predikat   string  `gorm:"column:predikat;size:255" json:"predikat"`
	Deskripsi  string  `gorm:"column:deskripsi;type:text" json:"deskripsi"`

	CreatedBy *uint `gorm:"column:created_by"`
	UpdatedBy *uint

	LastChangedBy *uint      `gorm:"column:last_changed_by"`
	LastChangedAt *time.Time `gorm:"column:last_changed_at"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	// Relasi
	Raport       Raport       `gorm:"foreignKey:RaportID;references:RaportID" json:"raport,omitempty"`
	TaKelasMapel TaKelasMapel `gorm:"foreignKey:TaKelasMapelID;references:TaKelasMapelID" json:"ta_kelas_mapel,omitempty"`
}

func (RaportNilai) TableName() string {
	return "tbl_raport_nilai"
}

type RaportNilaiAudit struct {
	AuditID uint `gorm:"primaryKey"`

	RaportNilaiID uint

	OldValue float64
	NewValue float64

	ChangedBy uint
	ChangedAt time.Time
}

func (RaportNilaiAudit) TableName() string {
	return "tbl_raport_nilai_audit"
}

type RaportEkskul struct {
	EkskulID   uint   `gorm:"column:ekskul_id;primaryKey;autoIncrement" json:"ekskul_id"`
	RaportID   uint   `gorm:"column:raport_id;not null" json:"raport_id"`
	NamaEkskul string `gorm:"column:nama_ekskul;size:100" json:"nama_ekskul"`
	Nilai      string `gorm:"column:nilai;size:20" json:"nilai"`
	Deskripsi  string `gorm:"column:deskripsi;type:text" json:"deskripsi"`

	// Relasi
	Raport Raport `gorm:"foreignKey:RaportID;references:RaportID" json:"raport,omitempty"`
}

func (RaportEkskul) TableName() string {
	return "tbl_raport_ekskul"
}
