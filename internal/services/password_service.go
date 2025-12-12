package services

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type PasswordService struct {
	bcryptCost int
}

func NewPasswordService() *PasswordService {
	return &PasswordService{
		bcryptCost: 12, // Good balance between security and performance
	}
}

// HashPassword hashes a password using bcrypt
func (s *PasswordService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), s.bcryptCost)
	return string(bytes), err
}

// VerifyPassword verifies a password against a hash
func (s *PasswordService) VerifyPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// GenerateToken generates a secure random token
func (s *PasswordService) GenerateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// IsTokenExpired checks if a token has expired
func (s *PasswordService) IsTokenExpired(expires *time.Time) bool {
	if expires == nil {
		return true
	}
	return time.Now().After(*expires)
}

