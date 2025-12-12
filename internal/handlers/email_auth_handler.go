package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/omarshah0/rest-api-with-social-auth/internal/config"
	"github.com/omarshah0/rest-api-with-social-auth/internal/middleware"
	"github.com/omarshah0/rest-api-with-social-auth/internal/models"
	"github.com/omarshah0/rest-api-with-social-auth/internal/services"
	"github.com/omarshah0/rest-api-with-social-auth/internal/utils"
)

type EmailAuthHandler struct {
	authService  *services.AuthService
	cookieConfig config.CookieConfig
}

func NewEmailAuthHandler(authService *services.AuthService, cookieConfig config.CookieConfig) *EmailAuthHandler {
	return &EmailAuthHandler{
		authService:  authService,
		cookieConfig: cookieConfig,
	}
}

// Register handles user registration with email/password
func (h *EmailAuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.UserRegister
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, utils.ErrorTypeBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if err := utils.ValidateStruct(req); err != nil {
		utils.SendError(w, http.StatusBadRequest, utils.ErrorTypeValidation, err.Error())
		return
	}

	// Register user
	user, err := h.authService.Register(r.Context(), &req)
	if err != nil {
		if err.Error() == "user with this email already exists" {
			utils.SendError(w, http.StatusConflict, utils.ErrorTypeConflict, err.Error())
			return
		}
		utils.SendError(w, http.StatusInternalServerError, utils.ErrorTypeInternalServer, "Failed to register user")
		return
	}

	response := map[string]interface{}{
		"user":    user,
		"message": "Registration successful. Please check your email to verify your account.",
	}

	utils.SendSuccess(w, http.StatusCreated, utils.ResponseTypeAuth, response, "User registered successfully")
}

// Login handles user login with email/password
func (h *EmailAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		models.UserLogin
		DeviceType string `json:"device_type" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, utils.ErrorTypeBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if err := utils.ValidateStruct(req); err != nil {
		utils.SendError(w, http.StatusBadRequest, utils.ErrorTypeValidation, err.Error())
		return
	}

	// Validate device type
	if req.DeviceType != "web" && req.DeviceType != "mobile" {
		utils.SendError(w, http.StatusBadRequest, utils.ErrorTypeBadRequest, "Invalid device_type. Must be 'web' or 'mobile'")
		return
	}

	// Authenticate user
	authResponse, err := h.authService.Login(r.Context(), &req.UserLogin, models.DeviceType(req.DeviceType))
	if err != nil {
		utils.SendError(w, http.StatusUnauthorized, utils.ErrorTypeUnauthorized, "Invalid email or password")
		return
	}

	// Set cookies (for web app)
	h.setAuthCookies(w, authResponse.AccessToken, authResponse.RefreshToken)

	utils.SendSuccess(w, http.StatusOK, utils.ResponseTypeAuth, authResponse, "Login successful")
}

// VerifyEmail handles email verification
func (h *EmailAuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		utils.SendError(w, http.StatusBadRequest, utils.ErrorTypeBadRequest, "Verification token is required")
		return
	}

	// Verify email
	err := h.authService.VerifyEmail(token)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, utils.ErrorTypeBadRequest, "Invalid or expired verification token")
		return
	}

	utils.SendSuccess(w, http.StatusOK, utils.ResponseTypeAction, nil, "Email verified successfully")
}

// ForgotPassword handles forgot password request
func (h *EmailAuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req models.ForgotPassword
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, utils.ErrorTypeBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if err := utils.ValidateStruct(req); err != nil {
		utils.SendError(w, http.StatusBadRequest, utils.ErrorTypeValidation, err.Error())
		return
	}

	// Process forgot password
	err := h.authService.ForgotPassword(r.Context(), &req)
	if err != nil {
		// Don't expose errors to prevent email enumeration
		utils.SendSuccess(w, http.StatusOK, utils.ResponseTypeAction, nil, "If the email exists, a password reset link has been sent")
		return
	}

	utils.SendSuccess(w, http.StatusOK, utils.ResponseTypeAction, nil, "If the email exists, a password reset link has been sent")
}

// ResetPassword handles password reset with token
func (h *EmailAuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req models.ResetPassword
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, utils.ErrorTypeBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if err := utils.ValidateStruct(req); err != nil {
		utils.SendError(w, http.StatusBadRequest, utils.ErrorTypeValidation, err.Error())
		return
	}

	// Reset password
	err := h.authService.ResetPassword(r.Context(), &req)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, utils.ErrorTypeBadRequest, err.Error())
		return
	}

	utils.SendSuccess(w, http.StatusOK, utils.ResponseTypeAction, nil, "Password reset successfully")
}

// ChangePassword handles password change for authenticated user
func (h *EmailAuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, utils.ErrorTypeUnauthorized, "Authentication required")
		return
	}

	var req models.PasswordChange
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, utils.ErrorTypeBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if err := utils.ValidateStruct(req); err != nil {
		utils.SendError(w, http.StatusBadRequest, utils.ErrorTypeValidation, err.Error())
		return
	}

	// Change password
	err := h.authService.ChangePassword(r.Context(), userID, &req)
	if err != nil {
		if err.Error() == "current password is incorrect" {
			utils.SendError(w, http.StatusBadRequest, utils.ErrorTypeBadRequest, err.Error())
			return
		}
		utils.SendError(w, http.StatusInternalServerError, utils.ErrorTypeInternalServer, "Failed to change password")
		return
	}

	utils.SendSuccess(w, http.StatusOK, utils.ResponseTypeAction, nil, "Password changed successfully. Please login again.")
}

// ResendVerification resends email verification link
func (h *EmailAuthHandler) ResendVerification(w http.ResponseWriter, r *http.Request) {
	var req models.ResendVerification
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, utils.ErrorTypeBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if err := utils.ValidateStruct(req); err != nil {
		utils.SendError(w, http.StatusBadRequest, utils.ErrorTypeValidation, err.Error())
		return
	}

	// Resend verification
	err := h.authService.ResendVerification(r.Context(), &req)
	if err != nil {
		if err.Error() == "email already verified" {
			utils.SendError(w, http.StatusBadRequest, utils.ErrorTypeBadRequest, err.Error())
			return
		}
		utils.SendError(w, http.StatusInternalServerError, utils.ErrorTypeInternalServer, "Failed to resend verification email")
		return
	}

	utils.SendSuccess(w, http.StatusOK, utils.ResponseTypeAction, nil, "Verification email sent successfully")
}

// Helper function to set auth cookies
func (h *EmailAuthHandler) setAuthCookies(w http.ResponseWriter, accessToken, refreshToken string) {
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

