-- Remove indexes
DROP INDEX IF EXISTS idx_trading_signals_free_for_all;
DROP INDEX IF EXISTS idx_trading_signals_duration_type;
DROP INDEX IF EXISTS idx_trading_signals_asset_class;

-- Remove new columns
ALTER TABLE trading_signals
DROP COLUMN IF EXISTS comments,
DROP COLUMN IF EXISTS free_for_all,
DROP COLUMN IF EXISTS duration_type,
DROP COLUMN IF EXISTS asset_class;

