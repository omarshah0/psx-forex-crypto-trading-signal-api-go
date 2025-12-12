-- Add new columns to trading_signals table
ALTER TABLE trading_signals
ADD COLUMN asset_class VARCHAR(20) CHECK (asset_class IN ('FOREX', 'CRYPTO', 'PSX')),
ADD COLUMN duration_type VARCHAR(20) CHECK (duration_type IN ('SHORT_TERM', 'LONG_TERM')),
ADD COLUMN free_for_all BOOLEAN DEFAULT false,
ADD COLUMN comments TEXT;

-- Create indexes for the new columns
CREATE INDEX idx_trading_signals_asset_class ON trading_signals(asset_class);
CREATE INDEX idx_trading_signals_duration_type ON trading_signals(duration_type);
CREATE INDEX idx_trading_signals_free_for_all ON trading_signals(free_for_all);

-- Make the new columns required (NOT NULL) after adding them
-- We set them nullable first to allow the ALTER TABLE to succeed on existing rows
-- If you have existing data, you may want to update it first before making them NOT NULL

