-- name: CreateAccount :one
INSERT INTO accounts (nickname, email, hashed_password)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetAccountByEmail :one
SELECT * FROM accounts
WHERE email = $1;
