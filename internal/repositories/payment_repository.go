package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/omarshah0/rest-api-with-social-auth/internal/models"
)

type PaymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

// Create creates a new payment record
func (r *PaymentRepository) Create(payment *models.PaymentCreate) (*models.Payment, error) {
	// Convert metadata map to JSON string
	var metadataJSON *string
	if payment.Metadata != nil {
		jsonBytes, err := json.Marshal(payment.Metadata)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal metadata: %w", err)
		}
		jsonStr := string(jsonBytes)
		metadataJSON = &jsonStr
	}

	query := `
		INSERT INTO payment_history (user_id, package_id, amount, payment_method, payment_status, transaction_id, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, user_id, package_id, amount, payment_method, payment_status, transaction_id, metadata, created_at
	`

	var newPayment models.Payment
	err := r.db.QueryRow(
		query,
		payment.UserID,
		payment.PackageID,
		payment.Amount,
		payment.PaymentMethod,
		payment.PaymentStatus,
		payment.TransactionID,
		metadataJSON,
	).Scan(
		&newPayment.ID,
		&newPayment.UserID,
		&newPayment.PackageID,
		&newPayment.Amount,
		&newPayment.PaymentMethod,
		&newPayment.PaymentStatus,
		&newPayment.TransactionID,
		&newPayment.Metadata,
		&newPayment.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	return &newPayment, nil
}

// GetByID retrieves a payment by ID
func (r *PaymentRepository) GetByID(id int64) (*models.Payment, error) {
	query := `
		SELECT id, user_id, package_id, amount, payment_method, payment_status, transaction_id, metadata, created_at
		FROM payment_history
		WHERE id = $1
	`

	var payment models.Payment
	err := r.db.QueryRow(query, id).Scan(
		&payment.ID,
		&payment.UserID,
		&payment.PackageID,
		&payment.Amount,
		&payment.PaymentMethod,
		&payment.PaymentStatus,
		&payment.TransactionID,
		&payment.Metadata,
		&payment.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	return &payment, nil
}

// GetByUserID retrieves all payments for a user
func (r *PaymentRepository) GetByUserID(userID int64, limit, offset int) ([]models.Payment, error) {
	query := `
		SELECT id, user_id, package_id, amount, payment_method, payment_status, transaction_id, metadata, created_at
		FROM payment_history
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get payments: %w", err)
	}
	defer rows.Close()

	var payments []models.Payment
	for rows.Next() {
		var payment models.Payment
		err := rows.Scan(
			&payment.ID,
			&payment.UserID,
			&payment.PackageID,
			&payment.Amount,
			&payment.PaymentMethod,
			&payment.PaymentStatus,
			&payment.TransactionID,
			&payment.Metadata,
			&payment.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan payment: %w", err)
		}
		payments = append(payments, payment)
	}

	return payments, nil
}

// GetByTransactionID retrieves a payment by transaction ID
func (r *PaymentRepository) GetByTransactionID(transactionID string) (*models.Payment, error) {
	query := `
		SELECT id, user_id, package_id, amount, payment_method, payment_status, transaction_id, metadata, created_at
		FROM payment_history
		WHERE transaction_id = $1
	`

	var payment models.Payment
	err := r.db.QueryRow(query, transactionID).Scan(
		&payment.ID,
		&payment.UserID,
		&payment.PackageID,
		&payment.Amount,
		&payment.PaymentMethod,
		&payment.PaymentStatus,
		&payment.TransactionID,
		&payment.Metadata,
		&payment.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get payment by transaction ID: %w", err)
	}

	return &payment, nil
}

// CountByUserID returns the total count of payments for a user
func (r *PaymentRepository) CountByUserID(userID int64) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM payment_history WHERE user_id = $1`
	err := r.db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count payments: %w", err)
	}
	return count, nil
}

// UpdateStatus updates the payment status
func (r *PaymentRepository) UpdateStatus(id int64, status models.PaymentStatus) error {
	query := `UPDATE payment_history SET payment_status = $1 WHERE id = $2`

	result, err := r.db.Exec(query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("payment not found")
	}

	return nil
}

