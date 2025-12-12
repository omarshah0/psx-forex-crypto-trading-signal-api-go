package repositories

import (
	"database/sql"
	"fmt"

	"github.com/omarshah0/rest-api-with-social-auth/internal/models"
)

type OAuthProviderRepository struct {
	db *sql.DB
}

func NewOAuthProviderRepository(db *sql.DB) *OAuthProviderRepository {
	return &OAuthProviderRepository{db: db}
}

// Create creates a new OAuth provider link
func (r *OAuthProviderRepository) Create(userID int64, provider models.OAuthProviderType, providerUserID, email string) (*models.OAuthProvider, error) {
	query := `
		INSERT INTO oauth_providers (user_id, provider, provider_user_id, email)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, provider, provider_user_id, email, created_at
	`

	var oauthProvider models.OAuthProvider
	err := r.db.QueryRow(query, userID, provider, providerUserID, email).Scan(
		&oauthProvider.ID,
		&oauthProvider.UserID,
		&oauthProvider.Provider,
		&oauthProvider.ProviderUserID,
		&oauthProvider.Email,
		&oauthProvider.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create oauth provider: %w", err)
	}

	return &oauthProvider, nil
}

// GetByProviderAndUserID retrieves an OAuth provider by provider and provider user ID
func (r *OAuthProviderRepository) GetByProviderAndUserID(provider models.OAuthProviderType, providerUserID string) (*models.OAuthProvider, error) {
	query := `
		SELECT id, user_id, provider, provider_user_id, email, created_at
		FROM oauth_providers
		WHERE provider = $1 AND provider_user_id = $2
	`

	var oauthProvider models.OAuthProvider
	err := r.db.QueryRow(query, provider, providerUserID).Scan(
		&oauthProvider.ID,
		&oauthProvider.UserID,
		&oauthProvider.Provider,
		&oauthProvider.ProviderUserID,
		&oauthProvider.Email,
		&oauthProvider.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get oauth provider: %w", err)
	}

	return &oauthProvider, nil
}

// GetByUserID retrieves all OAuth providers for a user
func (r *OAuthProviderRepository) GetByUserID(userID int64) ([]models.OAuthProvider, error) {
	query := `
		SELECT id, user_id, provider, provider_user_id, email, created_at
		FROM oauth_providers
		WHERE user_id = $1
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get oauth providers: %w", err)
	}
	defer rows.Close()

	var providers []models.OAuthProvider
	for rows.Next() {
		var provider models.OAuthProvider
		err := rows.Scan(
			&provider.ID,
			&provider.UserID,
			&provider.Provider,
			&provider.ProviderUserID,
			&provider.Email,
			&provider.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan oauth provider: %w", err)
		}
		providers = append(providers, provider)
	}

	return providers, nil
}

// Delete deletes an OAuth provider link
func (r *OAuthProviderRepository) Delete(id int64) error {
	query := `DELETE FROM oauth_providers WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete oauth provider: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("oauth provider not found")
	}

	return nil
}
