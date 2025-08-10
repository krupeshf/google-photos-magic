package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"krupesh.faldu/internal/domain"
)

const (
	tokenFile = "token.json"
)

// OAuthRepository implements the OAuthService interface
type OAuthRepository struct {
	config *oauth2.Config
}

// NewOAuthRepository creates a new instance of OAuthRepository
func NewOAuthRepository() (domain.OAuthService, error) {
	// Load OAuth2 config from credentials file
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		return nil, fmt.Errorf("unable to read credentials.json: %v", err)
	}

	// Configure OAuth2 scopes for Google Photos
	config, err := google.ConfigFromJSON(b,
		"https://www.googleapis.com/auth/photoslibrary.readonly.appcreateddata",
		"https://www.googleapis.com/auth/photoslibrary.appendonly",
		"https://www.googleapis.com/auth/photoslibrary.edit.appcreateddata")
	if err != nil {
		return nil, fmt.Errorf("unable to parse credentials.json: %v", err)
	}

	return &OAuthRepository{
		config: config,
	}, nil
}

// GetClient returns the OAuth2 configuration
func (r *OAuthRepository) GetClient() (*oauth2.Config, error) {
	return r.config, nil
}

// LoadToken loads the OAuth2 token from disk
func (r *OAuthRepository) LoadToken() (*oauth2.Token, error) {
	f, err := os.Open(tokenFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var tok oauth2.Token
	err = json.NewDecoder(f).Decode(&tok)
	return &tok, err
}

// SaveToken saves the OAuth2 token to disk
func (r *OAuthRepository) SaveToken(tok *oauth2.Token) error {
	f, err := os.Create(tokenFile)
	if err != nil {
		return fmt.Errorf("failed to create token file: %v", err)
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(tok)
}

// ExchangeCode exchanges an authorization code for an access token
func (r *OAuthRepository) ExchangeCode(code string) (*oauth2.Token, error) {
	return r.config.Exchange(context.Background(), code)
}

// GetAuthURL returns the authorization URL for the OAuth2 flow
func (r *OAuthRepository) GetAuthURL() string {
	return r.config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
}
