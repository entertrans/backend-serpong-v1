package siswa

import (
	"fmt"

	"github.com/entertrans/backend-bogor.git/internal/dto"
	"github.com/entertrans/backend-bogor.git/internal/model"
)

// UpdateSiswa untuk update data siswa berdasarkan NIS
func (c *siswaController) UpdateSiswa(nis string, req dto.UpdateSiswaRequest) error {
	fmt.Println("🚀 UPDATE SISWA:", nis)

	tx := c.db.Begin()

	// Update field yang dikirim (partial update)
	updates := c.buildSiswaUpdates(req)

	if len(updates) == 0 {
		tx.Rollback()
		return fmt.Errorf("tidak ada data yang akan diupdate")
	}

	// ✅ Update langsung dengan Where (pakai Model, bukan object siswa)
	if err := tx.Model(&model.Siswa{}).Where("siswa_nis = ?", nis).Updates(updates).Error; err != nil {
		tx.Rollback()
		return err
	}

	fmt.Printf("✅ SISWA UPDATED: %s\n", nis)
	return tx.Commit().Error
}

// UpdateOrangtua untuk update data orangtua berdasarkan NIS siswa
func (c *siswaController) UpdateOrangtua(nis string, req dto.UpdateOrangtuaRequest) error {
	fmt.Println("🚀 UPDATE ORANGTUA untuk siswa:", nis)

	tx := c.db.Begin()

	// Cek apakah siswa ada
	var siswa model.Siswa
	if err := tx.Where("siswa_nis = ?", nis).First(&siswa).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("siswa dengan NIS %s tidak ditemukan", nis)
	}

	// Update field yang dikirim
	updates := c.buildOrangtuaUpdates(req)

	if len(updates) == 0 {
		tx.Rollback()
		return fmt.Errorf("tidak ada data orangtua yang akan diupdate")
	}

	// UPSERT: Update jika ada, Create jika tidak ada
	var orangtua model.Orangtua
	result := tx.Where("siswa_nis = ?", nis).Assign(updates).FirstOrCreate(&orangtua, model.Orangtua{SiswaNIS: &nis})

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	fmt.Printf("✅ ORANGTUA UPDATED untuk siswa: %s\n", nis)
	return tx.Commit().Error
}

// Helper functions untuk UpdateSiswa
func (c *siswaController) buildSiswaUpdates(req dto.UpdateSiswaRequest) map[string]interface{} {
	updates := make(map[string]interface{})

	// Untuk string: jika nilai nil atau string kosong, set ke nil (NULL)
	if req.SiswaNama != nil {
		if *req.SiswaNama == "" {
			updates["siswa_nama"] = nil
		} else {
			updates["siswa_nama"] = *req.SiswaNama
		}
	}

	if req.SiswaNISN != nil {
		if *req.SiswaNISN == "" {
			updates["siswa_nisn"] = nil
		} else {
			updates["siswa_nisn"] = *req.SiswaNISN
		}
	}

	if req.SiswaJenkel != nil {
		if *req.SiswaJenkel == "" {
			updates["siswa_jenkel"] = nil
		} else {
			updates["siswa_jenkel"] = *req.SiswaJenkel
		}
	}

	if req.NoIjazah != nil {
		if *req.NoIjazah == "" {
			updates["no_ijazah"] = nil
		} else {
			updates["no_ijazah"] = *req.NoIjazah
		}
	}

	if req.SiswaTempat != nil {
		if *req.SiswaTempat == "" {
			updates["siswa_tempat"] = nil
		} else {
			updates["siswa_tempat"] = *req.SiswaTempat
		}
	}

	if req.SiswaTglLahir != nil {
		if *req.SiswaTglLahir == "" {
			updates["siswa_tgl_lahir"] = nil
		} else {
			updates["siswa_tgl_lahir"] = *req.SiswaTglLahir
		}
	}

	if req.SiswaAlamat != nil {
		if *req.SiswaAlamat == "" {
			updates["siswa_alamat"] = nil
		} else {
			updates["siswa_alamat"] = *req.SiswaAlamat
		}
	}

	if req.SiswaEmail != nil {
		if *req.SiswaEmail == "" {
			updates["siswa_email"] = nil
		} else {
			updates["siswa_email"] = *req.SiswaEmail
		}
	}

	if req.SiswaNoTelp != nil {
		if *req.SiswaNoTelp == "" {
			updates["siswa_no_telp"] = nil
		} else {
			updates["siswa_no_telp"] = *req.SiswaNoTelp
		}
	}

	if req.AnakKe != nil {
		// Untuk angka: jika 0, set ke NULL
		if *req.AnakKe == 0 {
			updates["anak_ke"] = nil
		} else {
			updates["anak_ke"] = *req.AnakKe
		}
	}

	if req.SiswaKelasID != nil {
		// Untuk ID: jika 0, set ke NULL
		if *req.SiswaKelasID == 0 {
			updates["siswa_kelas_id"] = nil
		} else {
			updates["siswa_kelas_id"] = *req.SiswaKelasID
		}
	}

	return updates
}

// Helper functions untuk UpdateOrangtua
func (c *siswaController) buildOrangtuaUpdates(req dto.UpdateOrangtuaRequest) map[string]interface{} {
	updates := make(map[string]interface{})

	// Data Ayah
	if req.AyahNama != nil {
		if *req.AyahNama == "" {
			updates["ayah_nama"] = nil
		} else {
			updates["ayah_nama"] = *req.AyahNama
		}
	}

	if req.AyahNik != nil {
		if *req.AyahNik == "" {
			updates["ayah_nik"] = nil
		} else {
			updates["ayah_nik"] = *req.AyahNik
		}
	}

	if req.AyahTempat != nil {
		if *req.AyahTempat == "" {
			updates["ayah_tempat"] = nil
		} else {
			updates["ayah_tempat"] = *req.AyahTempat
		}
	}

	if req.AyahTanggal != nil {
		if *req.AyahTanggal == "" {
			updates["ayah_tanggal"] = nil
		} else {
			updates["ayah_tanggal"] = *req.AyahTanggal
		}
	}

	if req.AyahPendidikan != nil {
		if *req.AyahPendidikan == "" {
			updates["ayah_pendidikan"] = nil
		} else {
			updates["ayah_pendidikan"] = *req.AyahPendidikan
		}
	}

	if req.AyahPekerjaan != nil {
		if *req.AyahPekerjaan == "" {
			updates["ayah_pekerjaan"] = nil
		} else {
			updates["ayah_pekerjaan"] = *req.AyahPekerjaan
		}
	}

	if req.AyahPenghasilan != nil {
		if *req.AyahPenghasilan == "" {
			updates["ayah_penghasilan"] = nil
		} else {
			updates["ayah_penghasilan"] = *req.AyahPenghasilan
		}
	}

	if req.NoTelpAyah != nil {
		if *req.NoTelpAyah == "" {
			updates["no_telp_ayah"] = nil
		} else {
			updates["no_telp_ayah"] = *req.NoTelpAyah
		}
	}

	if req.EmailAyah != nil {
		if *req.EmailAyah == "" {
			updates["email_ayah"] = nil
		} else {
			updates["email_ayah"] = *req.EmailAyah
		}
	}

	// Data Ibu
	if req.IbuNama != nil {
		if *req.IbuNama == "" {
			updates["ibu_nama"] = nil
		} else {
			updates["ibu_nama"] = *req.IbuNama
		}
	}

	if req.IbuNik != nil {
		if *req.IbuNik == "" {
			updates["ibu_nik"] = nil
		} else {
			updates["ibu_nik"] = *req.IbuNik
		}
	}

	if req.IbuTempat != nil {
		if *req.IbuTempat == "" {
			updates["ibu_tempat"] = nil
		} else {
			updates["ibu_tempat"] = *req.IbuTempat
		}
	}

	if req.IbuTanggal != nil {
		if *req.IbuTanggal == "" {
			updates["ibu_tanggal"] = nil
		} else {
			updates["ibu_tanggal"] = *req.IbuTanggal
		}
	}

	if req.IbuPendidikan != nil {
		if *req.IbuPendidikan == "" {
			updates["ibu_pendidikan"] = nil
		} else {
			updates["ibu_pendidikan"] = *req.IbuPendidikan
		}
	}

	if req.IbuPekerjaan != nil {
		if *req.IbuPekerjaan == "" {
			updates["ibu_pekerjaan"] = nil
		} else {
			updates["ibu_pekerjaan"] = *req.IbuPekerjaan
		}
	}

	if req.IbuPenghasilan != nil {
		if *req.IbuPenghasilan == "" {
			updates["ibu_penghasilan"] = nil
		} else {
			updates["ibu_penghasilan"] = *req.IbuPenghasilan
		}
	}

	if req.NoTelpIbu != nil {
		if *req.NoTelpIbu == "" {
			updates["no_telp_ibu"] = nil
		} else {
			updates["no_telp_ibu"] = *req.NoTelpIbu
		}
	}

	if req.EmailIbu != nil {
		if *req.EmailIbu == "" {
			updates["email_ibu"] = nil
		} else {
			updates["email_ibu"] = *req.EmailIbu
		}
	}

	// Data Wali
	if req.WaliNama != nil {
		if *req.WaliNama == "" {
			updates["wali_nama"] = nil
		} else {
			updates["wali_nama"] = *req.WaliNama
		}
	}

	if req.WaliNik != nil {
		if *req.WaliNik == "" {
			updates["wali_nik"] = nil
		} else {
			updates["wali_nik"] = *req.WaliNik
		}
	}

	if req.WaliTempat != nil {
		if *req.WaliTempat == "" {
			updates["wali_tempat"] = nil
		} else {
			updates["wali_tempat"] = *req.WaliTempat
		}
	}

	if req.WaliTanggal != nil {
		if *req.WaliTanggal == "" {
			updates["wali_tanggal"] = nil
		} else {
			updates["wali_tanggal"] = *req.WaliTanggal
		}
	}

	if req.WaliPendidikan != nil {
		if *req.WaliPendidikan == "" {
			updates["wali_pendidikan"] = nil
		} else {
			updates["wali_pendidikan"] = *req.WaliPendidikan
		}
	}

	if req.WaliPekerjaan != nil {
		if *req.WaliPekerjaan == "" {
			updates["wali_pekerjaan"] = nil
		} else {
			updates["wali_pekerjaan"] = *req.WaliPekerjaan
		}
	}

	if req.WaliPenghasilan != nil {
		if *req.WaliPenghasilan == "" {
			updates["wali_penghasilan"] = nil
		} else {
			updates["wali_penghasilan"] = *req.WaliPenghasilan
		}
	}

	if req.WaliAlamat != nil {
		if *req.WaliAlamat == "" {
			updates["wali_alamat"] = nil
		} else {
			updates["wali_alamat"] = *req.WaliAlamat
		}
	}

	if req.WaliNotelp != nil {
		if *req.WaliNotelp == "" {
			updates["wali_notelp"] = nil
		} else {
			updates["wali_notelp"] = *req.WaliNotelp
		}
	}

	return updates
}
