-- name: ListStudents :many
SELECT * FROM students ORDER BY id;

-- name: GetStudent :one
SELECT * FROM students WHERE id = $1;
