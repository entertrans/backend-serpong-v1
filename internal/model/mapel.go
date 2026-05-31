package model

import "time"

// model/mapel.go
type Mapel struct {
	KdMapel   uint      `gorm:"column:kd_mapel;primaryKey;autoIncrement" json:"kd_mapel"`
	KodeMapel string    `gorm:"column:kode_mapel;size:20" json:"kode_mapel"`
	NmMapel   string    `gorm:"column:nm_mapel;size:150;not null" json:"nm_mapel"`
	Kelompok  string    `gorm:"column:kelompok;type:enum('A','B','C');default:'A'" json:"kelompok"`
	Jenjang   string    `gorm:"column:jenjang;type:enum('SD','SMP','SMA');not null" json:"jenjang"`
	IsActive  *bool     `gorm:"column:is_active;default:1" json:"is_active"`
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	// Relasi ke guru_mapel
	GuruMapels []GuruMapel `gorm:"foreignKey:KdMapel;references:KdMapel" json:"guru_mapels"`
}

func (Mapel) TableName() string {
	return "tbl_mapel"
}

type GuruMapel struct {
	GuruMapelID uint   `gorm:"column:guru_mapel_id;primaryKey;autoIncrement" json:"guru_mapel_id"`
	GuruID      uint   `gorm:"column:guru_id;not null"`
	KdMapel     uint   `gorm:"column:kd_mapel;not null"`
	KelasID     uint   `gorm:"column:kelas_id;not null"`
	TahunAjaran string `gorm:"column:tahun_ajaran;size:20" json:"tahun_ajaran"`
	StatusAktif bool   `gorm:"column:status_aktif;default:true" json:"status_aktif"`

	// Relasi
	Guru  Guru  `gorm:"foreignKey:GuruID;references:GuruID" json:"guru"`
	Mapel Mapel `gorm:"foreignKey:KdMapel;references:KdMapel" json:"mapel"`
}

func (GuruMapel) TableName() string {
	return "tbl_guru_mapel"
}
