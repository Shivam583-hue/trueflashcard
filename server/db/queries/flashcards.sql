-- name: CreateFlashcard :one
INSERT INTO flashcards (deck_id, front, back, position)
SELECT
    d.id,
    $3,
    $4,
    COALESCE((SELECT max(fc.position) + 1 FROM flashcards fc WHERE fc.deck_id = d.id), 0)
FROM decks d
JOIN folders f ON f.id = d.folder_id
WHERE d.id = $1 AND f.owner_id = $2
RETURNING *;

-- name: GetFlashcard :one
SELECT fc.*
FROM flashcards fc
JOIN decks d ON d.id = fc.deck_id
JOIN folders f ON f.id = d.folder_id
WHERE fc.id = $1 AND f.owner_id = $2;

-- name: ListFlashcards :many
SELECT fc.*
FROM flashcards fc
JOIN decks d ON d.id = fc.deck_id
JOIN folders f ON f.id = d.folder_id
WHERE fc.deck_id = $1 AND f.owner_id = $2
ORDER BY fc.position ASC;

-- name: UpdateFlashcard :one
UPDATE flashcards fc
SET front = $3,
    back = $4,
    position = $5,
    updated_at = now()
FROM decks d
JOIN folders f ON f.id = d.folder_id
WHERE fc.deck_id = d.id AND fc.id = $1 AND f.owner_id = $2
RETURNING fc.*;

-- name: DeleteFlashcard :exec
DELETE FROM flashcards fc
USING decks d, folders f
WHERE fc.deck_id = d.id AND d.folder_id = f.id AND fc.id = $1 AND f.owner_id = $2;
