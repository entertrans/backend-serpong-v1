package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/entertrans/go-base-project.git/internal/config"
	authmodule "github.com/entertrans/go-base-project.git/internal/modules/auth"
	lampiranmodule "github.com/entertrans/go-base-project.git/internal/modules/lampiran"
	masterdatamodule "github.com/entertrans/go-base-project.git/internal/modules/masterdata"
	pingmodule "github.com/entertrans/go-base-project.git/internal/modules/ping"
	siswamodule "github.com/entertrans/go-base-project.git/internal/modules/siswa"
	"github.com/entertrans/go-base-project.git/internal/router"
	"github.com/entertrans/go-base-project.git/pkg/database"
	"github.com/entertrans/go-base-project.git/pkg/logger"
)

func main() {
	cfg := config.LoadConfig()
	logger.InitLogger(cfg.AppEnv)

	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	r := gin.Default()

	// CORS aman: jangan pakai cors.Default() untuk production
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(logger.LoggerMiddleware())

	// Register modules (tambah module baru cukup 1 baris di sini)
	router.RegisterModules(r, db, cfg,
		pingmodule.Register,
		authmodule.Register,
		siswamodule.Register,
		masterdatamodule.Register,
		lampiranmodule.Register,
		// kelasmodule.Register,
	)

	port := ":" + cfg.AppPort
	log.Printf("Server starting on port %s", port)
	if err := r.Run(port); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}
