package models

import (
	"time"
)

type OAuthProviderType string

const (
	ProviderGoogle   OAuthProviderType = "google"
	ProviderFacebook OAuthProviderType = "facebook"
)

type OAuthProvider struct {
	ID             int64             `json:"id" db:"id"`
	UserID         int64             `json:"user_id" db:"user_id"`
	Provider       OAuthProviderType `json:"provider" db:"provider" validate:"required,oneof=google facebook"`
	ProviderUserID string            `json:"provider_user_id" db:"provider_user_id" validate:"required"`
	Email          string            `json:"email" db:"email" validate:"required,email"`
	CreatedAt      time.Time         `json:"created_at" db:"created_at"`
}

