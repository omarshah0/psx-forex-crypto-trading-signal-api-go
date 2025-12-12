package repositories

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/omarshah0/rest-api-with-social-auth/internal/models"
)

type PackageRepository struct {
	db *sql.DB
}

func NewPackageRepository(db *sql.DB) *PackageRepository {
	return &PackageRepository{db: db}
}

// Create creates a new package
func (r *PackageRepository) Create(pkg *models.PackageCreate) (*models.Package, error) {
	query := `
		INSERT INTO packages (name, asset_class, duration_type, billing_cycle, duration_days, price, description)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, name, asset_class, duration_type, billing_cycle, duration_days, price, description, is_active, created_at, updated_at
	`

	var newPackage models.Package
	err := r.db.QueryRow(
		query,
		pkg.Name,
		pkg.AssetClass,
		pkg.DurationType,
		pkg.BillingCycle,
		pkg.DurationDays,
		pkg.Price,
		pkg.Description,
	).Scan(
		&newPackage.ID,
		&newPackage.Name,
		&newPackage.AssetClass,
		&newPackage.DurationType,
		&newPackage.BillingCycle,
		&newPackage.DurationDays,
		&newPackage.Price,
		&newPackage.Description,
		&newPackage.IsActive,
		&newPackage.CreatedAt,
		&newPackage.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create package: %w", err)
	}

	return &newPackage, nil
}

// GetByID retrieves a package by ID
func (r *PackageRepository) GetByID(id int64) (*models.Package, error) {
	query := `
		SELECT id, name, asset_class, duration_type, billing_cycle, duration_days, price, description, is_active, created_at, updated_at
		FROM packages
		WHERE id = $1
	`

	var pkg models.Package
	err := r.db.QueryRow(query, id).Scan(
		&pkg.ID,
		&pkg.Name,
		&pkg.AssetClass,
		&pkg.DurationType,
		&pkg.BillingCycle,
		&pkg.DurationDays,
		&pkg.Price,
		&pkg.Description,
		&pkg.IsActive,
		&pkg.CreatedAt,
		&pkg.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get package: %w", err)
	}

	return &pkg, nil
}

// GetAll retrieves all packages
func (r *PackageRepository) GetAll(activeOnly bool, limit, offset int) ([]models.Package, error) {
	query := `
		SELECT id, name, asset_class, duration_type, billing_cycle, duration_days, price, description, is_active, created_at, updated_at
		FROM packages
	`

	var args []interface{}
	argPosition := 1

	if activeOnly {
		query += fmt.Sprintf(" WHERE is_active = $%d", argPosition)
		args = append(args, true)
		argPosition++
	}

	query += " ORDER BY asset_class, duration_type, billing_cycle"

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argPosition)
		args = append(args, limit)
		argPosition++
	}

	if offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argPosition)
		args = append(args, offset)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get packages: %w", err)
	}
	defer rows.Close()

	var packages []models.Package
	for rows.Next() {
		var pkg models.Package
		err := rows.Scan(
			&pkg.ID,
			&pkg.Name,
			&pkg.AssetClass,
			&pkg.DurationType,
			&pkg.BillingCycle,
			&pkg.DurationDays,
			&pkg.Price,
			&pkg.Description,
			&pkg.IsActive,
			&pkg.CreatedAt,
			&pkg.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan package: %w", err)
		}
		packages = append(packages, pkg)
	}

	return packages, nil
}

// GetByIDs retrieves multiple packages by their IDs
func (r *PackageRepository) GetByIDs(ids []int64) ([]models.Package, error) {
	if len(ids) == 0 {
		return []models.Package{}, nil
	}

	// Build placeholders for IN clause
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`
		SELECT id, name, asset_class, duration_type, billing_cycle, duration_days, price, description, is_active, created_at, updated_at
		FROM packages
		WHERE id IN (%s)
	`, strings.Join(placeholders, ", "))

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get packages by IDs: %w", err)
	}
	defer rows.Close()

	var packages []models.Package
	for rows.Next() {
		var pkg models.Package
		err := rows.Scan(
			&pkg.ID,
			&pkg.Name,
			&pkg.AssetClass,
			&pkg.DurationType,
			&pkg.BillingCycle,
			&pkg.DurationDays,
			&pkg.Price,
			&pkg.Description,
			&pkg.IsActive,
			&pkg.CreatedAt,
			&pkg.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan package: %w", err)
		}
		packages = append(packages, pkg)
	}

	return packages, nil
}

// Update updates a package
func (r *PackageRepository) Update(id int64, update *models.PackageUpdate) (*models.Package, error) {
	var setClauses []string
	var args []interface{}
	argPosition := 1

	if update.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argPosition))
		args = append(args, *update.Name)
		argPosition++
	}
	if update.AssetClass != nil {
		setClauses = append(setClauses, fmt.Sprintf("asset_class = $%d", argPosition))
		args = append(args, *update.AssetClass)
		argPosition++
	}
	if update.DurationType != nil {
		setClauses = append(setClauses, fmt.Sprintf("duration_type = $%d", argPosition))
		args = append(args, *update.DurationType)
		argPosition++
	}
	if update.BillingCycle != nil {
		setClauses = append(setClauses, fmt.Sprintf("billing_cycle = $%d", argPosition))
		args = append(args, *update.BillingCycle)
		argPosition++
	}
	if update.DurationDays != nil {
		setClauses = append(setClauses, fmt.Sprintf("duration_days = $%d", argPosition))
		args = append(args, *update.DurationDays)
		argPosition++
	}
	if update.Price != nil {
		setClauses = append(setClauses, fmt.Sprintf("price = $%d", argPosition))
		args = append(args, *update.Price)
		argPosition++
	}
	if update.Description != nil {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", argPosition))
		args = append(args, *update.Description)
		argPosition++
	}
	if update.IsActive != nil {
		setClauses = append(setClauses, fmt.Sprintf("is_active = $%d", argPosition))
		args = append(args, *update.IsActive)
		argPosition++
	}

	if len(setClauses) == 0 {
		return r.GetByID(id)
	}

	setClauses = append(setClauses, "updated_at = CURRENT_TIMESTAMP")
	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE packages
		SET %s
		WHERE id = $%d
		RETURNING id, name, asset_class, duration_type, billing_cycle, duration_days, price, description, is_active, created_at, updated_at
	`, strings.Join(setClauses, ", "), argPosition)

	var pkg models.Package
	err := r.db.QueryRow(query, args...).Scan(
		&pkg.ID,
		&pkg.Name,
		&pkg.AssetClass,
		&pkg.DurationType,
		&pkg.BillingCycle,
		&pkg.DurationDays,
		&pkg.Price,
		&pkg.Description,
		&pkg.IsActive,
		&pkg.CreatedAt,
		&pkg.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("package not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update package: %w", err)
	}

	return &pkg, nil
}

// Delete deletes a package
func (r *PackageRepository) Delete(id int64) error {
	query := `DELETE FROM packages WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete package: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("package not found")
	}

	return nil
}

// Count returns the total count of packages
func (r *PackageRepository) Count(activeOnly bool) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM packages`
	if activeOnly {
		query += ` WHERE is_active = true`
	}
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count packages: %w", err)
	}
	return count, nil
}

