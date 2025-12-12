package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/omarshah0/rest-api-with-social-auth/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user (for OAuth - email is pre-verified)
func (r *UserRepository) Create(user *models.UserCreate) (*models.User, error) {
	query := `
		INSERT INTO users (email, name, profile_picture, blocked, email_verified)
		VALUES ($1, $2, $3, false, true)
		RETURNING id, email, name, profile_picture, blocked, email_verified, created_at, updated_at
	`

	var newUser models.User
	err := r.db.QueryRow(query, user.Email, user.Name, user.ProfilePicture).Scan(
		&newUser.ID,
		&newUser.Email,
		&newUser.Name,
		&newUser.ProfilePicture,
		&newUser.Blocked,
		&newUser.EmailVerified,
		&newUser.CreatedAt,
		&newUser.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &newUser, nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id int64) (*models.User, error) {
	query := `
		SELECT id, email, name, profile_picture, blocked, email_verified, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user models.User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.ProfilePicture,
		&user.Blocked,
		&user.EmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, email, name, profile_picture, blocked, email_verified, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user models.User
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.ProfilePicture,
		&user.Blocked,
		&user.EmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

// Update updates a user
func (r *UserRepository) Update(id int64, user *models.User) error {
	query := `
		UPDATE users
		SET email = $1, name = $2, profile_picture = $3, blocked = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
	`

	result, err := r.db.Exec(query, user.Email, user.Name, user.ProfilePicture, user.Blocked, id)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// Delete deletes a user
func (r *UserRepository) Delete(id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// CreateWithPassword creates a new user with password
func (r *UserRepository) CreateWithPassword(user *models.UserRegister, hashedPassword string) (*models.User, error) {
	query := `
		INSERT INTO users (email, name, hashed_password, email_verified, blocked)
		VALUES ($1, $2, $3, false, false)
		RETURNING id, email, name, profile_picture, hashed_password, email_verified, blocked, created_at, updated_at
	`

	var newUser models.User
	err := r.db.QueryRow(query, user.Email, user.Name, hashedPassword).Scan(
		&newUser.ID,
		&newUser.Email,
		&newUser.Name,
		&newUser.ProfilePicture,
		&newUser.HashedPassword,
		&newUser.EmailVerified,
		&newUser.Blocked,
		&newUser.CreatedAt,
		&newUser.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &newUser, nil
}

// GetByEmailWithPassword retrieves a user by email including password
func (r *UserRepository) GetByEmailWithPassword(email string) (*models.User, error) {
	query := `
		SELECT id, email, name, profile_picture, hashed_password, email_verified, verification_token, 
		       verification_token_expires, reset_token, reset_token_expires, blocked, 
		       created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user models.User
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.ProfilePicture,
		&user.HashedPassword,
		&user.EmailVerified,
		&user.VerificationToken,
		&user.VerificationTokenExpires,
		&user.ResetToken,
		&user.ResetTokenExpires,
		&user.Blocked,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

// SetVerificationToken sets the email verification token
func (r *UserRepository) SetVerificationToken(userID int64, token string, expires time.Time) error {
	query := `
		UPDATE users 
		SET verification_token = $1, verification_token_expires = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
	`

	_, err := r.db.Exec(query, token, expires, userID)
	if err != nil {
		return fmt.Errorf("failed to set verification token: %w", err)
	}

	return nil
}

// VerifyEmail marks email as verified and clears verification token
func (r *UserRepository) VerifyEmail(token string) error {
	query := `
		UPDATE users 
		SET email_verified = true, 
		    verification_token = NULL, 
		    verification_token_expires = NULL,
		    updated_at = CURRENT_TIMESTAMP
		WHERE verification_token = $1 
		  AND verification_token_expires > CURRENT_TIMESTAMP
	`

	result, err := r.db.Exec(query, token)
	if err != nil {
		return fmt.Errorf("failed to verify email: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("invalid or expired verification token")
	}

	return nil
}

// GetByVerificationToken retrieves a user by verification token
func (r *UserRepository) GetByVerificationToken(token string) (*models.User, error) {
	query := `
		SELECT id, email, name, profile_picture, email_verified, verification_token_expires, blocked, created_at, updated_at
		FROM users
		WHERE verification_token = $1
	`

	var user models.User
	err := r.db.QueryRow(query, token).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.ProfilePicture,
		&user.EmailVerified,
		&user.VerificationTokenExpires,
		&user.Blocked,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by verification token: %w", err)
	}

	return &user, nil
}

// SetResetToken sets the password reset token
func (r *UserRepository) SetResetToken(userID int64, token string, expires time.Time) error {
	query := `
		UPDATE users 
		SET reset_token = $1, reset_token_expires = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
	`

	_, err := r.db.Exec(query, token, expires, userID)
	if err != nil {
		return fmt.Errorf("failed to set reset token: %w", err)
	}

	return nil
}

// GetByResetToken retrieves a user by reset token
func (r *UserRepository) GetByResetToken(token string) (*models.User, error) {
	query := `
		SELECT id, email, name, profile_picture, hashed_password, reset_token_expires, blocked, created_at, updated_at
		FROM users
		WHERE reset_token = $1
	`

	var user models.User
	err := r.db.QueryRow(query, token).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.ProfilePicture,
		&user.HashedPassword,
		&user.ResetTokenExpires,
		&user.Blocked,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by reset token: %w", err)
	}

	return &user, nil
}

// UpdatePassword updates the user's password and clears reset token
func (r *UserRepository) UpdatePassword(userID int64, hashedPassword string) error {
	query := `
		UPDATE users 
		SET hashed_password = $1, 
		    reset_token = NULL, 
		    reset_token_expires = NULL,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	result, err := r.db.Exec(query, hashedPassword, userID)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// MarkEmailVerified marks a user's email as verified by user ID
func (r *UserRepository) MarkEmailVerified(userID int64) error {
	query := `
		UPDATE users 
		SET email_verified = true, 
		    verification_token = NULL, 
		    verification_token_expires = NULL,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	_, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to mark email as verified: %w", err)
	}

	return nil
}
