-- name: GetGame :one
SELECT * FROM games
WHERE id = $1 LIMIT 1;

-- name: CreateGame :one
INSERT INTO games (title) VALUES ($1)
RETURNING *;

-- name: ListGames :many
SELECT * FROM games;
