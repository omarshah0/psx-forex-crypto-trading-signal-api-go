package handlers

import (
	"net/http"

	"github.com/omarshah0/rest-api-with-social-auth/internal/middleware"
	"github.com/omarshah0/rest-api-with-social-auth/internal/services"
	"github.com/omarshah0/rest-api-with-social-auth/internal/utils"
)

type ProfileHandler struct {
	authService *services.AuthService
}

func NewProfileHandler(authService *services.AuthService) *ProfileHandler {
	return &ProfileHandler{
		authService: authService,
	}
}

// GetProfile retrieves the authenticated user's profile with linked accounts
func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		utils.SendError(w, http.StatusUnauthorized, utils.ErrorTypeUnauthorized, "Authentication required")
		return
	}

	// Get profile
	profile, err := h.authService.GetProfile(r.Context(), userID)
	if err != nil {
		if err.Error() == "user not found" {
			utils.SendError(w, http.StatusNotFound, utils.ErrorTypeNotFound, "User not found")
			return
		}
		utils.SendError(w, http.StatusInternalServerError, utils.ErrorTypeInternalServer, "Failed to retrieve profile")
		return
	}

	utils.SendSuccess(w, http.StatusOK, utils.ResponseTypeResource, profile, "Profile retrieved successfully")
}
