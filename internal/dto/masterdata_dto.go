package dto

type KelasListItem struct {
	KelasId   uint   `json:"kelas_id"`
	KelasNama string `json:"kelas_nama"`
	IsAlumni  bool   `json:"is_alumni"`
}

type SatelitListItem struct {
	SatelitId   uint   `json:"satelit_id"`
	SatelitNama string `json:"satelit_nama"`
}
type GuruAktifResponse struct {
	GuruID   uint   `json:"guru_id"`
	GuruNama string `json:"guru_nama"`
	GuruNIP  string `json:"guru_nip,omitempty"`
}

type MapelAktifResponse struct {
	KdMapel  uint   `json:"kd_mapel"`
	NmMapel  string `json:"nm_mapel"`
	Kelompok string `json:"kelompok"`
	Jenjang  string `json:"jenjang"`
}
