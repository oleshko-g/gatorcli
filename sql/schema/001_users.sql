-- +goose Up
CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  name VARCHAR NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS users;