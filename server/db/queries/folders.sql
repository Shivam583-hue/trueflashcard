-- name: CreateFolder :one
INSERT INTO folders (owner_id, name)
VALUES (@owner_id, @name)
RETURNING *;

-- name: GetFolder :one
SELECT * FROM folders
WHERE id = @id AND owner_id = @owner_id;

-- name: ListFolders :many
SELECT * FROM folders
WHERE owner_id = @owner_id
ORDER BY created_at ASC;

-- name: UpdateFolder :one
UPDATE folders
SET name = @name,
    updated_at = now()
WHERE id = @id AND owner_id = @owner_id
RETURNING *;

-- name: DeleteFolder :execrows
DELETE FROM folders
WHERE id = @id AND owner_id = @owner_id;
