-- +goose Up
CREATE TABLE refresh_tokens (
    id CHAR(36) PRIMARY KEY,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    token TEXT UNIQUE NOT NULL,
    revoked_at DATETIME DEFAULT NULL,
    expires_at DATETIME NOT NULL,
    user_id CHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

-- +goose StatementBegin
CREATE TRIGGER refresh_tokens_updated_at
AFTER UPDATE ON refresh_tokens
BEGIN
    UPDATE refresh_tokens SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose Down
DROP TRIGGER IF EXISTS refresh_tokens_updated_at;
DROP TABLE refresh_tokens;