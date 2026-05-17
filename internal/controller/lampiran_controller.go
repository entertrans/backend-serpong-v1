package controller

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/entertrans/backend-bogor.git/internal/config"
	"github.com/entertrans/backend-bogor.git/internal/dto"
	"github.com/entertrans/backend-bogor.git/internal/model"
)

var (
	ErrValidation       = errors.New("validation error")
	ErrSiswaNotFound    = errors.New("siswa not found")
	ErrLampiranNotFound = errors.New("lampiran not found")
)

func IsValidation(err error) bool       { return errors.Is(err, ErrValidation) }
func IsSiswaNotFound(err error) bool    { return errors.Is(err, ErrSiswaNotFound) }
func IsLampiranNotFound(err error) bool { return errors.Is(err, ErrLampiranNotFound) }

type LampiranController struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewLampiranController(db *gorm.DB, cfg *config.Config) *LampiranController {
	return &LampiranController{
		db:  db,
		cfg: cfg,
	}
}

/* =========================
   GET ALL BY NIS
========================= */

func (c *LampiranController) GetByNIS(
	ctx context.Context,
	nis string,
) (dto.SiswaLampiranResponse, error) {

	var siswa model.Siswa

	if err := c.db.WithContext(ctx).
		First(&siswa, "siswa_nis = ?", nis).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.SiswaLampiranResponse{}, ErrSiswaNotFound
		}

		return dto.SiswaLampiranResponse{}, err
	}

	var list []model.Lampiran

	if err := c.db.WithContext(ctx).
		Where("siswa_nis = ?", nis).
		Order("uploaded_at DESC").
		Find(&list).Error; err != nil {

		return dto.SiswaLampiranResponse{}, err
	}

	return dto.MapSiswaLampiranToResponse(siswa, list), nil
}

/* =========================
   UPLOAD FILE
========================= */

func (c *LampiranController) Upload(
	ctx context.Context,
	nis string,
	dokumenJenis string,
	fileName string,
	mime string,
	sizeBytes int64,
	reader io.Reader,
) (dto.LampiranResponse, error) {

	// validasi siswa
	if !c.siswaExists(ctx, nis) {
		return dto.LampiranResponse{}, ErrSiswaNotFound
	}

	// validasi jenis dokumen
	if err := validateDokumenJenis(dokumenJenis); err != nil {
		return dto.LampiranResponse{}, err
	}

	// validasi mime
	if err := validateMime(mime); err != nil {
		return dto.LampiranResponse{}, err
	}

	// validasi ukuran
	if err := validateSize(sizeBytes); err != nil {
		return dto.LampiranResponse{}, err
	}

	// buat nama file final
	finalFileName := buildLocalFileName(
		dokumenJenis,
		fileName,
	)

	// path folder siswa
	studentDir := filepath.Join(
		c.cfg.StoragePath,
		nis,
	)

	// buat folder jika belum ada
	if err := os.MkdirAll(studentDir, os.ModePerm); err != nil {
		return dto.LampiranResponse{}, err
	}

	// relative path untuk DB
	relativePath := filepath.Join(
		nis,
		finalFileName,
	)

	// full path fisik
	fullPath := filepath.Join(
		c.cfg.StoragePath,
		relativePath,
	)

	// simpan file ke disk
	dst, err := os.Create(fullPath)
	if err != nil {
		return dto.LampiranResponse{}, err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, reader); err != nil {
		return dto.LampiranResponse{}, err
	}

	now := time.Now()

	var out model.Lampiran
	var oldObjectKey string

	txErr := c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		var existing model.Lampiran

		qerr := tx.First(
			&existing,
			"siswa_nis = ? AND dokumen_jenis = ?",
			nis,
			dokumenJenis,
		).Error

		if qerr != nil && !errors.Is(qerr, gorm.ErrRecordNotFound) {
			return qerr
		}

		// replace existing
		if qerr == nil {

			// simpan file lama untuk dihapus nanti
			oldObjectKey = existing.ObjectKey

			existing.StorageProvider = "local"
			existing.ObjectKey = relativePath
			existing.FileName = finalFileName
			existing.MimeType = mime
			existing.SizeBytes = sizeBytes
			existing.UploadedAt = now

			if err := tx.Save(&existing).Error; err != nil {
				return err
			}

			out = existing

			return nil
		}

		// create new
		newRec := model.Lampiran{
			SiswaNIS:        nis,
			DokumenJenis:    dokumenJenis,
			StorageProvider: "local",
			ObjectKey:       relativePath,
			FileName:        finalFileName,
			MimeType:        mime,
			SizeBytes:       sizeBytes,
			UploadedAt:      now,
		}

		if err := tx.Create(&newRec).Error; err != nil {
			return err
		}

		out = newRec

		return nil
	})

	// rollback cleanup jika DB gagal
	if txErr != nil {

		// hapus file baru agar tidak orphan
		_ = os.Remove(fullPath)

		return dto.LampiranResponse{}, txErr
	}

	// delete file lama jika replace
	if oldObjectKey != "" && oldObjectKey != relativePath {

		oldPath := filepath.Join(
			c.cfg.StoragePath,
			oldObjectKey,
		)

		_ = os.Remove(oldPath)
	}

	return dto.MapLampiranToResponse(out), nil
}

/* =========================
   DOWNLOAD / PREVIEW FILE
========================= */

func (c *LampiranController) DownloadByID(
	ctx context.Context,
	idStr string,
) (string, io.ReadCloser, error) {

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return "", nil, fmt.Errorf("%w: invalid id", ErrValidation)
	}

	var l model.Lampiran

	if err := c.db.WithContext(ctx).
		First(&l, "id_lampiran = ?", id).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, ErrLampiranNotFound
		}

		return "", nil, err
	}

	fullPath := filepath.Join(
		c.cfg.StoragePath,
		l.ObjectKey,
	)

	file, err := os.Open(fullPath)
	if err != nil {
		return "", nil, err
	}

	contentType := l.MimeType

	if contentType == "" {
		contentType = "application/octet-stream"
	}

	return contentType, file, nil
}

/* =========================
   DELETE FILE
========================= */

func (c *LampiranController) DeleteByJenis(
	ctx context.Context,
	nis,
	jenis string,
) (dto.DeleteLampiranResponse, error) {

	if !c.siswaExists(ctx, nis) {
		return dto.DeleteLampiranResponse{}, ErrSiswaNotFound
	}

	if err := validateDokumenJenis(jenis); err != nil {
		return dto.DeleteLampiranResponse{}, err
	}

	var l model.Lampiran

	if err := c.db.WithContext(ctx).
		First(
			&l,
			"siswa_nis = ? AND dokumen_jenis = ?",
			nis,
			jenis,
		).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.DeleteLampiranResponse{}, ErrLampiranNotFound
		}

		return dto.DeleteLampiranResponse{}, err
	}

	// delete DB dulu
	if err := c.db.WithContext(ctx).
		Delete(&l).Error; err != nil {

		return dto.DeleteLampiranResponse{}, err
	}

	// delete physical file
	fullPath := filepath.Join(
		c.cfg.StoragePath,
		l.ObjectKey,
	)

	_ = os.Remove(fullPath)

	return dto.DeleteLampiranResponse{
		Message:      "lampiran deleted",
		SiswaNIS:     nis,
		DokumenJenis: jenis,
		ObjectKey:    l.ObjectKey,
	}, nil
}

/* =========================
   HELPERS
========================= */

func (c *LampiranController) siswaExists(
	ctx context.Context,
	nis string,
) bool {

	var count int64

	err := c.db.WithContext(ctx).
		Model(&model.Siswa{}).
		Where("siswa_nis = ?", nis).
		Count(&count).Error

	return err == nil && count > 0
}

func buildLocalFileName(
	jenis,
	original string,
) string {

	ext := strings.ToLower(
		filepath.Ext(original),
	)

	return jenis + ext
}

/* =========================
   VALIDATION
========================= */

var allowedDokumenJenis = map[string]struct{}{
	"profil-picture": {},
	"ktp-ayah":       {},
	"ktp-ibu":        {},
	"kk":             {},
	"ijazah":         {},
	"spd":            {},
	"akta":           {},
	"sp":             {},
}

func validateDokumenJenis(jenis string) error {

	jenis = strings.TrimSpace(jenis)

	if jenis == "" {
		return fmt.Errorf(
			"%w: dokumen_jenis is required",
			ErrValidation,
		)
	}

	if _, ok := allowedDokumenJenis[jenis]; !ok {
		return fmt.Errorf(
			"%w: dokumen_jenis not allowed",
			ErrValidation,
		)
	}

	return nil
}

func validateMime(mime string) error {

	mime = strings.ToLower(
		strings.TrimSpace(mime),
	)

	allowed := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
	}

	if !allowed[mime] {
		return fmt.Errorf(
			"%w: unsupported file type",
			ErrValidation,
		)
	}

	return nil
}

func validateSize(size int64) error {

	if size <= 0 {
		return fmt.Errorf(
			"%w: invalid file size",
			ErrValidation,
		)
	}

	// max 200kb
	if size > 200*1024 {
		return fmt.Errorf(
			"%w: max file size is 200kb",
			ErrValidation,
		)
	}

	return nil
}
