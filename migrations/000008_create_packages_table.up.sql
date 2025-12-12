CREATE TABLE IF NOT EXISTS packages (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    asset_class VARCHAR(20) NOT NULL CHECK (asset_class IN ('FOREX', 'CRYPTO', 'PSX')),
    duration_type VARCHAR(20) NOT NULL CHECK (duration_type IN ('SHORT_TERM', 'LONG_TERM')),
    billing_cycle VARCHAR(20) NOT NULL CHECK (billing_cycle IN ('MONTHLY', 'SIX_MONTHS', 'YEARLY')),
    duration_days INTEGER NOT NULL CHECK (duration_days > 0),
    price DECIMAL(10, 2) NOT NULL CHECK (price >= 0),
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(asset_class, duration_type, billing_cycle)
);

-- Create indexes
CREATE INDEX idx_packages_asset_class ON packages(asset_class);
CREATE INDEX idx_packages_duration_type ON packages(duration_type);
CREATE INDEX idx_packages_billing_cycle ON packages(billing_cycle);
CREATE INDEX idx_packages_is_active ON packages(is_active);

