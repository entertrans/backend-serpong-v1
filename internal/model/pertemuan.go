// internal/model/pertemuan.go
package model

import (
	"time"
)

type Pertemuan struct {
	PertemuanID    uint      `gorm:"column:pertemuan_id;primaryKey;autoIncrement" json:"pertemuan_id"`
	TaKelasMapelID uint      `gorm:"column:ta_kelas_mapel_id;not null" json:"ta_kelas_mapel_id"`
	PertemuanKe    int       `gorm:"column:pertemuan_ke;not null" json:"pertemuan_ke"`
	Tema           string    `gorm:"column:tema;size:255;not null" json:"tema"`
	Deskripsi      string    `gorm:"column:deskripsi;type:longtext" json:"deskripsi,omitempty"`
	Tanggal        time.Time `gorm:"column:tanggal;not null" json:"tanggal"`
	Status         string    `gorm:"column:status;type:enum('draft','dibuka','berlangsung','selesai');default:'draft'" json:"status"`
	CreatedBy      uint      `gorm:"column:created_by;not null" json:"created_by"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	// Relasi
	TaKelasMapel *TaKelasMapel      `gorm:"foreignKey:TaKelasMapelID;references:TaKelasMapelID" json:"ta_kelas_mapel,omitempty"`
	Guru         *Guru              `gorm:"foreignKey:CreatedBy;references:GuruID" json:"guru,omitempty"`
	Materi       *PertemuanMateri   `gorm:"foreignKey:PertemuanID;references:PertemuanID" json:"materi,omitempty"`
	Meeting      *PertemuanMeeting  `gorm:"foreignKey:PertemuanID;references:PertemuanID" json:"meeting,omitempty"`
	Absensi      []PertemuanAbsensi `gorm:"foreignKey:PertemuanID;references:PertemuanID" json:"absensi,omitempty"`
}

func (Pertemuan) TableName() string {
	return "tbl_ko_pertemuan"
}

type PertemuanMateri struct {
	MateriID    uint      `gorm:"column:materi_id;primaryKey;autoIncrement" json:"materi_id"`
	PertemuanID uint      `gorm:"column:pertemuan_id;not null" json:"pertemuan_id"`
	Judul       string    `gorm:"column:judul;size:255;not null" json:"judul"`
	Konten      string    `gorm:"column:konten;type:longtext;not null" json:"konten"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	// Relasi
	Pertemuan *Pertemuan `gorm:"foreignKey:PertemuanID;references:PertemuanID" json:"pertemuan,omitempty"`
}

func (PertemuanMateri) TableName() string {
	return "tbl_ko_pertemuan_materi"
}

type PertemuanMeeting struct {
	MeetingID       uint       `gorm:"column:meeting_id;primaryKey;autoIncrement" json:"meeting_id"`
	PertemuanID     uint       `gorm:"column:pertemuan_id;not null;unique" json:"pertemuan_id"`
	Provider        string     `gorm:"column:provider;type:enum('google_meet');default:'google_meet'" json:"provider"`
	CalendarEventID *string    `gorm:"column:calendar_event_id;size:255" json:"calendar_event_id,omitempty"`
	MeetingURL      *string    `gorm:"column:meeting_url;type:text" json:"meeting_url,omitempty"`
	StartedAt       *time.Time `gorm:"column:started_at" json:"started_at,omitempty"`
	EndedAt         *time.Time `gorm:"column:ended_at" json:"ended_at,omitempty"`
	CreatedAt       time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`

	// Relasi
	Pertemuan *Pertemuan `gorm:"foreignKey:PertemuanID;references:PertemuanID" json:"pertemuan,omitempty"`
}

func (PertemuanMeeting) TableName() string {
	return "tbl_ko_pertemuan_meeting"
}

type PertemuanAbsensi struct {
	AbsensiID   uint       `gorm:"column:absensi_id;primaryKey;autoIncrement" json:"absensi_id"`
	PertemuanID uint       `gorm:"column:pertemuan_id;not null" json:"pertemuan_id"`
	SiswaID     uint       `gorm:"column:siswa_id;not null" json:"siswa_id"`
	JoinTime    *time.Time `gorm:"column:join_time" json:"join_time,omitempty"`
	ViaLMS      bool       `gorm:"column:via_lms;default:false" json:"via_lms"`
	StatusFinal string     `gorm:"column:status_final;type:enum('belum_divalidasi','hadir','izin','sakit','alfa');default:'belum_divalidasi'" json:"status_final"`
	Catatan     *string    `gorm:"column:catatan;type:text" json:"catatan,omitempty"`
	ValidatedBy *uint      `gorm:"column:validated_by" json:"validated_by,omitempty"`
	ValidatedAt *time.Time `gorm:"column:validated_at" json:"validated_at,omitempty"`
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`

	// Relasi
	Pertemuan *Pertemuan `gorm:"foreignKey:PertemuanID;references:PertemuanID" json:"pertemuan,omitempty"`
	Siswa     *Siswa     `gorm:"foreignKey:SiswaID;references:SiswaID" json:"siswa,omitempty"`
	Validator *Guru      `gorm:"foreignKey:ValidatedBy;references:GuruID" json:"validator,omitempty"`
}

func (PertemuanAbsensi) TableName() string {
	return "tbl_ko_absensi_pertemuan"
}
