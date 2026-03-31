package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	AppEnv    string
	AppPort   string
	DBDriver  string
	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
	JWTSecret string
	// Google Drive OAuth
GDriveOAuthClientID     string
GDriveOAuthClientSecret string
GDriveTokenPath         string
GDriveFolderID          string

// 	GDriveCredentialsPath string
// GDriveFolderID        string

}

// Tambahkan DB sebagai field global
var (
	DB  *gorm.DB
	cfg *Config
)

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	cfg = &Config{
		AppEnv:    getEnv("APP_ENV", "development"),
		AppPort:   getEnv("APP_PORT", "8080"),
		DBDriver:  getEnv("DB_DRIVER", "mysql"),
		DBHost:    getEnv("DB_HOST", "localhost"),
		DBPort:    getEnv("DB_PORT", "3306"),
		DBUser:    getEnv("DB_USER", "root"),
		DBPass:    getEnv("DB_PASS", ""),
		DBName:    getEnv("DB_NAME", "go_base"),
		JWTSecret: getEnv("JWT_SECRET", "rahasia"),
		GDriveOAuthClientID:     getEnv("GDRIVE_OAUTH_CLIENT_ID", ""),
GDriveOAuthClientSecret: getEnv("GDRIVE_OAUTH_CLIENT_SECRET", ""),
GDriveTokenPath:         getEnv("GDRIVE_TOKEN_PATH", "storage/gdrive-token.json"),
GDriveFolderID:          getEnv("GDRIVE_FOLDER_ID", ""),

// 		GDriveCredentialsPath: getEnv("GDRIVE_CREDENTIALS_PATH", "storage/gdrive-service-account.json"),
// GDriveFolderID:        getEnv("GDRIVE_FOLDER_ID", ""),

	}

	return cfg
}

// Fungsi untuk inisialisasi database
func InitDatabase() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	// Konfigurasi GORM logger
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Warn,
			Colorful:      true,
		},
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connected successfully")
	return nil
}

// GetDB mengembalikan instance database
func GetDB() *gorm.DB {
	return DB
}

// GetConfig mengembalikan config
func GetConfig() *Config {
	return cfg
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}