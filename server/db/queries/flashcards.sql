-- name: CreateFlashcard :one
INSERT INTO flashcards (deck_id, front, back, position)
SELECT
    d.id,
    @front,
    @back,
    COALESCE((SELECT max(fc.position) + 1 FROM flashcards fc WHERE fc.deck_id = d.id), 0)
FROM decks d
JOIN folders f ON f.id = d.folder_id
WHERE d.id = @deck_id AND f.owner_id = @owner_id
RETURNING *;

-- name: GetFlashcard :one
SELECT fc.*
FROM flashcards fc
JOIN decks d ON d.id = fc.deck_id
JOIN folders f ON f.id = d.folder_id
WHERE fc.id = @id AND f.owner_id = @owner_id;

-- name: ListFlashcards :many
SELECT fc.*
FROM flashcards fc
JOIN decks d ON d.id = fc.deck_id
JOIN folders f ON f.id = d.folder_id
WHERE fc.deck_id = @deck_id AND f.owner_id = @owner_id
ORDER BY fc.position ASC;

-- name: UpdateFlashcard :one
UPDATE flashcards fc
SET front = @front,
    back = @back,
    position = @position,
    updated_at = now()
FROM decks d
JOIN folders f ON f.id = d.folder_id
WHERE fc.deck_id = d.id AND fc.id = @id AND f.owner_id = @owner_id
RETURNING fc.*;

-- name: DeleteFlashcard :execrows
DELETE FROM flashcards fc
USING decks d, folders f
WHERE fc.deck_id = d.id AND d.folder_id = f.id AND fc.id = @id AND f.owner_id = @owner_id;
