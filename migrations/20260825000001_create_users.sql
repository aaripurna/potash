-- +goose Up
CREATE TABLE users (
  id UUID PRIMARY KEY,
  email TEXT NOT NULL,
  name TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX users_email_unique ON users (lower(email));

-- +goose Down
DROP TABLE users;
