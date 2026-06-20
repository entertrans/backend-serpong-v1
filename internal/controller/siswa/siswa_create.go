package siswa

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/entertrans/backend-bogor.git/internal/dto"
	"github.com/entertrans/backend-bogor.git/internal/helper"
	"github.com/entertrans/backend-bogor.git/internal/model"
)

// CreateSiswa membuat siswa + user + orangtua dalam 1 transaksi
func (c *siswaController) CreateSiswa(req dto.CreateSiswaRequest) error {
	// fmt.Println("🚀 MASUK CONTROLLER - CREATE SISWA")

	tx := c.db.Begin()

	// ======================
	// 1. CEK & BUAT USER SISWA
	// ======================
	if req.SiswaNIS == "" {
		tx.Rollback()
		return fmt.Errorf("NIS wajib diisi")
	}

	email := req.SiswaNIS + "@siswa.sch.id"

	// Cek apakah user sudah ada
	var existingUser model.User
	if err := tx.Where("email = ?", email).First(&existingUser).Error; err == nil {
		tx.Rollback()
		return fmt.Errorf("user dengan NIS %s sudah terdaftar", req.SiswaNIS)
	}

	// Hash password (default = NIS)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.SiswaNIS), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Buat User
	user := model.User{
		Email:    email,
		Password: string(hashedPassword),
		Name:     req.SiswaNama,
		Role:     "siswa",
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return err
	}

	fmt.Printf("✅ USER CREATED: %s | Password default: %s\n", email, req.SiswaNIS)

	// ======================
	// 2. INSERT SISWA
	// ======================
	zero := 0
	bogor := 1 // ID Satelit Bogor (default)

	siswa := model.Siswa{
		SiswaNIS:             &req.SiswaNIS,
		SiswaNISN:            &req.SiswaNISN,
		SiswaNama:            &req.SiswaNama,
		SatelitID:            &bogor,
		SiswaJenkel:          &req.SiswaJenkel,
		SiswaTempat:          &req.SiswaTempat,
		SiswaTglLahir:        &req.SiswaTglLahir,
		SiswaAlamat:          &req.SiswaAlamat,
		SiswaEmail:           &req.SiswaEmail,
		SiswaNoTelp:          &req.SiswaNoTelp,
		SiswaNIK:             &req.SiswaNIK,
		SekolahAsal:          &req.SekolahAsal,
		SiswaKewarganegaraan: &req.SiswaKewarganegaraan,
		SiswaAgamaID:         helper.StringToIntPtr(req.SiswaAgamaID),
		SiswaKelasID:         helper.StringToIntPtr(req.SiswaKelasID),
		AnakKe:               helper.StringToIntPtr(req.AnakKe),
		NoIjazah:             &req.NoIjazah,
		UserID:               &user.ID,
		SoftDeleted:          &zero,
	}

	if err := tx.Create(&siswa).Error; err != nil {
		tx.Rollback()
		return err
	}

	// fmt.Println("✅ SISWA INSERTED")

	// ======================
	// 3. INSERT ORANGTUA (OPTIONAL)
	// ======================
	if err := c.createOrangtua(tx, req.SiswaNIS, req.Orangtua); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// createOrangtua helper untuk insert data orangtua
func (c *siswaController) createOrangtua(tx *gorm.DB, nis string, ortu dto.OrangtuaPayload) error {
	if ortu.AyahNama == nil && ortu.NoTelpAyah == nil {
		// fmt.Println("⏭️ SKIP INSERT ORTU (kosong)")
		return nil
	}

	// fmt.Println("👨‍👩‍👧 INSERT ORTU")

	ortuModel := model.Orangtua{
		SiswaNIS:   &nis,
		AyahNama:   ortu.AyahNama,
		AyahNotelp: ortu.NoTelpAyah,
		AyahNik:    ortu.AyahNik,
		IbuNama:    ortu.IbuNama,
		IbuNotelp:  ortu.NoTelpIbu,
		IbuNik:     ortu.IbuNik,
		WaliNama:   ortu.WaliNama,
		WaliNotelp: ortu.WaliNotelp,
		WaliAlamat: ortu.WaliAlamat,
		WaliNik:    ortu.WaliNik,
	}

	if err := tx.Create(&ortuModel).Error; err != nil {
		// fmt.Println("❌ GAGAL INSERT ORTU:", err)
		return err
	}

	// fmt.Println("✅ ORANGTUA INSERTED")
	return nil
}
