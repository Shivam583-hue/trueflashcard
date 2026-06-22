-- name: CreateFolder :one
INSERT INTO folders (owner_id, name)
VALUES ($1, $2)
RETURNING *;

-- name: GetFolder :one
SELECT * FROM folders
WHERE id = $1 AND owner_id = $2;

-- name: ListFolders :many
SELECT * FROM folders
WHERE owner_id = $1
ORDER BY created_at ASC;

-- name: UpdateFolder :one
UPDATE folders
SET name = $3,
    updated_at = now()
WHERE id = $1 AND owner_id = $2
RETURNING *;

-- name: DeleteFolder :exec
DELETE FROM folders
WHERE id = $1 AND owner_id = $2;
