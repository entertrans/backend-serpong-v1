package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"

)

type DriveOAuthCfg struct {
	ClientID     string
	ClientSecret string
	TokenPath    string
}

func NewDriveServiceOAuth(ctx context.Context, cfg DriveOAuthCfg) (*drive.Service, error) {
	if cfg.ClientID == "" || cfg.ClientSecret == "" {
		return nil, fmt.Errorf("missing GDRIVE_OAUTH_CLIENT_ID or GDRIVE_OAUTH_CLIENT_SECRET")
	}
	if cfg.TokenPath == "" {
		return nil, fmt.Errorf("missing GDRIVE_TOKEN_PATH")
	}

	oauthCfg := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Scopes:       []string{drive.DriveScope},
		Endpoint:     google.Endpoint,
		// Untuk flow CLI modern, gunakan localhost callback
		RedirectURL: "urn:ietf:wg:oauth:2.0:oob",
	}

	tok, err := tokenFromFile(cfg.TokenPath)
	if err != nil {
		// Token belum ada → login sekali
		tok, err = getTokenFromWeb(oauthCfg)
		if err != nil {
			return nil, err
		}
		if err := saveToken(cfg.TokenPath, tok); err != nil {
			return nil, err
		}
	}

	client := oauthCfg.Client(ctx, tok)
	return drive.NewService(ctx, option.WithHTTPClient(client))
}

func tokenFromFile(path string) (*oauth2.Token, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var tok oauth2.Token
	if err := json.NewDecoder(f).Decode(&tok); err != nil {
		return nil, err
	}
	return &tok, nil
}

func saveToken(path string, tok *oauth2.Token) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(tok)
}

// Flow: print URL, user open, paste code
func getTokenFromWeb(cfg *oauth2.Config) (*oauth2.Token, error) {
	authURL := cfg.AuthCodeURL("state-token", oauth2.AccessTypeOffline, oauth2.ApprovalForce)

	fmt.Println("Open this URL in your browser:")
	fmt.Println(authURL)
	fmt.Print("Paste authorization code here: ")

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		return nil, err
	}

	return cfg.Exchange(context.Background(), code)
}
