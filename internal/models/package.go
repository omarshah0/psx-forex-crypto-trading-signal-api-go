package models

import (
	"time"
)

type BillingCycle string

const (
	BillingCycleMonthly    BillingCycle = "MONTHLY"
	BillingCycleSixMonths  BillingCycle = "SIX_MONTHS"
	BillingCycleYearly     BillingCycle = "YEARLY"
)

type Package struct {
	ID           int64        `json:"id" db:"id"`
	Name         string       `json:"name" db:"name" validate:"required"`
	AssetClass   AssetClass   `json:"asset_class" db:"asset_class" validate:"required,oneof=FOREX CRYPTO PSX"`
	DurationType DurationType `json:"duration_type" db:"duration_type" validate:"required,oneof=SHORT_TERM LONG_TERM"`
	BillingCycle BillingCycle `json:"billing_cycle" db:"billing_cycle" validate:"required,oneof=MONTHLY SIX_MONTHS YEARLY"`
	DurationDays int          `json:"duration_days" db:"duration_days" validate:"required,gt=0"`
	Price        float64      `json:"price" db:"price" validate:"required,gte=0"`
	Description  *string      `json:"description,omitempty" db:"description"`
	IsActive     bool         `json:"is_active" db:"is_active"`
	CreatedAt    time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at" db:"updated_at"`
}

// PackageCreate represents the data needed to create a new package
type PackageCreate struct {
	Name         string       `json:"name" validate:"required"`
	AssetClass   AssetClass   `json:"asset_class" validate:"required,oneof=FOREX CRYPTO PSX"`
	DurationType DurationType `json:"duration_type" validate:"required,oneof=SHORT_TERM LONG_TERM"`
	BillingCycle BillingCycle `json:"billing_cycle" validate:"required,oneof=MONTHLY SIX_MONTHS YEARLY"`
	DurationDays int          `json:"duration_days" validate:"required,gt=0"`
	Price        float64      `json:"price" validate:"required,gte=0"`
	Description  *string      `json:"description,omitempty"`
}

// PackageUpdate represents the data needed to update a package
type PackageUpdate struct {
	Name         *string       `json:"name"`
	AssetClass   *AssetClass   `json:"asset_class" validate:"omitempty,oneof=FOREX CRYPTO PSX"`
	DurationType *DurationType `json:"duration_type" validate:"omitempty,oneof=SHORT_TERM LONG_TERM"`
	BillingCycle *BillingCycle `json:"billing_cycle" validate:"omitempty,oneof=MONTHLY SIX_MONTHS YEARLY"`
	DurationDays *int          `json:"duration_days" validate:"omitempty,gt=0"`
	Price        *float64      `json:"price" validate:"omitempty,gte=0"`
	Description  *string       `json:"description,omitempty"`
	IsActive     *bool         `json:"is_active"`
}

