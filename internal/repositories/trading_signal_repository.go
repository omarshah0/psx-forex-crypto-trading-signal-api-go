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
		INSERT INTO trading_signals (symbol, stop_loss_price, entry_price, take_profit_price, type, result, return, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, symbol, stop_loss_price, entry_price, take_profit_price, type, result, return, created_by, created_at, updated_at
	`

	var newSignal models.TradingSignal
	err := r.db.QueryRow(
		query,
		signal.Symbol,
		signal.StopLossPrice,
		signal.EntryPrice,
		signal.TakeProfitPrice,
		signal.Type,
		signal.Result,
		signal.Return,
		createdBy,
	).Scan(
		&newSignal.ID,
		&newSignal.Symbol,
		&newSignal.StopLossPrice,
		&newSignal.EntryPrice,
		&newSignal.TakeProfitPrice,
		&newSignal.Type,
		&newSignal.Result,
		&newSignal.Return,
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
		SELECT id, symbol, stop_loss_price, entry_price, take_profit_price, type, result, return, created_by, created_at, updated_at
		FROM trading_signals
		WHERE id = $1
	`

	var signal models.TradingSignal
	err := r.db.QueryRow(query, id).Scan(
		&signal.ID,
		&signal.Symbol,
		&signal.StopLossPrice,
		&signal.EntryPrice,
		&signal.TakeProfitPrice,
		&signal.Type,
		&signal.Result,
		&signal.Return,
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
		SELECT id, symbol, stop_loss_price, entry_price, take_profit_price, type, result, return, created_by, created_at, updated_at
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
			&signal.StopLossPrice,
			&signal.EntryPrice,
			&signal.TakeProfitPrice,
			&signal.Type,
			&signal.Result,
			&signal.Return,
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

	if len(setClauses) == 0 {
		return r.GetByID(id)
	}

	setClauses = append(setClauses, "updated_at = CURRENT_TIMESTAMP")
	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE trading_signals
		SET %s
		WHERE id = $%d
		RETURNING id, symbol, stop_loss_price, entry_price, take_profit_price, type, result, return, created_by, created_at, updated_at
	`, strings.Join(setClauses, ", "), argPosition)

	var signal models.TradingSignal
	err := r.db.QueryRow(query, args...).Scan(
		&signal.ID,
		&signal.Symbol,
		&signal.StopLossPrice,
		&signal.EntryPrice,
		&signal.TakeProfitPrice,
		&signal.Type,
		&signal.Result,
		&signal.Return,
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
