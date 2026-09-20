-- name: CreateSession :exec
INSERT INTO sessions (id, account_uuid, created_at, last_seen_at) VALUES ($1, $2, current_date, current_date);

-- name: ListSessions :many
SELECT * FROM sessions WHERE account_uuid = $1;

-- name: RevokeSession :exec
DELETE FROM sessions WHERE id = $1 AND account_uuid = $2;

-- name: TouchSession :exec
UPDATE sessions SET last_seen_at = current_date WHERE id = $1;
