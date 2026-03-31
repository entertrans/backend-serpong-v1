package dto

type OrangtuaPayload struct {
	AyahNama        *string `json:"ayah_nama"`
	AyahNik         *string `json:"ayah_nik"`
	AyahTempat      *string `json:"ayah_tempat"`
	AyahTanggal     *string `json:"ayah_tanggal"`
	AyahPendidikan  *string `json:"ayah_pendidikan"`
	AyahPekerjaan   *string `json:"ayah_pekerjaan"`
	AyahPenghasilan *string `json:"ayah_penghasilan"`
	NoTelpAyah      *string `json:"no_telp_ayah"`
	EmailAyah       *string `json:"email_ayah"`

	IbuNama        *string `json:"ibu_nama"`
	IbuNik         *string `json:"ibu_nik"`
	IbuTempat      *string `json:"ibu_tempat"`
	IbuTanggal     *string `json:"ibu_tanggal"`
	IbuPendidikan  *string `json:"ibu_pendidikan"`
	IbuPekerjaan   *string `json:"ibu_pekerjaan"`
	IbuPenghasilan *string `json:"ibu_penghasilan"`
	NoTelpIbu      *string `json:"no_telp_ibu"`
	EmailIbu       *string `json:"email_ibu"`

	WaliNama        *string `json:"wali_nama"`
	WaliNik         *string `json:"wali_nik"`
	WaliTempat      *string `json:"wali_tempat"`
	WaliTanggal     *string `json:"wali_tanggal"`
	WaliPendidikan  *string `json:"wali_pendidikan"`
	WaliPekerjaan   *string `json:"wali_pekerjaan"`
	WaliPenghasilan *string `json:"wali_penghasilan"`
	WaliAlamat      *string `json:"wali_alamat"`
	WaliNotelp      *string `json:"wali_notelp"`
}

type CreateSiswaRequest struct {
	SiswaNIS      *string `json:"siswa_nis"`
	SiswaNISN     *string `json:"siswa_nisn"`
	SiswaNama     *string `json:"siswa_nama"`
	SiswaJenkel   *string `json:"siswa_jenkel"`
	NoIjazah      *string `json:"no_ijazah"`
	SiswaTempat   *string `json:"siswa_tempat"`
	SiswaTglLahir *string `json:"siswa_tgl_lahir"`
	SiswaAlamat   *string `json:"siswa_alamat"`
	SiswaEmail    *string `json:"siswa_email"`
	SiswaNoTelp   *string `json:"siswa_no_telp"`
	AnakKe        *int    `json:"anak_ke"`
	SiswaKelasID  *int    `json:"siswa_kelas_id"`

	Orangtua OrangtuaPayload `json:"orangtua"`
}