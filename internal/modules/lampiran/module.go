package lampiranmodule

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/entertrans/backend-bogor.git/internal/config"
	lampirancontroller "github.com/entertrans/backend-bogor.git/internal/controller"
	lampiranhandler "github.com/entertrans/backend-bogor.git/internal/handler"
	"github.com/entertrans/backend-bogor.git/internal/middleware"
)

func Register(rg *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	ctrl := lampirancontroller.NewLampiranController(db, cfg)
	h := lampiranhandler.NewLampiranHandler(ctrl)

	lampiran := rg.Group("/lampiran")
	lampiran.Use(middleware.AuthMiddleware(cfg))
	{

		lampiran.GET("/:nis", h.GetByNIS)

		// Upload via backend (multipart/form-data)
		// POST /api/v1/lampiran/:nis/upload/:jenis
		lampiran.POST("/:nis/upload/:jenis", h.UploadByJenis)

		// Stream file (image/jpeg) from Google Drive
		// GET /api/v1/lampiran/file/:id/view
		lampiran.GET("/file/:id/view", h.ViewByID)

		// Delete record + delete Drive file
		lampiran.DELETE("/:nis/:jenis", h.DeleteByJenis)
	}
}
