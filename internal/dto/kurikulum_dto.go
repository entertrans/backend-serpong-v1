// internal/modules/rapor/dto/kurikulum_dto.go
package dto

// ==================== KURIKULUM DTO ====================

type TaKelasMapelItem struct {
	TaKelasMapelID uint   `json:"ta_kelas_mapel_id"`
	KdMapel        uint   `json:"kd_mapel"`
	NmMapel        string `json:"nm_mapel"`
	Kkm            int    `json:"kkm"`
	GuruID         *uint  `json:"guru_id"`
	GuruNama       string `json:"guru_nama"`
	Urutan         int    `json:"urutan"`
}

type KurikulumByKelasResponse struct {
	KelasID       uint               `json:"kelas_id"`
	KelasNama     string             `json:"kelas_nama"`
	WaliKelasID   *uint              `json:"wali_kelas_id"`
	WaliKelasNama string             `json:"wali_kelas_nama"`
	MapelList     []TaKelasMapelItem `json:"mapel_list"`
}

type SaveKurikulumRequest struct {
	TaID        uint                     `json:"ta_id" binding:"required"`
	KelasID     uint                     `json:"kelas_id" binding:"required"`
	WaliKelasID *uint                    `json:"wali_kelas_id"` // ← PERHATIKAN: binding:"required" dihapus
	MapelList   []SaveKurikulumMapelItem `json:"mapel_list" binding:"required"`
}

type SaveKurikulumMapelItem struct {
	KdMapel uint  `json:"kd_mapel" binding:"required"`
	GuruID  *uint `json:"guru_id"` // ← GuruID boleh null
	Urutan  int   `json:"urutan"`
}

type SaveKurikulumResponse struct {
	Message    string `json:"message"`
	TaID       uint   `json:"ta_id"`
	KelasID    uint   `json:"kelas_id"`
	TotalMapel int    `json:"total_mapel"`
}

type CopyKurikulumRequest struct {
	FromTaID  uint `json:"from_ta_id" binding:"required"`
	ToKelasID uint `json:"to_kelas_id" binding:"required"`
}

type CheckKurikulumResponse struct {
	BelumSetup []string `json:"belum_setup"`
	TotalKelas int      `json:"total_kelas"`
	SudahSetup int      `json:"sudah_setup"`
}

// GetKurikulumByGuruResponse - response untuk mapel by guru
type KurikulumByGuruResponse struct {
	TaID       uint              `json:"ta_id"`
	TaNama     string            `json:"ta_nama"` // Format: "2024/2025 - Semester Ganjil"
	KelasID    uint              `json:"kelas_id"`
	KelasNama  string            `json:"kelas_nama"`
	GuruID     uint              `json:"guru_id"`
	GuruNama   string            `json:"guru_nama"`
	GuruNIP    string            `json:"guru_nip"` // String, handle nil jadi string kosong
	TotalMapel int               `json:"total_mapel"`
	MapelList  []MapelByGuruItem `json:"mapel_list"`
	Message    string            `json:"message,omitempty"`
}

// MapelByGuruItem - item mapel yang diajar guru
type MapelByGuruItem struct {
	TaKelasMapelID uint   `json:"ta_kelas_mapel_id"`
	KdMapel        uint   `json:"kd_mapel"`
	NmMapel        string `json:"nm_mapel"`
	KKM            int    `json:"kkm"`
	Urutan         int    `json:"urutan"`
}
