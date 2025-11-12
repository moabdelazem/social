ALTER TABLE users
ADD COLUMN IF NOT EXISTS is_active BOOLEAN NOT NULL DEFAULT false;

-- Create an index for faster queries on active users
CREATE INDEX idx_users_is_active ON users (is_active);