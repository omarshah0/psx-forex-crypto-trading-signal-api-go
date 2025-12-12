-- Remove all seeded packages
DELETE FROM packages WHERE asset_class IN ('FOREX', 'CRYPTO', 'PSX');

