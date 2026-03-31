package model

import "time"

type Lampiran struct {
	IDLampiran      uint64    `json:"id_lampiran" gorm:"column:id_lampiran;primaryKey;autoIncrement"`
	SiswaNIS        string    `json:"siswa_nis" gorm:"column:siswa_nis;type:char(255);not null;index"`
	DokumenJenis    string    `json:"dokumen_jenis" gorm:"column:dokumen_jenis;size:50;not null"`

	StorageProvider string    `json:"storage_provider" gorm:"column:storage_provider;size:20;not null"`
	ObjectKey       string    `json:"object_key" gorm:"column:object_key;size:255;not null"`

	FileName  string `json:"file_name" gorm:"column:file_name;size:191"`
	MimeType  string `json:"mime_type" gorm:"column:mime_type;size:100"`
	SizeBytes int64  `json:"size_bytes" gorm:"column:size_bytes"`
	ETag      string `json:"etag" gorm:"column:etag;size:191"`

	UploadedAt time.Time `json:"uploaded_at" gorm:"column:uploaded_at;not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (Lampiran) TableName() string { return "tbl_lampiran" }
