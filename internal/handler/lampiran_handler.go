package handler

import (
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/entertrans/go-base-project.git/internal/controller"
	"github.com/gin-gonic/gin"
)

type LampiranHandler struct {
	ctrl *controller.LampiranController
}

func NewLampiranHandler(ctrl *controller.LampiranController) *LampiranHandler {
	return &LampiranHandler{ctrl: ctrl}
}

// GET /lampiran/:nis
func (h *LampiranHandler) GetByNIS(c *gin.Context) {
	nis := c.Param("nis")

	res, err := h.ctrl.GetByNIS(c.Request.Context(), nis)
	if err != nil {
		if controller.IsSiswaNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "siswa not found"})
			return
		}
		if controller.IsValidation(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// POST /lampiran/:nis/upload/:jenis  (multipart/form-data: file)
func (h *LampiranHandler) UploadByJenis(c *gin.Context) {
	nis := c.Param("nis")
	jenis := c.Param("jenis")

	fh, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required (form field: file)"})
		return
	}

	ext := strings.ToLower(filepath.Ext(fh.Filename))
	if ext != ".jpg" && ext != ".jpeg" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only .jpg/.jpeg files are allowed"})
		return
	}

	file, err := fh.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open uploaded file"})
		return
	}
	defer file.Close()

	mime := fh.Header.Get("Content-Type")
	if mime == "" {
		mime = "image/jpeg"
	}

	res, err := h.ctrl.Upload(c.Request.Context(), nis, jenis, fh.Filename, mime, fh.Size, file)
	if err != nil {
		if controller.IsSiswaNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "siswa not found"})
			return
		}
		if controller.IsValidation(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "upload success",
		"lampiran": res,
	})
}

// GET /lampiran/file/:id/view
func (h *LampiranHandler) ViewByID(c *gin.Context) {
	id := c.Param("id")

	contentType, body, err := h.ctrl.DownloadByID(c.Request.Context(), id)
	if err != nil {
		if controller.IsLampiranNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "lampiran not found"})
			return
		}
		if controller.IsValidation(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer body.Close()

	c.Header("Content-Type", contentType)
	c.Status(http.StatusOK)
	_, _ = io.Copy(c.Writer, body)
}

// DELETE /lampiran/:nis/:jenis
func (h *LampiranHandler) DeleteByJenis(c *gin.Context) {
	nis := c.Param("nis")
	jenis := c.Param("jenis")

	res, err := h.ctrl.DeleteByJenis(c.Request.Context(), nis, jenis)
	if err != nil {
		if controller.IsSiswaNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "siswa not found"})
			return
		}
		if controller.IsLampiranNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "lampiran not found"})
			return
		}
		if controller.IsValidation(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
