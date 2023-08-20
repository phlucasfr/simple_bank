-- name: CreateUser :one
INSERT INTO users (
    full_name,
    cpf_cnpj,
    email,
    password
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateUserPassword :one
UPDATE users
SET password = $2, last_updated = now()
WHERE id = $1
RETURNING *;

-- name: UpdateUserEmail :one
UPDATE users
SET email = $2, last_updated = now()
WHERE id = $1
RETURNING *;

-- name: UpdateUserIsMerchant :one
UPDATE users
SET is_merchant = CASE WHEN is_merchant = true THEN false ELSE true END, last_updated = now()
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;