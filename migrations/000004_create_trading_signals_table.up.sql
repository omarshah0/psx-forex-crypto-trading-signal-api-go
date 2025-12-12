CREATE TABLE IF NOT EXISTS trading_signals (
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(50) NOT NULL,
    stop_loss_price DECIMAL(20, 8) NOT NULL,
    entry_price DECIMAL(20, 8) NOT NULL,
    take_profit_price DECIMAL(20, 8) NOT NULL,
    type VARCHAR(10) NOT NULL CHECK (type IN ('LONG', 'SHORT')),
    result VARCHAR(20) CHECK (result IN ('WIN', 'LOSS', 'BREAKEVEN')),
    return DECIMAL(10, 2),
    created_by INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_trading_signals_symbol ON trading_signals(symbol);
CREATE INDEX idx_trading_signals_type ON trading_signals(type);
CREATE INDEX idx_trading_signals_result ON trading_signals(result);
CREATE INDEX idx_trading_signals_created_by ON trading_signals(created_by);
CREATE INDEX idx_trading_signals_created_at ON trading_signals(created_at DESC);

