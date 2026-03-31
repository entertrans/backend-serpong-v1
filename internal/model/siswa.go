package model

type Siswa struct {
	SiswaID       uint    `json:"siswa_id" gorm:"column:siswa_id"`
	SiswaNIS      *string `json:"siswa_nis" gorm:"column:siswa_nis"`
	SiswaNISN     *string `json:"siswa_nisn" gorm:"column:siswa_nisn"`
	NoIjazah      *string `json:"no_ijazah" gorm:"column:no_ijazah"`
	NIKSiswa      *string `json:"nik_siswa" gorm:"column:nik_siswa"`
	SiswaNama     *string `json:"siswa_nama" gorm:"column:siswa_nama"`
	SiswaJenkel   *string `json:"siswa_jenkel" gorm:"column:siswa_jenkel"`
	SiswaTempat   *string `json:"siswa_tempat" gorm:"column:siswa_tempat"`
	SiswaTglLahir *string `json:"siswa_tgl_lahir" gorm:"column:siswa_tgl_lahir"` // bisa diganti time.Time kalau pakai parsing tanggal
	SiswaAgamaID  *int    `json:"-" gorm:"column:siswa_agama_id"`
	SiswaAlamat   *string `json:"siswa_alamat" gorm:"column:siswa_alamat"`
	SiswaEmail    *string `json:"siswa_email" gorm:"column:siswa_email"`
	SiswaDokumen  *string `json:"siswa_dokumen" gorm:"column:siswa_dokumen"`
	SiswaNoTelp   *string `json:"siswa_no_telp" gorm:"column:siswa_no_telp"`
	SiswaKelasID  *int    `json:"siswa_kelas_id" gorm:"column:siswa_kelas_id"`
	SiswaPhoto    *string `json:"siswa_photo" gorm:"column:siswa_photo"`
	SoftDeleted   *int    `json:"soft_deleted" gorm:"column:soft_deleted"`
	TglKeluar     *string `json:"tgl_keluar" gorm:"column:tgl_keluar"`
	TglLulus      *string `json:"ta_lulus" gorm:"column:ta_lulus"`
	AnakKe        *int    `json:"anak_ke" gorm:"column:anak_ke"`
	SekolahAsal   *string `json:"sekolah_asal" gorm:"column:sekolah_asal"`
	SatelitID     *int    `json:"satelit" gorm:"column:satelit"`
	OC            *int    `json:"oc" gorm:"column:oc"`
	KC            *int    `json:"kc" gorm:"column:kc"`

	// Relasi: satu siswa punya satu ortu
	Orangtua Orangtua   `json:"orangtua" gorm:"foreignKey:SiswaNIS;references:SiswaNIS"`
	Agama    Agama      `json:"agama" gorm:"foreignKey:SiswaAgamaID;references:AgamaId"`
	Kelas    Kelas      `json:"kelas" gorm:"foreignKey:SiswaKelasID;references:KelasId"`
	Satelit  DtSatelit  `json:"Satelit" gorm:"foreignKey:SatelitID;references:SatelitId"`
	Lampiran []Lampiran `json:"lampiran" gorm:"foreignKey:SiswaNIS;references:SiswaNIS"`
}

func (Siswa) TableName() string {
	return "tbl_siswa"
}

type Agama struct {
	AgamaId   uint   `json:"agama_id" gorm:"column:agama_id"`
	AgamaNama string `json:"agama_nama" gorm:"column:agama_nama"`
}

func (Agama) TableName() string {
	return "tbl_agama"
}

type DtSatelit struct {
	SatelitId   uint   `json:"satelit_id" gorm:"column:satelit_id"`
	SatelitNama string `json:"satelit_nama" gorm:"column:satelit_nama"`
}

func (DtSatelit) TableName() string {
	return "tbl_satelit"
}

type Orangtua struct {
	OrtuID          uint    `json:"-" gorm:"column:ortu_id"`
	SiswaNIS        *string `json:"siswa_nis" gorm:"column:siswa_nis"`
	AyahNama        *string `json:"ayah_nama" gorm:"column:ayah_nama"`
	AyahNik         *string `json:"ayah_nik" gorm:"column:ayah_nik"`
	AyahTempat      *string `json:"ayah_tempat" gorm:"column:ayah_tempat"`
	AyahTanggal     *string `json:"ayah_tanggal" gorm:"column:ayah_tanggal"`
	AyahPendidikan  *string `json:"ayah_pendidikan" gorm:"column:ayah_pendidikan"`
	AyahPekerjaan   *string `json:"ayah_pekerjaan" gorm:"column:ayah_pekerjaan"`
	AyahPenghasilan *string `json:"ayah_penghasilan" gorm:"column:ayah_penghasilan"`
	AyahNotelp      *string `json:"no_telp_ayah" gorm:"column:no_telp_ayah"`
	AyahEmail       *string `json:"email_ayah" gorm:"column:email_ayah"`
	IbuNama         *string `json:"ibu_nama" gorm:"column:ibu_nama"`
	IbuNik          *string `json:"ibu_nik" gorm:"column:ibu_nik"`
	IbuTempat       *string `json:"ibu_tempat" gorm:"column:ibu_tempat"`
	IbuTanggal      *string `json:"ibu_tanggal" gorm:"column:ibu_tanggal"`
	IbuPendidikan   *string `json:"ibu_pendidikan" gorm:"column:ibu_pendidikan"`
	IbuPekerjaan    *string `json:"ibu_pekerjaan" gorm:"column:ibu_pekerjaan"`
	IbuPenghasilan  *string `json:"ibu_penghasilan" gorm:"column:ibu_penghasilan"`
	IbuNotelp       *string `json:"no_telp_ibu" gorm:"column:no_telp_ibu"`
	IbuEmail        *string `json:"email_ibu" gorm:"column:email_ibu"`
	WaliNama        *string `json:"wali_nama" gorm:"column:wali_nama"`
	WaliNik         *string `json:"wali_nik" gorm:"column:wali_nik"`
	WaliTempat      *string `json:"wali_tempat" gorm:"column:wali_tempat"`
	WaliTanggal     *string `json:"wali_tanggal" gorm:"column:wali_tanggal"`
	WaliPendidikan  *string `json:"wali_pendidikan" gorm:"column:wali_pendidikan"`
	WaliPekerjaan   *string `json:"wali_pekerjaan" gorm:"column:wali_pekerjaan"`
	WaliPenghasilan *string `json:"wali_penghasilan" gorm:"column:wali_penghasilan"`
	WaliAlamat      *string `json:"wali_alamat" gorm:"column:wali_alamat"`
	WaliNotelp      *string `json:"wali_notelp" gorm:"column:wali_notelp"`
}

func (Orangtua) TableName() string {
	return "tbl_orangtua"
}