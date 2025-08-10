package usecase

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

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

// CompleteAuthenticationWithServer automatically completes OAuth2 flow using a local server
func (uc *OAuthUseCase) CompleteAuthenticationWithServer() error {
	log.Printf("Starting OAuth2 flow with local server...")

	// Generate a random state for security
	state := "random-state-" + fmt.Sprintf("%d", time.Now().Unix())

	// Get the authorization URL with the state
	authURL := uc.oauthService.GetAuthURLWithState(state)

	// Create a channel to receive the authorization code
	codeChan := make(chan string, 1)
	errChan := make(chan error, 1)

	// Start local server to capture the callback
	server := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Handle OAuth callback
			if r.URL.Path == "/oauth2callback" {
				query := r.URL.Query()

				// Check if there's an error
				if err := query.Get("error"); err != "" {
					errChan <- fmt.Errorf("OAuth error: %s", err)
					return
				}

				// Verify state parameter
				if receivedState := query.Get("state"); receivedState != state {
					errChan <- fmt.Errorf("invalid state parameter")
					return
				}

				// Get the authorization code
				code := query.Get("code")
				if code == "" {
					errChan <- fmt.Errorf("no authorization code received")
					return
				}

				// Send success response to browser
				w.Header().Set("Content-Type", "text/html")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`
					<html>
						<body>
							<h1>Authorization Successful!</h1>
							<p>You can close this window now.</p>
							<script>window.close();</script>
						</body>
					</html>
				`))

				// Send the code through the channel
				codeChan <- code
			} else {
				http.NotFound(w, r)
			}
		}),
	}

	// Start the server in a goroutine
	go func() {
		log.Printf("Starting local server on http://localhost:8080")
		log.Printf("Visit this URL in your browser to authorize:")
		log.Printf("%s", authURL)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("server error: %v", err)
		}
	}()

	// Wait for the authorization code or an error
	select {
	case code := <-codeChan:
		// Shutdown the server gracefully
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx)

		// Complete the authentication
		return uc.CompleteAuthentication(code)

	case err := <-errChan:
		// Shutdown the server gracefully
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx)
		return err

	case <-time.After(10 * time.Minute):
		// Timeout after 10 minutes
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx)
		return fmt.Errorf("OAuth flow timed out")
	}
}

// GetAuthURL returns the authorization URL for the OAuth2 flow
func (uc *OAuthUseCase) GetAuthURL() string {
	return uc.oauthService.GetAuthURL()
}

// LoadToken loads the OAuth token from storage
func (uc *OAuthUseCase) LoadToken() (*oauth2.Token, error) {
	return uc.oauthService.LoadToken()
}
