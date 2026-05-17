// internal/handler/avatar_handler.go
package handler

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/entertrans/backend-bogor.git/internal/config"
	"github.com/gin-gonic/gin"
)

type AvatarHandler struct {
	cfg *config.Config
}

func NewAvatarHandler(cfg *config.Config) *AvatarHandler {
	return &AvatarHandler{cfg: cfg}
}

// GET /avatar/:nis
func (h *AvatarHandler) ServeAvatar(c *gin.Context) {
	nis := c.Param("nis")

	// Path ke file profil-picture
	avatarPath := h.getAvatarPath(nis)

	// Cek file exist
	if info, err := os.Stat(avatarPath); err == nil && !info.IsDir() {
		// Serve file dengan cache header
		c.Header("Cache-Control", "public, max-age=86400") // cache 1 hari
		c.Header("Content-Type", "image/jpeg")
		c.File(avatarPath)
		return
	}

	// FALLBACK: serve default avatar
	h.serveDefaultAvatar(c)
}

func (h *AvatarHandler) getAvatarPath(nis string) string {
	// Daftar ekstensi yang mungkin
	extensions := []string{".jpg", ".jpeg", ".png", ".JPG", ".JPEG", ".PNG"}

	for _, ext := range extensions {
		path := filepath.Join(h.cfg.StoragePath, nis, "profil-picture"+ext)
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// Return default path (akan trigger fallback)
	return filepath.Join(h.cfg.StoragePath, nis, "profil-picture.jpg")
}

func (h *AvatarHandler) serveDefaultAvatar(c *gin.Context) {
	// Coba beberapa kemungkinan path default
	defaultPaths := []string{
		filepath.Join(h.cfg.StoragePath, "defaults", "avatar.jpg"),
		filepath.Join(h.cfg.StoragePath, "defaults", "avatar.png"),
		filepath.Join("storage", "defaults", "avatar.jpg"), // fallback
		filepath.Join(".", "storage", "defaults", "avatar.jpg"),
	}

	for _, defaultPath := range defaultPaths {
		if info, err := os.Stat(defaultPath); err == nil && !info.IsDir() {
			// Cache default avatar lebih lama (1 minggu)
			c.Header("Cache-Control", "public, max-age=604800")
			c.Header("Content-Type", "image/jpeg")
			c.File(defaultPath)
			return
		}
	}

	// Ultimate fallback: return 204 No Content (bukan 404)
	// Biar browser tidak request ulang terus
	c.Header("Cache-Control", "public, max-age=3600") // cache 1 jam
	c.Status(http.StatusNoContent)
}
