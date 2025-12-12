package models

import (
	"encoding/json"
	"time"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "PENDING"
	PaymentStatusCompleted PaymentStatus = "COMPLETED"
	PaymentStatusFailed    PaymentStatus = "FAILED"
	PaymentStatusRefunded  PaymentStatus = "REFUNDED"
)

type Payment struct {
	ID            int64         `json:"id" db:"id"`
	UserID        int64         `json:"user_id" db:"user_id"`
	PackageID     int64         `json:"package_id" db:"package_id"`
	Amount        float64       `json:"amount" db:"amount"`
	PaymentMethod *string       `json:"payment_method,omitempty" db:"payment_method"`
	PaymentStatus PaymentStatus `json:"payment_status" db:"payment_status"`
	TransactionID *string       `json:"transaction_id,omitempty" db:"transaction_id"`
	Metadata      *string       `json:"metadata,omitempty" db:"metadata"` // JSONB stored as string
	CreatedAt     time.Time     `json:"created_at" db:"created_at"`
}

// PaymentWithPackage represents a payment with its package details
type PaymentWithPackage struct {
	Payment
	Package *Package `json:"package,omitempty"`
}

// PaymentCreate represents the data needed to create a payment record
type PaymentCreate struct {
	UserID        int64                  `json:"user_id" validate:"required,gt=0"`
	PackageID     int64                  `json:"package_id" validate:"required,gt=0"`
	Amount        float64                `json:"amount" validate:"required,gte=0"`
	PaymentMethod *string                `json:"payment_method,omitempty"`
	PaymentStatus PaymentStatus          `json:"payment_status" validate:"required,oneof=PENDING COMPLETED FAILED REFUNDED"`
	TransactionID *string                `json:"transaction_id,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// GetMetadata parses the JSONB metadata field
func (p *Payment) GetMetadata() (map[string]interface{}, error) {
	if p.Metadata == nil || *p.Metadata == "" {
		return nil, nil
	}
	var metadata map[string]interface{}
	if err := json.Unmarshal([]byte(*p.Metadata), &metadata); err != nil {
		return nil, err
	}
	return metadata, nil
}

