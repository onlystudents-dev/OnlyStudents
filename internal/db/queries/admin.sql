-- name: GetAdminDebugData :one
SELECT
  (SELECT COUNT(*) FROM accounts) AS account_count,
  (SELECT COUNT(*) FROM schools) AS school_count,
  (SELECT COUNT(*) FROM students) AS student_count,
  (SELECT COUNT(*) FROM teachers) AS teacher_count,
  (SELECT COUNT(*) FROM guardians) AS guardian_count;

-- name: AdminListSchools :many
SELECT name, city, address_line, zip_code, id FROM schools;
