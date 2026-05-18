-- +goose Up
ALTER TABLE tldrs
ADD COLUMN flag TEXT NOT NULL DEFAULT 'Mild';

-- +goose Down
ALTER TABLE tldrs
DROP COLUMN IF EXISTS flag;
