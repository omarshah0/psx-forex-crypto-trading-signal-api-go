package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/omarshah0/rest-api-with-social-auth/internal/config"
	"github.com/omarshah0/rest-api-with-social-auth/internal/middleware"
	"github.com/omarshah0/rest-api-with-social-auth/internal/models"
	"github.com/omarshah0/rest-api-with-social-auth/internal/services"
	"github.com/omarshah0/rest-api-with-social-auth/internal/utils"
)

type AuthHandler struct {
	authService  *services.AuthService
	oauthService *services.OAuthService
	cookieConfig config.CookieConfig
}

func NewAuthHandler(authService *services.AuthService, oauthService *services.OAuthService, cookieConfig config.CookieConfig) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		oauthService: oauthService,
		cookieConfig: cookieConfig,
	}
}

// Refresh token handler
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var refreshToken string

	// Try to get refresh token from request body
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err == nil && body.RefreshToken != "" {
		refreshToken = body.RefreshToken
	}

	// If not in body, try cookie
	if refreshToken == "" {
		cookie, err := r.Cookie("refresh_token")
		if err != nil {
			utils.SendError(w, http.StatusUnauthorized, utils.ErrorTypeUnauthorized, "Refresh token required")
			return
		}
		refreshToken = cookie.Value
	}

	// Refresh tokens
	authResponse, err := h.authService.RefreshTokens(refreshToken)
	if err != nil {
		utils.SendError(w, http.StatusUnauthorized, utils.ErrorTypeUnauthorized, "Invalid or expired refresh token")
		return
	}

	// Set new cookies
	h.setAuthCookies(w, authResponse.AccessToken, authResponse.RefreshToken)

	utils.SendSuccess(w, http.StatusOK, utils.ResponseTypeAuth, authResponse, "Tokens refreshed successfully")
}

// Logout handler
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, utils.ErrorTypeUnauthorized, "Authentication required")
		return
	}

	// Logout user (revoke refresh token)
	if err := h.authService.Logout(userID); err != nil {
		utils.SendError(w, http.StatusInternalServerError, utils.ErrorTypeInternalServer, "Failed to logout")
		return
	}

	// Clear cookies
	h.clearAuthCookies(w)

	utils.SendSuccess(w, http.StatusOK, utils.ResponseTypeAction, nil, "Logged out successfully")
}

// Helper functions

func (h *AuthHandler) setAuthCookies(w http.ResponseWriter, accessToken, refreshToken string) {
	// Access token cookie (15 minutes)
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		MaxAge:   15 * 60,
		HttpOnly: h.cookieConfig.HTTPOnly,
		Secure:   h.cookieConfig.Secure,
		SameSite: http.SameSiteLaxMode,
		Domain:   h.cookieConfig.Domain,
	})

	// Refresh token cookie (7 days default)
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		MaxAge:   7 * 24 * 60 * 60,
		HttpOnly: h.cookieConfig.HTTPOnly,
		Secure:   h.cookieConfig.Secure,
		SameSite: http.SameSiteLaxMode,
		Domain:   h.cookieConfig.Domain,
	})
}

func (h *AuthHandler) clearAuthCookies(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
	})
}

// VerifyGoogleIDToken verifies Google ID token (for React Native/Expo)
func (h *AuthHandler) VerifyGoogleIDToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		IDToken string `json:"id_token" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, utils.ErrorTypeBadRequest, "Invalid request body")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.SendError(w, http.StatusBadRequest, utils.ErrorTypeValidation, err.Error())
		return
	}

	// Verify ID token with Google
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get("https://oauth2.googleapis.com/tokeninfo?id_token=" + req.IDToken)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, utils.ErrorTypeInternalServer, "Failed to verify token with Google")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Try to read error message from response
		body, _ := io.ReadAll(resp.Body)
		var errorResp struct {
			Error            string `json:"error"`
			ErrorDescription string `json:"error_description"`
		}
		if err := json.Unmarshal(body, &errorResp); err == nil && errorResp.Error != "" {
			utils.SendError(w, http.StatusUnauthorized, utils.ErrorTypeUnauthorized, "Invalid ID token: "+errorResp.ErrorDescription)
			return
		}
		utils.SendError(w, http.StatusUnauthorized, utils.ErrorTypeUnauthorized, "Invalid or expired ID token")
		return
	}

	var tokenInfo struct {
		Sub     string `json:"sub"` // User ID
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"` // Profile picture URL
		Aud     string `json:"aud"`     // Client ID
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenInfo); err != nil {
		utils.SendError(w, http.StatusInternalServerError, utils.ErrorTypeInternalServer, "Failed to parse token info")
		return
	}

	// Create OAuthUserInfo from token
	var picture *string
	if tokenInfo.Picture != "" {
		picture = &tokenInfo.Picture
	}

	userInfo := &services.OAuthUserInfo{
		ID:      tokenInfo.Sub,
		Email:   tokenInfo.Email,
		Name:    tokenInfo.Name,
		Picture: picture,
	}

	// Create or login user using the same OAuth flow
	authResponse, err := h.authService.AuthenticateWithOAuthUserInfo(r.Context(), models.ProviderGoogle, userInfo)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, utils.ErrorTypeInternalServer, "Authentication failed")
		return
	}

	utils.SendSuccess(w, http.StatusOK, utils.ResponseTypeAuth, authResponse, "Authentication successful")
}

// VerifyFacebookAccessToken verifies Facebook access token (for React Native/Expo)
func (h *AuthHandler) VerifyFacebookAccessToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AccessToken string `json:"access_token" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, utils.ErrorTypeBadRequest, "Invalid request body")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.SendError(w, http.StatusBadRequest, utils.ErrorTypeValidation, err.Error())
		return
	}

	// Verify access token with Facebook (including picture)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get("https://graph.facebook.com/me?fields=id,name,email,picture&access_token=" + req.AccessToken)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, utils.ErrorTypeInternalServer, "Failed to verify token")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		utils.SendError(w, http.StatusUnauthorized, utils.ErrorTypeUnauthorized, "Invalid access token")
		return
	}

	// Facebook returns picture as a nested object
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

	if err := json.NewDecoder(resp.Body).Decode(&fbResponse); err != nil {
		utils.SendError(w, http.StatusInternalServerError, utils.ErrorTypeInternalServer, "Failed to parse user info")
		return
	}

	// Extract picture URL if available
	var picture *string
	if fbResponse.Picture.Data.URL != "" {
		picture = &fbResponse.Picture.Data.URL
	}

	userInfo := services.OAuthUserInfo{
		ID:      fbResponse.ID,
		Email:   fbResponse.Email,
		Name:    fbResponse.Name,
		Picture: picture,
	}

	// Create or login user using the same OAuth flow
	authResponse, err := h.authService.AuthenticateWithOAuthUserInfo(r.Context(), models.ProviderFacebook, &userInfo)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, utils.ErrorTypeInternalServer, "Authentication failed")
		return
	}

	utils.SendSuccess(w, http.StatusOK, utils.ResponseTypeAuth, authResponse, "Authentication successful")
}

// ExchangeGoogleCode exchanges Google authorization code for JWT tokens (for React Web/Mobile)
func (h *AuthHandler) ExchangeGoogleCode(w http.ResponseWriter, r *http.Request) {
	if !h.oauthService.IsGoogleEnabled() {
		utils.SendError(w, http.StatusForbidden, utils.ErrorTypeForbidden, "Google OAuth is not enabled")
		return
	}

	var req struct {
		Code string `json:"code" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, utils.ErrorTypeBadRequest, "Invalid request body")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.SendError(w, http.StatusBadRequest, utils.ErrorTypeValidation, err.Error())
		return
	}

	// Authenticate with OAuth using the authorization code
	authResponse, err := h.authService.AuthenticateWithOAuth(r.Context(), models.ProviderGoogle, req.Code)
	if err != nil {
		utils.SendError(w, http.StatusUnauthorized, utils.ErrorTypeUnauthorized, "Authentication failed: "+err.Error())
		return
	}

	utils.SendSuccess(w, http.StatusOK, utils.ResponseTypeAuth, authResponse, "Authentication successful")
}

// ExchangeFacebookCode exchanges Facebook authorization code for JWT tokens (for React Web/Mobile)
func (h *AuthHandler) ExchangeFacebookCode(w http.ResponseWriter, r *http.Request) {
	if !h.oauthService.IsFacebookEnabled() {
		utils.SendError(w, http.StatusForbidden, utils.ErrorTypeForbidden, "Facebook OAuth is not enabled")
		return
	}

	var req struct {
		Code string `json:"code" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, utils.ErrorTypeBadRequest, "Invalid request body")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.SendError(w, http.StatusBadRequest, utils.ErrorTypeValidation, err.Error())
		return
	}

	// Authenticate with OAuth using the authorization code
	authResponse, err := h.authService.AuthenticateWithOAuth(r.Context(), models.ProviderFacebook, req.Code)
	if err != nil {
		utils.SendError(w, http.StatusUnauthorized, utils.ErrorTypeUnauthorized, "Authentication failed: "+err.Error())
		return
	}

	utils.SendSuccess(w, http.StatusOK, utils.ResponseTypeAuth, authResponse, "Authentication successful")
}
