package repositories

import (
	"database/sql"
	"fmt"

	"github.com/omarshah0/rest-api-with-social-auth/internal/models"
)

type AdminRepository struct {
	db *sql.DB
}

func NewAdminRepository(db *sql.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

// GetByUserID retrieves an admin by user ID
func (r *AdminRepository) GetByUserID(userID int64) (*models.Admin, error) {
	query := `
		SELECT id, user_id, created_at, updated_at
		FROM admins
		WHERE user_id = $1
	`

	var admin models.Admin
	err := r.db.QueryRow(query, userID).Scan(
		&admin.ID,
		&admin.UserID,
		&admin.CreatedAt,
		&admin.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get admin: %w", err)
	}

	return &admin, nil
}

// IsAdmin checks if a user is an admin
func (r *AdminRepository) IsAdmin(userID int64) (bool, error) {
	admin, err := r.GetByUserID(userID)
	if err != nil {
		return false, err
	}
	return admin != nil, nil
}
