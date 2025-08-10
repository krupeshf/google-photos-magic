package domain

import (
	"golang.org/x/oauth2"
)

// OAuthService defines the interface for OAuth operations
type OAuthService interface {
	GetClient() (*oauth2.Config, error)
	LoadToken() (*oauth2.Token, error)
	SaveToken(tok *oauth2.Token) error
	ExchangeCode(code string) (*oauth2.Token, error)
	GetAuthURL() string
}
