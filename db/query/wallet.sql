-- name: CreateWallet :one
INSERT INTO wallets (
    owner,
    balance,
    currency
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetWallet :one
SELECT * FROM wallets
WHERE id = $1 LIMIT 1;

-- name: GetWalletForUpdate :one
SELECT * FROM wallets
WHERE id = $1 
LIMIT 1 
FOR NO KEY UPDATE;

-- name: ListWallets :many
SELECT * FROM wallets
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateWallet :one
UPDATE wallets
SET balance = $2
WHERE id = $1
RETURNING *;

-- name: AddWalletBalance :one
UPDATE wallets
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteWallet :exec
DELETE FROM wallets WHERE id = $1;