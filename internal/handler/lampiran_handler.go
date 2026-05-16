package handler

import (
	"io"
	"net/http"

	"github.com/entertrans/go-base-project.git/internal/controller"
	"github.com/gin-gonic/gin"
)

type LampiranHandler struct {
	ctrl *controller.LampiranController
}

func NewLampiranHandler(
	ctrl *controller.LampiranController,
) *LampiranHandler {

	return &LampiranHandler{
		ctrl: ctrl,
	}
}

/* =========================
   GET ALL
========================= */

func (h *LampiranHandler) GetByNIS(c *gin.Context) {

	nis := c.Param("nis")

	res, err := h.ctrl.GetByNIS(
		c.Request.Context(),
		nis,
	)

	if err != nil {

		if controller.IsSiswaNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "siswa not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, res)
}

/* =========================
   UPLOAD
========================= */

func (h *LampiranHandler) UploadByJenis(c *gin.Context) {

	nis := c.Param("nis")
	jenis := c.Param("jenis")

	fh, err := c.FormFile("file")
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "file is required",
		})

		return
	}

	file, err := fh.Open()
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to open file",
		})

		return
	}
	defer file.Close()

	// baca 512 byte pertama untuk detect mime asli
	buffer := make([]byte, 512)

	_, err = file.Read(buffer)
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read file",
		})

		return
	}

	// reset pointer file
	_, _ = file.Seek(0, io.SeekStart)

	mime := http.DetectContentType(buffer)

	res, err := h.ctrl.Upload(
		c.Request.Context(),
		nis,
		jenis,
		fh.Filename,
		mime,
		fh.Size,
		file,
	)

	if err != nil {

		if controller.IsValidation(err) {

			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		if controller.IsSiswaNotFound(err) {

			c.JSON(http.StatusNotFound, gin.H{
				"error": "siswa not found",
			})

			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "upload success",
		"lampiran": res,
	})
}

/* =========================
   VIEW FILE
========================= */

func (h *LampiranHandler) ViewByID(c *gin.Context) {

	id := c.Param("id")

	contentType, body, err := h.ctrl.DownloadByID(
		c.Request.Context(),
		id,
	)

	if err != nil {

		if controller.IsLampiranNotFound(err) {

			c.JSON(http.StatusNotFound, gin.H{
				"error": "lampiran not found",
			})

			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}
	defer body.Close()

	c.Header("Content-Type", contentType)

	c.Status(http.StatusOK)

	_, _ = io.Copy(c.Writer, body)
}

/* =========================
   DELETE
========================= */

func (h *LampiranHandler) DeleteByJenis(c *gin.Context) {

	nis := c.Param("nis")
	jenis := c.Param("jenis")

	res, err := h.ctrl.DeleteByJenis(
		c.Request.Context(),
		nis,
		jenis,
	)

	if err != nil {

		if controller.IsLampiranNotFound(err) {

			c.JSON(http.StatusNotFound, gin.H{
				"error": "lampiran not found",
			})

			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, res)
}
