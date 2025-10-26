-- +goose Up
-- +goose StatementBegin
CREATE TABLE games (
    id uuid PRIMARY KEY default gen_random_uuid(),
    title text NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE games;
-- +goose StatementEnd
