-- name: CreateCache :one
INSERT INTO caches (
  owner,
  title, 
  link
) VALUES (
  $1, $2, $3
) RETURNING *;
   
-- name: GetCache :one
SELECT * FROM caches
WHERE id = $1 LIMIT 1;

-- name: ListCaches :many
SELECT * FROM caches
WHERE owner =$1
-- ORDER BY id
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;

-- name: UpdateCache :one
UPDATE caches
SET title = $2,
    link = $3
WHERE id = $1
RETURNING *;

-- name: DeleteCache :exec
DELETE FROM caches
WHERE id = $1;