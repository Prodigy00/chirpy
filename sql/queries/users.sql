-- name: CreateUser :one
INSERT INTO  users(id, email,created_at, updated_at)
VALUES (
        gen_random_uuid(),
        $1,
        NOW(),
        NOW()
)
RETURNING *;