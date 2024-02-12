-- name: CreateCache :one
INSERT INTO caches (
  owner,
  title, 
  content,
  is_public
) VALUES (
  $1, $2, $3, $4
) RETURNING *;
   
-- name: GetCache :one
SELECT * FROM caches
WHERE id = $1 LIMIT 1;

-- name: ListCaches :many
SELECT * FROM caches
WHERE owner =$1
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;

-- name: UpdateCache :one
UPDATE caches
SET title = $2,
    content = $3,
    is_public = $4
WHERE id = $1
RETURNING *;

-- name: DeleteCache :exec
DELETE FROM caches
WHERE id = $1;

-- name: ListPublicCaches :many
SELECT c.*
FROM caches c
WHERE c.is_public = TRUE
LIMIT $1
OFFSET $2;

-- name: ListPublicCachesByTags :many
SELECT c.*
FROM caches c
JOIN cache_tags ct ON c.id = ct.cache_id
JOIN tags t ON ct.tag_id = t.tag_id
WHERE c.is_public = TRUE
AND t.tag_id = ANY(sqlc.arg(tag_ids)::int[])
LIMIT $1
OFFSET $2;









