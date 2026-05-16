// model/ta_kelas_mapel.go
package model

import "time"

type TaKelasMapel struct {
	TaKelasMapelID uint      `gorm:"column:ta_kelas_mapel_id;primaryKey;autoIncrement" json:"ta_kelas_mapel_id"`
	TaID           uint      `gorm:"column:ta_id;not null" json:"ta_id"`
	KelasID        uint      `gorm:"column:kelas_id;not null" json:"kelas_id"`
	KdMapel        uint      `gorm:"column:kd_mapel;not null" json:"kd_mapel"`
	GuruID         *uint     `gorm:"column:guru_id" json:"guru_id,omitempty"`
	Urutan         int       `gorm:"column:urutan;default:0" json:"urutan"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`

	// Relasi - UBAH MENJADI POINTER
	TahunAjaran *TahunAjaran `gorm:"foreignKey:TaID;references:TaID" json:"tahun_ajaran,omitempty"`
	Kelas       *Kelas       `gorm:"foreignKey:KelasID;references:KelasId" json:"kelas,omitempty"`
	Mapel       *Mapel       `gorm:"foreignKey:KdMapel;references:KdMapel" json:"mapel,omitempty"`
	Guru        *Guru        `gorm:"foreignKey:GuruID;references:GuruID" json:"guru,omitempty"`
}

func (TaKelasMapel) TableName() string {
	return "tbl_ta_kelas_mapel"
}
