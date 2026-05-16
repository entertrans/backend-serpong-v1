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
	// SEMUA FIELD PAKAI STRING (bukan *string atau number)
	SiswaNIS             string `json:"siswa_nis"`         // ← ubah jadi string
	SiswaNISN            string `json:"siswa_nisn"`        // ← ubah jadi string
	SiswaNama            string `json:"siswa_nama"`        // ← ubah jadi string
	SiswaJenkel          string `json:"siswa_jenkel"`      // ← ubah jadi string
	SekolahAsal          string `json:"sekolah_asal"`      // ← ubah jadi string
	NoIjazah             string `json:"no_ijazah"`         // ← ubah jadi string
	SiswaNIK             string `json:"siswa_nik"`         // ← ubah jadi string
	SiswaAgamaID         string `json:"siswa_agama"`       // ← ubah jadi string
	SiswaTempat          string `json:"siswa_tempat"`      // ← ubah jadi string
	SiswaTglLahir        string `json:"siswa_tgl_lahir"`   // ← ubah jadi string
	SiswaAlamat          string `json:"siswa_alamat"`      // ← ubah jadi string
	SiswaEmail           string `json:"siswa_email"`       // ← ubah jadi string
	SiswaKewarganegaraan string `json:"siswa_negara_asal"` // ← ubah jadi string
	SiswaNoTelp          string `json:"siswa_no_telp"`     // ← ubah jadi string
	AnakKe               string `json:"anak_ke"`           // ← ubah jadi int (bukan pointer)
	SiswaKelasID         string `json:"siswa_kelas_id"`    // ← ubah jadi int (bukan pointer)

	Orangtua OrangtuaPayload `json:"orangtua"`
}