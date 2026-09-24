-- name: ListTeachers :many
SELECT * FROM teachers ORDER BY id;

-- name: GetTeacher :one
SELECT * FROM teachers WHERE id = $1;

-- name: CreateTeacher :one
INSERT INTO teachers (phone_number,username,birth_first_name,birth_last_name,birth_date,birth_city,birth_country,permament_address,temporary_address,first_name,last_name) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id;
