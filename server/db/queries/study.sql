-- name: ListDueCards :many
SELECT fc.*
FROM card_review_states crs
JOIN flashcards fc ON fc.id = crs.card_id
WHERE crs.owner_id = @owner_id AND crs.due_at <= @as_of
ORDER BY crs.due_at ASC
LIMIT @lim;

-- name: ListDueCardsInDeck :many
SELECT fc.*
FROM card_review_states crs
JOIN flashcards fc ON fc.id = crs.card_id
WHERE crs.owner_id = @owner_id AND crs.due_at <= @as_of AND fc.deck_id = @deck_id
ORDER BY crs.due_at ASC
LIMIT @lim;

-- name: ListNewCards :many
SELECT fc.*
FROM flashcards fc
JOIN decks d ON d.id = fc.deck_id
JOIN folders f ON f.id = d.folder_id
LEFT JOIN card_review_states crs ON crs.card_id = fc.id
WHERE f.owner_id = @owner_id AND crs.card_id IS NULL
ORDER BY fc.created_at ASC, fc.position ASC
LIMIT @lim;

-- name: ListNewCardsInDeck :many
SELECT fc.*
FROM flashcards fc
JOIN decks d ON d.id = fc.deck_id
JOIN folders f ON f.id = d.folder_id
LEFT JOIN card_review_states crs ON crs.card_id = fc.id
WHERE f.owner_id = @owner_id AND crs.card_id IS NULL AND fc.deck_id = @deck_id
ORDER BY fc.position ASC
LIMIT @lim;

-- name: GetCardReviewState :one
SELECT *
FROM card_review_states
WHERE card_id = @card_id AND owner_id = @owner_id;

-- name: UpsertCardReviewState :one
INSERT INTO card_review_states (
    card_id, owner_id, state, due_at, stability, difficulty, reps, lapses, last_reviewed_at
) VALUES (
    @card_id, @owner_id, @state, @due_at, @stability, @difficulty, @reps, @lapses, @last_reviewed_at
)
ON CONFLICT (card_id) DO UPDATE SET
    state = EXCLUDED.state,
    due_at = EXCLUDED.due_at,
    stability = EXCLUDED.stability,
    difficulty = EXCLUDED.difficulty,
    reps = EXCLUDED.reps,
    lapses = EXCLUDED.lapses,
    last_reviewed_at = EXCLUDED.last_reviewed_at,
    updated_at = now()
RETURNING *;

-- name: InsertCardReview :exec
INSERT INTO card_reviews (
    card_id, owner_id, rating, state_before, elapsed_days, stability, difficulty, scheduled_days
) VALUES (
    @card_id, @owner_id, @rating, @state_before, @elapsed_days, @stability, @difficulty, @scheduled_days
);

-- name: CountDueTotal :one
SELECT count(*)::int
FROM card_review_states
WHERE owner_id = @owner_id AND due_at <= @as_of;

-- name: CountNewTotal :one
SELECT count(*)::int
FROM flashcards fc
JOIN decks d ON d.id = fc.deck_id
JOIN folders f ON f.id = d.folder_id
LEFT JOIN card_review_states crs ON crs.card_id = fc.id
WHERE f.owner_id = @owner_id AND crs.card_id IS NULL;

-- name: CountDueByDeck :many
SELECT fc.deck_id, count(*)::int AS due
FROM card_review_states crs
JOIN flashcards fc ON fc.id = crs.card_id
WHERE crs.owner_id = @owner_id AND crs.due_at <= @as_of
GROUP BY fc.deck_id;

-- name: CountNewByDeck :many
SELECT fc.deck_id, count(*)::int AS new
FROM flashcards fc
JOIN decks d ON d.id = fc.deck_id
JOIN folders f ON f.id = d.folder_id
LEFT JOIN card_review_states crs ON crs.card_id = fc.id
WHERE f.owner_id = @owner_id AND crs.card_id IS NULL
GROUP BY fc.deck_id;

-- name: CountReviewsSince :one
SELECT count(*)::int
FROM card_reviews
WHERE owner_id = @owner_id AND reviewed_at >= @since;

-- name: ListReviewDays :many
SELECT DISTINCT (reviewed_at AT TIME ZONE 'UTC')::date AS day
FROM card_reviews
WHERE owner_id = @owner_id
ORDER BY day DESC
LIMIT 400;
