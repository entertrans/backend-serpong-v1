package controller

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"google.golang.org/api/drive/v3"
	"gorm.io/gorm"

	"github.com/entertrans/go-base-project.git/internal/config"
	"github.com/entertrans/go-base-project.git/internal/dto"
	"github.com/entertrans/go-base-project.git/internal/model"
	"github.com/entertrans/go-base-project.git/internal/service"
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
	return &LampiranController{db: db, cfg: cfg}
}

/* =========================
   Public methods
========================= */

func (c *LampiranController) GetByNIS(ctx context.Context, nis string) (dto.SiswaLampiranResponse, error) {
	var siswa model.Siswa
	if err := c.db.WithContext(ctx).First(&siswa, "siswa_nis = ?", nis).Error; err != nil {
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

// Upload JPG via backend -> Google Drive -> upsert DB
func (c *LampiranController) Upload(
	ctx context.Context,
	nis string,
	dokumenJenis string,
	fileName string,
	mime string,
	sizeBytes int64,
	reader io.Reader,
) (dto.LampiranResponse, error) {

	if !c.siswaExists(ctx, nis) {
		return dto.LampiranResponse{}, ErrSiswaNotFound
	}

	if err := validateDokumenJenis(dokumenJenis); err != nil {
		return dto.LampiranResponse{}, err
	}
	if err := validateJpegOnly(mime); err != nil {
		return dto.LampiranResponse{}, err
	}
	if err := validateSize(sizeBytes); err != nil {
		return dto.LampiranResponse{}, err
	}
	if strings.TrimSpace(fileName) == "" {
		fileName = dokumenJenis + ".jpg"
	}

	driveFileID, err := c.uploadToDrive(ctx, reader, buildDriveFileName(nis, dokumenJenis, fileName), mime)
	if err != nil {
		return dto.LampiranResponse{}, err
	}

	now := time.Now()
	var out model.Lampiran
	var oldDriveID string

	txErr := c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing model.Lampiran
		qerr := tx.First(&existing, "siswa_nis = ? AND dokumen_jenis = ?", nis, dokumenJenis).Error
		if qerr != nil && !errors.Is(qerr, gorm.ErrRecordNotFound) {
			return qerr
		}

		if qerr == nil {
			// replace
			oldDriveID = existing.ObjectKey

			existing.StorageProvider = "gdrive"
			existing.ObjectKey = driveFileID
			existing.FileName = fileName
			existing.MimeType = mime
			existing.SizeBytes = sizeBytes
			existing.ETag = "" // tidak dipakai di drive
			existing.UploadedAt = now

			if uerr := tx.Save(&existing).Error; uerr != nil {
				return uerr
			}
			out = existing
			return nil
		}

		newRec := model.Lampiran{
			SiswaNIS:        nis,
			DokumenJenis:    dokumenJenis,
			StorageProvider: "gdrive",
			ObjectKey:       driveFileID,
			FileName:        fileName,
			MimeType:        mime,
			SizeBytes:       sizeBytes,
			ETag:            "",
			UploadedAt:      now,
		}
		if cerr := tx.Create(&newRec).Error; cerr != nil {
			return cerr
		}
		out = newRec
		return nil
	})

	if txErr != nil {
		// best-effort cleanup file yang barusan diupload agar tidak jadi sampah
		_ = c.deleteDriveFileBestEffort(ctx, driveFileID)
		return dto.LampiranResponse{}, txErr
	}

	// best-effort delete file lama jika replace
	if oldDriveID != "" && oldDriveID != driveFileID {
		_ = c.deleteDriveFileBestEffort(ctx, oldDriveID)
	}

	return dto.MapLampiranToResponse(out), nil
}

// Stream file by lampiran id (untuk handler menulis ke ResponseWriter)
func (c *LampiranController) DownloadByID(ctx context.Context, idStr string) (contentType string, body io.ReadCloser, err error) {
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return "", nil, fmt.Errorf("%w: invalid id", ErrValidation)
	}

	var l model.Lampiran
	if err := c.db.WithContext(ctx).First(&l, "id_lampiran = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, ErrLampiranNotFound
		}
		return "", nil, err
	}

	if l.StorageProvider != "gdrive" {
		return "", nil, fmt.Errorf("%w: unsupported storage_provider", ErrValidation)
	}

	srv, err := service.NewDriveServiceOAuth(ctx, service.DriveOAuthCfg{
	ClientID:     c.cfg.GDriveOAuthClientID,
	ClientSecret: c.cfg.GDriveOAuthClientSecret,
	TokenPath:    c.cfg.GDriveTokenPath,
})
	if err != nil {
		return "", nil, err
	}

	// Ambil metadata untuk content-type (optional)
	meta, _ := srv.Files.Get(l.ObjectKey).Fields("mimeType").Do()

	resp, err := srv.Files.Get(l.ObjectKey).Download()
	if err != nil {
		return "", nil, err
	}

	ct := l.MimeType
	if ct == "" && meta != nil && meta.MimeType != "" {
		ct = meta.MimeType
	}
	if ct == "" {
		ct = "image/jpeg"
	}

	return ct, resp.Body, nil
}

func (c *LampiranController) DeleteByJenis(ctx context.Context, nis, jenis string) (dto.DeleteLampiranResponse, error) {
	if !c.siswaExists(ctx, nis) {
		return dto.DeleteLampiranResponse{}, ErrSiswaNotFound
	}
	if err := validateDokumenJenis(jenis); err != nil {
		return dto.DeleteLampiranResponse{}, err
	}

	var l model.Lampiran
	if err := c.db.WithContext(ctx).First(&l, "siswa_nis = ? AND dokumen_jenis = ?", nis, jenis).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.DeleteLampiranResponse{}, ErrLampiranNotFound
		}
		return dto.DeleteLampiranResponse{}, err
	}

	driveID := l.ObjectKey

	if err := c.db.WithContext(ctx).Delete(&l).Error; err != nil {
		return dto.DeleteLampiranResponse{}, err
	}

	_ = c.deleteDriveFileBestEffort(ctx, driveID)

	return dto.DeleteLampiranResponse{
		Message:      "lampiran deleted",
		SiswaNIS:     nis,
		DokumenJenis: jenis,
		ObjectKey:    driveID,
	}, nil
}

/* =========================
   Internal helpers
========================= */

func (c *LampiranController) siswaExists(ctx context.Context, nis string) bool {
	var count int64
	err := c.db.WithContext(ctx).Model(&model.Siswa{}).
		Where("siswa_nis = ?", nis).
		Count(&count).Error
	return err == nil && count > 0
}

func (c *LampiranController) uploadToDrive(ctx context.Context, file io.Reader, fileName, mime string) (string, error) {
	if strings.TrimSpace(c.cfg.GDriveFolderID) == "" {
  return "", fmt.Errorf("%w: GDRIVE_FOLDER_ID is empty", ErrValidation)
}

	srv, err := service.NewDriveServiceOAuth(ctx, service.DriveOAuthCfg{
	ClientID:     c.cfg.GDriveOAuthClientID,
	ClientSecret: c.cfg.GDriveOAuthClientSecret,
	TokenPath:    c.cfg.GDriveTokenPath,
})
	if err != nil {
		return "", err
	}

	f := &drive.File{
		Name:    fileName,
		Parents: []string{c.cfg.GDriveFolderID},
	}

	res, err := srv.Files.Create(f).
		Media(file).
		SupportsAllDrives(true).
		Fields("id").
		Do()
	if err != nil {
		return "", err
	}
	return res.Id, nil
}

func (c *LampiranController) deleteDriveFileBestEffort(ctx context.Context, fileID string) error {
	if strings.TrimSpace(fileID) == "" {
		return nil
	}
	srv, err := service.NewDriveServiceOAuth(ctx, service.DriveOAuthCfg{
	ClientID:     c.cfg.GDriveOAuthClientID,
	ClientSecret: c.cfg.GDriveOAuthClientSecret,
	TokenPath:    c.cfg.GDriveTokenPath,
})
	if err != nil {
		return err
	}
	return srv.Files.Delete(fileID).SupportsAllDrives(true).Do()
}

/* =========================
   Validation
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
		return fmt.Errorf("%w: dokumen_jenis is required", ErrValidation)
	}
	if _, ok := allowedDokumenJenis[jenis]; !ok {
		return fmt.Errorf("%w: dokumen_jenis not allowed", ErrValidation)
	}
	return nil
}

func validateJpegOnly(mime string) error {
	m := strings.ToLower(strings.TrimSpace(mime))
	if m != "image/jpeg" && m != "image/jpg" {
		return fmt.Errorf("%w: only image/jpeg is allowed", ErrValidation)
	}
	return nil
}

func validateSize(size int64) error {
	if size <= 0 {
		return fmt.Errorf("%w: size_bytes must be > 0", ErrValidation)
	}
	if size > 2*1024*1024 {
		return fmt.Errorf("%w: max file size is 2MB", ErrValidation)
	}
	return nil
}

func buildDriveFileName(nis, jenis, original string) string {
	// aman & konsisten: nis_jenis_timestamp_original
	base := strings.ReplaceAll(original, " ", "_")
	return fmt.Sprintf("%s_%s_%d_%s", nis, jenis, time.Now().Unix(), base)
}