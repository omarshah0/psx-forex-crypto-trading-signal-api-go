package models

import (
	"time"
)

type Subscription struct {
	ID           int64     `json:"id" db:"id"`
	UserID       int64     `json:"user_id" db:"user_id"`
	PackageID    int64     `json:"package_id" db:"package_id"`
	PricePaid    float64   `json:"price_paid" db:"price_paid"`
	SubscribedAt time.Time `json:"subscribed_at" db:"subscribed_at"`
	ExpiresAt    time.Time `json:"expires_at" db:"expires_at"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// SubscriptionWithPackage represents a subscription with its package details
type SubscriptionWithPackage struct {
	Subscription
	Package *Package `json:"package,omitempty"`
}

// SubscribeRequest represents a request to subscribe to packages
type SubscribeRequest struct {
	PackageIDs []int64 `json:"package_ids" validate:"required,min=1,dive,gt=0"`
}

// SubscribeResponse represents the response after subscribing
type SubscribeResponse struct {
	Subscriptions []SubscriptionWithPackage `json:"subscriptions"`
	TotalAmount   float64                   `json:"total_amount"`
	Message       string                    `json:"message"`
}

// CheckAccessRequest represents a request to check access
type CheckAccessRequest struct {
	AssetClass   AssetClass   `json:"asset_class" validate:"required,oneof=FOREX CRYPTO PSX"`
	DurationType DurationType `json:"duration_type" validate:"required,oneof=SHORT_TERM LONG_TERM"`
}

// CheckAccessResponse represents the response for access check
type CheckAccessResponse struct {
	HasAccess  bool       `json:"has_access"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	PackageID  *int64     `json:"package_id,omitempty"`
	PricePaid  *float64   `json:"price_paid,omitempty"`
}

