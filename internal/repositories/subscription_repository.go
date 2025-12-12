package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/omarshah0/rest-api-with-social-auth/internal/models"
)

type SubscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

// Create creates a new subscription
func (r *SubscriptionRepository) Create(userID, packageID int64, pricePaid float64, expiresAt time.Time) (*models.Subscription, error) {
	query := `
		INSERT INTO user_subscriptions (user_id, package_id, price_paid, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, package_id, price_paid, subscribed_at, expires_at, is_active, created_at, updated_at
	`

	var subscription models.Subscription
	err := r.db.QueryRow(query, userID, packageID, pricePaid, expiresAt).Scan(
		&subscription.ID,
		&subscription.UserID,
		&subscription.PackageID,
		&subscription.PricePaid,
		&subscription.SubscribedAt,
		&subscription.ExpiresAt,
		&subscription.IsActive,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	return &subscription, nil
}

// GetByID retrieves a subscription by ID
func (r *SubscriptionRepository) GetByID(id int64) (*models.Subscription, error) {
	query := `
		SELECT id, user_id, package_id, price_paid, subscribed_at, expires_at, is_active, created_at, updated_at
		FROM user_subscriptions
		WHERE id = $1
	`

	var subscription models.Subscription
	err := r.db.QueryRow(query, id).Scan(
		&subscription.ID,
		&subscription.UserID,
		&subscription.PackageID,
		&subscription.PricePaid,
		&subscription.SubscribedAt,
		&subscription.ExpiresAt,
		&subscription.IsActive,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	return &subscription, nil
}

// GetActiveByUserID retrieves all active subscriptions for a user
func (r *SubscriptionRepository) GetActiveByUserID(userID int64) ([]models.Subscription, error) {
	query := `
		SELECT id, user_id, package_id, price_paid, subscribed_at, expires_at, is_active, created_at, updated_at
		FROM user_subscriptions
		WHERE user_id = $1 AND is_active = true AND expires_at > CURRENT_TIMESTAMP
		ORDER BY expires_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active subscriptions: %w", err)
	}
	defer rows.Close()

	var subscriptions []models.Subscription
	for rows.Next() {
		var subscription models.Subscription
		err := rows.Scan(
			&subscription.ID,
			&subscription.UserID,
			&subscription.PackageID,
			&subscription.PricePaid,
			&subscription.SubscribedAt,
			&subscription.ExpiresAt,
			&subscription.IsActive,
			&subscription.CreatedAt,
			&subscription.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}
		subscriptions = append(subscriptions, subscription)
	}

	return subscriptions, nil
}

// GetAllByUserID retrieves all subscriptions for a user (active and expired)
func (r *SubscriptionRepository) GetAllByUserID(userID int64, limit, offset int) ([]models.Subscription, error) {
	query := `
		SELECT id, user_id, package_id, price_paid, subscribed_at, expires_at, is_active, created_at, updated_at
		FROM user_subscriptions
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscriptions: %w", err)
	}
	defer rows.Close()

	var subscriptions []models.Subscription
	for rows.Next() {
		var subscription models.Subscription
		err := rows.Scan(
			&subscription.ID,
			&subscription.UserID,
			&subscription.PackageID,
			&subscription.PricePaid,
			&subscription.SubscribedAt,
			&subscription.ExpiresAt,
			&subscription.IsActive,
			&subscription.CreatedAt,
			&subscription.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}
		subscriptions = append(subscriptions, subscription)
	}

	return subscriptions, nil
}

// CheckAccess checks if user has active subscription for specific asset class and duration type
func (r *SubscriptionRepository) CheckAccess(userID int64, assetClass models.AssetClass, durationType models.DurationType) (*models.Subscription, error) {
	query := `
		SELECT us.id, us.user_id, us.package_id, us.price_paid, us.subscribed_at, us.expires_at, us.is_active, us.created_at, us.updated_at
		FROM user_subscriptions us
		JOIN packages p ON us.package_id = p.id
		WHERE us.user_id = $1
		AND us.is_active = true
		AND us.expires_at > CURRENT_TIMESTAMP
		AND p.asset_class = $2
		AND p.duration_type = $3
		ORDER BY us.expires_at DESC
		LIMIT 1
	`

	var subscription models.Subscription
	err := r.db.QueryRow(query, userID, assetClass, durationType).Scan(
		&subscription.ID,
		&subscription.UserID,
		&subscription.PackageID,
		&subscription.PricePaid,
		&subscription.SubscribedAt,
		&subscription.ExpiresAt,
		&subscription.IsActive,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to check access: %w", err)
	}

	return &subscription, nil
}

// DeactivateExpired deactivates all expired subscriptions
func (r *SubscriptionRepository) DeactivateExpired() (int64, error) {
	query := `
		UPDATE user_subscriptions
		SET is_active = false, updated_at = CURRENT_TIMESTAMP
		WHERE is_active = true AND expires_at <= CURRENT_TIMESTAMP
	`

	result, err := r.db.Exec(query)
	if err != nil {
		return 0, fmt.Errorf("failed to deactivate expired subscriptions: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rows, nil
}

// CountByUserID returns the total count of subscriptions for a user
func (r *SubscriptionRepository) CountByUserID(userID int64) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM user_subscriptions WHERE user_id = $1`
	err := r.db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count subscriptions: %w", err)
	}
	return count, nil
}

