package middleware

import (
	"net/http"

	"github.com/omarshah0/rest-api-with-social-auth/internal/repositories"
	"github.com/omarshah0/rest-api-with-social-auth/internal/utils"
)

type AdminMiddleware struct {
	adminRepo *repositories.AdminRepository
}

func NewAdminMiddleware(adminRepo *repositories.AdminRepository) *AdminMiddleware {
	return &AdminMiddleware{adminRepo: adminRepo}
}

// RequireAdmin middleware checks if the authenticated user is an admin
func (m *AdminMiddleware) RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from context (set by AuthMiddleware)
		userID, ok := GetUserIDFromContext(r.Context())
		if !ok {
			utils.SendError(w, http.StatusUnauthorized, utils.ErrorTypeUnauthorized, "Authentication required")
			return
		}

		// Check if user is admin
		isAdmin, err := m.adminRepo.IsAdmin(userID)
		if err != nil {
			utils.SendError(w, http.StatusInternalServerError, utils.ErrorTypeInternalServer, "Failed to verify admin status")
			return
		}

		if !isAdmin {
			utils.SendError(w, http.StatusForbidden, utils.ErrorTypeForbidden, "Admin access required")
			return
		}

		next.ServeHTTP(w, r)
	})
}

