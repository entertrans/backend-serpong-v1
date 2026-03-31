package dto

import "github.com/entertrans/go-base-project.git/internal/model"

type SiswaResponse struct {
	SiswaID      uint    `json:"siswa_id"`
	SiswaNIS     *string `json:"siswa_nis"`
	SiswaNISN    *string `json:"siswa_nisn"`
	SiswaNama    *string `json:"siswa_nama"`
	SiswaJenkel  *string `json:"siswa_jenkel"`
	NoIjazah     *string `json:"no_ijazah"`
	SiswaTempat  *string `json:"siswa_tempat"`
	SiswaTglLahir *string `json:"siswa_tgl_lahir"`
	SiswaAlamat  *string `json:"siswa_alamat"`
	SiswaEmail   *string `json:"siswa_email"`
	SiswaNoTelp  *string `json:"siswa_no_telp"`
	SiswaDokumen *string `json:"siswa_dokumen"`
	TglKeluar    *string `json:"tgl_keluar"`
	TglLulus     *string `json:"ta_lulus"`
	AnakKe       *int    `json:"anak_ke"`
	SiswaKelasID *int    `json:"siswa_kelas_id"`
	SoftDeleted  *int    `json:"soft_deleted"`
	Kelas        model.Kelas `json:"kelas"`
	Satelit      model.DtSatelit `json:"satelit"`
	Orangtua    model.Orangtua `json:"orangtua"`
	Lampiran    []model.Lampiran `json:"lampiran"`
}