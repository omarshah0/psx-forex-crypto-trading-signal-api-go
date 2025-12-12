package repositories

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/omarshah0/rest-api-with-social-auth/internal/models"
)

type TradingSignalRepository struct {
	db *sql.DB
}

func NewTradingSignalRepository(db *sql.DB) *TradingSignalRepository {
	return &TradingSignalRepository{db: db}
}

// Create creates a new trading signal
func (r *TradingSignalRepository) Create(signal *models.TradingSignalCreate, createdBy int64) (*models.TradingSignal, error) {
	query := `
		INSERT INTO trading_signals (symbol, asset_class, duration_type, stop_loss_price, entry_price, take_profit_price, type, result, return, free_for_all, comments, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, symbol, asset_class, duration_type, stop_loss_price, entry_price, take_profit_price, type, result, return, free_for_all, comments, created_by, created_at, updated_at
	`

	var newSignal models.TradingSignal
	err := r.db.QueryRow(
		query,
		signal.Symbol,
		signal.AssetClass,
		signal.DurationType,
		signal.StopLossPrice,
		signal.EntryPrice,
		signal.TakeProfitPrice,
		signal.Type,
		signal.Result,
		signal.Return,
		signal.FreeForAll,
		signal.Comments,
		createdBy,
	).Scan(
		&newSignal.ID,
		&newSignal.Symbol,
		&newSignal.AssetClass,
		&newSignal.DurationType,
		&newSignal.StopLossPrice,
		&newSignal.EntryPrice,
		&newSignal.TakeProfitPrice,
		&newSignal.Type,
		&newSignal.Result,
		&newSignal.Return,
		&newSignal.FreeForAll,
		&newSignal.Comments,
		&newSignal.CreatedBy,
		&newSignal.CreatedAt,
		&newSignal.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create trading signal: %w", err)
	}

	return &newSignal, nil
}

// GetByID retrieves a trading signal by ID
func (r *TradingSignalRepository) GetByID(id int64) (*models.TradingSignal, error) {
	query := `
		SELECT id, symbol, asset_class, duration_type, stop_loss_price, entry_price, take_profit_price, type, result, return, free_for_all, comments, created_by, created_at, updated_at
		FROM trading_signals
		WHERE id = $1
	`

	var signal models.TradingSignal
	err := r.db.QueryRow(query, id).Scan(
		&signal.ID,
		&signal.Symbol,
		&signal.AssetClass,
		&signal.DurationType,
		&signal.StopLossPrice,
		&signal.EntryPrice,
		&signal.TakeProfitPrice,
		&signal.Type,
		&signal.Result,
		&signal.Return,
		&signal.FreeForAll,
		&signal.Comments,
		&signal.CreatedBy,
		&signal.CreatedAt,
		&signal.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get trading signal: %w", err)
	}

	return &signal, nil
}

// GetAll retrieves all trading signals with optional filtering
func (r *TradingSignalRepository) GetAll(limit, offset int) ([]models.TradingSignal, error) {
	query := `
		SELECT id, symbol, asset_class, duration_type, stop_loss_price, entry_price, take_profit_price, type, result, return, free_for_all, comments, created_by, created_at, updated_at
		FROM trading_signals
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get trading signals: %w", err)
	}
	defer rows.Close()

	var signals []models.TradingSignal
	for rows.Next() {
		var signal models.TradingSignal
		err := rows.Scan(
			&signal.ID,
			&signal.Symbol,
			&signal.AssetClass,
			&signal.DurationType,
			&signal.StopLossPrice,
			&signal.EntryPrice,
			&signal.TakeProfitPrice,
			&signal.Type,
			&signal.Result,
			&signal.Return,
			&signal.FreeForAll,
			&signal.Comments,
			&signal.CreatedBy,
			&signal.CreatedAt,
			&signal.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan trading signal: %w", err)
		}
		signals = append(signals, signal)
	}

	return signals, nil
}

// Update updates a trading signal
func (r *TradingSignalRepository) Update(id int64, update *models.TradingSignalUpdate) (*models.TradingSignal, error) {
	// Build dynamic update query
	var setClauses []string
	var args []interface{}
	argPosition := 1

	if update.Symbol != nil {
		setClauses = append(setClauses, fmt.Sprintf("symbol = $%d", argPosition))
		args = append(args, *update.Symbol)
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
	if update.StopLossPrice != nil {
		setClauses = append(setClauses, fmt.Sprintf("stop_loss_price = $%d", argPosition))
		args = append(args, *update.StopLossPrice)
		argPosition++
	}
	if update.EntryPrice != nil {
		setClauses = append(setClauses, fmt.Sprintf("entry_price = $%d", argPosition))
		args = append(args, *update.EntryPrice)
		argPosition++
	}
	if update.TakeProfitPrice != nil {
		setClauses = append(setClauses, fmt.Sprintf("take_profit_price = $%d", argPosition))
		args = append(args, *update.TakeProfitPrice)
		argPosition++
	}
	if update.Type != nil {
		setClauses = append(setClauses, fmt.Sprintf("type = $%d", argPosition))
		args = append(args, *update.Type)
		argPosition++
	}
	if update.Result != nil {
		setClauses = append(setClauses, fmt.Sprintf("result = $%d", argPosition))
		args = append(args, *update.Result)
		argPosition++
	}
	if update.Return != nil {
		setClauses = append(setClauses, fmt.Sprintf("return = $%d", argPosition))
		args = append(args, *update.Return)
		argPosition++
	}
	if update.FreeForAll != nil {
		setClauses = append(setClauses, fmt.Sprintf("free_for_all = $%d", argPosition))
		args = append(args, *update.FreeForAll)
		argPosition++
	}
	if update.Comments != nil {
		setClauses = append(setClauses, fmt.Sprintf("comments = $%d", argPosition))
		args = append(args, *update.Comments)
		argPosition++
	}

	if len(setClauses) == 0 {
		return r.GetByID(id)
	}

	setClauses = append(setClauses, "updated_at = CURRENT_TIMESTAMP")
	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE trading_signals
		SET %s
		WHERE id = $%d
		RETURNING id, symbol, asset_class, duration_type, stop_loss_price, entry_price, take_profit_price, type, result, return, free_for_all, comments, created_by, created_at, updated_at
	`, strings.Join(setClauses, ", "), argPosition)

	var signal models.TradingSignal
	err := r.db.QueryRow(query, args...).Scan(
		&signal.ID,
		&signal.Symbol,
		&signal.AssetClass,
		&signal.DurationType,
		&signal.StopLossPrice,
		&signal.EntryPrice,
		&signal.TakeProfitPrice,
		&signal.Type,
		&signal.Result,
		&signal.Return,
		&signal.FreeForAll,
		&signal.Comments,
		&signal.CreatedBy,
		&signal.CreatedAt,
		&signal.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("trading signal not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update trading signal: %w", err)
	}

	return &signal, nil
}

// Delete deletes a trading signal
func (r *TradingSignalRepository) Delete(id int64) error {
	query := `DELETE FROM trading_signals WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete trading signal: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("trading signal not found")
	}

	return nil
}

// Count returns the total count of trading signals
func (r *TradingSignalRepository) Count() (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM trading_signals`
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count trading signals: %w", err)
	}
	return count, nil
}

// GetSignalsForUser retrieves trading signals visible to a specific user based on their subscriptions
func (r *TradingSignalRepository) GetSignalsForUser(userID int64, limit, offset int) ([]models.TradingSignal, error) {
	query := `
		SELECT DISTINCT ts.id, ts.symbol, ts.asset_class, ts.duration_type, ts.stop_loss_price, ts.entry_price, 
			ts.take_profit_price, ts.type, ts.result, ts.return, ts.free_for_all, ts.comments, ts.created_by, ts.created_at, ts.updated_at
		FROM trading_signals ts
		WHERE ts.free_for_all = true
		OR EXISTS (
			SELECT 1 FROM user_subscriptions us
			JOIN packages p ON us.package_id = p.id
			WHERE us.user_id = $1
			AND us.is_active = true
			AND us.expires_at > CURRENT_TIMESTAMP
			AND p.asset_class = ts.asset_class
			AND p.duration_type = ts.duration_type
		)
		ORDER BY ts.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get trading signals for user: %w", err)
	}
	defer rows.Close()

	var signals []models.TradingSignal
	for rows.Next() {
		var signal models.TradingSignal
		err := rows.Scan(
			&signal.ID,
			&signal.Symbol,
			&signal.AssetClass,
			&signal.DurationType,
			&signal.StopLossPrice,
			&signal.EntryPrice,
			&signal.TakeProfitPrice,
			&signal.Type,
			&signal.Result,
			&signal.Return,
			&signal.FreeForAll,
			&signal.Comments,
			&signal.CreatedBy,
			&signal.CreatedAt,
			&signal.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan trading signal: %w", err)
		}
		signals = append(signals, signal)
	}

	return signals, nil
}

// CountForUser returns the total count of trading signals visible to a specific user
func (r *TradingSignalRepository) CountForUser(userID int64) (int64, error) {
	var count int64
	query := `
		SELECT COUNT(DISTINCT ts.id)
		FROM trading_signals ts
		WHERE ts.free_for_all = true
		OR EXISTS (
			SELECT 1 FROM user_subscriptions us
			JOIN packages p ON us.package_id = p.id
			WHERE us.user_id = $1
			AND us.is_active = true
			AND us.expires_at > CURRENT_TIMESTAMP
			AND p.asset_class = ts.asset_class
			AND p.duration_type = ts.duration_type
		)
	`
	err := r.db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count trading signals for user: %w", err)
	}
	return count, nil
}

// CheckUserAccess checks if a user has access to a specific signal
func (r *TradingSignalRepository) CheckUserAccess(userID, signalID int64) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM trading_signals ts
			WHERE ts.id = $2
			AND (
				ts.free_for_all = true
				OR EXISTS (
					SELECT 1 FROM user_subscriptions us
					JOIN packages p ON us.package_id = p.id
					WHERE us.user_id = $1
					AND us.is_active = true
					AND us.expires_at > CURRENT_TIMESTAMP
					AND p.asset_class = ts.asset_class
					AND p.duration_type = ts.duration_type
				)
			)
		)
	`
	var hasAccess bool
	err := r.db.QueryRow(query, userID, signalID).Scan(&hasAccess)
	if err != nil {
		return false, fmt.Errorf("failed to check user access: %w", err)
	}
	return hasAccess, nil
}
