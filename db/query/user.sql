-- name: CreateUser :one
INSERT INTO users (
    username,
    full_name,
    cpf_cnpj,
    email,
    hashed_password
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY username
LIMIT $1
OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET 
    hashed_password = $2,
    email = $3,
    is_merchant = $4,
    password_changed_at = $5,
    last_updated = now()
WHERE username = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE username = $1;