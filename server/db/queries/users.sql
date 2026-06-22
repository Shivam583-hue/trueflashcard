-- name: CreateUser :one
INSERT INTO users (google_subject, email, display_name)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByGoogleSubject :one
SELECT * FROM users
WHERE google_subject = $1;

-- name: UpsertUserByGoogleSubject :one
INSERT INTO users (google_subject, email, display_name)
VALUES (@google_subject, @email, @display_name)
ON CONFLICT (google_subject)
DO UPDATE SET
    email = EXCLUDED.email,
    display_name = EXCLUDED.display_name,
    updated_at = now()
RETURNING *;

-- name: UpdateUserProfile :one
UPDATE users
SET email = $2,
    display_name = $3,
    updated_at = now()
WHERE id = $1
RETURNING *;
