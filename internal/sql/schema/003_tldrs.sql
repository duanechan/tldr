-- +goose Up
CREATE TABLE tldrs (
    id CHAR(36) PRIMARY KEY,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    user_id CHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

-- +goose StatementBegin
CREATE TRIGGER tldrs_updated_at
AFTER UPDATE ON tldrs
BEGIN
    UPDATE tldrs SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose Down
DROP TRIGGER IF EXISTS tldrs_updated_at;
DROP TABLE tldrs;