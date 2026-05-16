package dto

import "github.com/entertrans/go-base-project.git/internal/model"

type SiswaResponse struct {
	SiswaID       uint               `json:"siswa_id"`
	SiswaNIS      *string            `json:"siswa_nis"`
	SiswaNISN     *string            `json:"siswa_nisn"`
	SiswaNama     *string            `json:"siswa_nama"`
	SiswaJenkel   *string            `json:"siswa_jenkel"`
	NoIjazah      *string            `json:"no_ijazah"`
	SiswaTempat   *string            `json:"siswa_tempat"`
	SiswaTglLahir *string            `json:"siswa_tgl_lahir"`
	SiswaAlamat   *string            `json:"siswa_alamat"`
	SiswaEmail    *string            `json:"siswa_email"`
	SiswaNoTelp   *string            `json:"siswa_no_telp"`
	SiswaDokumen  *string            `json:"siswa_dokumen"`
	TglKeluar     *string            `json:"tgl_keluar"`
	TglLulus      *string            `json:"ta_lulus"`
	AnakKe        *int               `json:"anak_ke"`
	SiswaKelasID  *int               `json:"siswa_kelas_id"`
	SoftDeleted   *int               `json:"soft_deleted"`
	Kelas         model.Kelas        `json:"kelas"`
	Satelit       model.DtSatelit    `json:"satelit"`
	Orangtua      model.Orangtua     `json:"orangtua"`
	Lampiran      []LampiranResponse `json:"lampiran"` // 👈 PAKAI YANG SUDAH ADA
}

type UpdateSiswaRequest struct {
	SiswaNama     *string `json:"siswa_nama,omitempty"`
	SiswaNISN     *string `json:"siswa_nisn,omitempty"`
	SiswaJenkel   *string `json:"siswa_jenkel,omitempty"`
	NoIjazah      *string `json:"no_ijazah,omitempty"`
	SiswaTempat   *string `json:"siswa_tempat,omitempty"`
	SiswaTglLahir *string `json:"siswa_tgl_lahir,omitempty"`
	SiswaAlamat   *string `json:"siswa_alamat,omitempty"`
	SiswaEmail    *string `json:"siswa_email,omitempty"`
	SiswaNoTelp   *string `json:"siswa_no_telp,omitempty"`
	AnakKe        *int    `json:"anak_ke,omitempty"`
	SiswaKelasID  *int    `json:"siswa_kelas_id,omitempty"`
}
type UpdateOrangtuaRequest struct {
	AyahNama        *string `json:"ayah_nama,omitempty"`
	AyahNik         *string `json:"ayah_nik,omitempty"`
	AyahTempat      *string `json:"ayah_tempat,omitempty"`
	AyahTanggal     *string `json:"ayah_tanggal,omitempty"`
	AyahPendidikan  *string `json:"ayah_pendidikan,omitempty"`
	AyahPekerjaan   *string `json:"ayah_pekerjaan,omitempty"`
	AyahPenghasilan *string `json:"ayah_penghasilan,omitempty"`
	NoTelpAyah      *string `json:"no_telp_ayah,omitempty"`
	EmailAyah       *string `json:"email_ayah,omitempty"`

	IbuNama        *string `json:"ibu_nama,omitempty"`
	IbuNik         *string `json:"ibu_nik,omitempty"`
	IbuTempat      *string `json:"ibu_tempat,omitempty"`
	IbuTanggal     *string `json:"ibu_tanggal,omitempty"`
	IbuPendidikan  *string `json:"ibu_pendidikan,omitempty"`
	IbuPekerjaan   *string `json:"ibu_pekerjaan,omitempty"`
	IbuPenghasilan *string `json:"ibu_penghasilan,omitempty"`
	NoTelpIbu      *string `json:"no_telp_ibu,omitempty"`
	EmailIbu       *string `json:"email_ibu,omitempty"`

	WaliNama        *string `json:"wali_nama,omitempty"`
	WaliNik         *string `json:"wali_nik,omitempty"`
	WaliTempat      *string `json:"wali_tempat,omitempty"`
	WaliTanggal     *string `json:"wali_tanggal,omitempty"`
	WaliPendidikan  *string `json:"wali_pendidikan,omitempty"`
	WaliPekerjaan   *string `json:"wali_pekerjaan,omitempty"`
	WaliPenghasilan *string `json:"wali_penghasilan,omitempty"`
	WaliAlamat      *string `json:"wali_alamat,omitempty"`
	WaliNotelp      *string `json:"wali_notelp,omitempty"`
}
