-- Add password field (nullable since OAuth users won't have passwords)
ALTER TABLE users ADD COLUMN hashed_password VARCHAR(255);

-- Add email verification fields
ALTER TABLE users ADD COLUMN email_verified BOOLEAN DEFAULT FALSE;
ALTER TABLE users ADD COLUMN verification_token VARCHAR(255);
ALTER TABLE users ADD COLUMN verification_token_expires TIMESTAMP;

-- Add password reset fields
ALTER TABLE users ADD COLUMN reset_token VARCHAR(255);
ALTER TABLE users ADD COLUMN reset_token_expires TIMESTAMP;

-- Indexes for tokens (for faster lookups)
CREATE INDEX idx_users_verification_token ON users(verification_token);
CREATE INDEX idx_users_reset_token ON users(reset_token);
CREATE INDEX idx_users_email_verified ON users(email_verified);

-- OAuth users are auto-verified
UPDATE users SET email_verified = TRUE WHERE hashed_password IS NULL;

