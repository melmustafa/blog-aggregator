-- +goose Up
CREATE TABLE posts(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
	updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
	title TEXT NOT NULL,
	url TEXT UNIQUE NOT NULL,
	description TEXT,
	published_at TIMESTAMPTZ,
	feed_id UUID,
	CONSTRAINT fk_feed
	FOREIGN KEY (feed_id)
	REFERENCES feeds(id)
	ON UPDATE CASCADE
	ON DELETE SET NULL
);

-- +goose Down
DROP TABLE posts;
