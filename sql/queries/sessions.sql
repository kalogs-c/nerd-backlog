-- name: CreateSession :exec
INSERT INTO sessions (token, account_id, expires_at)
VALUES ($1, $2, $3);

-- name: GetSessionAccountID :one
SELECT account_id FROM sessions
WHERE token = $1
  AND expires_at > now();

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE token = $1;
