package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/entertrans/backend-bogor.git/internal/config"
	authmodule "github.com/entertrans/backend-bogor.git/internal/modules/auth"
	avatarmodule "github.com/entertrans/backend-bogor.git/internal/modules/avatar"

	financemodule "github.com/entertrans/backend-bogor.git/internal/modules/finance"
	lampiranmodule "github.com/entertrans/backend-bogor.git/internal/modules/lampiran"
	masterdatamodule "github.com/entertrans/backend-bogor.git/internal/modules/masterdata"
	pingmodule "github.com/entertrans/backend-bogor.git/internal/modules/ping"
	rapor "github.com/entertrans/backend-bogor.git/internal/modules/rapor"
	siswamodule "github.com/entertrans/backend-bogor.git/internal/modules/siswa"
	"github.com/entertrans/backend-bogor.git/internal/router"
	"github.com/entertrans/backend-bogor.git/pkg/database"
	"github.com/entertrans/backend-bogor.git/pkg/logger"
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
	// r.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"http://localhost:5173"},
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
	// 	AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
	// 	AllowCredentials: true,
	// 	MaxAge:           12 * time.Hour,
	// }))
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
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
		avatarmodule.Register,
		rapor.Register,
		financemodule.Register,
		// kelasmodule.Register,
	)

	port := ":" + cfg.AppPort
	log.Printf("Server starting on port %s", port)
	if err := r.Run(port); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}
