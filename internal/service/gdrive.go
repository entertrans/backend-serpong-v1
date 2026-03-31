package service

import (
	"context"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

func NewDriveService(ctx context.Context, credPath string) (*drive.Service, error) {
	b, err := os.ReadFile(credPath)
	if err != nil {
		return nil, err
	}

	conf, err := google.JWTConfigFromJSON(b, drive.DriveScope)
	if err != nil {
		return nil, err
	}

	client := conf.Client(ctx)
	return drive.New(client)
}
