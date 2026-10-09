-- name: GetAdminDebugData :one
SELECT
  (SELECT COUNT(*) FROM accounts) AS account_count,
  (SELECT COUNT(*) FROM schools) AS school_count,
  (SELECT COUNT(*) FROM students) AS student_count,
  (SELECT COUNT(*) FROM teachers) AS teacher_count,
  (SELECT COUNT(*) FROM guardians) AS guardian_count;

-- name: AdminListSchools :many
SELECT name, city, address_line, zip_code, id FROM schools;

-- name: AdminListClasses :many
SELECT name, id FROM classes WHERE school_id = $1;

-- name: AdminCheckSchoolExists :execrows
SELECT 1 FROM schools WHERE id = $1;

-- name: AdminCheckClassExists :execrows
SELECT 1 FROM classes WHERE id = $1 AND school_id = $2;
