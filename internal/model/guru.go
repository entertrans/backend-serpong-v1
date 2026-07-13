package model

import "time"

// internal/model/guru.go
type Guru struct {
	GuruID          uint       `gorm:"column:guru_id;primaryKey;autoIncrement" json:"guru_id"`
	UserID          *uint      `gorm:"column:user_id;index" json:"user_id,omitempty"`
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

	// Relasi
	User       *User       `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	GuruMapels []GuruMapel `gorm:"foreignKey:GuruID;references:GuruID" json:"guru_mapels"`
}

func (Guru) TableName() string {
	return "tbl_guru"
}
