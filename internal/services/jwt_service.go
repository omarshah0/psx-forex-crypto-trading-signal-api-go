package services

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/omarshah0/rest-api-with-social-auth/internal/config"
	"github.com/omarshah0/rest-api-with-social-auth/internal/database"
	"github.com/omarshah0/rest-api-with-social-auth/internal/models"
)

type JWTService struct {
	config   *config.JWTConfig
	redisDB  *database.RedisDB
}

type Claims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func NewJWTService(cfg *config.JWTConfig, redisDB *database.RedisDB) *JWTService {
	return &JWTService{
		config:  cfg,
		redisDB: redisDB,
	}
}

// GenerateAccessToken generates a new access token
func (s *JWTService) GenerateAccessToken(userID int64, email string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.AccessExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    s.config.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.AccessSecret))
}

// GenerateRefreshToken generates a new refresh token and stores it in Redis with device type
func (s *JWTService) GenerateRefreshToken(userID int64, email string, deviceType models.DeviceType) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.RefreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    s.config.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.RefreshSecret))
	if err != nil {
		return "", err
	}

	// Store refresh token in Redis with device type
	ctx := context.Background()
	key := fmt.Sprintf("refresh_token:%d:%s", userID, deviceType)
	err = s.redisDB.Set(ctx, key, tokenString, s.config.RefreshExpiry)
	if err != nil {
		return "", fmt.Errorf("failed to store refresh token: %w", err)
	}

	return tokenString, nil
}

// ValidateAccessToken validates an access token and returns the claims
func (s *JWTService) ValidateAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.AccessSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// ValidateRefreshToken validates a refresh token and returns the claims
func (s *JWTService) ValidateRefreshToken(tokenString string, deviceType models.DeviceType) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.RefreshSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// Verify token exists in Redis for the specific device type
		ctx := context.Background()
		key := fmt.Sprintf("refresh_token:%d:%s", claims.UserID, deviceType)
		storedToken, err := s.redisDB.Get(ctx, key)
		if err != nil || storedToken != tokenString {
			return nil, fmt.Errorf("invalid or expired refresh token")
		}
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// RevokeRefreshToken revokes a refresh token for a specific device by removing it from Redis
func (s *JWTService) RevokeRefreshToken(userID int64, deviceType models.DeviceType) error {
	ctx := context.Background()
	key := fmt.Sprintf("refresh_token:%d:%s", userID, deviceType)
	return s.redisDB.Delete(ctx, key)
}

// RevokeAllRefreshTokens revokes all refresh tokens for a user (all devices)
func (s *JWTService) RevokeAllRefreshTokens(userID int64) error {
	ctx := context.Background()
	
	// Delete both web and mobile tokens
	keyWeb := fmt.Sprintf("refresh_token:%d:web", userID)
	keyMobile := fmt.Sprintf("refresh_token:%d:mobile", userID)
	
	// Try to delete both, but don't fail if one doesn't exist
	errWeb := s.redisDB.Delete(ctx, keyWeb)
	errMobile := s.redisDB.Delete(ctx, keyMobile)
	
	// Return error only if both deletions failed
	if errWeb != nil && errMobile != nil {
		return fmt.Errorf("failed to revoke tokens: web error: %v, mobile error: %v", errWeb, errMobile)
	}
	
	return nil
}

// RefreshTokens generates new access and refresh tokens for a specific device
func (s *JWTService) RefreshTokens(refreshToken string, deviceType models.DeviceType) (string, string, error) {
	claims, err := s.ValidateRefreshToken(refreshToken, deviceType)
	if err != nil {
		return "", "", err
	}

	// Revoke old refresh token for this device
	if err := s.RevokeRefreshToken(claims.UserID, deviceType); err != nil {
		return "", "", fmt.Errorf("failed to revoke old refresh token: %w", err)
	}

	// Generate new tokens
	accessToken, err := s.GenerateAccessToken(claims.UserID, claims.Email)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, err := s.GenerateRefreshToken(claims.UserID, claims.Email, deviceType)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return accessToken, newRefreshToken, nil
}

