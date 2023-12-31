-- name: CreateUser :one
INSERT INTO users (
  id,
  email,
  name
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET 
  name = COALESCE(sqlc.narg(name), name)
WHERE
  id = sqlc.narg(id)
RETURNING *;
