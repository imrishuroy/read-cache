-- name: CreateCache :one
INSERT INTO caches (
    title, 
    link
) VALUES (
  $1, $2
) RETURNING *;
   
-- name: GetCache :one
SELECT * FROM caches
WHERE id = $1 LIMIT 1;

-- name: ListCaches :many
SELECT * FROM caches
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateCache :one
UPDATE caches
SET title = $2,
    link = $3
WHERE id = $1
RETURNING *;

-- name: DeleteCache :exec
DELETE FROM caches
WHERE id = $1;