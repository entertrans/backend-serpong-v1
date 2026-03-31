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
