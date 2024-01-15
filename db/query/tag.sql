-- name: CreateTag :one
INSERT INTO tags (
  tag_name
) VALUES (
  $1
) RETURNING *;

-- name: ListTags :many
SELECT * FROM tags;

-- name: DeleteTagFromCacheTagsTable :exec
DELETE FROM cache_tags 
WHERE tag_id = $1;

-- name: DeleteTagFromTagsTable :exec
DELETE FROM tags
WHERE tag_id = $1;

-- name: AddTagToCache :one
INSERT INTO cache_tags (
  cache_id,
  tag_id
) VALUES (
  $1, $2
) RETURNING *;

-- name: SubscribeTag :one
INSERT INTO user_tags (
  user_id,
  tag_id
) VALUES (
  $1, $2
) RETURNING *;

-- name: UnsubscribeTag :exec
DELETE FROM user_tags
WHERE user_id = $1 AND tag_id = $2;

-- name: ListUserSubscriptions :many
SELECT t.*
FROM user_tags ut
JOIN tags t ON ut.tag_id = t.tag_id
WHERE ut.user_id =$1;

-- name: ListCacheTags :many
SELECT t.*
FROM cache_tags ct
JOIN tags t ON ct.tag_id = t.tag_id
WHERE ct.cache_id =$1;

