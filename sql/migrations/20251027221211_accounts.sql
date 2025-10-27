-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS accounts (
    id uuid PRIMARY KEY default gen_random_uuid(),
    nickname text NOT NULL,
    email text NOT NULL UNIQUE,
    hashed_password text NOT NULL,
    inserted_at timestamp WITH TIME ZONE DEFAULT now(),
    updated_at timestamp WITH TIME ZONE DEFAULT now(),
    deleted_at timestamp WITH TIME ZONE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS accounts;
-- +goose StatementEnd
