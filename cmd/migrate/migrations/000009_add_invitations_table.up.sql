CREATE TABLE IF NOT EXISTS user_invitations (
    token TEXT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    expiry TIMESTAMP NOT NULL DEFAULT (NOW() + INTERVAL '7 days'),
    created_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_user_invitations_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX idx_user_invitations_user_id ON user_invitations (user_id);

CREATE INDEX idx_user_invitations_expiry ON user_invitations (expiry);