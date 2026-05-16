// internal/model/tahun_ajaran.go
package model

import "time"

type TahunAjaran struct {
	TaID        uint       `gorm:"column:ta_id;primaryKey;autoIncrement" json:"ta_id"`
	TahunAjaran string     `gorm:"column:tahun_ajaran;size:20;not null" json:"tahun_ajaran"`
	Semester    string     `gorm:"column:semester;type:enum('1','2');not null" json:"semester"`
	Status      string     `gorm:"column:status;type:enum('draft','aktif','selesai');default:'draft'" json:"status"`
	PublishDate *time.Time `gorm:"column:publish_date;type:date" json:"publish_date"` // PASTIKAN ADA
	IsActive    bool       `gorm:"column:is_active;default:false" json:"is_active"`
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`

	TaKelasMapels []TaKelasMapel `gorm:"foreignKey:TaID;references:TaID" json:"ta_kelas_mapels,omitempty"`
}

func (TahunAjaran) TableName() string {
	return "tbl_tahun_ajaran"
}
