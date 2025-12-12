package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
	googleapi "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

type OAuthService struct {
	googleConfig    *oauth2.Config
	facebookConfig  *oauth2.Config
	googleEnabled   bool
	facebookEnabled bool
}

type OAuthUserInfo struct {
	ID      string  `json:"id"`
	Email   string  `json:"email"`
	Name    string  `json:"name"`
	Picture *string `json:"picture,omitempty"`
}

func NewOAuthService(googleClientID, googleClientSecret, googleRedirectURL string, googleEnabled bool,
	facebookClientID, facebookClientSecret, facebookRedirectURL string, facebookEnabled bool) *OAuthService {

	var googleConfig *oauth2.Config
	if googleEnabled {
		googleConfig = &oauth2.Config{
			ClientID:     googleClientID,
			ClientSecret: googleClientSecret,
			RedirectURL:  googleRedirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		}
	}

	var facebookConfig *oauth2.Config
	if facebookEnabled {
		facebookConfig = &oauth2.Config{
			ClientID:     facebookClientID,
			ClientSecret: facebookClientSecret,
			RedirectURL:  facebookRedirectURL,
			Scopes:       []string{"email", "public_profile"},
			Endpoint:     facebook.Endpoint,
		}
	}

	return &OAuthService{
		googleConfig:    googleConfig,
		facebookConfig:  facebookConfig,
		googleEnabled:   googleEnabled,
		facebookEnabled: facebookEnabled,
	}
}

// GetGoogleAuthURL returns the Google OAuth authorization URL
func (s *OAuthService) GetGoogleAuthURL(state string) (string, error) {
	if !s.googleEnabled || s.googleConfig == nil {
		return "", fmt.Errorf("google oauth is not enabled")
	}
	return s.googleConfig.AuthCodeURL(state, oauth2.AccessTypeOffline), nil
}

// GetFacebookAuthURL returns the Facebook OAuth authorization URL
func (s *OAuthService) GetFacebookAuthURL(state string) (string, error) {
	if !s.facebookEnabled || s.facebookConfig == nil {
		return "", fmt.Errorf("facebook oauth is not enabled")
	}
	return s.facebookConfig.AuthCodeURL(state), nil
}

// ExchangeGoogleCode exchanges the authorization code for user info
func (s *OAuthService) ExchangeGoogleCode(ctx context.Context, code string) (*OAuthUserInfo, error) {
	if !s.googleEnabled || s.googleConfig == nil {
		return nil, fmt.Errorf("google oauth is not enabled")
	}

	token, err := s.googleConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Get user info using the Google OAuth2 API
	oauth2Service, err := googleapi.NewService(ctx, option.WithTokenSource(s.googleConfig.TokenSource(ctx, token)))
	if err != nil {
		return nil, fmt.Errorf("failed to create oauth2 service: %w", err)
	}

	userInfo, err := oauth2Service.Userinfo.Get().Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	// Get profile picture if available
	var picture *string
	if userInfo.Picture != "" {
		picture = &userInfo.Picture
	}

	return &OAuthUserInfo{
		ID:      userInfo.Id,
		Email:   userInfo.Email,
		Name:    userInfo.Name,
		Picture: picture,
	}, nil
}

// ExchangeFacebookCode exchanges the authorization code for user info
func (s *OAuthService) ExchangeFacebookCode(ctx context.Context, code string) (*OAuthUserInfo, error) {
	if !s.facebookEnabled || s.facebookConfig == nil {
		return nil, fmt.Errorf("facebook oauth is not enabled")
	}

	token, err := s.facebookConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Get user info from Facebook (including picture)
	client := s.facebookConfig.Client(ctx, token)
	resp, err := client.Get("https://graph.facebook.com/me?fields=id,name,email,picture")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Facebook returns picture as a nested object, so we need a custom struct
	var fbResponse struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture struct {
			Data struct {
				URL string `json:"url"`
			} `json:"data"`
		} `json:"picture"`
	}

	if err := json.Unmarshal(body, &fbResponse); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	// Extract picture URL if available
	var picture *string
	if fbResponse.Picture.Data.URL != "" {
		picture = &fbResponse.Picture.Data.URL
	}

	return &OAuthUserInfo{
		ID:      fbResponse.ID,
		Email:   fbResponse.Email,
		Name:    fbResponse.Name,
		Picture: picture,
	}, nil
}

// IsGoogleEnabled returns whether Google OAuth is enabled
func (s *OAuthService) IsGoogleEnabled() bool {
	return s.googleEnabled
}

// IsFacebookEnabled returns whether Facebook OAuth is enabled
func (s *OAuthService) IsFacebookEnabled() bool {
	return s.facebookEnabled
}
