-- name: ListSchools :many
SELECT * FROM schools ORDER BY id;

-- name: GetSchool :one
SELECT * FROM schools WHERE id = $1;
