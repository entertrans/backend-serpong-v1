package model

import "time"

type Kelas struct {
	KelasId   uint   `gorm:"column:kelas_id;primaryKey;autoIncrement" json:"kelas_id"`
	KelasNama string `gorm:"column:kelas_nama" json:"kelas_nama"`

	KelasMapels []KelasMapel `gorm:"foreignKey:KelasID;references:KelasId"`
}

func (Kelas) TableName() string {
	return "tbl_kelas"
}

type KelasMapel struct {
	ID      uint `gorm:"column:id_kelas_mapel;primaryKey;autoIncrement" json:"id"`
	KelasID uint `gorm:"column:kelas_id;not null" json:"kelas_id"`
	KdMapel uint `gorm:"column:kd_mapel;not null" json:"kd_mapel"`

	// Relasi
	Kelas Kelas `gorm:"foreignKey:KelasID;references:KelasId"`
	Mapel Mapel `gorm:"foreignKey:KdMapel;references:KdMapel"`
}

func (KelasMapel) TableName() string {
	return "tbl_kelas_mapel"
}

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

type Guru struct {
	GuruID          uint       `gorm:"column:guru_id;primaryKey;autoIncrement" json:"guru_id"`
	GuruNIP         *string    `gorm:"column:guru_nip;unique;size:50" json:"guru_nip,omitempty"`
	GuruNUPTK       *string    `gorm:"column:guru_nuptk;unique;size:50" json:"guru_nuptk,omitempty"`
	GuruNama        string     `gorm:"column:guru_nama;size:150;not null" json:"guru_nama"`
	GuruJenkel      string     `gorm:"column:guru_jenkel;type:enum('L','P');default:'L'" json:"guru_jenkel"`
	GuruTempatLahir *string    `gorm:"column:guru_tempat_lahir;size:100" json:"guru_tempat_lahir,omitempty"`
	GuruTglLahir    *time.Time `gorm:"column:guru_tgl_lahir" json:"guru_tgl_lahir,omitempty"`
	GuruAgamaID     *int       `gorm:"column:guru_agama_id" json:"guru_agama_id,omitempty"`
	GuruAlamat      *string    `gorm:"column:guru_alamat;type:text" json:"guru_alamat,omitempty"`
	GuruEmail       *string    `gorm:"column:guru_email;size:100" json:"guru_email,omitempty"`
	GuruNoTelp      *string    `gorm:"column:guru_no_telp;size:20" json:"guru_no_telp,omitempty"`
	GuruPhoto       *string    `gorm:"column:guru_photo;size:255" json:"guru_photo,omitempty"`
	StatusAktif     bool       `gorm:"column:status_aktif;default:true" json:"status_aktif"`
	CreatedAt       time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	// Relasi ke guru_mapel
	GuruMapels []GuruMapel `gorm:"foreignKey:GuruID;references:GuruID" json:"guru_mapels"`
}

func (Guru) TableName() string {
	return "tbl_guru"
}
