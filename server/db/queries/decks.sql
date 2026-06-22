-- name: CreateDeck :one
INSERT INTO decks (folder_id, name, description)
SELECT f.id, @name, @description
FROM folders f
WHERE f.id = @folder_id AND f.owner_id = @owner_id
RETURNING *;

-- name: GetDeck :one
SELECT
    d.*,
    (SELECT count(*) FROM flashcards fc WHERE fc.deck_id = d.id)::int AS card_count
FROM decks d
JOIN folders f ON f.id = d.folder_id
WHERE d.id = @id AND f.owner_id = @owner_id;

-- name: ListDecks :many
SELECT
    d.*,
    (SELECT count(*) FROM flashcards fc WHERE fc.deck_id = d.id)::int AS card_count
FROM decks d
JOIN folders f ON f.id = d.folder_id
WHERE d.folder_id = @folder_id AND f.owner_id = @owner_id
ORDER BY d.created_at ASC;

-- name: UpdateDeck :one
UPDATE decks d
SET name = @name,
    description = @description,
    updated_at = now()
FROM folders f
WHERE d.folder_id = f.id AND d.id = @id AND f.owner_id = @owner_id
RETURNING
    d.*,
    (SELECT count(*) FROM flashcards fc WHERE fc.deck_id = d.id)::int AS card_count;

-- name: DeleteDeck :execrows
DELETE FROM decks d
USING folders f
WHERE d.folder_id = f.id AND d.id = @id AND f.owner_id = @owner_id;
