-- name: GetMessages :many
SELECT * FROM messages ORDER BY created_at ASC;