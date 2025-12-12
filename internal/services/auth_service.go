package services

import (
	"context"
	"fmt"
	"time"

	"github.com/omarshah0/rest-api-with-social-auth/internal/models"
	"github.com/omarshah0/rest-api-with-social-auth/internal/repositories"
)

type AuthService struct {
	userRepo        *repositories.UserRepository
	oauthRepo       *repositories.OAuthProviderRepository
	adminRepo       *repositories.AdminRepository
	jwtService      *JWTService
	oauthService    *OAuthService
	passwordService *PasswordService
	emailService    *EmailService
}

func NewAuthService(
	userRepo *repositories.UserRepository,
	oauthRepo *repositories.OAuthProviderRepository,
	adminRepo *repositories.AdminRepository,
	jwtService *JWTService,
	oauthService *OAuthService,
	passwordService *PasswordService,
	emailService *EmailService,
) *AuthService {
	return &AuthService{
		userRepo:        userRepo,
		oauthRepo:       oauthRepo,
		adminRepo:       adminRepo,
		jwtService:      jwtService,
		oauthService:    oauthService,
		passwordService: passwordService,
		emailService:    emailService,
	}
}

type AuthResponse struct {
	User         *models.User `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	IsAdmin      bool         `json:"is_admin"`
}

// AuthenticateWithOAuth authenticates or creates a user via OAuth
func (s *AuthService) AuthenticateWithOAuth(ctx context.Context, provider models.OAuthProviderType, code string) (*AuthResponse, error) {
	var userInfo *OAuthUserInfo
	var err error

	// Exchange code for user info based on provider
	switch provider {
	case models.ProviderGoogle:
		if !s.oauthService.IsGoogleEnabled() {
			return nil, fmt.Errorf("google oauth is not enabled")
		}
		userInfo, err = s.oauthService.ExchangeGoogleCode(ctx, code)
	case models.ProviderFacebook:
		if !s.oauthService.IsFacebookEnabled() {
			return nil, fmt.Errorf("facebook oauth is not enabled")
		}
		userInfo, err = s.oauthService.ExchangeFacebookCode(ctx, code)
	default:
		return nil, fmt.Errorf("unsupported oauth provider")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user info from provider: %w", err)
	}

	// Check if OAuth provider is already linked
	oauthProvider, err := s.oauthRepo.GetByProviderAndUserID(provider, userInfo.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check oauth provider: %w", err)
	}

	var user *models.User

	if oauthProvider != nil {
		// OAuth provider exists, get the user
		user, err = s.userRepo.GetByID(oauthProvider.UserID)
		if err != nil {
			return nil, fmt.Errorf("failed to get user: %w", err)
		}
		if user == nil {
			return nil, fmt.Errorf("user not found")
		}
	} else {
		// Check if user exists by email (for account linking)
		user, err = s.userRepo.GetByEmail(userInfo.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to check user by email: %w", err)
		}

		if user == nil {
			// Create new user (OAuth users have verified emails)
			userCreate := &models.UserCreate{
				Email:          userInfo.Email,
				Name:           userInfo.Name,
				ProfilePicture: userInfo.Picture,
			}
			user, err = s.userRepo.Create(userCreate)
			if err != nil {
				return nil, fmt.Errorf("failed to create user: %w", err)
			}
		} else {
			// User exists - update profile picture if provided and not already set
			if userInfo.Picture != nil && user.ProfilePicture == nil {
				user.ProfilePicture = userInfo.Picture
				err = s.userRepo.Update(user.ID, user)
				if err != nil {
					// Log but don't fail
					fmt.Printf("Warning: Failed to update profile picture for OAuth user: %v\n", err)
				}
			}

			if !user.EmailVerified {
				// If user exists but email not verified, verify it now since OAuth provider confirmed it
				// This happens when user registered with email/password but didn't verify yet
				// Now they're linking OAuth which proves email ownership
				err = s.userRepo.MarkEmailVerified(user.ID)
				if err != nil {
					// Log but don't fail - just continue
					fmt.Printf("Warning: Failed to auto-verify email for OAuth user: %v\n", err)
				}
			}
			// Refresh user object to get updated data
			user, _ = s.userRepo.GetByID(user.ID)
		}

		// Link OAuth provider to user
		_, err = s.oauthRepo.Create(user.ID, provider, userInfo.ID, userInfo.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to link oauth provider: %w", err)
		}
	}

	// Check if user is blocked
	if user.Blocked {
		return nil, fmt.Errorf("user account is blocked")
	}

	// Check if user is admin
	isAdmin, err := s.adminRepo.IsAdmin(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check admin status: %w", err)
	}

	// Generate tokens
	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		IsAdmin:      isAdmin,
	}, nil
}

// RefreshTokens refreshes the access and refresh tokens
func (s *AuthService) RefreshTokens(refreshToken string) (*AuthResponse, error) {
	accessToken, newRefreshToken, err := s.jwtService.RefreshTokens(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh tokens: %w", err)
	}

	// Validate the new access token to get user info
	claims, err := s.jwtService.ValidateAccessToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to validate access token: %w", err)
	}

	// Get user
	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Check if user is blocked
	if user.Blocked {
		return nil, fmt.Errorf("user account is blocked")
	}

	// Check if user is admin
	isAdmin, err := s.adminRepo.IsAdmin(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check admin status: %w", err)
	}

	return &AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		IsAdmin:      isAdmin,
	}, nil
}

// Logout logs out a user by revoking their refresh token
func (s *AuthService) Logout(userID int64) error {
	return s.jwtService.RevokeRefreshToken(userID)
}

// AuthenticateWithOAuthUserInfo authenticates or creates a user with OAuthUserInfo (for mobile ID token verification)
func (s *AuthService) AuthenticateWithOAuthUserInfo(ctx context.Context, provider models.OAuthProviderType, userInfo *OAuthUserInfo) (*AuthResponse, error) {
	// Check if OAuth provider is already linked
	oauthProvider, err := s.oauthRepo.GetByProviderAndUserID(provider, userInfo.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check oauth provider: %w", err)
	}

	var user *models.User

	if oauthProvider != nil {
		// OAuth provider exists, get the user
		user, err = s.userRepo.GetByID(oauthProvider.UserID)
		if err != nil {
			return nil, fmt.Errorf("failed to get user: %w", err)
		}
		if user == nil {
			return nil, fmt.Errorf("user not found")
		}
	} else {
		// Check if user exists by email (for account linking)
		user, err = s.userRepo.GetByEmail(userInfo.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to check user by email: %w", err)
		}

		if user == nil {
			// Create new user (OAuth users have verified emails)
			userCreate := &models.UserCreate{
				Email:          userInfo.Email,
				Name:           userInfo.Name,
				ProfilePicture: userInfo.Picture,
			}
			user, err = s.userRepo.Create(userCreate)
			if err != nil {
				return nil, fmt.Errorf("failed to create user: %w", err)
			}
		} else {
			// User exists - update profile picture if provided and not already set
			if userInfo.Picture != nil && user.ProfilePicture == nil {
				user.ProfilePicture = userInfo.Picture
				err = s.userRepo.Update(user.ID, user)
				if err != nil {
					// Log but don't fail
					fmt.Printf("Warning: Failed to update profile picture for OAuth user: %v\n", err)
				}
			}

			if !user.EmailVerified {
				// If user exists but email not verified, verify it now since OAuth provider confirmed it
				err = s.userRepo.MarkEmailVerified(user.ID)
				if err != nil {
					// Log but don't fail
					fmt.Printf("Warning: Failed to auto-verify email for OAuth user: %v\n", err)
				}
			}
			// Refresh user object to get updated data
			user, _ = s.userRepo.GetByID(user.ID)
		}

		// Link OAuth provider to user
		_, err = s.oauthRepo.Create(user.ID, provider, userInfo.ID, userInfo.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to link oauth provider: %w", err)
		}
	}

	// Check if user is blocked
	if user.Blocked {
		return nil, fmt.Errorf("user account is blocked")
	}

	// Check if user is admin
	isAdmin, err := s.adminRepo.IsAdmin(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check admin status: %w", err)
	}

	// Generate tokens
	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		IsAdmin:      isAdmin,
	}, nil
}

// Register registers a new user with email and password
func (s *AuthService) Register(ctx context.Context, req *models.UserRegister) (*models.User, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return nil, fmt.Errorf("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := s.passwordService.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user, err := s.userRepo.CreateWithPassword(req, hashedPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate verification token
	token, err := s.passwordService.GenerateToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate verification token: %w", err)
	}

	// Set token expiry (24 hours)
	expires := time.Now().Add(24 * time.Hour)
	err = s.userRepo.SetVerificationToken(user.ID, token, expires)
	if err != nil {
		return nil, fmt.Errorf("failed to set verification token: %w", err)
	}

	// Send verification email
	err = s.emailService.SendVerificationEmail(user.Email, user.Name, token)
	if err != nil {
		// Log error but don't fail registration
		fmt.Printf("Failed to send verification email: %v\n", err)
	}

	return user, nil
}

// Login authenticates a user with email and password
func (s *AuthService) Login(ctx context.Context, req *models.UserLogin) (*AuthResponse, error) {
	// Get user with password
	user, err := s.userRepo.GetByEmailWithPassword(req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil || user.HashedPassword == nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Verify password
	err = s.passwordService.VerifyPassword(req.Password, *user.HashedPassword)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Check if user is blocked
	if user.Blocked {
		return nil, fmt.Errorf("user account is blocked")
	}

	// Optional: Check if email is verified (uncomment to enforce)
	// if !user.EmailVerified {
	// 	return nil, fmt.Errorf("please verify your email address")
	// }

	// Check if user is admin
	isAdmin, err := s.adminRepo.IsAdmin(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check admin status: %w", err)
	}

	// Generate tokens
	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		IsAdmin:      isAdmin,
	}, nil
}

// VerifyEmail verifies a user's email address
func (s *AuthService) VerifyEmail(token string) error {
	return s.userRepo.VerifyEmail(token)
}

// ForgotPassword initiates password reset process
func (s *AuthService) ForgotPassword(ctx context.Context, req *models.ForgotPassword) error {
	// Get user
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Don't reveal if email exists or not (security best practice)
	if user == nil {
		return nil
	}

	// Generate reset token
	token, err := s.passwordService.GenerateToken()
	if err != nil {
		return fmt.Errorf("failed to generate reset token: %w", err)
	}

	// Set token expiry (1 hour)
	expires := time.Now().Add(1 * time.Hour)
	err = s.userRepo.SetResetToken(user.ID, token, expires)
	if err != nil {
		return fmt.Errorf("failed to set reset token: %w", err)
	}

	// Send reset email
	err = s.emailService.SendPasswordResetEmail(user.Email, user.Name, token)
	if err != nil {
		return fmt.Errorf("failed to send reset email: %w", err)
	}

	return nil
}

// ResetPassword resets user password with token
func (s *AuthService) ResetPassword(ctx context.Context, req *models.ResetPassword) error {
	// Get user by reset token
	user, err := s.userRepo.GetByResetToken(req.Token)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("invalid or expired reset token")
	}

	// Check if token is expired
	if s.passwordService.IsTokenExpired(user.ResetTokenExpires) {
		return fmt.Errorf("reset token has expired")
	}

	// Hash new password
	hashedPassword, err := s.passwordService.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	err = s.userRepo.UpdatePassword(user.ID, hashedPassword)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Revoke all refresh tokens (force logout from all devices)
	err = s.jwtService.RevokeRefreshToken(user.ID)
	if err != nil {
		// Log but don't fail
		fmt.Printf("Failed to revoke refresh tokens: %v\n", err)
	}

	// Send confirmation email
	err = s.emailService.SendPasswordChangedEmail(user.Email, user.Name)
	if err != nil {
		// Log but don't fail
		fmt.Printf("Failed to send password changed email: %v\n", err)
	}

	return nil
}

// ChangePassword changes user password (authenticated)
func (s *AuthService) ChangePassword(ctx context.Context, userID int64, req *models.PasswordChange) error {
	// Get user with password
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	// Get user with password to verify old password
	userWithPassword, err := s.userRepo.GetByEmailWithPassword(user.Email)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if userWithPassword.HashedPassword == nil {
		return fmt.Errorf("password authentication not set up for this account")
	}

	// Verify old password
	err = s.passwordService.VerifyPassword(req.OldPassword, *userWithPassword.HashedPassword)
	if err != nil {
		return fmt.Errorf("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := s.passwordService.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Update password
	err = s.userRepo.UpdatePassword(userID, hashedPassword)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Revoke all refresh tokens (force logout from all devices)
	err = s.jwtService.RevokeRefreshToken(userID)
	if err != nil {
		// Log but don't fail
		fmt.Printf("Failed to revoke refresh tokens: %v\n", err)
	}

	// Send confirmation email
	err = s.emailService.SendPasswordChangedEmail(user.Email, user.Name)
	if err != nil {
		// Log but don't fail
		fmt.Printf("Failed to send password changed email: %v\n", err)
	}

	return nil
}

// ResendVerification resends verification email
func (s *AuthService) ResendVerification(ctx context.Context, req *models.ResendVerification) error {
	// Get user
	user, err := s.userRepo.GetByEmailWithPassword(req.Email)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	// Check if already verified
	if user.EmailVerified {
		return fmt.Errorf("email already verified")
	}

	// Generate new verification token
	token, err := s.passwordService.GenerateToken()
	if err != nil {
		return fmt.Errorf("failed to generate verification token: %w", err)
	}

	// Set token expiry (24 hours)
	expires := time.Now().Add(24 * time.Hour)
	err = s.userRepo.SetVerificationToken(user.ID, token, expires)
	if err != nil {
		return fmt.Errorf("failed to set verification token: %w", err)
	}

	// Send verification email
	err = s.emailService.SendVerificationEmail(user.Email, user.Name, token)
	if err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	return nil
}

// GetProfile retrieves user profile with linked accounts
func (s *AuthService) GetProfile(ctx context.Context, userID int64) (*models.ProfileResponse, error) {
	// Get user by ID
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Get all OAuth providers for the user
	oauthProviders, err := s.oauthRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get oauth providers: %w", err)
	}

	// Build linked accounts
	linkedAccounts := models.LinkedAccounts{
		Email:       user.Email,
		PasswordSet: user.HashedPassword != nil,
		Google:      nil,
		Facebook:    nil,
	}

	// Find Google and Facebook providers
	for _, provider := range oauthProviders {
		if provider.Provider == models.ProviderGoogle {
			linkedAccounts.Google = &models.LinkedAccount{
				Provider: string(provider.Provider),
				LinkedAt: provider.CreatedAt,
			}
		} else if provider.Provider == models.ProviderFacebook {
			linkedAccounts.Facebook = &models.LinkedAccount{
				Provider: string(provider.Provider),
				LinkedAt: provider.CreatedAt,
			}
		}
	}

	return &models.ProfileResponse{
		User:           user,
		LinkedAccounts: linkedAccounts,
	}, nil
}
