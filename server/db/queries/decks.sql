-- name: CreateDeck :one
INSERT INTO decks (folder_id, name, description)
SELECT f.id, $3, $4
FROM folders f
WHERE f.id = $1 AND f.owner_id = $2
RETURNING *;

-- name: GetDeck :one
SELECT
    d.*,
    (SELECT count(*) FROM flashcards fc WHERE fc.deck_id = d.id)::int AS card_count
FROM decks d
JOIN folders f ON f.id = d.folder_id
WHERE d.id = $1 AND f.owner_id = $2;

-- name: ListDecks :many
SELECT
    d.*,
    (SELECT count(*) FROM flashcards fc WHERE fc.deck_id = d.id)::int AS card_count
FROM decks d
JOIN folders f ON f.id = d.folder_id
WHERE d.folder_id = $1 AND f.owner_id = $2
ORDER BY d.created_at ASC;

-- name: UpdateDeck :one
UPDATE decks d
SET name = $3,
    description = $4,
    updated_at = now()
FROM folders f
WHERE d.folder_id = f.id AND d.id = $1 AND f.owner_id = $2
RETURNING d.*;

-- name: DeleteDeck :exec
DELETE FROM decks d
USING folders f
WHERE d.folder_id = f.id AND d.id = $1 AND f.owner_id = $2;
