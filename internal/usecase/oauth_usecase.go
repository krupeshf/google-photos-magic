package usecase

import (
	"log"

	"golang.org/x/oauth2"
	"krupesh.faldu/internal/domain"
)

// OAuthUseCase implements the business logic for OAuth operations
type OAuthUseCase struct {
	oauthService domain.OAuthService
}

// NewOAuthUseCase creates a new instance of OAuthUseCase
func NewOAuthUseCase(oauthService domain.OAuthService) *OAuthUseCase {
	return &OAuthUseCase{
		oauthService: oauthService,
	}
}

// AuthenticateClient handles the OAuth2 authentication flow
func (uc *OAuthUseCase) AuthenticateClient() (*oauth2.Config, error) {
	log.Printf("Starting OAuth2 authentication...")

	config, err := uc.oauthService.GetClient()
	if err != nil {
		log.Printf("Failed to get OAuth config: %v", err)
		return nil, err
	}

	// Try to load existing token
	token, err := uc.oauthService.LoadToken()
	if err != nil {
		log.Printf("No existing token found, starting OAuth flow...")
		return config, nil
	}

	// Validate token
	if token.Valid() {
		log.Printf("Valid token found, authentication successful")
		return config, nil
	}

	log.Printf("Token expired, starting OAuth flow...")
	return config, nil
}

// CompleteAuthentication completes the OAuth2 flow with the authorization code
func (uc *OAuthUseCase) CompleteAuthentication(code string) error {
	log.Printf("Completing OAuth2 authentication with code...")

	token, err := uc.oauthService.ExchangeCode(code)
	if err != nil {
		log.Printf("Failed to exchange code for token: %v", err)
		return err
	}

	err = uc.oauthService.SaveToken(token)
	if err != nil {
		log.Printf("Failed to save token: %v", err)
		return err
	}

	log.Printf("Authentication completed successfully")
	return nil
}

// GetAuthURL returns the authorization URL for the OAuth2 flow
func (uc *OAuthUseCase) GetAuthURL() string {
	return uc.oauthService.GetAuthURL()
}

// LoadToken loads the OAuth token from storage
func (uc *OAuthUseCase) LoadToken() (*oauth2.Token, error) {
	return uc.oauthService.LoadToken()
}
