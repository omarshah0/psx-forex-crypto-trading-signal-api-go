-- Seed initial packages: 3 asset classes × 2 durations × 3 billing cycles = 18 packages

-- FOREX SHORT_TERM packages
INSERT INTO packages (name, asset_class, duration_type, billing_cycle, duration_days, price, description) VALUES
('Forex Short Term - Monthly', 'FOREX', 'SHORT_TERM', 'MONTHLY', 30, 10.00, 'Access to Forex day trading signals for 1 month'),
('Forex Short Term - 6 Months', 'FOREX', 'SHORT_TERM', 'SIX_MONTHS', 180, 50.00, 'Access to Forex day trading signals for 6 months'),
('Forex Short Term - Yearly', 'FOREX', 'SHORT_TERM', 'YEARLY', 365, 80.00, 'Access to Forex day trading signals for 1 year');

-- FOREX LONG_TERM packages
INSERT INTO packages (name, asset_class, duration_type, billing_cycle, duration_days, price, description) VALUES
('Forex Long Term - Monthly', 'FOREX', 'LONG_TERM', 'MONTHLY', 30, 15.00, 'Access to Forex swing trading signals for 1 month'),
('Forex Long Term - 6 Months', 'FOREX', 'LONG_TERM', 'SIX_MONTHS', 180, 75.00, 'Access to Forex swing trading signals for 6 months'),
('Forex Long Term - Yearly', 'FOREX', 'LONG_TERM', 'YEARLY', 365, 120.00, 'Access to Forex swing trading signals for 1 year');

-- CRYPTO SHORT_TERM packages
INSERT INTO packages (name, asset_class, duration_type, billing_cycle, duration_days, price, description) VALUES
('Crypto Short Term - Monthly', 'CRYPTO', 'SHORT_TERM', 'MONTHLY', 30, 8.00, 'Access to Cryptocurrency day trading signals for 1 month'),
('Crypto Short Term - 6 Months', 'CRYPTO', 'SHORT_TERM', 'SIX_MONTHS', 180, 40.00, 'Access to Cryptocurrency day trading signals for 6 months'),
('Crypto Short Term - Yearly', 'CRYPTO', 'SHORT_TERM', 'YEARLY', 365, 65.00, 'Access to Cryptocurrency day trading signals for 1 year');

-- CRYPTO LONG_TERM packages
INSERT INTO packages (name, asset_class, duration_type, billing_cycle, duration_days, price, description) VALUES
('Crypto Long Term - Monthly', 'CRYPTO', 'LONG_TERM', 'MONTHLY', 30, 12.00, 'Access to Cryptocurrency swing trading signals for 1 month'),
('Crypto Long Term - 6 Months', 'CRYPTO', 'LONG_TERM', 'SIX_MONTHS', 180, 60.00, 'Access to Cryptocurrency swing trading signals for 6 months'),
('Crypto Long Term - Yearly', 'CRYPTO', 'LONG_TERM', 'YEARLY', 365, 95.00, 'Access to Cryptocurrency swing trading signals for 1 year');

-- PSX SHORT_TERM packages
INSERT INTO packages (name, asset_class, duration_type, billing_cycle, duration_days, price, description) VALUES
('PSX Short Term - Monthly', 'PSX', 'SHORT_TERM', 'MONTHLY', 30, 5.00, 'Access to Pakistan Stock Exchange day trading signals for 1 month'),
('PSX Short Term - 6 Months', 'PSX', 'SHORT_TERM', 'SIX_MONTHS', 180, 25.00, 'Access to Pakistan Stock Exchange day trading signals for 6 months'),
('PSX Short Term - Yearly', 'PSX', 'SHORT_TERM', 'YEARLY', 365, 40.00, 'Access to Pakistan Stock Exchange day trading signals for 1 year');

-- PSX LONG_TERM packages
INSERT INTO packages (name, asset_class, duration_type, billing_cycle, duration_days, price, description) VALUES
('PSX Long Term - Monthly', 'PSX', 'LONG_TERM', 'MONTHLY', 30, 10.00, 'Access to Pakistan Stock Exchange swing trading signals for 1 month'),
('PSX Long Term - 6 Months', 'PSX', 'LONG_TERM', 'SIX_MONTHS', 180, 50.00, 'Access to Pakistan Stock Exchange swing trading signals for 6 months'),
('PSX Long Term - Yearly', 'PSX', 'LONG_TERM', 'YEARLY', 365, 80.00, 'Access to Pakistan Stock Exchange swing trading signals for 1 year');

