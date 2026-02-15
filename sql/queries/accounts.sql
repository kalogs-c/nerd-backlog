-- name: CreateAccount :one
INSERT INTO accounts (nickname, email, hashed_password)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetAccountByEmail :one
SELECT * FROM accounts
WHERE email = $1;

-- name: StoreRefreshToken :exec
INSERT INTO refresh_tokens (token, account_id, expires_at)
VALUES ($1, $2, $3)
ON CONFLICT (account_id) DO UPDATE SET token = EXCLUDED.token, expires_at = EXCLUDED.expires_at;

-- name: DeleteRefreshTokenByAccountID :exec
DELETE FROM refresh_tokens
WHERE account_id = $1;
