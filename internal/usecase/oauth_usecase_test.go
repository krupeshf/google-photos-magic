package usecase

import (
	"testing"
	"time"

	"golang.org/x/oauth2"
)

// MockOAuthService is a mock implementation for testing
type MockOAuthService struct {
	config     *oauth2.Config
	token      *oauth2.Token
	err        error
	authURL    string
	stateValue string
}

func (m *MockOAuthService) GetClient() (*oauth2.Config, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.config, nil
}

func (m *MockOAuthService) LoadToken() (*oauth2.Token, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.token, nil
}

func (m *MockOAuthService) SaveToken(tok *oauth2.Token) error {
	if m.err != nil {
		return m.err
	}
	m.token = tok
	return nil
}

func (m *MockOAuthService) ExchangeCode(code string) (*oauth2.Token, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &oauth2.Token{
		AccessToken: "mock-access-token",
		TokenType:   "Bearer",
		Expiry:      time.Now().Add(1 * time.Hour),
	}, nil
}

func (m *MockOAuthService) GetAuthURL() string {
	return m.authURL
}

func (m *MockOAuthService) GetAuthURLWithState(state string) string {
	m.stateValue = state
	return m.authURL + "?state=" + state
}

func TestOAuthUseCase_CompleteAuthentication(t *testing.T) {
	// Arrange
	mockService := &MockOAuthService{}
	useCase := NewOAuthUseCase(mockService)
	code := "test-auth-code"

	// Act
	err := useCase.CompleteAuthentication(code)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if mockService.token == nil {
		t.Error("Expected token to be saved")
	}

	if mockService.token.AccessToken != "mock-access-token" {
		t.Errorf("Expected access token 'mock-access-token', got '%s'", mockService.token.AccessToken)
	}
}

func TestOAuthUseCase_GetAuthURL(t *testing.T) {
	// Arrange
	expectedURL := "https://accounts.google.com/oauth/authorize"
	mockService := &MockOAuthService{authURL: expectedURL}
	useCase := NewOAuthUseCase(mockService)

	// Act
	url := useCase.GetAuthURL()

	// Assert
	if url != expectedURL {
		t.Errorf("Expected URL '%s', got '%s'", expectedURL, url)
	}
}

func TestOAuthUseCase_LoadToken(t *testing.T) {
	// Arrange
	expectedToken := &oauth2.Token{AccessToken: "test-token"}
	mockService := &MockOAuthService{token: expectedToken}
	useCase := NewOAuthUseCase(mockService)

	// Act
	token, err := useCase.LoadToken()

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if token != expectedToken {
		t.Errorf("Expected token %v, got %v", expectedToken, token)
	}
}

func TestOAuthUseCase_AuthenticateClient_WithValidToken(t *testing.T) {
	// Arrange
	config := &oauth2.Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/photoslibrary.readonly"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/oauth/authorize",
			TokenURL: "https://oauth2.googleapis.com/token",
		},
	}

	validToken := &oauth2.Token{
		AccessToken: "valid-token",
		Expiry:      time.Now().Add(1 * time.Hour),
	}

	mockService := &MockOAuthService{
		config: config,
		token:  validToken,
	}

	useCase := NewOAuthUseCase(mockService)

	// Act
	resultConfig, err := useCase.AuthenticateClient()

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if resultConfig != config {
		t.Errorf("Expected config %v, got %v", config, resultConfig)
	}
}

func TestOAuthUseCase_AuthenticateClient_WithExpiredToken(t *testing.T) {
	// Arrange
	config := &oauth2.Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/photoslibrary.readonly"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/oauth/authorize",
			TokenURL: "https://oauth2.googleapis.com/token",
		},
	}

	expiredToken := &oauth2.Token{
		AccessToken: "expired-token",
		Expiry:      time.Now().Add(-1 * time.Hour), // Expired
	}

	mockService := &MockOAuthService{
		config: config,
		token:  expiredToken,
	}

	useCase := NewOAuthUseCase(mockService)

	// Act
	resultConfig, err := useCase.AuthenticateClient()

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if resultConfig != config {
		t.Errorf("Expected config %v, got %v", config, resultConfig)
	}
}
