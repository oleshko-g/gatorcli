-- +goose Up
CREATE TABLE IF NOT EXISTS feed_follows (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  feed_url VARCHAR REFERENCES feeds (url) ON DELETE CASCADE,
  user_id UUID REFERENCES users (id) ON DELETE CASCADE,

  UNIQUE (user_id, feed_url)
);

-- +goose Down
DROP TABLE IF EXISTS feed_follows;