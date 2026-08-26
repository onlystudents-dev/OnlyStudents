-- name: ListGuardians :many
SELECT * FROM guardians ORDER BY id;

-- name: GetGuardian :one
SELECT * FROM guardians WHERE id = $1;
