-- +goose up
ALTER TABLE users
ADD COLUMN api_key VARCHAR(64) NOT NULL DEFAULT ENCODE(SHA256(random()::text::bytea), 'hex');

-- +goose Down
ALTER TABLE users
DROP COLUMN api_key;

