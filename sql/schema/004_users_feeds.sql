-- +goose Up
CREATE TABLE feed_follows(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id UUID NOT NULL,
	feed_id UUID NOT NULL,
	created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
	updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
	CONSTRAINT user_feed UNIQUE(user_id, feed_id),
	CONSTRAINT fk_user
	FOREIGN KEY (user_id)
	REFERENCES users(id)
	ON DELETE CASCADE 
	ON UPDATE CASCADE,
	CONSTRAINT fk_feed
	FOREIGN KEY (feed_id)
	REFERENCES feeds(id)
	ON DELETE CASCADE
	ON UPDATE CASCADE
);

-- +goose Down
DROP TABLE feed_follows;
