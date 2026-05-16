// internal/modules/rapor/dto/penilaian_dto.go
package dto

// ==================== DTO UNTUK SISWA ====================

type SiswaByKelasResponse struct {
	SiswaNIS      string `json:"siswa_nis"`
	SiswaNama     string `json:"siswa_nama"`
	SiswaNISN     string `json:"siswa_nisn"`
	StatusPublish bool   `json:"status_publish"` // ← TAMBAHKAN
}

// ==================== DTO UNTUK NILAI MAPEL ====================

type NilaiMapelRequest struct {
	TaID           uint             `json:"ta_id" binding:"required"`
	KelasID        uint             `json:"kelas_id" binding:"required"`
	TaKelasMapelID uint             `json:"ta_kelas_mapel_id" binding:"required"`
	NilaiList      []NilaiMapelItem `json:"nilai_list" binding:"required"`
}

type NilaiMapelItem struct {
	SiswaNIS   string  `json:"siswa_nis" binding:"required"`
	NilaiAngka float64 `json:"nilai_angka" binding:"min=0,max=100"`
	Deskripsi  string  `json:"deskripsi"`
}

type NilaiMapelResponse struct {
	Message        string `json:"message"`
	TotalData      int    `json:"total_data"`
	TaKelasMapelID uint   `json:"ta_kelas_mapel_id"`
}

// ==================== DTO UNTUK NILAI MAPEL (GET) ====================

type GetNilaiMapelResponse struct {
	SiswaNIS   string  `json:"siswa_nis"`
	SiswaNama  string  `json:"siswa_nama"`
	NilaiAngka float64 `json:"nilai_angka"`
	Deskripsi  string  `json:"deskripsi"`
}

// ==================== DTO UNTUK ABSENSI & CATATAN ====================

type AbsensiRequest struct {
	TaID        uint          `json:"ta_id" binding:"required"`
	KelasID     uint          `json:"kelas_id" binding:"required"`
	AbsensiList []AbsensiItem `json:"absensi_list" binding:"required"`
}

type AbsensiItem struct {
	SiswaNIS         string `json:"siswa_nis" binding:"required"`
	Sakit            int    `json:"sakit" binding:"min=0"`
	Izin             int    `json:"izin" binding:"min=0"`
	Alpha            int    `json:"alpha" binding:"min=0"`
	CatatanWaliKelas string `json:"catatan_wali_kelas"`
}

type AbsensiResponse struct {
	Message   string `json:"message"`
	TotalData int    `json:"total_data"`
}

// ==================== DTO UNTUK EKSKUL ====================

type EkskulRequest struct {
	TaID       uint         `json:"ta_id" binding:"required"`
	KelasID    uint         `json:"kelas_id" binding:"required"`
	EkskulList []EkskulItem `json:"ekskul_list" binding:"required"`
}

type EkskulItem struct {
	SiswaNIS   string `json:"siswa_nis" binding:"required"`
	NamaEkskul string `json:"nama_ekskul" binding:"required"`
	Nilai      string `json:"nilai"`
	Deskripsi  string `json:"deskripsi"`
}

type EkskulResponse struct {
	Message   string `json:"message"`
	TotalData int    `json:"total_data"`
}

// ==================== DTO UNTUK RAPORT ====================

type RaportResponse struct {
	RaportID         uint               `json:"raport_id"`
	KelasID          uint               `json:"kelas_id"`
	SiswaNIS         string             `json:"siswa_nis"`
	SiswaNama        string             `json:"siswa_nama"`
	SiswaNISN        string             `json:"siswa_nisn"`
	KelasNama        string             `json:"kelas_nama"`
	TahunAjaran      string             `json:"tahun_ajaran"`
	PublishDate      *string            `json:"publish_date"`
	Semester         string             `json:"semester"`
	CatatanWaliKelas string             `json:"catatan_wali_kelas"`
	Sakit            int                `json:"sakit"`
	Izin             int                `json:"izin"`
	Alpha            int                `json:"alpha"`
	NilaiMapel       []NilaiMapelDetail `json:"nilai_mapel"`
	NilaiEkskul      []EkskulDetail     `json:"nilai_ekskul"`
	WaliKelas        string             `json:"wali_kelas,omitempty"`
}

// dto/raport.go
type NilaiMapelDetail struct {
	KdMapel    uint    `json:"kd_mapel"`
	NmMapel    string  `json:"nm_mapel"`
	Kelompok   string  `json:"kelompok"` // tambahkan ini
	NilaiAngka float64 `json:"nilai_angka"`
	NilaiHuruf string  `json:"nilai_huruf"`
	Predikat   string  `json:"predikat"`
	Deskripsi  string  `json:"deskripsi"`
}

type EkskulDetail struct {
	NamaEkskul string `json:"nama_ekskul"`
	Nilai      string `json:"nilai"`
	Deskripsi  string `json:"deskripsi"`
}

type EditNilaiPerSiswaRequest struct {
	RaportID  uint                    `json:"raport_id"`
	NilaiList []EditNilaiPerSiswaItem `json:"nilai_list"`
}

type EditNilaiPerSiswaItem struct {
	TaKelasMapelID uint    `json:"ta_kelas_mapel_id"`
	NilaiAngka     float64 `json:"nilai_angka"`
	Deskripsi      string  `json:"deskripsi"`
}

type UpdateStatusPublishRequest struct {
	StatusPublish int `json:"status_publish"` // hapus binding dulu
}
