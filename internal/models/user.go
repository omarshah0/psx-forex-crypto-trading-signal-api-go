package models

import (
	"time"
)

type User struct {
	ID                       int64      `json:"id" db:"id"`
	Email                    string     `json:"email" db:"email" validate:"required,email"`
	Name                     string     `json:"name" db:"name" validate:"required"`
	ProfilePicture           *string    `json:"profile_picture,omitempty" db:"profile_picture"`
	HashedPassword           *string    `json:"-" db:"hashed_password"` // Never expose in JSON
	EmailVerified            bool       `json:"email_verified" db:"email_verified"`
	VerificationToken        *string    `json:"-" db:"verification_token"`
	VerificationTokenExpires *time.Time `json:"-" db:"verification_token_expires"`
	ResetToken               *string    `json:"-" db:"reset_token"`
	ResetTokenExpires        *time.Time `json:"-" db:"reset_token_expires"`
	Blocked                  bool       `json:"blocked" db:"blocked"`
	CreatedAt                time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at" db:"updated_at"`
}

// UserCreate represents the data needed to create a new user
type UserCreate struct {
	Email          string  `json:"email" validate:"required,email"`
	Name           string  `json:"name" validate:"required"`
	ProfilePicture *string `json:"profile_picture,omitempty"`
}

// UserRegister represents registration with email/password
type UserRegister struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=2"`
	Password string `json:"password" validate:"required,min=8"`
}

// UserLogin represents login credentials
type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// PasswordChange represents password change request
type PasswordChange struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// ForgotPassword represents forgot password request
type ForgotPassword struct {
	Email string `json:"email" validate:"required,email"`
}

// ResetPassword represents password reset with token
type ResetPassword struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// ResendVerification represents resend verification email request
type ResendVerification struct {
	Email string `json:"email" validate:"required,email"`
}

// LinkedAccount represents a simplified OAuth provider link
type LinkedAccount struct {
	Provider string    `json:"provider"`
	LinkedAt time.Time `json:"linked_at"`
}

// LinkedAccounts represents all linked authentication methods for a user
type LinkedAccounts struct {
	Email       string         `json:"email"`
	PasswordSet bool           `json:"password_set"`
	Google      *LinkedAccount `json:"google"`
	Facebook    *LinkedAccount `json:"facebook"`
}

// ProfileResponse represents the complete profile data including linked accounts
type ProfileResponse struct {
	User           *User          `json:"user"`
	LinkedAccounts LinkedAccounts `json:"linked_accounts"`
}
