-- name: GetTeacher :one
SELECT * FROM teachers WHERE id = $1;

-- name: CreateTeacher :one
INSERT INTO teachers (phone_number,birth_first_name,birth_last_name,birth_date,birth_city,birth_country,permament_address,temporary_address,first_name,last_name) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id;

-- name: ListTeacherSchoolMemberships :many
SELECT *, school.name AS school_name FROM teacher_school ts
JOIN schools school ON school.id = ts.school_id
WHERE teacher_id = $1;