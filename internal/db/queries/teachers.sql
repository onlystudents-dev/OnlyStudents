-- name: ListTeachers :many
SELECT * FROM teachers ORDER BY id;

-- name: GetTeacher :one
SELECT * FROM teachers WHERE id = $1;
