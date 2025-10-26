-- name: GetGame :one
SELECT * FROM games
WHERE id = $1;

-- name: CreateGame :one
INSERT INTO games (title) VALUES ($1)
RETURNING *;

-- name: ListGames :many
SELECT * FROM games;

-- name: DeleteGameByID :exec
DELETE FROM games
WHERE id = $1;
