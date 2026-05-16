// internal/model/ta_kelas_wali.go
package model

import "time"

type TaKelasWali struct {
	TaKelasWaliID uint      `gorm:"column:ta_kelas_wali_id;primaryKey;autoIncrement" json:"ta_kelas_wali_id"`
	TaID          uint      `gorm:"column:ta_id;not null" json:"ta_id"`
	KelasID       uint      `gorm:"column:kelas_id;not null" json:"kelas_id"`
	GuruID        uint      `gorm:"column:guru_id;not null" json:"guru_id"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	// Relasi - UBAH MENJADI POINTER
	TahunAjaran *TahunAjaran `gorm:"foreignKey:TaID;references:TaID" json:"tahun_ajaran,omitempty"`
	Kelas       *Kelas       `gorm:"foreignKey:KelasID;references:KelasId" json:"kelas,omitempty"`
	Guru        *Guru        `gorm:"foreignKey:GuruID;references:GuruID" json:"guru,omitempty"`
}

func (TaKelasWali) TableName() string {
	return "tbl_ta_kelas_wali"
}
