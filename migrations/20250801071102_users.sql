-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
user_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
username VARCHAR(64) NOT NULL UNIQUE,
password_hash VARCHAR(255) NOT NULL,
created_at timestamp default now() NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
