-- +goose Up
CREATE TABLE IF NOT EXISTS posts (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title VARCHAR NOT NULL,
    url VARCHAR NOT NULL UNIQUE,
    description VARCHAR NOT NULL,
    published_at TIMESTAMP NOT NULL,
    feed_id UUID REFERENCES feeds (id) ON DELETE CASCADE
);

-- +goose Down
drop table if exists posts;