-- name: ListGuardians :many
SELECT * FROM guardians ORDER BY id;

-- name: GetGuardian :one
SELECT * FROM guardians WHERE id = $1;

-- name: CreateGuardian :one
INSERT INTO guardians (phone_number,first_name,last_name,birth_first_name,birth_last_name,birth_date,birth_city,birth_country,permament_address,temporary_address) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id;
