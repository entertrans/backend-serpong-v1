package main

import (
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"

	"github.com/entertrans/go-base-project.git/internal/config"
	"github.com/entertrans/go-base-project.git/internal/service"

)

func main() {
	_ = godotenv.Load()

	cfg := config.LoadConfig()

	srv, err := service.NewDriveServiceOAuth(context.Background(), service.DriveOAuthCfg{
		ClientID:     cfg.GDriveOAuthClientID,
		ClientSecret: cfg.GDriveOAuthClientSecret,
		TokenPath:    cfg.GDriveTokenPath,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Ping ringan supaya jelas token valid
	about, err := srv.About.Get().Fields("user,storageQuota").Do()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("OAuth OK. Token saved to:", cfg.GDriveTokenPath)
	if about.User != nil {
		fmt.Println("Logged in as:", about.User.EmailAddress)
	}
}
