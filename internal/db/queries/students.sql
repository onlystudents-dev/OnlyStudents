-- name: ListStudents :many
SELECT * FROM students ORDER BY id;

-- name: GetStudent :one
SELECT * FROM students WHERE id = $1;

-- name: CreateStudent :one
INSERT INTO students (id_number,school_id,phone_number,first_name,last_name,birth_first_name,birth_last_name,birth_date,birth_city,birth_country,mother_birth_first_name,mother_birth_last_name,classes_id,permament_address,temporary_address,tax_number,ssn_number,bank_name,iban_owner,iban_number,document_type,document_number) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22) RETURNING id;
